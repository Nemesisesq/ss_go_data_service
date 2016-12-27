package streamsavvy

import (
	"encoding/json"
	"net/http"
	"os"

	com "github.com/nemesisesq/ss_data_service/common"
	"fmt"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")

	result := Search(query)

	encoder := json.NewEncoder(w)
	encoder.Encode(result)

}

func Search(query string) interface{} {

	client := &http.Client{}
	url := fmt.Sprintf("%v/search", os.Getenv("SS_DJANGO_DATA_SERVICE"))

	req, err := http.NewRequest("GET", url, nil)

	com.Check(err)

	params := map[string]string{
		"q": query,
	}

	com.BuildQuery(req, params)

	res, err := client.Do(req)

	com.Check(err)

	var temp interface{}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&temp)

	com.Check(err)
	return temp
}
