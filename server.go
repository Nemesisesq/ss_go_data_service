package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "fmt"
    "log"
)


func main () {
    r := mux.NewRouter()
    r.HandleFunc("/", index).Methods("GET")
    //http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}
