package dao

import (
	"os"
	"github.com/nemesisesq/ss_data_service/common"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	"github.com/nemesisesq/ss_data_service/gracenote"
	"gopkg.in/mgo.v2/bson"
	log"github.com/Sirupsen/logrus"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"fmt"
)

type Sport gracenote.Sport
type Team map[string]interface{}
type CollegeTeam map[string]interface{}

type NeoData [][]interface{}

func (t Team) isATeam() string {
	return "pro"
}
func (c CollegeTeam) isATeam() string {
	return "amature"
}

type SportsTeams []map[string]interface{}


var bolt_url = os.Getenv("NEO4JBOLT")

func GetOrgs(sportId string) []map[string]interface{}{
	conn := conn()
	defer conn.Close()

	cypher_query := `MATCH (o)-[:REPRESENTS]->(s)
			WHERE s.gracenote_sport_id = {id}
	 		RETURN o
	 		`
	params := map[string]interface{}{"id":sportId}

	data, metadata := QueryNeo(cypher_query, params)
	log.Info(metadata)

	orgs := []map[string]interface{}{}
	for _, val := range data {

		o := &map[string]interface{}{}
		the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
		common.Check(err)
		err = bson.Unmarshal(the_bson, &o)
		common.Check(err)
		orgs = append(orgs, *o)

	}
	return orgs
}

func GetTeams(sportId string) SportsTeams {
	conn := conn()
	defer conn.Close()
	pro_cypher_query := getquery("Team")
	college_cypher_query := getquery("CollegeTeam")

	//TODO if logic should be implemented at soem point in time
	qs := []string{pro_cypher_query, college_cypher_query}

	params := map[string]interface{}{"id":sportId}

	data_slice := []NeoData{}

	for _, r := range qs {

		data, metadata := QueryNeo(r, params)
		log.Info(metadata)

		data_slice = append(data_slice, data)
	}

	sportTeams := SportsTeams{}
	for indx, data := range data_slice {
		for _, val := range data {

			if indx == 1 {
				tc := CollegeTeam{}
				the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
				common.Check(err)
				//log.Info(the_bson)
				err = bson.Unmarshal(the_bson, &tc)
				common.Check(err)
				sportTeams = append(sportTeams, tc)
			} else {
				tc := Team{}

				the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
				common.Check(err)
				//log.Info(string(the_bson[:]))
				err = bson.Unmarshal(the_bson, &tc)
				common.Check(err)
				sportTeams = append(sportTeams, tc)
			}

		}
	}

	return sportTeams

}
func conn() bolt.Conn {
	driver := bolt.NewDriver()
	bolt_url = os.Getenv("NEO4JBOLT")
	conn, err := driver.OpenNeo(bolt_url)
	common.Check(err)
	return conn
}
func getquery(nodeName string) string {
	cypher_query := fmt.Sprintf(`MATCH (t:%v {sports_id:{id}}) RETURN t`, nodeName)
	return cypher_query
}

func GetSport(p map[string]interface{}) []Sport {
	driver := bolt.NewDriver()
	bolt_url = os.Getenv("NEO4JBOLT")
	log.Info(bolt_url)

	conn, err := driver.OpenNeo(bolt_url)
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
