package gracenote

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/square/go-jose.v1/json"
)

type Organization struct {
	OrganizationId   string            `json:"organizationId" bson:"organization_id"`
	OrganizationName string            `json:"organizationName" bson:"organizationName"`
	PreferredImage   map[string]string `json:"preferredImage" bson:"preferredImage"`
}

type Sport struct {
	SportsId      string         `json:"sportsId" bson:"sportsId"`
	SportsName    string         `json:"sportsName" bson:"sportsName"`
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
		//rmvSpace := strings.NewReplacer(" ", "", ".", "_", "'", "", "-", "", "/", "")
		//sport_node_name := rmvSpace.Replace(val.SportsName)
		cypher_query := `CREATE (n:Sport {sport_name:{sport_name}, gracenote_sport_id:{sports_id}})`

		params := map[string]interface{}{"sport_name": val.SportsName, "sports_id": val.SportsId}

		ProcessCypher(conn, cypher_query, params)

		for _, org := range val.Organizations {

			//org_node_name := rmvSpace.Replace(org.OrganizationName)

			cypher_query := `CREATE (n:Org {organization:{org_name}, gracenote_organization_id:{org_id}})`
			params := map[string]interface{}{"org_name": org.OrganizationName, "org_id": org.OrganizationId}
			ProcessCypher(conn, cypher_query, params)

			_, _ = conn.ExecNeo(`CREATE INDEX ON:Sport(sports_id) CREATE INDEX ON:Org(org_id)`, nil)
			cypher_query = `MATCH (s:Sport {gracenote_sport_id: {s}})
					MATCH (o:Org {gracenote_organization_id: {o}})
					CREATE (o)-[:B]->(s)
					`
			result, err := conn.ExecNeo(cypher_query, map[string]interface{}{"s":val.SportsId, "o":org.OrganizationId})


			logrus.Info(result)
			common.Check(err)
		}

	}
}

func ProcessCypher(conn bolt.Conn, cypher_template string, params map[string]interface{}) {
	stmt, err := conn.PrepareNeo(cypher_template)

	common.Check(err)

	result, err := stmt.ExecNeo(params)
	common.Check(err)
	logrus.Info(result)

	//numResult, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("CREATED ROWS: %d\n", numResult)

	// Closing the statment will also close the rows
	stmt.Close()
}

func GetSport(sportsId string) {
	sClient := &http.Client{}
	url := fmt.Sprintf("%v/%v", SportsUri, sportsId)
	req, err := http.NewRequest("GET", url, nil)

	common.Check(err)

	params := map[string]string{
		"api_key":    ApiKey,
		"includeOrg": "true",
		"imageSize":  "md",
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
