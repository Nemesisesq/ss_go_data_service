package main

import (
	"fmt"
	//"log"
	"net/http"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/joho/godotenv"
	com "github.com/nemesisesq/ss_data_service/common"
	dbase "github.com/nemesisesq/ss_data_service/database"
	edr "github.com/nemesisesq/ss_data_service/email_data_service"
	gnote "github.com/nemesisesq/ss_data_service/gracenote"
	"github.com/nemesisesq/ss_data_service/middleware"
	pop "github.com/nemesisesq/ss_data_service/popularity"
	serv_proc "github.com/nemesisesq/ss_data_service/service_processor"
	"github.com/nemesisesq/ss_data_service/socket"
	ss "github.com/nemesisesq/ss_data_service/streamsavvy"
	"github.com/nemesisesq/ss_data_service/timers"
	"github.com/newrelic/go-agent"
	"github.com/rs/cors"
)

func main() {
	//configure new relic
	config := newrelic.NewConfig("Your App Name", "baa40a4680d3d03079bb6f7bfbc9130934bf33e0")
	app, err := newrelic.NewApplication(config)

	com.Check(err)

	log.SetFormatter(&log.JSONFormatter{})
	//Handle port environment variables for local and remote

	err = godotenv.Load()

	com.Check(err)

	port := com.GetPort()

	// Create Redis Client
	redis_url := os.Getenv("REDISCLOUD_URL")

	u, err := url.Parse(redis_url)

	com.Check(err)

	pass, b := u.User.Password()

	if !b {
		pass = ""
	}

	com.Check(err)

	n := negroni.Classic()

	dbAccessor := dbase.DBStartup()
	n.Use(nigronimgosession.NewDatabase(dbAccessor).Middleware())

	cacheAccessor, err := middleware.NewCacheAccessor(u.Host, pass, 0)
	n.Use(middleware.NewRedisClient(*cacheAccessor).Middleware())

	messengerAccessor, err := middleware.NewRabbitMQAccesor(os.Getenv("RABBITMQ_BIGWIG_TX_URL"), os.Getenv("RABBITMQ_BIGWIG_RX_URL"))
	n.Use(middleware.NewRabbitMQConnection(*messengerAccessor).Middleware())

	r := mux.NewRouter()

	r.HandleFunc("/echo", socket.EchoHandler)
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/epis", ss.HandleEpisodeSocket))

	r.HandleFunc(newrelic.WrapHandleFunc(app, "/popular", pop.GetPopularityScore))
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/episodes", ss.GetEpisodes)).Methods("GET")
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/live-streaming-service", serv_proc.GetLiveStreamingServices))
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/on-demand-streaming-service", serv_proc.GetOnDemandServices))
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/gracenote/lineup-airings/{lat}/{long}", gnote.GetLineupAirings))

	r.HandleFunc("/data", edr.EmailDataHandler).Methods("POST")
	r.HandleFunc("/update", pop.UpdatePopularShows).Methods("GET")
	r.HandleFunc("/favorites", ss.GetFavorites)
	r.HandleFunc("/favorites/add", ss.AddContentToFavorites)
	r.HandleFunc("/favorites/remove", ss.RemoveContentFromFavorites).Methods("DELETE")
	r.HandleFunc("/favorites/delete_all/test", ss.DeleteTestFavorites).Methods("DELETE")
	r.HandleFunc("/fff", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "1")
	})
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	//r.HandleFunc("/stop-ticker", func(w http.ResponseWriter, r *http.Request) {close(quit)})
	//r.HandleFunc("/test/{email}", testHandler).Methods("GET")

	//Socket

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	n.Use(c)
	n.UseHandler(r)

	//timers

	if os.Getenv("DEBUG") != "true" {

		timers.GraceNoteListingTimer()
		timers.GuideboxEpisodeTimer()
		timers.PopularityTimer()
	}


	//
	//ticker := time.NewTicker(25 * time.Minute)
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			log.Println("ticker fired")
	//			gnote.RefreshListings()
	//		case <-quit:
	//			log.Println("ticker Stoping")
	//			ticker.Stop()
	//			return
	//		}
	//	}
	//
	//	log.Println("Cleaning up!!")
	//}()

	fmt.Println(fmt.Sprintf("listening on port :%s", port))
	log.Fatal(http.ListenAndServe(":"+port, n))

}
