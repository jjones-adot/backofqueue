package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var ExpiryTimeStampFormatted = time.Now().UTC().Add(-time.Hour * 1).Format("2006-01-02 15:04:05")

//Inrix Back of Queue Struct Definition

type BackofQueue struct {
	Doctype       string    `json:"docType"`
	Copyright     string    `json:"copyright"`
	Versionnumber string    `json:"versionNumber"`
	Createddate   time.Time `json:"createdDate"`
	Statusid      int       `json:"statusId"`
	Statustext    string    `json:"statusText"`
	Responseid    string    `json:"responseId"`
	Result        struct {
		Xdincidents []struct {
			ID       int    `json:"id"`
			Version  string `json:"version"`
			Type     string `json:"type"`
			Severity string `json:"severity"`
			Geometry struct {
				Type        string   `json:"type"`
				Coordinates []string `json:"coordinates"`
			} `json:"geometry"`
			Impacting string `json:"impacting"`
			Status    string `json:"status"`
			Messages  struct {
				Alertcmessagecodes []struct {
					Eventcode string `json:"eventCode"`
					Level     string `json:"level"`
				} `json:"alertCMessageCodes"`
				Inrixmessage []struct {
					Inrixcode      string `json:"inrixCode"`
					Type           string `json:"type"`
					Quantifierdata string `json:"quantifierData"`
					Quantifiertype string `json:"quantifierType"`
				} `json:"inrixMessage"`
			} `json:"messages"`
			Location struct {
				Countrycode   string `json:"countryCode"`
				Direction     string `json:"direction"`
				Bidirectional string `json:"biDirectional"`
				Segments      []struct {
					Type   string `json:"type"`
					Offset string `json:"offset"`
					Code   string `json:"code"`
				} `json:"segments"`
			} `json:"location"`
			Schedule struct {
				Planned             string    `json:"planned"`
				Advancewarning      string    `json:"advanceWarning"`
				Occurrencestarttime time.Time `json:"occurrenceStartTime"`
				Occurrenceendtime   time.Time `json:"occurrenceEndTime"`
				Descriptions        struct {
					Lang string `json:"lang"`
					Desc string `json:"desc"`
				} `json:"descriptions"`
			} `json:"schedule"`
			Descriptions []struct {
				Type string `json:"type"`
				Lang string `json:"lang"`
				Desc string `json:"desc"`
			} `json:"descriptions"`
			Parameterizeddescription struct {
				Eventcode  string `json:"eventCode"`
				Eventtext  string `json:"eventText"`
				Roadname   string `json:"roadName"`
				Direction  string `json:"direction"`
				Crossroad1 string `json:"crossroad1"`
				Crossroad2 string `json:"crossroad2"`
				Position1  string `json:"position1"`
				Position2  string `json:"position2"`
			} `json:"parameterizedDescription,omitempty"`
			Head struct {
				Geometry struct {
					Type        string   `json:"type"`
					Coordinates []string `json:"coordinates"`
				} `json:"geometry"`
			} `json:"head"`
			Tail []struct {
				Geometry struct {
					Type        string   `json:"type"`
					Coordinates []string `json:"coordinates"`
				} `json:"geometry"`
			} `json:"tail"`
			Lastdetourpoints []struct {
				Geometry struct {
					Type        string   `json:"type"`
					Coordinates []string `json:"coordinates"`
				} `json:"geometry"`
			} `json:"lastDetourPoints"`
			Dlrs struct {
				Type     string `json:"type"`
				Segments []struct {
					ID     string `json:"id"`
					Offset string `json:"offset"`
				} `json:"segments"`
			} `json:"dlrs"`
			Rds struct {
				Alertcmessage         string `json:"alertcMessage"`
				Direction             string `json:"direction"`
				Extent                string `json:"extent"`
				Duration              string `json:"duration"`
				Diversion             string `json:"diversion"`
				Directionalitychanged string `json:"directionalityChanged"`
				Eventcode             []struct {
					Code    string `json:"code"`
					Primary string `json:"primary"`
				} `json:"eventCode"`
			} `json:"rds"`
			Delayimpact struct {
				Fromtypicalminutes  string `json:"fromTypicalMinutes"`
				Fromfreeflowminutes string `json:"fromFreeFlowMinutes"`
				Fromnas             string `json:"fromNas"`
				Distance            string `json:"distance"`
				Abnormal            string `json:"abnormal"`
			} `json:"delayImpact"`
		} `json:"XDIncidents"`
	} `json:"result"`
}



// Inrix XML Token Struct Definition

type Inrix struct {
	XMLName       xml.Name `xml:"Inrix"`
	Text          string   `xml:",chardata"`
	DocType       string   `xml:"docType,attr"`
	Copyright     string   `xml:"copyright,attr"`
	VersionNumber string   `xml:"versionNumber,attr"`
	CreatedDate   string   `xml:"createdDate,attr"`
	StatusId      string   `xml:"statusId,attr"`
	StatusText    string   `xml:"statusText,attr"`
	ResponseId    string   `xml:"responseId,attr"`
	AuthResponse  struct {
		Text      string `xml:",chardata"`
		AuthToken struct {
			Text   string `xml:",chardata"`
			Expiry string `xml:"expiry,attr"`
		} `xml:"AuthToken"`
		ServerPath  string `xml:"ServerPath"`
		ServerPaths struct {
			Text       string `xml:",chardata"`
			ServerPath []struct {
				Text   string `xml:",chardata"`
				Type   string `xml:"type,attr"`
				Region string `xml:"region,attr"`
			} `xml:"ServerPath"`
		} `xml:"ServerPaths"`
	} `xml:"AuthResponse"`
}


// Token Refresh Call

func current_token() string {

	resp, err := http.Get("http://na.api.inrix.com/Traffic/Inrix.ashx?Action=GetSecurityToken&vendorId=1922770974&consumerId=3efa617e-0310-4d97-81c6-ec8cd5d37b06")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	var Response Inrix
	if err := xml.Unmarshal(body, &Response); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	token := (Response.AuthResponse.AuthToken.Text)
	ExpirytimeString := Response.AuthResponse.AuthToken.Expiry
	ExpiryTimeStamp, err := time.Parse(time.RFC3339, ExpirytimeString)
	ExpiryTimeStampFormatted = ExpiryTimeStamp.Format("2006-01-02 15:04:05")
	
	return token

}

func main() {

	var tokenc string

	for true {
		CurrentUTCStampFormatted := time.Now().UTC().Format("2006-01-02 15:04:05")
		
		if ExpiryTimeStampFormatted < CurrentUTCStampFormatted {

			tokenc = current_token()
			
		} else {
			fmt.Println("0")
		}
		resp, err := http.Get("https://na-api.inrix.com/Traffic/Inrix.ashx?action=GetXDIncidentsInBox&units=1&locale=en-US&corner1=33.925066|-112.803268&corner2=32.743371|-111.178080&incidentSource=INRIXonly&incidentType=Flow&geometryTolerance=10&incidentoutputfields=all&locrefmethod=XD&token=" + tokenc + "&compress=true&format=json")

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		msg := BackofQueue{}
		jsonErr := json.Unmarshal(body, &msg)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		for _, s := range msg.Result.Xdincidents {

			for _, a := range s.Descriptions {
				fmt.Println(a.Desc + "\n")
			}
		}

		time.Sleep(time.Second * 61)

	}
}
