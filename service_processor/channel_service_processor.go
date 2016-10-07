package service_processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	com "github.com/nemesisesq/ss_data_service/common"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
)

type RawPayload struct {
	CallLetters    string `json:"CallLetters"`
	DisplayName    string `json:"DisplayName"`
	SourceLongName string `json:"SourceLongName"`
	SourceId       string `json:"SourceId"`
}

type StreamingSource struct {
	Source        string `json:"source"`
	MatchedSource string `json:"matched_source"`
	DisplayName   string `json:"display_name"`
	Id            int    `json:"id"`
	DeepLinks     Links  `json:"deep_links"`
}

//type StreamingSources struct {
//	TypeName []*StreamingSource
//}

type PP struct {
	StreamingSources []StreamingSource `json:"streamingServices"`
}

type Links struct {
	DeepLinks []string
	AppStore  string
}

type ViewingWindows struct {
	PayPerView []StreamingSource `json:"pay_per_view"`
	Binge      []StreamingSource `json:"binge"`
	Live       []StreamingSource `json:"live"`
	OnDemand   []StreamingSource `json:"on_demand"`
	Misc       []StreamingSource `json:"misc"`
}

func GetOnDemandServices(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	url := fmt.Sprintf("%s/detail_sources", os.Getenv("NODE_DATA_SERVICE"))
	fmt.Println("URL created")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	fmt.Println("request made")
	com.Check(err)

	//v := &ViewingWindows{}
	v := make(map[string]interface{})
	decoder := json.NewDecoder(response.Body)
	fmt.Println("Decoding boddy")
	fmt.Println(response.Body)
	fmt.Println("Response from Node server unmarshalled to map interface made")
	err = decoder.Decode(&v)
	com.Check(err)

	fToExt := []string{"on_demand", "binge", "pay_per_view"}

	ss_slice := []StreamingSource{}

	for _, fldNm := range fToExt {
		//fmt.Println(reflect.TypeOf(v[fldNm]))=
		if t, ok := v[fldNm]; ok {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				fmt.Println(s.Index(i))
				data := s.Index(i).Interface().(map[string]interface{})
				fmt.Println(data)

				newSS := &StreamingSource{}

				jsonData, err := json.Marshal(data)

				err = json.Unmarshal(jsonData, newSS)

				com.Check(err)

				newSS.MatchedSource = newSS.Source

				//filter out amazon and goole

				googleRegExp, err := regexp.CompilePOSIX(`google`)
				com.Check(err)
				amazonRegExp, err := regexp.CompilePOSIX(`amazon_buy`)
				com.Check(err)

				if !googleRegExp.Match([]byte(newSS.Source)) && !amazonRegExp.Match([]byte(newSS.Source)) {

					ss_slice = append(ss_slice, *newSS)
				}

			}
		}
	}
	for idx, val := range ss_slice {
		streamSource := MatchDeepLinks(&val)
		ss_slice[idx] = *streamSource
	}

	ODPayload := &PP{}

	for _, val := range ss_slice {
		if CheckIfLowestTier(ss_slice, val.Source) {
			ODPayload.StreamingSources = append(ODPayload.StreamingSources, val)
		}
	}

	ODPayload.RemoveDuplicates()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ODPayload.StreamingSources)

}

func CheckIfLowestTier(ss_slice []StreamingSource, source string) bool {

	source_slice := &PP{}

	source_slice.StreamingSources = ss_slice

	x := true

	switch source {
	case "sling_blue":
		x = true
	case "sling_orange":

		x = !source_slice.CheckStreamingServicesForSource("sling_blue")
	case "sling_blue_orange":
		x = !source_slice.CheckStreamingServicesForSource("sling_orange")
	case "sony_vue_slim":
		x = true
	case "sony_vue_core":
		x = !source_slice.CheckStreamingServicesForSource("sony_vue_slim")
	case "sony_vue_elite":
		x = !source_slice.CheckStreamingServicesForSource("sony_vue_core")
	}
	return x
}

func (a PP) CheckStreamingServicesForSource(source string) bool {
	for _, val := range a.StreamingSources {
		if val.Source == source {
			return true
		}
	}

	return false
}

func (a *PP) RemoveDuplicates() {
	m := map[string]bool{}

	for _, i := range a.StreamingSources {
		if _, seen := m[i.Source]; !seen {
			a.StreamingSources[len(m)] = i
			m[i.Source] = true
		}
	}

	a.StreamingSources = a.StreamingSources[:len(m)]
}

