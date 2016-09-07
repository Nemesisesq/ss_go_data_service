package main

import (
    "net/http"
    "encoding/json"
    "os"
    "fmt"
    "bytes"
    "io/ioutil"
)

type PayloadDict struct {
    CallLetters    string `json:"CallLetters"`
    DisplayName    string `json:"DisplayName"`
    SourceLongName string `json:"SourceLongName"`
    SourceId       string `json:"SourceId"`
}


type PaylooadTarget struct {

}

func GetLiveStreamingServices(w http.ResponseWriter, r *http.Request) {

    pyld := &PayloadDict{}

    decoder := json.NewDecoder(r.Body)
    err:= decoder.Decode(&pyld)

    check(err)

    url := fmt.Sprintf("%s/guide",os.Getenv("NODE_DATA_SERVICE"))

    buf:= new(bytes.Buffer)

    json.NewEncoder(buf).Encode(pyld)

    res, err := http.Post(url, "application/json", buf)

    check(err)

    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)

    w.Header().Set("Content-Type", "applicaation/json")
    fmt.Fprintf(w, string(body))
}
