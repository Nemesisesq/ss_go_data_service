package main

import (
	"fmt"
	//"log"
	"net/http"
	"os"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	com "github.com/nemesisesq/ss_data_service/common"
	dbase "github.com/nemesisesq/ss_data_service/database"
	edr "github.com/nemesisesq/ss_data_service/email_data_service"
	gnote "github.com/nemesisesq/ss_data_service/gracenote"
	"github.com/nemesisesq/ss_data_service/middleware"
	pop "github.com/nemesisesq/ss_data_service/popularity"
	serv_proc "github.com/nemesisesq/ss_data_service/service_processor"
	"github.com/nemesisesq/ss_data_service/socket"
	ss "github.com/nemesisesq/ss_data_service/streamsavvy"
	//"github.com/nemesisesq/ss_data_service/timers"
	//"github.com/rs/cors"
	"net/url"
	"github.com/urfave/negroni"
	"github.com/nemesisesq/ss_data_service/timers"
)

func main() {

	//configure new relic

	log.SetFormatter(&log.JSONFormatter{})
	//Handle port environment variables for local and remote

	//err = godotenv.Load()

	//com.Check(err)

	port := com.GetPort()

	// Create Redis Client
	var redis_url string
	if os.Getenv("REDIS_URL") != "" {

		redis_url = os.Getenv("REDIS_URL")
	} else {
		redis_url = os.Getenv("R_PORT")
	}

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0], " : ", pair[1])
	}

	u, err := url.Parse(redis_url)

	log.Info("u here", u)
	com.Check(err)

	//pass, b := u.User.Password()
	//
	//if !b {
	//	pass = ""
	//}

	com.Check(err)

	n := negroni.New()

	n.Use(negroni.NewLogger())

	dbAccessor := dbase.DBStartup()
	n.Use(nigronimgosession.NewDatabase(dbAccessor).Middleware())

	cacheAccessor, err := middleware.NewCacheAccessor(u.Host, "", 0)
	com.Check(err)
	n.Use(middleware.NewRedisClient(*cacheAccessor).Middleware())

	//TODO fix these urls for AWS ElasticBeanStalk
	var tx_url string
	var rx_url string
	if os.Getenv("RABBITMQ_URL") != "" {
		log.Info("*$*$*$*$*$*$*$*$*")
		tx_url = os.Getenv("RABBITMQ_URL")
		rx_url = os.Getenv("RABBITMQ_URL")
	} else {
		rx_url = fmt.Sprintf("amqp://%v:%v", os.Getenv("RABBITMQ_1_PORT_5672_TCP_ADDR"), os.Getenv("RABBITMQ_1_PORT_5672_TCP_PORT"))
		tx_url = fmt.Sprintf("amqp://%v:%v", os.Getenv("RABBITMQ_1_PORT_5671_TCP_ADDR"), os.Getenv("RABBITMQ_1_PORT_5672_TCP_PORT"))

	}

	log.Info(tx_url, " ", rx_url)

	messengerAccessor, err := middleware.NewRabbitMQAccesor(tx_url, rx_url)
	n.Use(middleware.NewRabbitMQConnection(*messengerAccessor).Middleware())

	r := mux.NewRouter()

	socketRouter := mux.NewRouter().PathPrefix("/sock").Subrouter().StrictSlash(true)
	socketRouter.HandleFunc("/echo", socket.EchoHandler)
	socketRouter.HandleFunc("/epis", ss.HandleEpisodeSocket)
	socketRouter.HandleFunc("/reco", ss.HandleRecomendations)

	r.PathPrefix("/sock").Handler(negroni.New(
		nigronimgosession.NewDatabase(dbAccessor).Middleware(),
		middleware.NewRedisClient(*cacheAccessor).Middleware(),
		middleware.NewRabbitMQConnection(*messengerAccessor).Middleware(),
		middleware.CleanupMiddleware(),
		negroni.Wrap(socketRouter),
	))

	r.HandleFunc("/popular", pop.GetPopularityScore)
	r.HandleFunc("/episodes", ss.GetEpisodes).Methods("GET")
	r.HandleFunc("/live-streaming-service", serv_proc.GetLiveStreamingServices)
	r.HandleFunc("/on-demand-streaming-service", serv_proc.GetOnDemandServices)
	r.HandleFunc("/gracenote/lineup-airings/{lat}/{long}", gnote.GetLineupAirings)

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

	//c := cors.New(cors.Options{
	//	AllowedOrigins: []string{"*"},
	//})
	//
	//n.Use(c)
	n.UseHandler(r)

	//timers

	//if os.Getenv("DEBUG") != "true" {

		timers.GraceNoteListingTimer()
		timers.GuideboxEpisodeTimer()
		timers.PopularityTimer()
	//}

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
	log.Fatal(http.ListenAndServe(":" + port, n))

}
