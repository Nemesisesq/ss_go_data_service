package email_data_service

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

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
    Code    int
    Message string
}

func EmailDataHandler(w http.ResponseWriter, r *http.Request) {
    db := context.Get(r, "db").(*mgo.Database)

    c := db.C("ed_records")

    if r.Method == "GET" {

        vars := mux.Vars(r)

        email := vars["email"]

        result := &EDRecord{}
        err := c.Find(bson.M{"email": email}).One(result)

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

        err := decoder.Decode(&t)

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
