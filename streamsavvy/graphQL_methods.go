package streamsavvy

import (
	"os"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/compose/transporter/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/nemesisesq/ss_data_service/common"
	"github.com/nemesisesq/ss_data_service/dao"
	dbase "github.com/nemesisesq/ss_data_service/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ProcessTeamForFavorites(userId, email, team_brand_id, label string, favorite bool) bool {

	driver := bolt.NewDriver()
	bolt_url := os.Getenv("NEO4JBOLT")
	log.Info(bolt_url)

	conn, err := driver.OpenNeo(bolt_url)
	defer conn.Close()
	common.Check(err)

	var cypher_query string

	if favorite {

		cypher_query = `
			MATCH (t)
			WHERE t.team_brand_id={team_brand_id}
			MERGE (u:User {email: {email}, user_id: {user_id}})
			MERGE (u)-[:LIKES]->(t)

			`
	} else {

		cypher_query = `
			MATCH (t)
			WHERE t.team_brand_id={team_brand_id}
			MERGE (u:User {email: {email}, user_id: {user_id}})
			MERGE (u)-[rel:LIKES]->(t)
			DELETE rel

			`
	}
	params := map[string]interface{}{"email": email, "user_id": userId, "team_brand_id": team_brand_id}

	//stmt, err := conn.PrepareNeo(cypher_query)


	dao.ProcessCypher(conn, cypher_query, params)
	if favorite {
		logrus.WithField("team", team_brand_id).Info("was favortied")
	} else {

		logrus.WithField("team", team_brand_id).Info("was un-favortied")
	}

	return favorite


}

func ProcessShowForFavorites(userId, email string, guidebox_id int, favorite bool) bool {
	favorites := &Favorites{}

	var _sesh, err = mgo.Dial(os.Getenv("MONGODB_URI"))

	common.Check(err)
	dbname := dbase.GetDatabase()
	db := _sesh.DB(dbname)
	contentDb := _sesh.DB("content")

	collection := *db.C("favorites")

	err = collection.Find(bson.M{"user.user_id": userId}).One(&favorites)

	if err != nil {
		user := &User{Email: email, UserId: userId}
		favorites.User = *user

	}

	content := &Content{}

	err = contentDb.C("shows").Find(bson.M{"guidebox_data.id": guidebox_id}).One(&content)

	common.Check(err)

	collection.Find(bson.M{"user.user_id": userId}).One(&favorites)

	driver := bolt.NewDriver()
	bolt_url := os.Getenv("NEO4JBOLT")
	log.Info(bolt_url)

	conn, err := driver.OpenNeo(bolt_url)
	defer conn.Close()
	common.Check(err)

	var cypher_query string
	if favorite {

		cypher_query = `
		MATCH (c:Content)<-[d:DETAIL]-()
			 WHERE d.mongo_id = {id}
		MERGE (u:User {email: {email}, user_id: {user_id}})

			 MERGE (u)-[:LIKES]->(c)
			`

	} else {

		cypher_query = `MATCH (u:User {email: {email}, user_id: {user_id}})
			 MATCH (c:Content)<-[d:DETAIL]-()
			 WHERE d.mongo_id = {id}
			 MATCH (u)-[rel:LIKES]->(c)
			 DELETE rel

			`

	}

	params := map[string]interface{}{"email": email, "user_id": userId, "id": fmt.Sprintf(`%x`, string(content.Id))}

	dao.ProcessCypher(conn, cypher_query, params)

	if favorite {
		logrus.WithField("show", content.Title).Info("was favortied")
	} else {

		logrus.WithField("show", content.Title).Info("was un-favortied")
	}

	return favorite
}

func ProcessSportForFavorites(userId, email  string ,sportId int, fav bool) bool {

	//favorites := &Favorites{}

	//var _sesh, err = mgo.Dial(os.Getenv("MONGODB_URI"))

	//common.Check(err)
	//_sesh := GetMongoSession()
	//dbname := dbase.GetDatabase()
	//db := _sesh.DB(dbname)

	//collection := *db.C("favorites")

	//err = collection.Find(bson.M{"user.user_id": userId}).One(&favorites)

	//if err != nil {
	//	user := &User{Email: email, UserId: userId}
	//	favorites.User = *user
	//
	//}

	driver := bolt.NewDriver()
	bolt_url := os.Getenv("NEO4JBOLT")
	log.Info(bolt_url)

	conn, err := driver.OpenNeo(bolt_url)
	defer conn.Close()
	common.Check(err)

	var cypher_query string
	if fav {

		cypher_query = `
			 MATCH (s:Sport {gracenote_sport_id:{sportsId}})
		     	 MERGE (u:User {email: {email}, user_id: {user_id}})
			 MERGE (u)-[:LIKES]->(s)

			`
	} else {

		cypher_query = `
			 MATCH (s:Sport {gracenote_sport_id:{sportsId}})
			 MERGE (u:User {email: {email}, user_id: {user_id}})
			 MERGE (u)-[rel:LIKES]->(s)
			 DELETE rel

			`
	}
	params := map[string]interface{}{"email": email, "user_id": userId, "sportsId": fmt.Sprint(sportId)}

	//stmt, err := conn.PrepareNeo(cypher_query)

	dao.ProcessCypher(conn, cypher_query, params)

	if fav {
		logrus.WithField("sport", sportId).Info("was favortied")
	} else {

		logrus.WithField("sport", sportId).Info("was un-favortied")
	}

	return fav
}
