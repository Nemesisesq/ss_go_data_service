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
type Team gracenote.Team
type CollegeTeam gracenote.CollegeTeam

func(t Team)isATeam() string {
	return "pro"
}
func (c CollegeTeam)isATeam() string {
	return "amature"
}

type SportsTeam interface {
	isATeam() string
}

var bolt_url = os.Getenv("NEO4JBOLT")

func GetTeams(sportId string) []SportsTeam {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(bolt_url)
	defer conn.Close()
	common.Check(err)

	cypher_query := `MATCH (t:Team {sports_id:{id}})
			MATCH (c:CollegeTeam {sports_id:{id})
			RETURN t,c
	`
	params := map[string]interface{}{"id":sportId}

	data, metadata := QueryNeo(cypher_query, params)
	log.Info(metadata)
	sportTeams := []SportsTeam{}

	for _, val := range data {
		var tc SportsTeam

		if true {
			tc = CollegeTeam{}
		} else {
			tc = Team{}
		}


		the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
		common.Check(err)
		err = bson.Unmarshal(the_bson, &tc)
		common.Check(err)
		sportTeams = append(sportTeams, tc)

	}

	return sportTeams

}

func GetSport(p map[string]interface{}) []Sport {
	driver := bolt.NewDriver()
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
