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

type PrefferedImage struct {
	Width    string `json:"width" bson:"width"`
	Height   string `json:"height" bson:"height"`
	Uri      string `json:"uri" bson:"uri"`
	Category string `json:"category" bson:"category"`
	Primary  string `json:"primary" bson:"primary"`
	Tier     string `json:"tier" bson:"tier"`
}
type University  struct {
	UniversityId   string        `json:"universityId" bson:"universityId"`
	UniversityName string        `json:"universityName" bson:"universityName"`
	NickName       string         `json:"nickName" bson:"nickName"`
	PrefferedImage PrefferedImage `json:"prefferedImage" bson:"prefferedImage"`
}

type Universities []University

func (unis Universities) SaveUniversities() {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)
	for _, val := range unis {
		cypher_query := `
			MERGE (u:University {gracenote_university_id:{uni_id}, name:{uni_name}, nickname:{uni_nick}, img:{uri}})
		`
		params := map[string]interface{}{"uni_id": val.UniversityId, "uni_name": val.UniversityName, "uni_nick":val.NickName, "img":val.PrefferedImage.Uri}
		ProcessCypher(conn, cypher_query, params)
	}

}

func (sl SportsList) SaveSportsList() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	driver := bolt.NewDriver()
	logrus.Info(os.Getenv("NEO4JBOLT"))
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)
	for _, val := range sl {
		cypher_query := `MERGE (s:Sport {sport_name:{sport_name}, gracenote_sport_id:{sports_id}})`
		params := map[string]interface{}{"sport_name": val.SportsName, "sports_id": val.SportsId, }
		ProcessCypher(conn, cypher_query, params)
		for _, org := range val.Organizations {
			cypher_query := `MERGE (o:Org {organization:{org_name}, gracenote_organization_id:{org_id}, preffered_image:{img}})`
			params := map[string]interface{}{"org_name": org.OrganizationName, "org_id": org.OrganizationId, "img":fmt.Sprintf("%v", org.PreferredImage)}
			ProcessCypher(conn, cypher_query, params)
			cypher_query = `MATCH (s:Sport {gracenote_sport_id: {s}})
					MATCH (o:Org {gracenote_organization_id: {o}})
					MERGE (o)-[:REPRESENTS]->(s)
					`
			params = map[string]interface{}{"s":val.SportsId, "o":org.OrganizationId}
			ProcessCypher(conn, cypher_query, params)
		}
	}
}

func ProcessCypher(conn bolt.Conn, cypher_template string, params map[string]interface{}) {
	stmt, err := conn.PrepareNeo(cypher_template)
	common.Check(err)
	result, err := stmt.ExecNeo(params)
	common.Check(err)
	logrus.Info(result)
	if val, ok := result.RowsAffected(); ok {
		//numResult, err := result.RowsAffected()
		common.Check(err)
		fmt.Printf("CREATED ROWS: %d\n", val)
	}
	//Closing the statment will also close the rows
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

func GetUniversities() {
	sClient := &http.Client{}
	url := fmt.Sprintf("%v/%v", SportsUri, "universities")
	req, err := http.NewRequest("GET", url, nil)
	common.Check(err)
	params := map[string]string{
		"api_key":    ApiKey,
		"imageSize":  "md",
	}
	common.BuildQuery(req, params)
	res, err := sClient.Do(req)
	common.Check(err)
	decoder := json.NewDecoder(res.Body)
	unis := &Universities{}
	err = decoder.Decode(unis)
	common.Check(err)
	unis.SaveUniversities()
}
