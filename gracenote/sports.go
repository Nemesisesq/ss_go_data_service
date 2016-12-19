package gracenote

import (
	"net/http"
	"fmt"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/square/go-jose.v1/json"
	bolt"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
	"github.com/Sirupsen/logrus"
	"strings"
)

type Organization struct {
	OrganizationId   string `json:"organizationId" bson:"organization_id"`
	OrganizationName string `json:"organizationName" bson:"organizationName"`
	PreferredImage   map[string]string `json:"preferredImage" bson:"preferredImage"`
}

type Sport struct {
	SportsId      string `json:"sportsId" bson:"sportsId"`
	SportsName    string `json:"sportsName" bson:"sportsName"`
	Organizations []Organization `json:"organizations" bson:"organizations"`
}

type SportsList []Sport

func (sl SportsList) SaveSportsList() {

	logrus.SetFormatter(&logrus.JSONFormatter{})

	driver := bolt.NewDriver()
	logrus.Info(os.Getenv("NEO4JBOLT"))

	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))

	common.Check(err)

	for _, val := range sl {
		logrus.Info(val.SportsName, " sports name")
		rmvSpace := strings.NewReplacer(" ", "", ".", "_", "'", "", "-", "", "/", "")
		sport_node_name := rmvSpace.Replace(val.SportsName)
		cypher_query := fmt.Sprintf(`CREATE (%v:Sport {sport_name:{sport_name}, gracenote_sport_id:{sports_id}})`, sport_node_name)

		logrus.Info(sport_node_name, val.SportsName, val.SportsId)
		params := map[string]interface{}{"sport_node_name": sport_node_name, "sport_name": val.SportsName, "sports_id": val.SportsId}

		ProcessCypher(conn, cypher_query, params)

		for _, org := range val.Organizations {

			org_node_name := rmvSpace.Replace(org.OrganizationName)

			cypher_query := fmt.Sprintf(`CREATE (%v:Organization {organization:{org_name}, gracenote_organization_id:{org_id}})`, org_node_name)
			params := map[string]interface{}{"sport_node_name": sport_node_name, "org_node_name": org_node_name, "org_name": org.OrganizationName, "org_id": org.OrganizationId}
			ProcessCypher(conn, cypher_query, params)

			cypher_query = `CREATE ({sport_node_name})-[:BELONGS_TO]->({org_node_name})`
			params = map[string]interface{}{"sport_node_name": sport_node_name, "org_node_name": org_node_name}
			ProcessCypher(conn, cypher_query, params)
		}

	}
}


func ProcessCypher(conn bolt.Conn, cypher_template string, params map[string]interface{}) {
	stmt, err := conn.PrepareNeo(cypher_template)

	common.Check(err)

	result, err := stmt.ExecNeo(params)
	common.Check(err)

	numResult, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("CREATED ROWS: %d\n", numResult)

	// Closing the statment will also close the rows
	stmt.Close()
}

func GetSport(sportsId string) {
	sClient := &http.Client{}
	url := fmt.Sprintf("%v/%v", SportsUri, sportsId)
	req, err := http.NewRequest("GET", url, nil)

	common.Check(err)

	params := map[string]string{
		"api_key": ApiKey,
		"includeOrg" : "true",
		"imageSize" : "md",
	}

	common.BuildQuery(req, params)

	res, err := sClient.Do(req)
	common.Check(err)

	decoder := json.NewDecoder(res.Body)

	sportsList := SportsList{}

	err = decoder.Decode(&sportsList)

	common.Check(err)

	sportsList.SaveSportsList()

}


