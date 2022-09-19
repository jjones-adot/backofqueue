import geocoding1 as gc
#from pprint import pprint as pp
import urllib, json
import json
import sys
import math

import ssl
ssl._create_default_https_context = ssl._create_unverified_context

import urllib.parse
urlEncode = urllib.parse.urlencode # 3.7 encoder

import urllib.request
urlOpen = urllib.request.urlopen # 3.7 http request
    
geometry_service = "https://gis.dot.state.az/rhgeocode/rest/services/Utilities/Geometry/GeometryServer"
def boq_war_loc(x,y):
    slowdown_point = (-111.0853171, 33.3038192)
    #slowdown_point = (x,y)
    centerline_offset = 15  ## units in meters

    rg_params = gc.ReverseGeocodeParams(x=slowdown_point[0], y=slowdown_point[1], inSR=4326, snappedRouteCandidateSubtypes=[70,71,80,81,90,91])   ## constrain snapping behavior to snap only on State Highway mainlines - exclude Ramps, Frontages, Local Routes
    rg_result = gc.reverseGeocode(rg_params)

    ##pp(rg_result)
    snapped_routeId     = rg_result['locations'][0]['geocode_onRoad']['geocode_onRoad_ATIS']       ## the RouteId in the Roads and Highways Database of the snapped route
    snapped_cardinality = rg_result['locations'][0]['geocode_onRoad']['geocode_onRoadCard']        ## the cardinality of the route
    snapped_measure     = rg_result['locations'][0]['geocode_routePoint']['geocode_XY_WM']['m']    ## the Roads and Highways Database managed measure value of the route at the snapped location

                     
    for i, demarcation_distance in enumerate([1,2,3]):
        directional_multiplier = 1 if snapped_cardinality == 'N' else -1
                          
        params = gc.GeocodeParams(routeId=snapped_routeId,
                              fromId=snapped_measure,   ## here we are using the native route measure as our referent feature
                              fromOffset=directional_multiplier*demarcation_distance,  ## roads and highways measures alway increase in the direction of increasing mileposts - not direction of travel
                              fromType=-3,              ## here we are telling the Geocoder to interpret fromId as a measure value, not milepost
                              returnGeometry=True)
                          
        g_result = gc.geocode(params)
##      pp(g_result)

        ## the Geocoder returns coordinates in Web Mercator (a coordinate systems that facilitates web mapping)
        ## Web Mercator is a planar coordinate system, as opposed to WGS84 (lat/long) which is a spherical coordinate system
        ## we will do are coordinate calculations in Web Mercator as their easier in planar systems, the native unit in Web Mercator in meter
        snapped_point_wm      = g_result['locations'][0]['geocode_fromLocation']['geocode_XY_WM']     ## the coordinates in web mercator coordinate system (wkid=102100)
        snapped_point_bearing = g_result['locations'][0]['geocode_fromLocation']['geocode_bearing']

        x, y = snapped_point_wm['x'], snapped_point_wm['y']
        p1_x, p1_y = x + centerline_offset * math.cos(snapped_point_bearing + math.pi / 2), y + centerline_offset * math.sin(snapped_point_bearing + math.pi / 2)
        p2_x, p2_y = x + centerline_offset * math.cos(snapped_point_bearing - math.pi / 2), y + centerline_offset * math.sin(snapped_point_bearing - math.pi / 2)
        ## now we will use a Geometry REST service to project the coordinates to WGS84
        proj_parameters = {'geometryType' : 'esriGeometryPoint', 'geometries': "{},{},{},{}".format(p1_x, p1_y, p2_x, p2_y), 'inSR': 102100, 'outSR': 4326, 'f': 'json'}   ## project to wgs84 coordinate system (wkid=4326)
        response = urlOpen(geometry_service+'/Project', urlEncode(proj_parameters).encode('utf-8'))
        response = json.loads(response.read().decode('utf-8'))

        if i == 1:
            p11_x, p11_y = response['geometries'][0]['x'], response['geometries'][0]['y']
            p22_x, p22_y = response['geometries'][1]['x'], response['geometries'][1]['y']
        elif i == 2:
            p3_x, p3_y = response['geometries'][0]['x'], response['geometries'][0]['y']
            p4_x, p4_y = response['geometries'][1]['x'], response['geometries'][1]['y']
        else :
            p5_x, p5_y = response['geometries'][0]['x'], response['geometries'][0]['y']
            p6_x, p6_y = response['geometries'][1]['x'], response['geometries'][1]['y']
        ##print("Coordinates of Warning {} located {} miles upstream of the reported location ({},{}) are: ({},{}))".format(i, demarcation_distance, slowdown_point[0], slowdown_point[1], x, y))

    return p11_x,p11_y,p22_x,p22_y,p3_x,p3_y,p4_x,p4_y,p5_x,p5_y,p6_x,p6_y

def main(argv):

    print(json.dumps(list(boq_war_loc(argv[0],argv[1]))))

if __name__ == '__main__':
    main(sys.argv[1:])


                              
    

