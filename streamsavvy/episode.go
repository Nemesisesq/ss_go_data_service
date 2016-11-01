package streamsavvy

import (
	"fmt"
	"net/http"
	//"github.com/gorilla/context"
	"encoding/json"

	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/redis.v4"
	"time"
)

type Episode struct {
	Id                         int         `json:"id"`
	Tvdb                       int         `json:"tvdb"`
	ContentType                string      `json:"content_type"`
	IsShadow                   int         `json:"is_shadow"`
	AlternateTvdb              []string    `json:"alternate_tvdb"`
	ImdbId                     string      `json:"imdb_id"`
	SeasonNumber               int         `json:"season_number"`
	EpisodeNumber              int         `json:"episode_number"`
	ShowId                     int         `json:"show_id"`
	Themoviedb                 int         `json:"themoviedb"`
	Special                    int         `json:"special"`
	FirstAired                 string      `json:"first_aired"`
	Title                      string      `json:"title"`
	OriginalTitle              string      `json:"original_title"`
	AlternateTitles            []string    `json:"alternate_titles"`
	Overview                   string      `json:"overview"`
	Duration                   int         `json:"duration"`
	ProductionCode             string      `json:"production_code"`
	Thumbnail208X117           string      `json:"thumbnail_208x117"`
	Thumbnail304X171           string      `json:"thumbnail_304x171"`
	Thumbnail400X225           string      `json:"thumbnail_400x225"`
	Thumbnail608X342           string      `json:"thumbnail_608x342"`
	FreeWebSources             interface{} `json:"free_web_sources"`
	FreeIosSources             interface{} `json:"free_ios_sources"`
	SubscriptionWebSources     interface{} `json:"subscription_web_sources"`
	SubscriptionIosSources     interface{} `json:"subscription_ios_sources"`
	SubscriptionAndroidSources interface{} `json:"subscription_android_sources"`
	PurchaseWebSources         interface{} `json:"purchase_web_sources"`
	PurchaseIosSources         interface{} `json:"purchase_ios_sources"`
	PurchaseAndroidSources     interface{} `json:"purchase_android_sources"`
}
type GuideBoxEpisodes struct {
	GuideboxId    string    `bson:"guidebox_id"`
	Results       []interface{} `json:"results"`
	TotalResults  int       `json:"total_results" bson:"total_results"`
	TotalReturned int       `json:"total_returned" bson:"total_returned"`
}
it 
func GetEpisodes(w http.ResponseWriter, r *http.Request) {

	guideboxId := r.URL.Query().Get("guidebox_id")

	epi := &GuideBoxEpisodes{}

	client := r.Context().Value("redis_client").(*redis.Client)

	val, err := client.Get(guideboxId).Result()

	if err == redis.Nil {

		total_results, episode_list := epi.GetAllEpisodes(0, 25, guideboxId)

		if total_results > len(episode_list) {
			for total_results > len(episode_list) {
				_, res := epi.GetAllEpisodes(len(episode_list), 25, guideboxId)
				episode_list = append(episode_list, res...)
			}

		}

		print(total_results)
		epi.Results = episode_list

		epi.GuideboxId = guideboxId
		epi.TotalResults = total_results
		epi.TotalReturned = len(episode_list)

		val, err := json.Marshal(epi)
		com.Check(err)

		timeout := time.Hour * 24 * 3
		err = client.Set(guideboxId, val, timeout).Err()
		com.Check(err)
	} else {
		json.Unmarshal([]byte(val), &epi)
	}

	json.NewEncoder(w).Encode(epi)

}

func (gbe GuideBoxEpisodes) GetAllEpisodes(start int, chunk int, guideboxId string) (total_results int, epiList []interface{}) {

	//TODO Logging

	client := &http.Client{}

	baseUrl := "https://api-public.guidebox.com/v1.43/US/%v/show/%v/episodes/all/%v/%v/all/all/true"

	apiKey := "rKWvTOuKvqzFbORmekPyhkYMGinuxgxM"

	url := fmt.Sprintf(baseUrl, apiKey, guideboxId, start, chunk)
	fmt.Println(url)

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
