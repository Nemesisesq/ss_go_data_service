package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "log"
    "os"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    //"encoding/json"
    "encoding/json"
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
    Email          string `json: "email" bson:"e"`
    Fingerprint    int    `json : "fingerprint"`
    Browser        string `json : browser`
    BrowserVersion string `json : browserVersion`
    Device         string `json: device`
    DeviceType     string `json: deviceType`
    DeviceVendor   string `json: deviceVendor`
    Time           int    `json: time`
    TimeZone       string `json: timeZone`
    Platform       string `json: platform`
    Package        map[string]interface{} `json: package.data`
}

type EDPriv struct {
    Password string `json: "password"`
}

type EDWhole struct {
    EDRecord `, bson:"inline"`
    EDPriv `, bson:"inline`
}

type ResponseStatus struct {
    Code int
    Message string
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    mongo_uri := os.Getenv("MONGODB_URI")

    session, err := mgo.Dial(mongo_uri)

    if err != nil {
        //panic(err)
    }

    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    db_name := ""

    if mongo_uri != "" {
        db_name = "heroku_8c97bzpr"
    } else {
        db_name = "test"
    }

    c := session.DB(db_name).C("ed_records")

    if r.Method == "GET" {

        vars := mux.Vars(r)

        email := vars["email"]

        result := &EDRecord{}
        err = c.Find(bson.M{"email": email}).One(result)

        if err != nil {
            //log.Fatal(err)
            http.NotFound(w, r)
        }

        res, err := json.Marshal(result)

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            //return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(res)

    }

    if r.Method == "POST" {

        decoder := json.NewDecoder(r.Body)

        //fmt.Println(r.Body)

        var t EDRecord

        err = decoder.Decode(&t)

        if err != nil {
            log.Fatal(err)
            //fmt.Println("an error happend here 1st")
            fmt.Print(err)
            http.NotFound(w, r)
        }

        //fmt.Println(string(t.Package))

        err = c.Insert(t)

        if err != nil {
            log.Fatal(err)
            //fmt.Println("an error happend here 2nd")
            fmt.Print(err)
            http.NotFound(w, r)
        }

        res, err := json.Marshal(&ResponseStatus{200, "Data saved sucessfully!"})

        if err != nil {
            log.Fatal(err)
            //fmt.Println("an error happend here 2nd")
            fmt.Print(err)
            http.NotFound(w, r)
        }

        w.Write(res)
    }

}
