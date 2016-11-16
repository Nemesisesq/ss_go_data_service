package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
)

func Check(e error) {
	if e != nil {
		GetLogger().Debug(e)
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

func HomePage(w http.ResponseWriter, r *http.Request) {
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

func GetLogger() *log.Logger {
	log.SetFormatter(&log.JSONFormatter{})
	return log.New()
}

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

//Returns true if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

//Returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

//Returns true if all of the strings in the slice satisfy the predicate f.
func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

//Returns a new slice containing all strings in the slice that satisfy the predicate f.
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

//Returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
