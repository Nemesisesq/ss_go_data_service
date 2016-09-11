package popularity

import (
    "encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/context"
	com "github.com/nemesisesq/ss_data_service/common"
)

type Results struct {
	Results []PopShow `bson:"result"`
}

type ShowName struct {
	Name string `json:"name"`
}

type PopShow struct {
	PosterPath         string        `json:"poster_path"`
	Popularity         float32       `json:"popularity"`
	TheMovieDatabaseId int32         `json:"id"`
	BackdropPath       string        `json:"backdrop_path"`
	VoteAverage        float32       `json:"vote_average"`
	Overview           string        `json:"overview"`
	FirstAirDate       string        `json:"first_air_date"`
	OriginCountry      []interface{} `json:"origin_country"`
	GenreIds           []interface{} `json:"genere_ids"`
	OriginalLanguge    string        `json:"original_language"`
	VoteCount          int           `json:"vote_count"`
	Name               string        `json:"name"`
	OriginalName       string        `json:"original_name"`
}

func UpdatePopularShows(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "db").(*mgo.Database)

	col := com.GetCollection()

	c := db.C(col)

	//time.Sleep(1 * time.Minute)

	for i := 1; i <= 1000; i++ {
		GetPopularShows(i, "1995-01-01", c)

	}
}

func GetPopularShows(page int, air_date string, c *mgo.Collection) {
	url := "http://api.themoviedb.org/3/discover/tv?api_key=186e0e756acb157c80d75708e227cf25&sort_by=popularity.desc&page=%d&first_air_date.gte=%s"

	url = fmt.Sprintf(url, page, air_date)

	res, err := http.Get(url)

	defer res.Body.Close()

	if err != nil {
		log.Panic(err)
	}
	decoder := json.NewDecoder(res.Body)

	t := &Results{}

	err = decoder.Decode(&t)

	if err != nil {
		log.Panic(err)
	}

	for _, elem := range t.Results {

		err = c.Insert(elem)

		if err != nil {
			print("I'm here")
			log.Panic(err)
		}

		fmt.Println(fmt.Sprintf("%s saved \n", elem.Name))
	}

	time.Sleep(2500 * time.Millisecond)

}



func GetPopularityScore(w http.ResponseWriter, r *http.Request) {

	t := &ShowName{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(t)

	if err != nil {
		fmt.Println(err)
	}

	//db := GetDatabase()

	db := context.Get(r, "db").(*mgo.Database)

	col := com.GetCollection()
	//defer session.Close()

	c := db.C(col)
	show := &PopShow{}

	err = c.Find(t).One(&show)

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(show)

}
