package streamsavvy

import (
	"fmt"
	"net/http"
	//"github.com/gorilla/context"
	"encoding/json"

	"net/url"
	"os"
	"time"

	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	com "github.com/nemesisesq/ss_data_service/common"
	"github.com/nemesisesq/ss_data_service/middleware"
	"github.com/streadway/amqp"
	"gopkg.in/redis.v5"
	//"github.com/aws/aws-sdk-go/aws/client"
)

type Episode struct {
	Id                         interface{}              `json:"id"`
	Tvdb                       interface{}              `json:"tvdb"`
	ContentType                interface{}              `json:"content_type"`
	IsShadow                   interface{}              `json:"is_shadow"`
	AlternateTvdb              []interface{}            `json:"alternate_tvdb"`
	ImdbId                     interface{}              `json:"imdb_id"`
	SeasonNumber               interface{}              `json:"season_number"`
	EpisodeNumber              interface{}              `json:"episode_number"`
	ShowId                     interface{}              `json:"show_id"`
	Themoviedb                 interface{}              `json:"themoviedb"`
	Special                    interface{}              `json:"special"`
	FirstAired                 string                   `json:"first_aired"`
	Title                      string                   `json:"title"`
	OriginalTitle              string                   `json:"original_title"`
	AlternateTitles            []string                 `json:"alternate_titles"`
	Overview                   string                   `json:"overview"`
	Duration                   interface{}              `json:"duration"`
	ProductionCode             interface{}              `json:"production_code"`
	Thumbnail208X117           string                   `json:"thumbnail_208x117"`
	Thumbnail304X171           string                   `json:"thumbnail_304x171"`
	Thumbnail400X225           string                   `json:"thumbnail_400x225"`
	Thumbnail608X342           string                   `json:"thumbnail_608x342"`
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleEpisodeSocket(w http.ResponseWriter, r *http.Request) {

	//timeout := time.NewTicker(10 * time.Second)
	conn, err := upgrader.Upgrade(w, r, nil)
	com.Check(err)

	//epiChan := make(chan []interface{}, 100)

	client := r.Context().Value("redis_client").(*redis.Client)

	rmqc := r.Context().Value("rabbitmq").(middleware.RMQCH)

	//timeout := time.NewTicker(20 * time.Minute)
	//	stop := make(chan bool)

	//go func(){
	//	select {
	//	case <- timeout.C:
	//		conn.Close()
	//		//stop <- true
	//	}
	//}()

	//wg := sync.WaitGroup{}

	for {
		//timeout.Stop()
		//timeout = time.NewTicker(20 * time.Minute)

		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Error(err)
			conn.Close()
			return
		}
		com.Check(err)

		guideboxId := string(p[:])

		epi := &GuideBoxEpisodes{}

		val, err := client.Get(guideboxId).Result()

		if err == redis.Nil || val == "" {

			rx_q, err := rmqc.RX.QueueDeclare(
				"episodes",
				false,
				false,
				false,
				false,
				nil,
			)
			com.Check(err)

			/*
				Here we set the time out to close the socket connection after there is no more
				Episodes to send
			*/
			go func() {

				msgs, err := rmqc.RX.Consume(
					rx_q.Name, // queue
					"", // consumer
					true, // auto-ack
					false, // exclusive
					false, // no-local
					false, // no-wait
					nil, // args
				)

				com.Check(err)
				x := 1
				for {

					select {
					case d := <-msgs:
						if d.Body != nil {

							//log.Info(string(d.Body[:]))
							log.Info("sending to client ", x * 12)
							err = conn.WriteMessage(messageType, d.Body)
							x += 1
							//if err != nil {
							//	conn.Close()
							//	stop <- true
							//	return
							//}
						}
					//case <-stop:
					//	conn.Close()
					//	return
					//
					}
				}
			}()

			log.Info("Getting Initial")
			//wg.Add(1)
			total_results, episode_list := epi.GetEpisodes(0, 12, guideboxId)
			log.Info("Got Initial")
			tx_q, err := rmqc.TX.QueueDeclare(
				"episodes",
				false,
				false,
				false,
				false,
				nil,
			)
			com.Check(err)

			for i := 1; (i * 12) <= total_results; i += 1 {
				s := i * 12

				log.Print("############## %v #################", s)
				//wg.Add(1)
				go func(s int, guideboxId string) {
					start := time.Now()
					log.Debug("sending ")
					_, res := epi.GetEpisodes(s, 12, guideboxId)

					response, err := json.Marshal(res)
					com.Check(err)

					err = rmqc.TX.Publish(
						"", // exchange
						tx_q.Name, // routing key
						false, // mandatory
						false, // immediate
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        response,
						})

					com.Check(err)

					//err = conn.WriteMessage(messageType, response)
					//wg.Done()
					//timeout = time.NewTicker(1 * time.Second)

					log.WithFields(log.Fields{

						"Starting At": s,
						"start time":  start,
					}).Info(time.Since(start))
				}(s, guideboxId)

				// Be careful with select statements with out defaults

			}

			response, err := json.Marshal(episode_list)

			com.Check(err)
			err = conn.WriteMessage(messageType, response)
			com.Check(err)
			//wg.Done()

		} else {

			log.Info(fmt.Sprintf("%v found in cache", guideboxId))
			json.Unmarshal([]byte(val), &epi)
			response, err := json.Marshal(epi.Results)
			com.Check(err)
			err = conn.WriteMessage(messageType, response)

		}
	}
}

func GetEpisodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	guideboxId := r.URL.Query().Get("guidebox_id")

	epi := &GuideBoxEpisodes{}

	client := r.Context().Value("redis_client").(*redis.Client)

	//val, _ := client.Get(guideboxId).Result()
	//ttl, _ := client.TTL(guideboxId).Result()

	//log.Info(fmt.Sprintf("the redis error is %v", err))
	//log.Info(fmt.Sprintf("the value is %v", val))

	if true {

		log.Info(fmt.Sprintf("Getting %v, not present in cache", guideboxId))

		episode_list, total_results := epi.GetAllEpisodes(guideboxId)

		log.Info(fmt.Sprintf("Got %v, not present in cache", guideboxId))

		epi.CacheEpisode(total_results, episode_list, guideboxId, *client)

		epi.Results = episode_list
		//log.Info(fmt.Sprintf("this is the value of epi %v", epi))

	} else {
		//log.Info("checking TTL", reflect.TypeOf(ttl))
		//if ttl < time.Hour*12 {
		//	log.Info(fmt.Sprintf("refreshing %v", guideboxId))
		//	go epi.RefreshEpisodes(guideboxId, *client)
		//}
		//
		//log.Info(fmt.Sprintf("%v found in cache", guideboxId))
		//json.Unmarshal([]byte(val), &epi)
	}

	epi.Results = CleanUpDeepLinks(epi.Results)

	json.NewEncoder(w).Encode(epi)

}

