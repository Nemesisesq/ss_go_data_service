package gracenote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sort"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v5"
	//"regexp"
	"regexp"
	"sync"
)

const format = "2006-01-02T15:04Z"

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
	Stations
}

type Guide struct {
	Location
	ZipCode string
	Lineups []Lineup
}

type Program struct {
	TMSID            string                 `json:"tmsId"`
	RootId           string                 `json:"rootId"`
	SeriesId         string                 `json:"seriesId"`
	SubType          string                 `json:"subType"`
	Title            string                 `json:"title"`
	EpisodeTitle     string                 `json:"episodeTitle"`
	ReleaseYear      int                    `json:"releaseYear"`
	ReleaseDate      string                 `json:"releaseDate"`
	OrigAirDate      string                 `json:"origAirDate"`
	TitleLang        string                 `json:"titleLang"`
	DescriptionLang  string                 `json:"descriptionLang"`
	EntityType       string                 `json:"entityType"`
	Genres           []string               `json:"genres"`
	ShortDescription string                 `json:"shortDescription"`
	PreferredImage   map[string]interface{} `json:"preferredImage"`
}

type StationMetaData struct {
	StationIdPrimary  string `json:"stationId_primary" json:"stationId_primary"`
	CallsignPrimary   string `json:"callsign_primary" bson:"callsign_primary"`
	CallsignSecondary string `json:"callsign_secondary" bson:"callsign_secondary"`
	DefaultRank       string `json:"default_rank" bson:"default_rank"`
}

type Station struct {
	StationId         string                 `json:"stationId"`
	CallSign          string                 `json:"callSign"`
	AffiliateCallSign string                 `json:"affiliateCallSign" bson:"affiliateCallSign"`
	Channel           string                 `json:"channel"`
	PreferredImage    map[string]interface{} `json:"preferredImage"`
	Airings           []Airing               `json:"airings"`
	DefaultRank       int                    `json:"default_rank"`
}

type Stations []Station

func (slice Stations) Len() int {
	return len(slice)
}

func (slice Stations) Less(i, j int) bool {
	return slice[i].DefaultRank < slice[j].DefaultRank
}

func (slice Stations) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Airing struct {
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
	Duration  int      `json:"duration"`
	Channels  []string `json:"channels"`
	StationId string   `json:"stationId"`
	Program   Program  `json:"program"`
}

func GetLineupAirings(w http.ResponseWriter, r *http.Request) {

	guideObj := &Guide{}
	vars := mux.Vars(r)
	guideObj.Lat = vars["lat"]
	guideObj.Long = vars["long"]
	//guideObj.CheckLineUpsForGeoCoords()
	guideObj.SetZipCode()
	guideObj.GetLineups(r)
	lineups := guideObj.GetTVGrid(r)
	stations := GetCombinedGrid(lineups)
	stations = guideObj.FilterAirings(stations, r)
	stations = RemoveDuplicates(stations)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(stations)

	com.Check(err)
}

func RemoveDuplicates(stations Stations) (dedupedStations Stations) {

	seen := map[string]bool{}

	for _, station := range stations {
		channelNumber, _ := strconv.Atoi(station.Channel)
		rank := strconv.Itoa(station.DefaultRank)

		if !seen[station.StationId] && channelNumber < 8000 && !seen[rank] {
			dedupedStations = append(dedupedStations, station)
			seen[station.StationId] = true
			seen[rank] = true
		}
	}

	return dedupedStations
}

func GetCombinedGrid(lineups []Lineup) (combinedStations Stations) {

	for _, lineup := range lineups {

		for _, station := range lineup.Stations {

			combinedStations = append(combinedStations, station)

		}
	}
	return combinedStations
}

