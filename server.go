package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "log"
    "os"
)


func main () {
    port := os.Getenv("PORT")
    if port == "" {
    port = "8080"
    }


    r := mux.NewRouter()
    r.HandleFunc("/", index).Methods("GET")
    //http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":" + port, r))
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}
