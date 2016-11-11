package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	ss "github.com/nemesisesq/ss_data_service/streamsavvy"
	"github.com/rs/cors"
	"net/url"
)

func main() {

	//Handle port environment variables for local and remote

	err := godotenv.Load()

	com.Check(err)

	port := com.GetPort()

	dbAccessor := dbase.DBStartup()

	// Create Redis Client
	redis_url := os.Getenv("REDISCLOUD_URL")

	u, err := url.Parse(redis_url)

	com.Check(err)

	pass, b := u.User.Password()

	if !b {
		pass = ""
	}
	//rURL := fmt.Sprintf("%v://%v",u.Scheme, u.Host)
	cacheAccessor, err := middleware.NewCacheAccessor(u.Host, pass, 0)
	com.Check(err)

	n := negroni.Classic()
	n.Use(nigronimgosession.NewDatabase(dbAccessor).Middleware())

	x := middleware.NewRedisClient(*cacheAccessor)
	n.Use(x.Middleware())

	r := mux.NewRouter()
	r.HandleFunc("/", com.Index).Methods("GET")
	r.HandleFunc("/data", edr.EmailDataHandler).Methods("POST")
	r.HandleFunc("/update", pop.UpdatePopularShows).Methods("GET")
	r.HandleFunc("/popular", pop.GetPopularityScore).Methods("POST")
	r.HandleFunc("/live-streaming-service", serv_proc.GetLiveStreamingServices).Methods("POST")
	r.HandleFunc("/on-demand-streaming-service", serv_proc.GetOnDemandServices).Methods("POST")
	r.HandleFunc("/gracenote/lineup-airings/{lat}/{long}", gnote.GetLineupAirings)
	r.HandleFunc("/favorites", ss.GetFavorites)
	r.HandleFunc("/favorites/add", ss.AddContentToFavorites)
	r.HandleFunc("/favorites/remove", ss.RemoveContentFromFavorites).Methods("DELETE")
	r.HandleFunc("/favorites/delete_all/test", ss.DeleteTestFavorites).Methods("DELETE")
	r.HandleFunc("/episodes", ss.GetEpisodes).Methods("GET")
	r.HandleFunc("/fff", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "1") })
	//r.HandleFunc("/test/{email}", testHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	n.Use(c)
	n.UseHandler(r)

	fmt.Println(fmt.Sprintf("listening on port :%s", port))
	log.Fatal(http.ListenAndServe(":"+port, n))

}
