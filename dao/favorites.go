package dao

import (
	"github.com/Sirupsen/logrus"
	"github.com/compose/transporter/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"github.com/nemesisesq/ss_data_service/common"
)


func GetFavorties(userId string) (favs []map[string]interface{}) {
	conn := conn()
	defer conn.Close()



	cypher_query := `MATCH (u:User {user_id: {user_id}})
			 MATCH (u)-[:LIKES]->(f)
			 RETURN f
			 `


	params := map[string]interface{}{"user_id": userId}

	d, m :=  QueryNeo(cypher_query, params)

	logrus.Info(m)

	for _, val := range d {
		temp := map[string]interface{}{}
		node := val[0].(graph.Node)
		the_bson, err := bson.Marshal(node.Properties)

		common.Check(err)

		bson.Unmarshal(the_bson, &temp)

		if node.Labels[0] == "Content" {
			temp["name"] = temp["title"]
		}

		switch node.Labels[0] {
		case "Content":
			temp["name"] = temp["title"]
		case "Sport":
			temp["name"] = temp["sport_name"]
		}

		favs = append(favs, temp)

	}


	 return favs
}
