package gracenote

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	com"github.com/nemesisesq/ss_data_service/common"
	"encoding/json"
	"time"
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
	FormattedAddress  string `json:"formatted_address"`
	Geometry
	PlaceId           string `json:"place_id"`
	Types             []string `json:"types"`
}

type AddressComponent struct {
	LongName  string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string `json:"types"`
}

type Lineup struct {
	Type     string `json:"type"`
	Device   string `json:"device"`
	LineupId string `json:"lineupId"`
	Name     string `json:"name"`
	Location string `json:"location"`
	MSO      map[string]interface{} `json:"mso"`
}

type Guide struct {
	Location
	ZipCode  string
	Lineups  []Lineup
	Stations []Station
}



type Program struct {
	TMSID            string `json:"tmsId"`
	RootId           string `json:"rootId"`
	SeriesId         string `json:"seriesId"`
	SubType          string `json:"subType"`
	Title            string `json:"title"`
	ReleaseYear      string `json:"releaseYear"`
	ReleaseDate      string `json:"releaseDate"`
	OrigAirDate      string `json:"origAirDate"`
	TitleLang        string `json:"titleLang"`
	DescriptionLang  string `json:"descriptionLang"`
	EntityType       string `json:"entityType"`
	Genres           []string `json:"genres"`
	ShortDescription string `json:"shortDescription"`
}

type Airing struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Duration  int `json:"duration"`
	Channels  []string `json:"channels"`
	StationId string `json:"stationId"`
	Program
}

type Station struct {
	StationId      string `json:"stationId"`
	CallSign       string `json:"callSign"`
	Channel        string `json:"channel"`
	PreferredImage map[string]interface{} `json:"preferredImage"`
	Airings        []Airing `json:"airings"`
}

func GetLineupAirings(w http.ResponseWriter, r *http.Request) {

	guideObj := &Guide{}

	vars := mux.Vars(r)

	guideObj.Lat = vars["lat"]
	guideObj.Long = vars["long"]

	guideObj.SetZipCode()
	guideObj.GetLineups()
	guideObj.GetTVGrid()
}

func (g *Guide) GetTVGrid() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", LineupsUri + "/USA-ECHOST-DEFAULT/grid", nil)

	com.Check(err)

	curr_time := time.Now().Format(time.RFC3339)

	params := map[string]string{
		"api_key": ApiKey,
		"starDateTime": curr_time,
		//"lineupId" : "USA-ECHOST-DEFAULT",
		"size" : "Basic",
		"imageSize" : "Sm",
		"excludeChannels": "music, ppv, adult",
		"enhancedCallSign": "true",
	}

	BuildQuery(req, params)

	res, err := client.Do(req)
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&g.Stations)

	com.Check(err)
}

func BuildQuery(r *http.Request, m map[string]string) {
	q := r.URL.Query()

	for key, val := range m {
		q.Add(key, val)
	}
	r.URL.RawQuery = q.Encode()
}

func (g *Guide) GetLineups() {
	client := &http.Client{}
	//res, err := client.Get(LineupsUri)

	req, err := http.NewRequest("GET", LineupsUri, nil)

	params := map[string]string{"country": "USA", "postalCode": g.ZipCode, "api_key": ApiKey}

	BuildQuery(req, params)

	res, err := client.Do(req)

	defer res.Body.Close()

	com.Check(err)

	decoder := json.NewDecoder(res.Body)

	//var x []Lineup

	err = decoder.Decode(&g.Lineups)

	com.Check(err)

}

func (g *Guide) SetZipCode() {

	req, err := http.NewRequest("GET", GeoCodeUri, nil )

	com.Check(err)

	params := map[string]string{
		"latlng": fmt.Sprintf("%s,%s",g.Lat, g.Long),
		"sensor" : "true",
	}

	BuildQuery(req, params)

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