func (lineup Lineup) GetFreshTVListingsGrid() []byte {
	log.SetFormatter(&log.JSONFormatter{})
	iClient := &http.Client{}
	url := fmt.Sprintf("%v/%v/grid", LineupsUri, lineup.LineupId)
	req, err := http.NewRequest("GET", url, nil)
	fmt.Println(lineup.LineupId)

	com.Check(err)

	fmt.Println("\n\nnow", time.Now())
	start_time := time.Now().Format(format)
	fmt.Println("\nstart", start_time)
	end_time := time.Now().Add(time.Minute * 30).Format(format)
	params := map[string]string{
		"api_key":       ApiKey,
		"startDateTime": start_time,
		"endDateTime":   end_time,

		"imageAspectTV":    "16x9",
		"size":             "Detailed",
		"imageSize":        "Sm",
		"excludeChannels":  "ppv,adult",
		"enhancedCallSign": "true",
	}
	com.BuildQuery(req, params)

	log.Debug(req)

	res, err := iClient.Do(req)
	com.Check(err)

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&lineup.Stations)
	com.Check(err)

	the_json, err := json.Marshal(&lineup.Stations)

	com.Check(err)

	return the_json
}
func (g *Guide) GetTVGrid(r *http.Request) (lineups []Lineup) {

	var wg sync.WaitGroup

	log.SetFormatter(&log.JSONFormatter{})
	rc := r.Context().Value("redis_client").(*redis.Client)

	lchan := make(chan Lineup, 2)

	for _, lineup := range g.Lineups {

		if IsRightLinup(lineup) {
			wg.Add(1)
			go func(lineup Lineup, wg *sync.WaitGroup, rc *redis.Client, lchan chan Lineup) {
				log.Info("Getting ", lineup.LineupId)
				val, err := rc.Get(lineup.LineupId).Result()

				if err == redis.Nil {
					the_json := lineup.GetFreshTVListingsGrid()
					timeout := time.Hour * 5
					rc.Set(lineup.LineupId, the_json, timeout)
					err = json.Unmarshal(the_json, &lineup.Stations)
					com.Check(err)

					//lineups = append(lineups, lineup)
					lchan <- lineup
					wg.Done()

				} else {
					log.Info("Redis Value Found for ", lineup.LineupId)
					json.Unmarshal([]byte(val), &lineup.Stations)

					//lineups = append(lineups, lineup)

					lchan <- lineup
					wg.Done()
				}
			}(lineup, &wg, rc, lchan)
		}
	}

	wg.Wait()

	close(lchan)

	for x := range lchan {
		lineups = append(lineups, x)
	}
	return lineups

}

func IsRightLinup(lineup Lineup) bool {
	uverseMatch, _ := regexp.Match("U-verse", []byte(lineup.Name))

	return lineup.LineupId == "USA-ECHOST-DEFAULT" || uverseMatch
}

func (g *Guide) GetLineups(r *http.Request) {

	db := r.Context().Value("db").(mgo.Database)
	c := db.C("lineups")

	//pipeline := []bson.M{
	//	{"$match": bson.M{"zip_code": g.ZipCode}},
	//	{"$or": []bson.M{
	//		{"lineup_id": "USA-ECHOST-DEFAULT"},
	//		{"name": `/U-verse/i`},
	//	}},
	//}

	query := *c.Find(bson.M{"zip_code": g.ZipCode})
	count, err := query.Count()

	com.Check(err)

	if count > 0 {

		err := query.All(&g.Lineups)
		com.Check(err)

		log.WithField("count", count).Info("Lineups were found in the database")

		return
		//err := pipe.All(&g.Lineups)

		//TODO do some stuff here we would want to return all the lineups for a zipcode evenrtually or crtain lineups based on query

		//return lineups
	}

	iClient := &http.Client{}
	req, err := http.NewRequest("GET", LineupsUri, nil)
	params := map[string]string{"country": "USA", "postalCode": g.ZipCode, "api_key": ApiKey}
	com.BuildQuery(req, params)
	log.WithFields(log.Fields{
		"request url": req.URL.Path,
		"postal code": g.ZipCode,
		"API key":     ApiKey,
	}).Info()
	res, err := iClient.Do(req)
	log.WithField("request status", res.Status).Info()

	//log.Info("got reponse from call to lineups")
	defer res.Body.Close()

	com.Check(err)

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&g.Lineups)

	com.Check(err)

	for _, val := range g.Lineups {
		val.ZipCode = g.ZipCode
	}

	//TODO Do something to pick the correct lineup here
	//fmt.Println(g.Lineups)

	for _, l := range g.Lineups {
		l.ZipCode = g.ZipCode
		err = c.Insert(l)
		com.Check(err)
	}

	index := mgo.Index{
		Key:        []string{"zip_code"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	}
	err = c.EnsureIndex(index)
	com.Check(err)
}

func (g *Guide) SetZipCode() {
	geoCode := GetGeoCodeFromGoogle(g.Lat, g.Long)

	g.ZipCode = ExtractZipFromGeoCode(geoCode)

}

