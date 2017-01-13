package streamsavvy

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//New endpoint for graphQL following shows

func GetFavorites(w http.ResponseWriter, r *http.Request) {

	user := &User{}

	favorites := &Favorites{}

	query := r.URL.Query()

	user.Email = query["email"][0]
	user.UserName = query["name"][0]
	user.UserId = r.Header["User-Id"][0]

	db := r.Context().Value("db").(mgo.Database)

	collection := *db.C("favorites")

	collection.Find(bson.M{"user.user_id": user.UserId}).One(&favorites)

	json.NewEncoder(w).Encode(&favorites.ContentList)
}

func DeleteTestFavorites(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(mgo.Database)

	c := *db.C("favorites")

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

func RemoveContentFromFavorites(w http.ResponseWriter, r *http.Request) {
	content := &Content{}

	favorites := &Favorites{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(content)

	common.Check(err)

	userId := r.Header["User-Id"][0]

	db := r.Context().Value("db").(mgo.Database)
	c := *db.C("favorites")

	delQuery := c.Find(bson.M{"user.user_id": userId})

	delQuery.One(&favorites)

	newList := favorites.ContentList

	for idx, value := range favorites.ContentList {
		if value.Title == content.Title {
			if len(newList) == 1 {
				newList = []Content{}
			} else {

				newList = append(newList[:idx], newList[idx+1:]...)
			}
		}
	}

	favorites.ContentList = newList

	c.Update(delQuery, favorites)

	res, err := json.Marshal(&common.ResponseStatus{200, "Item deleted sucessfully!"})

	if err != nil {
		log.Fatal(err)
		//fmt.Println("an error happend here 2nd")
		fmt.Print(err)
		http.NotFound(w, r)
	}

	w.Write(res)

}

func AddContentToFavorites(w http.ResponseWriter, r *http.Request) {
	authUser := &User{}

	favorites := &Favorites{}

	content := &Content{}

	userId := r.Header["User-Id"][0]

	authUser.UserId = userId

	//authUser.UserName = "test"

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&content)

	common.Check(err)

	db := r.Context().Value("db").(mgo.Database)

	c := *db.C("favorites")

	query := c.Find(bson.M{"user.user_id": userId})

	if dbCount, _ := query.Count(); dbCount == 0 {
		favorites.User = *authUser

		favorites.ContentList = append(favorites.ContentList, *content)
		err = c.Insert(favorites)
		common.Check(err)

	} else {
		query.One(&favorites)
		favorites.ContentList = append(favorites.ContentList, *content)
		colQuery := bson.M{"user.user_id": userId}
		err = c.Update(colQuery, favorites)
		common.Check(err)
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
