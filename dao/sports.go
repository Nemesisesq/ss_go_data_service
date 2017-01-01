package dao

import (
	"os"
	"github.com/nemesisesq/ss_data_service/common"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	"github.com/nemesisesq/ss_data_service/gracenote"
	"gopkg.in/mgo.v2/bson"
	log"github.com/Sirupsen/logrus"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
)

type Sport gracenote.Sport

func GetSport(p map[string]interface{}) []Sport {
	driver := bolt.NewDriver()
	log.Info(os.Getenv("NEO4JBOLT"))

	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)

	cypher_query := `MATCH (s:Sport) RETURN s`
	params := map[string]interface{}{}

	if id, ok := p["id"].(string); ok {
		cypher_query = `MATCH (s:Sport {gracenote_sport_id:{id}) RETURN s`
		params["id"] = id
	}

	//stmt, err := conn.PrepareNeo(cypher_query)

	data, metadata := QueryNeo(cypher_query, params)

	log.Info(metadata)

	sports := []Sport{}
	for _, val := range data {

		s := &Sport{}
		the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
		common.Check(err)
		err = bson.Unmarshal(the_bson, &s)
		common.Check(err)
		sports = append(sports, *s)

	}
	return sports
}