func ExtractZipFromGeoCode(geoCode GeoCode) (zip string) {


	for _, result := range geoCode.Results {
		for _, component := range result.AddressComponents {
			for _, t := range component.Types {
				if t == "postal_code" {

					zip = component.LongName
				}
			}
		}
	}

	return zip
}

func GetGeoCodeFromGoogle(lat, long string) GeoCode {

	req, err := http.NewRequest("GET", GeoCodeUri, nil)
	com.Check(err)
	params := map[string]string{
		"latlng": fmt.Sprintf("%s,%s", lat, long),
		"sensor": "true",
	}
	com.BuildQuery(req, params)
	iClient := &http.Client{}
	res, err := iClient.Do(req)
	defer res.Body.Close()

	com.Check(err)

	geoCode := &GeoCode{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&geoCode)

	com.Check(err)

	return *geoCode
}

func (g *Guide) FilterAirings(stations Stations, r *http.Request) (filteredStations Stations) {

	db := r.Context().Value("db").(mgo.Database)
	col := db.C("live_streaming_services")

	filtered := make(chan Station)

	log.Print("entering")
	go stations.Process(col, filtered)

	for s := range filtered {
		filteredStations = append(filteredStations, s)
	}

	//fmt.Println(len(filteredStations))

	sort.Sort(filteredStations)

	log.Print("exiting")
	return filteredStations
}

func (stations Stations) Process(col *mgo.Collection, filtered chan Station) {
	var md []StationMetaData
	var wg sync.WaitGroup

	err := col.Find(nil).All(&md)

	com.Check(err)
	//log.Println(len(md))

	for _, station := range stations {

		wg.Add(1)
		go func(station Station, wg *sync.WaitGroup) {

			//start := time.Now()
			for _, record := range md {
				//r, _ := json.MarshalIndent(record, "","\t")
				//fmt.Println(string(r))
				if record.CallsignPrimary == station.CallSign && record.CallsignPrimary != "" {
					station.DefaultRank, err = strconv.Atoi(record.DefaultRank)
					//log.Printf("passing station to filterd channel %v", record.DefaultRank)

					filtered <- station
				} else if record.CallsignSecondary == station.CallSign && record.CallsignSecondary != "" {
					station.DefaultRank, err = strconv.Atoi(record.DefaultRank)
					//log.Printf("passing station to filterd channel %v", record.DefaultRank)

					filtered <- station
				} else if record.CallsignPrimary == station.AffiliateCallSign && record.CallsignPrimary != "" {
					station.DefaultRank, err = strconv.Atoi(record.DefaultRank)
					//log.Printf("passing station to filterd channel %v", record.DefaultRank)

					filtered <- station
				} else if record.CallsignSecondary == station.AffiliateCallSign && record.CallsignSecondary != "" {
					station.DefaultRank, err = strconv.Atoi(record.DefaultRank)
					//log.Printf("passing station to filterd channel %v", record.DefaultRank)

					filtered <- station
				}

			}

			//log.Printf("records loop duration %v", time.Since(start))
			wg.Done()

			//query := []bson.M{}
			//
			//query = append(query, bson.M{"stationId_primary": station.StationId})
			//query = append(query, bson.M{"stationId_second": station.StationId})
			//if station.CallSign != "" {
			//	query = append(query, bson.M{"callsign_primary": station.CallSign})
			//	query = append(query, bson.M{"callsign_secondary": station.CallSign})
			//}
			//
			//if station.AffiliateCallSign != "" {
			//	query = append(query, bson.M{"callsign_secondary": station.AffiliateCallSign})
			//	query = append(query, bson.M{"callsign_primary": station.AffiliateCallSign})
			//}
			//count, _ := col.Find(bson.M{"$or": query}).Count()
			//if count > 0 {
			//	//wg.Add(1)
			//	md := &StationMetaData{}
			//	err := col.Find(bson.M{"$or": query}).One(&md)
			//	com.Check(err)
			//
			//	station.DefaultRank, err = strconv.Atoi(md.DefaultRank)
			//	filtered <- station
			//	wg.Done()
			//} else {
			//	wg.Done()
			//}

		}(station, &wg)

	}

	log.Println("waiting")
	wg.Wait()
	log.Println("Done Waiting")
	close(filtered)

}
