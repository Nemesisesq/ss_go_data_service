package gracenote

import "fmt"



const BaseUri = "http://data.tmsapi.com/v1.1"

var LineupsUri = fmt.Sprintf("%v/lineups", BaseUri)

const ApiKey = "3w8hvfmfxjuwgvbqkahrss35"

const GeoCodeUri = "http://maps.googleapis.com/maps/api/geocode/json?"

var SportsUri = fmt.Sprintf("%v/sports", BaseUri)