func GetLiveStreamingServices(w http.ResponseWriter, r *http.Request) {

	print(`hello world`)

	rawPayload := &RawPayload{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rawPayload)

	com.Check(err)

	url := fmt.Sprintf("%s/guide", os.Getenv("NODE_DATA_SERVICE"))
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(rawPayload)
	res, err := http.Post(url, "application/json", buf)
	defer res.Body.Close()

	com.Check(err)

	processedPayloads := &PP{}
	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&processedPayloads)
	com.Check(err)

	for i, sS := range processedPayloads.StreamingSources {
		//fmt.Println(sS.Source)

		sS.MatchedSource = sS.Source

		streamSource := MatchDeepLinks(&sS)

		processedPayloads.StreamingSources[i] = *streamSource

	}

	NewPP := &PP{}

	for _, val := range processedPayloads.StreamingSources {
		if CheckIfLowestTier(processedPayloads.StreamingSources, val.Source) {
			NewPP.StreamingSources = append(NewPP.StreamingSources, val)
		}
	}

	seen := map[string]bool{}
	resSS := []StreamingSource{}
	for i, val := range NewPP.StreamingSources {

		switch {

		case GSM(`sling`, val.Source):

			if seen["sling"] {

			} else {

				NewPP.StreamingSources[i].Source = "sling"
				resSS = append(resSS, NewPP.StreamingSources[i])
				seen["sling"] = true
			}
		case GSM(`vue`, val.Source):
			if seen["ps_vue"] {

			} else {

				NewPP.StreamingSources[i].Source = "ps_vue"
				resSS = append(resSS, NewPP.StreamingSources[i])
				seen["ps_vue"] = true
			}

		case GSM(`ota`, val.Source):
			continue

		default:
			resSS = append(resSS, NewPP.StreamingSources[i])

		}

	}

	NewPP.StreamingSources = resSS

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewPP)
}

func MatchDeepLinks(sS *StreamingSource) *StreamingSource {
	deepLinkMap := GetDeepLinks()

	switch {

	case GSM(`youtube`, sS.MatchedSource):
		sS.MatchedSource = "youtube"
	case GSM(`sling`, sS.MatchedSource):
		sS.MatchedSource = "sling_tv"
	case GSM(`(vue|sony|playstation)`, sS.MatchedSource):
		sS.MatchedSource = "playstation_vue"
	case GSM(`hulu`, sS.MatchedSource):
		sS.MatchedSource = "hulu"
	}

	sS.DeepLinks = deepLinkMap[sS.MatchedSource]

	return sS
}

func GSM(key string, source string) bool {
	re, err := regexp.Compile(key)
	match := re.Match([]byte(source))
	com.Check(err)

	return match

}

