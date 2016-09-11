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
