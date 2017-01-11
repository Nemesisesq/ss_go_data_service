package main

import (
	"fmt"
	//"log"
	"net/http"
	"os"

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

	"github.com/graphql-go/graphql-go-handler"
	"github.com/nemesisesq/ss_data_service/graphqlApi"
	"github.com/nemesisesq/ss_data_service/strand"
	"github.com/nemesisesq/ss_data_service/timers"
	"github.com/urfave/negroni"
)

func main() {

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

	//for _, e := range os.Environ() {
	//	pair := strings.Split(e, "=")
	//	fmt.Println(pair[0], " : ", pair[1])
	//}

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

	//Auth0
	/* this passes linting*/

	//jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
	//	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	//		secret := os.Getenv("AUTH0_CLIENT_SECRET")
	//		if secret == "" {
	//			return nil, errors.New("AUTH0_CLIENT_SECRET is not set")
	//		}
	//		return secret, nil
	//	},
	//})

	//func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	//respondJson("All good. You only get this message if you're authenticated", w)
	//}

	//n.UseFunc(jwtMiddleware.HandlerWithNext)

	var url string
	if os.Getenv("RABBITMQ_URL") != "" {
		log.Info("*$*$*$*$*$*$*$*$*")
		//tx_url = os.Getenv("RABBITMQ_URL")
		//rx_url = os.Getenv("RABBITMQ_URL")
		url = os.Getenv("RABBITMQ_URL")
	} else {
		//rx_url = fmt.Sprintf("amqp://%v:%v", os.Getenv("RABBITMQ_1_PORT_5672_TCP_ADDR"), os.Getenv("RABBITMQ_1_PORT_5672_TCP_PORT"))
		//tx_url = fmt.Sprintf("amqp://%v:%v", os.Getenv("RABBITMQ_1_PORT_5671_TCP_ADDR"), os.Getenv("RABBITMQ_1_PORT_5672_TCP_PORT"))
		url = fmt.Sprintf("amqp://%v:%v", os.Getenv("RABBITMQ_1_PORT_5671_TCP_ADDR"), os.Getenv("RABBITMQ_1_PORT_5672_TCP_PORT"))

	}

	log.Info(url)

	messengerAccessor, err := middleware.NewRabbitMQAccesor(url)
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

	/*GraphQL Endpoint*/

	h := handler.New(&handler.Config{
		Schema: graphqlApi.Schema(),
		Pretty: true,
	})
	r.Handle("/graphql", h)

	r.HandleFunc("/search", ss.SearchHandler).Methods("GET")
	r.HandleFunc("/strand", strand.ServePage)

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
	if os.Getenv("DEV") == "True" {

	} else {

		timers.GraceNoteListingTimer()
		timers.GuideboxEpisodeTimer()
		timers.PopularityTimer()
	}
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
	log.Fatal(http.ListenAndServe(":"+port, n))

}