func CleanUpDeepLinks(epi_list []interface{}) []interface{} {
	//res := []Episode{}

	for indx, val := range epi_list {
		x_epi := &Episode{}

		the_json, err := json.Marshal(val)
		com.Check(err)
		err = json.Unmarshal(the_json, &x_epi)
		com.Check(err)
		for indx, val := range x_epi.SubscriptionIosSources {
			if val["source"].(string) == "hulu_with_showtime" {
				//log.WithField("length of sources before", len(x_epi.SubscriptionIosSources)).Info()
				x_epi.SubscriptionIosSources = append(x_epi.SubscriptionIosSources[:indx], x_epi.SubscriptionIosSources[indx + 1:]...)
				//log.WithField("length of sources after", len(x_epi.SubscriptionIosSources)).Info()
			}

		}

		the_json, _ = json.Marshal(x_epi)
		err = json.Unmarshal(the_json, &val)

		//if indx == 0{

		//pretty, _ := json.MarshalIndent(val,"", "\t")
		//fmt.Println(string(pretty[:]))
		//}
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

	err = client.Set(guideboxId, val, 0).Err()
	com.Check(err)

	json.Unmarshal([]byte(val), &epi)
}

func (epi GuideBoxEpisodes) RefreshEpisodes(guideboxId string, client redis.Client) {
	episode_list, total_results := epi.GetAllEpisodes(guideboxId)
	epi.CacheEpisode(total_results, episode_list, guideboxId, client)

}

func (epi GuideBoxEpisodes) GetAllEpisodes(guideboxId string) (episode_list []interface{}, total_results int) {
	//startingIndexList := []int{}
	//log.Info("Im in")

	total_results, episode_list = epi.GetEpisodes(0, 25, guideboxId)
	//log.Info("Im out")

	chanBuffer := total_results / 25
	episodeListChan := make(chan []interface{}, chanBuffer)

	wg := sync.WaitGroup{}

	for i := 1; (i * 25) <= total_results; i++ {
		s := i * 25
		//log.Printf("getting episodes starting with %v", s)
		wg.Add(1)
		go func(s int, guideboxId string, wg *sync.WaitGroup, c chan []interface{}) {
			_, res := epi.GetEpisodes(s, 25, guideboxId)

			//log.Printf("sending reuslts for %v to chan", s)
			c <- res
			wg.Done()
			//log.Println(len(res))
			//fmt.Println(wg)
		}(s, guideboxId, &wg, episodeListChan)
		//time.Sleep(25 * time.Millisecond)
	}

	//log.Info("waiting")
	wg.Wait()
	//log.Info("done waiting ")
	close(episodeListChan)

	for e := range episodeListChan {

		episode_list = append(episode_list, e...)

	}

	//
	//if total_results > len(episode_list) {
	//	go func() {
	//		for total_results > len(episode_list) {
	//			_, res := epi.GetEpisodes(len(episode_list), 25, guideboxId)
	//			episode_list = append(episode_list, res...)
	//		}
	//	}()
	//
	//}

	//print(total_results)
	epi.Results = episode_list

	return episode_list, total_results
}

func (gbe GuideBoxEpisodes) GetEpisodes(start int, chunk int, guideboxId string) (total_results int, epiList []interface{}) {

	//TODO Logging
	resChan := make(chan *http.Response)

	go func() {

		client := &http.Client{}

		baseUrl := "https://api-public.guidebox.com/v1.43/US/%v/show/%v/episodes/all/%v/%v/all/all/true"

		apiKey := "rKWvTOuKvqzFbORmekPyhkYMGinuxgxM"

		gbox_url := fmt.Sprintf(baseUrl, apiKey, guideboxId, start, chunk)

		fmt.Println(gbox_url)

		gbox_req, err := http.NewRequest("GET", gbox_url, nil)

		com.Check(err)

		params := map[string]string{
			"reverse_ordering": "true",
		}

		com.BuildQuery(gbox_req, params)

		res, err := client.Do(gbox_req)

		resChan <- res

	}()

	//com.Check(err)
	res := <- resChan
	decoder := json.NewDecoder(res.Body)

	err := decoder.Decode(&gbe)
	com.Check(err)

	total_results = gbe.TotalResults
	epiList = gbe.Results

	return total_results, epiList
}

func RefreshEpisodes() {
	popularShowList := []Content{}
	epis := &GuideBoxEpisodes{}

	redis_url := os.Getenv("REDIS_URL")

	u, err := url.Parse(redis_url)

	com.Check(err)

	pass, b := u.User.Password()

	if !b {
		pass = ""
	}
	//rURL := fmt.Sprintf("%v://%v",u.Scheme, u.Host)
	client := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: pass,
		DB:       0,
	})

	pong, err := client.Ping().Result()
	com.Check(err)

	fmt.Printf("redis %v", pong)

	popshows_url := fmt.Sprintf("%v/popular_shows", os.Getenv("SS_MASTER"))

	res, err := http.Get(popshows_url)

	com.Check(err)

	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&popularShowList)

	for _, show := range popularShowList {

		id := show.GuideboxData["id"].(string)

		log.Printf("getting episodes for &v", id)

		episode_list, total_results := epis.GetAllEpisodes(id)
		//
		epis.CacheEpisode(total_results, episode_list, id, *client)

	}

}
