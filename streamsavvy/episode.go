package streamsavvy

import (
	"fmt"
	"net/http"
	//"github.com/gorilla/context"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	com "github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/redis.v5"
	"reflect"
	"time"
)

type Episode struct {
	Id                         interface{}         `json:"id"`
	Tvdb                       interface{}         `json:"tvdb"`
	ContentType                string      `json:"content_type"`
	IsShadow                   interface{}         `json:"is_shadow"`
	AlternateTvdb              []string    `json:"alternate_tvdb"`
	ImdbId                     string      `json:"imdb_id"`
	SeasonNumber               interface{}         `json:"season_number"`
	EpisodeNumber              interface{}         `json:"episode_number"`
	ShowId                     interface{}         `json:"show_id"`
	Themoviedb                 interface{}         `json:"themoviedb"`
	Special                    interface{}         `json:"special"`
	FirstAired                 string      `json:"first_aired"`
	Title                      string      `json:"title"`
	OriginalTitle              string      `json:"original_title"`
	AlternateTitles            []string    `json:"alternate_titles"`
	Overview                   string      `json:"overview"`
	Duration                   interface{}         `json:"duration"`
	ProductionCode             string      `json:"production_code"`
	Thumbnail208X117           string      `json:"thumbnail_208x117"`
	Thumbnail304X171           string      `json:"thumbnail_304x171"`
	Thumbnail400X225           string      `json:"thumbnail_400x225"`
	Thumbnail608X342           string      `json:"thumbnail_608x342"`
	FreeWebSources             []map[string]interface{} `json:"free_web_sources"`
	FreeIosSources             []map[string]interface{} `json:"free_ios_sources"`
	SubscriptionWebSources     []map[string]interface{} `json:"subscription_web_sources"`
	SubscriptionIosSources     []map[string]interface{} `json:"subscription_ios_sources"`
	SubscriptionAndroidSources []map[string]interface{} `json:"subscription_android_sources"`
	PurchaseWebSources         []map[string]interface{} `json:"purchase_web_sources"`
	PurchaseIosSources         []map[string]interface{} `json:"purchase_ios_sources"`
	PurchaseAndroidSources     []map[string]interface{} `json:"purchase_android_sources"`
}
type GuideBoxEpisodes struct {
	GuideboxId    string        `bson:"guidebox_id"`
	Results       []interface{} `json:"results"`
	TotalResults  int           `json:"total_results" bson:"total_results"`
	TotalReturned int           `json:"total_returned" bson:"total_returned"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func GetEpisodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	guideboxId := r.URL.Query().Get("guidebox_id")

	epi := &GuideBoxEpisodes{}

	client := r.Context().Value("redis_client").(*redis.Client)

	val, err := client.Get(guideboxId).Result()
	ttl, err := client.TTL(guideboxId).Result()

	//log.Info(fmt.Sprintf("the redis error is %v", err))
	//log.Info(fmt.Sprintf("the value is %v", val))

	if err == redis.Nil || len(val) == 0 {

		log.Info(fmt.Sprintf("Getting %v, not present in cache", guideboxId))

		episode_list, total_results := epi.GetAllEpisodes(guideboxId)


		epi.CacheEpisode(total_results, episode_list, guideboxId, *client)

		epi.Results = episode_list
		//log.Info(fmt.Sprintf("this is the value of epi %v", epi))

	} else {
		log.Info("checking TTL", reflect.TypeOf(ttl))
		if ttl < time.Hour * 12 {
			log.Info(fmt.Sprintf("refreshing %v", guideboxId))
			go epi.RefreshEpisodes(guideboxId, *client)
		}

		log.Info(fmt.Sprintf("%v found in cache", guideboxId))
		json.Unmarshal([]byte(val), &epi)
	}

	epi.Results = CleanUpDeepLinks(epi.Results)

	json.NewEncoder(w).Encode(epi)

}

func CleanUpDeepLinks(epi_list []interface{}) []interface{} {
	//res := []Episode{}

	for indx, val := range epi_list {
		x_epi := &Episode{}

		the_json, _ := json.Marshal(val)
		err := json.Unmarshal(the_json, &x_epi)
		com.Check(err)
		for indx, val := range x_epi.SubscriptionIosSources {
			if val["source"].(string) == "hulu_with_showtime" {
				log.WithField("length of sources before", len(x_epi.SubscriptionIosSources)).Info()
				x_epi.SubscriptionIosSources = append(x_epi.SubscriptionIosSources[:indx], x_epi.SubscriptionIosSources[indx + 1:]...)
				log.WithField("length of sources after", len(x_epi.SubscriptionIosSources)).Info()
			}

		}

		the_json, _ = json.Marshal(x_epi)
		err = json.Unmarshal(the_json, &val)

		if indx == 0{

			pretty, _ := json.MarshalIndent(val,"", "\t")
			fmt.Println(string(pretty[:]))
		}
		com.Check(err)

		epi_list[indx] = val


	}


	return epi_list
}

func (epi GuideBoxEpisodes) CacheEpisode(total_results int, episode_list []interface{}, guideboxId string, client redis.Client) {
	//print(total_results)
	epi.Results = episode_list

	//log.Info(fmt.Sprintf("length of episode list %v", len(episode_list)))


	epi.GuideboxId = guideboxId
	epi.TotalResults = total_results
	epi.TotalReturned = len(episode_list)

	val, err := json.Marshal(epi)
	com.Check(err)

	timeout := time.Hour * 24 * 3
	err = client.Set(guideboxId, val, timeout).Err()
	com.Check(err)

	json.Unmarshal([]byte(val), &epi)
}

func (epi GuideBoxEpisodes) RefreshEpisodes(guideboxId string, client redis.Client) {
	episode_list, total_results := epi.GetAllEpisodes(guideboxId)
	epi.CacheEpisode(total_results, episode_list, guideboxId, client)

}

func (epi GuideBoxEpisodes) GetAllEpisodes(guideboxId string) (episode_list []interface{}, total_results int) {

	total_results, episode_list = epi.GetEpisodes(0, 25, guideboxId)

	if total_results > len(episode_list) {
		for total_results > len(episode_list) {
			_, res := epi.GetEpisodes(len(episode_list), 25, guideboxId)
			episode_list = append(episode_list, res...)
		}

	}

	print(total_results)
	epi.Results = episode_list

	return episode_list, total_results
}

func (gbe GuideBoxEpisodes) GetEpisodes(start int, chunk int, guideboxId string) (total_results int, epiList []interface{}) {

	//TODO Logging

	client := &http.Client{}

	baseUrl := "https://api-public.guidebox.com/v1.43/US/%v/show/%v/episodes/all/%v/%v/all/all/true"

	apiKey := "rKWvTOuKvqzFbORmekPyhkYMGinuxgxM"

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
