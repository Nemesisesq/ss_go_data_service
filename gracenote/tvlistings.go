package gracenote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GeoCode struct {
	Results []Result `json:"results"`
}

type Location struct {
	Lat  string
	Long string
}

type Viewport struct {
	Northeast Location `json:"northeast"`
	Southeast Location `json:"southeast"`
}

type Geometry struct {
	Location
	LocationType string `json:"location_type"`
	Viewport
}

type Result struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry
	PlaceId           string   `json:"place_id"`
	Types             []string `json:"types"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Lineup struct {
	ZipCode  string                 `json:"zip_code" bson:"zip_code"`
	Type     string                 `json:"type" bson:"type"`
	Device   string                 `json:"device" bson:"device"`
	LineupId string                 `json:"lineupId" bson:"lineup_id"`
	Name     string                 `json:"name" bson:"name"`
	Location string                 `json:"location" bson:"location"`
	MSO      map[string]interface{} `json:"mso" bson:"mso"`
	Stations []Station
}

type Guide struct {
	Location
	ZipCode string
	Lineups []Lineup
}

type Program struct {
	TMSID            string   `json:"tmsId"`
	RootId           string   `json:"rootId"`
	SeriesId         string   `json:"seriesId"`
	SubType          string   `json:"subType"`
	Title            string   `json:"title"`
	ReleaseYear      int      `json:"releaseYear"`
	ReleaseDate      string   `json:"releaseDate"`
	OrigAirDate      string   `json:"origAirDate"`
	TitleLang        string   `json:"titleLang"`
	DescriptionLang  string   `json:"descriptionLang"`
	EntityType       string   `json:"entityType"`
	Genres           []string `json:"genres"`
	ShortDescription string   `json:"shortDescription"`
}

type Station struct {
	StationId      string `json:"stationId"`
	CallSign       string `json:"callSign"`
	AffiliateCallSign string `json:"affiliateCallSign" bson:"affiliateCallSign`
	Channel        string `json:"channel"`
	PreferredImage map[string]interface{} `json:"preferredImage"`
	Airings        []Airing `json:"airings"`
}

type Airing struct {
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	Duration  int      `json:"duration"`
	Channels  []string `json:"channels"`
	StationId string   `json:"stationId"`
	Program   Program `json:"program"`
}

func GetLineupAirings(w http.ResponseWriter, r *http.Request) {

	guideObj := &Guide{}

	vars := mux.Vars(r)

	guideObj.Lat = vars["lat"]
	guideObj.Long = vars["long"]

	guideObj.CheckLineUpsForGeoCoords()

	guideObj.SetZipCode()
	lineup := guideObj.GetLineups(r)
	stations := guideObj.GetTVGrid(r, lineup)


	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(stations)

	com.Check(err)
}

func (g *Guide) GetTVGrid(r *http.Request, lineup Lineup) []Station {
	db := context.Get(r, "db").(*mgo.Database)

	c := db.C("lineups")

	l := &Lineup{}

	query := c.Find(bson.M{"lineup_id": lineup.LineupId})

	count, err := query.Count()

	com.Check(err)

	if count > 0 {

		query.One(&l)
		return l.Stations
	}

	client := &http.Client{}
	//TODO Actually set the correct Lineup id in the URL here
	//url := fmt.Sprintf("%v/%v/grid", LineupsUri, lineup.LineupId)
	url := fmt.Sprintf("%v/%v/grid", LineupsUri, "USA-TX42500-X")
	req, err := http.NewRequest("GET", url, nil)

	com.Check(err)

	curr_time := time.Now().Format(time.RFC3339)

	params := map[string]string{
		"api_key":      ApiKey,
		"starDateTime": curr_time,
		//"lineupId" : "USA-ECHOST-DEFAULT",
		//"imageSize":"Md",
		"imageAspectTV":"16x9",
		"size":             "Basic",
		"imageSize":        "Sm",
		"excludeChannels":  "music, ppv, adult",
		"enhancedCallSign": "true",
	}

	com.BuildQuery(req, params)

	res, err := client.Do(req)

	//body, err := ioutil.ReadAll(res.Body)

	//if err != nil {
	//	log.Fatalf("ERROR: %s", err)
	//}
	//
	//fmt.Printf("%s", body)

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&lineup.Stations)

	com.Check(err)

	return lineup.Stations

}


func (g *Guide) CheckLineUpsForGeoCoords() {
	//TODO check the geo coordinates for
}

func (g *Guide) GetLineups(r *http.Request) (lineup Lineup) {

	db := context.Get(r, "db").(*mgo.Database)

	c := db.C("lineups")

	query := c.Find(bson.M{"zip_code": g.ZipCode})
	count, err := query.Count()

	com.Check(err)

	if count > 0 {
		//TODO do some stuff here we would want to return all the lineups for a zipcode evenrtually
		query.One(&lineup)

		return lineup
	}

	client := &http.Client{}
	//res, err := client.Get(LineupsUri)

	req, err := http.NewRequest("GET", LineupsUri, nil)

	params := map[string]string{"country": "USA", "postalCode": g.ZipCode, "api_key": ApiKey}

	com.BuildQuery(req, params)

	res, err := client.Do(req)

	defer res.Body.Close()

	com.Check(err)

	decoder := json.NewDecoder(res.Body)

	//var x []Lineup

	err = decoder.Decode(&g.Lineups)

	com.Check(err)

	//TODO Do something to pick the correct lineup here

	//save the New Lineups

	c.Insert(g.Lineups)

	return g.Lineups[0]

}

func (g *Guide) SetZipCode() {

	req, err := http.NewRequest("GET", GeoCodeUri, nil)

	com.Check(err)

	params := map[string]string{
		"latlng": fmt.Sprintf("%s,%s", g.Lat, g.Long),
		"sensor": "true",
	}

	com.BuildQuery(req, params)

	client := &http.Client{}

	res, err := client.Do(req)

	defer res.Body.Close()

	com.Check(err)

	geoCode := &GeoCode{}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&geoCode)

	com.Check(err)

	for _, result := range geoCode.Results {
		for _, component := range result.AddressComponents {
			for _, t := range component.Types {
				if t == "postal_code" {
					g.ZipCode = component.LongName
					return
				}
			}
		}
	}
}
