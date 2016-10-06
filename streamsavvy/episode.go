package streamsavvy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	//"github.com/gorilla/mux"
	"encoding/json"

	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)

type Episode struct {
	Id               int      `json:"id"`
	Tvdb             int      `json:"tvdb"`
	ContentType      string   `json:"content_type"`
	IsShadow         int      `json:"is_shadow"`
	AlternateTvdb    []string `json:"alternate_tvdb"`
	ImdbId           string   `json:"imdb_id"`
	SeasonNumber     int      `json:"season_number"`
	EpisodeNumber    int      `json:"episode_number"`
	ShowId           int      `json:"show_id"`
	Themoviedb       int      `json:"themoviedb"`
	Special          int      `json:"special"`
	FirstAired       string   `json:"first_aired"`
	Title            string   `json:"title"`
	OriginalTitle    string   `json:"original_title"`
	AlternateTitles  []string `json:"alternate_titles"`
	Overview         string   `json:"overview"`
	Duration         int      `json:"duration"`
	ProductionCode   string   `json:"production_code"`
	Thumbnail208X117 string   `json:"thumbnail_208x117"`
	Thumbnail304X171 string   `json:"thumbnail_304x171"`
	Thumbnail400X225 string   `json:"thumbnail_400x225"`
	Thumbnail608X342 string   `json:"thumbnail_608x342"`
}
type GuideBoxEpisodes struct {
	GuideboxId    string    `bson:"guidebox_id"`
	Results       []Episode `json:"results" bson:"list"`
	TotalResults  int   `json:"total_results" bson:"total_results"`
	TotalReturned int    `json:"total_returned" bson:"total_returned"`
}

func GetEpisodes(w http.ResponseWriter, r *http.Request) {

	guideboxId := r.URL.Query().Get("guidebox_id")
	db := context.Get(r, "db").(*mgo.Database)

	c := db.C("episodes")

	query := c.Find(bson.M{"guidebox_id": guideboxId})

	count, err := query.Count()

	epi := &GuideBoxEpisodes{}

	com.Check(err)

	if count > 0 {
		query.One(&epi)
	} else {


		total_results, episode_list :=epi.GetAllEpisodes(0, 25, guideboxId)


		if total_results > len(episode_list) {
			for total_results > len(episode_list) {
				_, res := epi.GetAllEpisodes(len(episode_list), 25, guideboxId)
				episode_list = append(episode_list,res...)
			}

			epi.Results = episode_list

		}

		epi.GuideboxId = guideboxId
		epi.TotalResults = total_results
		epi.TotalReturned = len(episode_list)

		c.Insert(*epi)

	}

	json.NewEncoder(w).Encode(epi)

}


func (gbe GuideBoxEpisodes) GetAllEpisodes(start int, chunk int, guideboxId string) (total_results int, epiList []Episode) {


	client := &http.Client{}

	baseUrl := "https://api-public.guidebox.com/v1.43/US/%v/show/%v/episodes/all/%v/%v/all/all/true"

	apiKey := "rKWvTOuKvqzFbORmekPyhkYMGinuxgxM"
	//TODO Actually set the correct Lineup id in the URL here
	//url := fmt.Sprintf("%v/%v/grid", LineupsUri, lineup.LineupId)
	url := fmt.Sprintf(baseUrl, apiKey, guideboxId, start, chunk)

	req, err := http.NewRequest("GET", url, nil)

	com.Check(err)

	params := map[string]string{
		"reverse_ordering": "true",
	}

	com.BuildQuery(req, params)

	res, err := client.Do(req)

	com.Check(err)

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&gbe)

	com.Check(err)

	total_results = gbe.TotalResults
	epiList = gbe.Results

	return total_results, epiList
}
