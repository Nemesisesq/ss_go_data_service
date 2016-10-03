package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
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
