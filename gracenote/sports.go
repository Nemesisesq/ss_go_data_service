package gracenote

import (
	"net/http"
	"fmt"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/square/go-jose.v1/json"
	bolt"github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"os"
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
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))

	common.Check(err)

	for _, val := range sl {
		cypher_query := `CREATE	({sport_name}:Sport {sport_name:{sport_name}, gracenote_sport_id: {sports_id}})`

		params := map[string]interface{}{"sport_name": val.SportsName, "sports_id": val.SportsId}

		ProcessCypher(conn, cypher_query, params)

		for _, val := range val.Organizations {
			cypher_query := `CREATE	({org_name}:Organization {organization: {org_name}, gracenote_organization_id: {org_id} })
			CREATE ({sport_name})-[:BELONGS_TO]->({org_name)
			`
			params := map[string]interface{}{"org_name": val.OrganizationName, "org_id": val.OrganizationId}
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

	decoder := json.NewDecoder(res)

	sportsList := SportsList{}

	err = decoder.Decode(&sportsList)

	common.Check(err)

}

