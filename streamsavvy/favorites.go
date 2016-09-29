package streamsavvy

import (
	"net/http"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"github.com/nemesisesq/ss_data_service/common"
	"log"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func GetTestFavorites(w http.ResponseWriter, r *http.Request) {

	testUser := &User{}

	favorites := &Favorites{}

	testUser.UserName = "test"


	db := context.Get(r, "db").(*mgo.Database)

	collection := db.C("favorites")

	collection.Find(bson.M{"user_uuid" : 999 }).One(&favorites)

	json.NewEncoder(w).Encode(&favorites.ContentList)
}

func DeleteTestFavorites(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "db").(*mgo.Database)
	c := db.C("favorites")
	c.RemoveAll("")

	res, err := json.Marshal(&common.ResponseStatus{200, "All Test Data Deleted"})

	if err != nil {
		log.Fatal(err)
		//fmt.Println("an error happend here 2nd")
		fmt.Print(err)
		http.NotFound(w, r)
	}

	w.Write(res)



}

func AddContentToTestFavorites(w http.ResponseWriter, r *http.Request) {
 	testUser := &User{}

	favorites := &Favorites{}

	content := &Content{}

	testUser.UserName = "test"

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&content)

	common.Check(err)

	db := context.Get(r, "db").(*mgo.Database)

	c := db.C("favorites")

	query := c.Find(bson.M{"user_uuid" : 999})

	if dbCount, _ := query.Count(); dbCount == 0 {
		favorites.User = *testUser
		favorites.UserUUID = 999

		favorites.ContentList = append(favorites.ContentList, *content)
		err = c.Insert(favorites)

	} else {
		query.One(&favorites)
		favorites.ContentList = append(favorites.ContentList, *content)
		colQuery := bson.M{"user_uuid" : 999}
		err  = c.Update(colQuery, favorites)
	}

	res, err := json.Marshal(&common.ResponseStatus{200, "Data saved sucessfully!"})

	if err != nil {
		log.Fatal(err)
		//fmt.Println("an error happend here 2nd")
		fmt.Print(err)
		http.NotFound(w, r)
	}

	w.Write(res)
}