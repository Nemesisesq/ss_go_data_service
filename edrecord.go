package main

import (
  "log"
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "os"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

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
