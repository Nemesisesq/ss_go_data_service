package streamsavvy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	//"github.com/gorilla/mux"
	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"gopkg.in/mgo.v2"
)

type Episode struct {
	Id               int    `json:"id"`
	Tvdb             int `json:"tvdb"`
	ContentType      string `json:"content_type"`
	IsShadow         int `json:"is_shadow"`
	AlternateTvdb    []string `json:"alternate_tvdb"`
	ImdbId           string `json:"imdb_id"`
	SeasonNumber     int `json:"season_number"`
	EpisodeNumber    int `json:"episode_number"`
	ShowId           int `json:"show_id"`
	Themoviedb       int `json:"themoviedb"`
	Special          int `json:"special"`
	FirstAired       string `json:"first_aired"`
	Title            string `json:"title"`
	OriginalTitle    string `json:"original_title"`
	AlternateTitles  []string `json:"alternate_titles"`
	Overview         string `json:"overview"`
	Duration         int `json:"duration"`
	ProductionCode   string `json:"production_code"`
	Thumbnail208X117 string `json:"thumbnail_208x117"`
	Thumbnail304X171 string `json:"thumbnail_304x171"`
	Thumbnail400X225 string `json:"thumbnail_400x225"`
	Thumbnail608X342 string `json:"thumbnail_608x342"`
}
type Episodes struct {

	GuideboxId string `bson:"guidebox_id"`
	List []Episode `json:"results" bson:"list"`
}

func GetEpisodes(w http.ResponseWriter, r *http.Request) {

	guideboxId := r.URL.Query().Get("guidebox_id")
	db := context.Get(r, "db").(*mgo.Database)

	c := db.C("episodes")

	query := c.Find(bson.M{"guidebox_id": guideboxId})

	count, err := query.Count()

	epi := &Episodes{}

	com.Check(err)

	if count > 0 {
		query.One(&epi)
	} else {
		client := &http.Client{}

		baseUrl := "https://api-public.guidebox.com/v1.43/US/%v/show/%v/episodes/all/0/100/all/all/true"

		apiKey := "rKWvTOuKvqzFbORmekPyhkYMGinuxgxM"
		//TODO Actually set the correct Lineup id in the URL here
		//url := fmt.Sprintf("%v/%v/grid", LineupsUri, lineup.LineupId)
		url := fmt.Sprintf(baseUrl, apiKey,guideboxId)

		req, err := http.NewRequest("GET", url, nil)

		com.Check(err)

		params := map[string]string{
			"reverse_ordering": "true",
		}

		com.BuildQuery(req, params)

		res, err := client.Do(req)

		decoder := json.NewDecoder(res.Body)

		err = decoder.Decode(&epi)

		com.Check(err)

		epi.GuideboxId = guideboxId

		c.Insert(*epi)

	}

	json.NewEncoder(w).Encode(epi)

}
