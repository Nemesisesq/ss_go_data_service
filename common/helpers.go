package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	"github.com/joho/godotenv"
	"log"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func BuildQuery(r *http.Request, m map[string]string) {
	q := r.URL.Query()

	for key, val := range m {
		q.Add(key, val)
	}
	r.URL.RawQuery = q.Encode()
}


func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func ReadJSON(r *http.Request, p ...interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil || data == nil {
		return err
	}
	for _, v := range p {
		err := json.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}
	return nil
}

type ResponseStatus struct {
	Code    int
	Message string
}

func AnnounceMongoConnection(dbURL, dbName, dbColl string) {
	fmt.Println("Connecting to MongoDB: ", dbURL)
	fmt.Println("Database Name: ", dbName)
	fmt.Println("Collection Name: ", dbColl)
}

func GetPort() string {
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

	return port
}