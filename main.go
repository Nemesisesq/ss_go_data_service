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
	"github.com/rs/cors"
	com "github.com/nemesisesq/ss_data_service/common"
	pop "github.com/nemesisesq/ss_data_service/popularity"
	edr "github.com/nemesisesq/ss_data_service/email_data_service"
	serv_proc"github.com/nemesisesq/ss_data_service/service_processor"
	gnote "github.com/nemesisesq/ss_data_service/gracenote"
	ss "github.com/nemesisesq/ss_data_service/streamsavvy"
)

func main() {

	//Handle port environment variables for local and remote

	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		port = os.Getenv("PORT")
	}
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Use the MongoDB `DATABASE_URL` from the env
	dbURL := os.Getenv("MONGODB_URI")
	fmt.Println(dbURL)
	// Use the MongoDB `DATABASE_NAME` from the env
	dbName := com.GetDatabase()
	fmt.Println(dbName)
	// Set the MongoDB collection name
	dbColl := com.GetCollection()

	fmt.Println("Connecting to MongoDB: ", dbURL)
	fmt.Println("Database Name: ", dbName)
	fmt.Println("Collection Name: ", dbColl)

	// Creating the database accessor here.
	// Pointer to this database accessor will be passed to the middleware.
	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)

	com.Check(err)

	n := negroni.Classic()
	n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())

	r := mux.NewRouter()
	r.HandleFunc("/", com.Index).Methods("GET")
	r.HandleFunc("/data", edr.EmailDataHandler).Methods("POST")
	r.HandleFunc("/update", pop.UpdatePopularShows).Methods("GET")
	r.HandleFunc("/popular", pop.GetPopularityScore).Methods("POST")
	r.HandleFunc("/live-streaming-service", serv_proc.GetLiveStreamingServices).Methods("POST")
	r.HandleFunc("/on-demand-streaming-service", serv_proc.GetOnDemandServices).Methods("POST")
	r.HandleFunc("/gracenote/lineup-airings/{lat}/{long}", gnote.GetLineupAirings)
	r.HandleFunc("/favorites/test", ss.GetTestFavorites)
	r.HandleFunc("/favorites/add/test", ss.AddContentToTestFavorites)
	//r.HandleFunc("/test/{email}", testHandler).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	n.Use(c)
	n.UseHandler(r)

	fmt.Println(fmt.Sprintf("listening on port :%s", port))
	log.Fatal(http.ListenAndServe(":"+port, n))

}