func GetDeepLinks() map[string]Links {
	deepLinks := make(map[string]Links)

	deepLinks["hulu"] = Links{
		DeepLinks: []string{"fb40582213222://", "hulu://"},
		AppStore:  "https://itunes.apple.com/us/app/hulu/id376510438?mt=8",
	}

	deepLinks["netflix"] = Links{
		DeepLinks: []string{"fb163114453728333://"},
		AppStore:  "https://itunes.apple.com/us/app/netflix/id363590051?mt=8",
	}

	deepLinks["sling_tv"] = Links{
		DeepLinks: []string{"slingtv://"},
		AppStore:  "https://itunes.apple.com/us/app/sling-tv-live-and-on-demand/id945077360?mt=8",
	}
	deepLinks["hbo_now"] = Links{
		DeepLinks: []string{"hbonow://"},
		AppStore:  "https://itunes.apple.com/us/app/hbo-now/id971265422?mt=8",
	}
	deepLinks["playstation_vue"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "https://itunes.apple.com/us/app/playstation-vue-mobile/id957987596?mt=8",
	}
	deepLinks["showtime"] = Links{
		DeepLinks: []string{"fb457474864412373://", "showtime://"},
		AppStore:  "https://itunes.apple.com/us/app/showtime/id996246479?mt=8",
	}
	deepLinks["amazon_instant_video"] = Links{
		DeepLinks: []string{"aiv://"},
		AppStore:  "https://itunes.apple.com/us/app/amazon-video/id545519333?mt=8",
	}

	deepLinks["amazon_prime"] = Links{
		DeepLinks: []string{"aiv://"},
		AppStore:  "https://itunes.apple.com/us/app/amazon-video/id545519333?mt=8",
	}
	deepLinks["cbs"] = Links{
		DeepLinks: []string{"cbs-svod://", "cbs-sv://"},
		AppStore:  "https://itunes.apple.com/us/app/cbs-watch-full-episodes-on/id530168168?mt=8",
	}
	deepLinks["twitter"] = Links{
		DeepLinks: []string{"twitter://"},
		AppStore:  "https://itunes.apple.com/us/app/twitter/id333903271?mt=8",
	}
	deepLinks["starz"] = Links{
		DeepLinks: []string{"starz://", "fb385790948110676://"},
		AppStore:  "https://itunes.apple.com/us/app/starz/id550221096?mt=8",
	}
	deepLinks["nbc"] = Links{
		DeepLinks: []string{"nbctve://"},
		AppStore:  "https://itunes.apple.com/us/app/nbc-watch-now-stream-full/id442839435?mt=8",
	}
	deepLinks["the_cw"] = Links{
		DeepLinks: []string{"fb111598788905376comcwfullepisodesios://", "fb391138331009052comcwfullepisodesios://", "cwtv://"},
		AppStore:  "https://itunes.apple.com/us/app/the-cw/id491730359?mt=8",
	}
	deepLinks["cw_seed"] = Links{
		DeepLinks: []string{"cwseed://", "cwseed-pvn://"},
		AppStore:  "https://itunes.apple.com/us/app/cw-seed/id967093677?mt=8",
	}
	deepLinks["seeso"] = Links{
		DeepLinks: []string{"fb1505747749718971://", "seeso://"},
		AppStore:  "https://itunes.apple.com/us/app/seeso/id1053181816?mt=8",
	}
	deepLinks["acorn_tv"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "https://itunes.apple.com/us/app/acorn-tv-the-best-british-tv/id896014310?mt=8",
	}
	deepLinks["history"] = Links{
		DeepLinks: []string{"fb271893612990351://", "historyplus://"},
		AppStore:  "https://itunes.apple.com/us/app/history/id576009463?mt=8",
	}
	deepLinks["history_vault"] = Links{
		DeepLinks: []string{"historyvault://"},
		AppStore:  "https://itunes.apple.com/us/app/history-vault/id1076619087?mt=8",
	}
	deepLinks["twitch"] = Links{
		DeepLinks: []string{"ttv://", "twitch://"},
		AppStore:  "https://itunes.apple.com/us/app/twitch/id460177396?mt=8",
	}
	deepLinks["machinima"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "nil",
	}
	deepLinks["fubo_tv"] = Links{
		DeepLinks: []string{"fb162788457265492://", "tv.fubo.mobile://", "a09gnk7Cg77AOFuNemS5Kek3qT8HiJ87N1://", "a0RTpfOxSkv1LOIgHlxo2YmxLqUfaSvjaG://", "a0d6YiOzgcOnC305cKkBZoydAu62K1z7Ly://"},
		AppStore:  "https://itunes.apple.com/us/app/fubotv-live/id905401434?mt=8",
	}
	deepLinks["crunchyroll"] = Links{
		DeepLinks: []string{"fb56424855326://", "crunchyroll://"},
		AppStore:  "https://itunes.apple.com/us/app/crunchyroll-everything-anime/id329913454?mt=8",
	}
	deepLinks["pbs_kids"] = Links{
		DeepLinks: []string{"fb151570254902333://", "pbskidsvideo://"},
		AppStore:  "https://itunes.apple.com/us/app/pbs-kids-video/id435138734?mt=8",
	}
	deepLinks["tubi_tv"] = Links{
		DeepLinks: []string{"fb205962049613862://", "tubitv://"},
		AppStore:  "https://itunes.apple.com/us/app/tubi-tv-stream-free-movies/id886445756?mt=8",
	}
	deepLinks["crackle"] = Links{
		DeepLinks: []string{"fb91018702399://", "crackle://"},
		AppStore:  "https://itunes.apple.com/us/app/crackle-movies-tv/id377951542?mt=8",
	}
	deepLinks["newsy"] = Links{
		DeepLinks: []string{"fb396724197178595://", "newsy://"},
		AppStore:  "https://itunes.apple.com/us/app/newsy-video-news/id330879884?mt=8",
	}
	deepLinks["mlb_tv"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "nil",
	}
	deepLinks["nba_league_pass"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "nil",
	}
	deepLinks["youtube"] = Links{
		DeepLinks: []string{"youtube://watch?v=<video-id here>", "youtube://"},
		AppStore:  "https://itunes.apple.com/us/app/youtube/id544007664?mt=8",
	}
	deepLinks["vudu"] = Links{
		DeepLinks: []string{"vuduiosplayer://"},
		AppStore:  "https://itunes.apple.com/us/app/vudu-movies-tv/id487285735?mt=8",
	}
	deepLinks["itunes"] = Links{
		DeepLinks: []string{"nil"},
		AppStore: "https://itunes.apple.com",
	}
	deepLinks["google_play"] = Links{
		DeepLinks: []string{"nil"},
		AppStore:  "https://itunes.apple.com/us/app/google-play-movies-tv/id746894884?mt=8",
	}
	return deepLinks
}

//
