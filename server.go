package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "log"
    "os"
    //"encoding/json"
    "github.com/urfave/negroni"
    "github.com/rs/cors"
    //"golang.org/x/net/icmp"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
    })

    r := mux.NewRouter()
    r.HandleFunc("/", index).Methods("GET")
    r.HandleFunc("/data", testHandler).Methods("POST")
    //r.HandleFunc("/test/{email}", testHandler).Methods("GET")
    n := negroni.Classic()

    n.Use(c)
    n.UseHandler(r)

    log.Fatal(http.ListenAndServe(":" + port, n))

    fmt.Println("application is running @ http://localhost:8080")

}


//// DATABASE Mongo connection


type EDRecord struct {
    Email          string `json:"email" bson:"e"`
    Fingerprint    int    `json:"fingerprint"`
    Browser        string `json:"browser"`
    BrowserVersion string `json:"browserVersion"`
    Device         string `json:"device"`
    DeviceType     string `json:"deviceType"`
    DeviceVendor   string `json:"deviceVendor"`
    Time           int    `json:"time"`
    TimeZone       string `json:"timeZone"`
    Platform       string `json:"platform"`
    Package        map[string]interface{} `json:"package.data"`
}

type EDPriv struct {
    Password string `json:"password"`
}

type EDWhole struct {
    EDRecord `bson:"inline"`
    EDPriv `bson:"inline"`
}

type ResponseStatus struct {
    Code int
    Message string
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}

