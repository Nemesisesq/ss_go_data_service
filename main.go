package main

import (
    "fmt"
    "net/http"
    "os"
    "log"

    "github.com/codegangsta/negroni"
    //"github.com/gorilla/context"
    nigronimgosession "github.com/joeljames/nigroni-mgo-session"
    //mgo "gopkg.in/mgo.v2"
    //"github.com/joho/godotenv"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "github.com/joho/godotenv"
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
	dbName := GetDatabase()
    fmt.Println(dbName)
	// Set the MongoDB collection name
	dbColl := GetCollection()

    fmt.Println("Connecting to MongoDB: ", dbURL)
    fmt.Println("Database Name: ", dbName)
    fmt.Println("Collection Name: ", dbColl)

    // Creating the database accessor here.
    // Pointer to this database accessor will be passed to the middleware.
    dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)

    check(err)

    n := negroni.Classic()
    n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())

    r := mux.NewRouter()
    r.HandleFunc("/", index).Methods("GET")
    r.HandleFunc("/data", emailDataHandler).Methods("POST")
    r.HandleFunc("/update", UpdatePopularShows).Methods("GET")
    r.HandleFunc("/popular", GetPopularityScore).Methods("POST")
    r.HandleFunc("/live-streaming-service", GetLiveStreamingServices).Methods("POST")
    //r.HandleFunc("/test/{email}", testHandler).Methods("GET")


    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
    })

    n.Use(c)
    n.UseHandler(r)

    fmt.Println(fmt.Sprintf("listening on port :%s", port))
    log.Fatal(http.ListenAndServe(":" + port, n))

}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}

