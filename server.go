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
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r := mux.NewRouter()
    r.HandleFunc("/", index).Methods("GET")
    r.HandleFunc("/test/", testHandler).Methods("POST")
    r.HandleFunc("/test/{email}", testHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":" + port, r))

    fmt.Println("application is running @ http://localhost:8080")

}


//// DATABASE Mongo connection


type EDRecord struct {
    Email    string `json: "email"`
    DeviceID []string `json : "deviceId"`
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    mongo_uri := os.Getenv("MONGODB_URI")

    session, err := mgo.Dial(mongo_uri)

    if err != nil {
        panic(err)
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

        //m := r.URL.Path[len("/test/"):]

        //fmt.Fprint(w, m)

        //fmt.Println(m)

        //if m == nil {
        //    http.NotFound(w, r)
        //}

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

        fmt.Println(r.Body)

        var t EDRecord

        err = decoder.Decode(&t)

        if err != nil {
            log.Fatal(err)
            fmt.Print(err)
        }

        fmt.Println(t)

        err = c.Insert(t)

        if err != nil {
            log.Fatal(err)
            fmt.Print(err)
        }

        fmt.Fprint(w, "OKAY!")
    }

}
