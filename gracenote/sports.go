package gracenote

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"github.com/nemesisesq/ss_data_service/common"
	"gopkg.in/square/go-jose.v1/json"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Team struct {
	TeamBrandId    string `json:"teamBrandId" bson:"team_brand_id"`
	SportsId       string `json:"sportsId" bson:"sports_id"`
	TeamBrandName  string `json:"teamBrandName" bson:"team_brand_name"`
	Nickname       string `json:"nickname" bson:"nickname"`
	ProperName     string `json:"properName" bson:"proper_name"`
	Abbreviation   string `json:"abbreviation" bson:"abbreviation"`
	PreferredImage PreferredImage `json:"preferredImage" bson:"preferred_image"`
	Img            string        `bson:"img"`
}

type CollegeTeam struct {
	Team
	University     University `json:"university" bson:"university"`
	PreferredImage PreferredImage `json:"preferredImage" bson:"preferred_image"`
	Img            string        `bson:"img"`
}

type CollegeTeams []CollegeTeam

type Teams []Team

type Organization struct {
	OrganizationId   interface{}      `json:"organizationId" bson:"gracenote_organization_id"`
	OrganizationName string      `json:"organizationName" bson:"organization"`
	PreferredImage   PreferredImage `json:"preferredImage" bson:"preferred_image"`
	Img              string        `bson:"img"`
	Teams            `json:"teams"`
}

type Sport struct {
	SportsId      string         `json:"sportsId" bson:"sportsId"`
	SportsName    string         `json:"sportsName" bson:"sportsName"`
	Organizations []Organization `json:"organizations" bson:"organizations"`
}

type SportsList []Sport

type PreferredImage struct {
	Width    string `json:"width" bson:"width"`
	Height   string `json:"height" bson:"height"`
	Uri      string `json:"uri" bson:"uri"`
	Category string `json:"category" bson:"category"`
	Primary  string `json:"primary" bson:"primary"`
	Tier     string `json:"tier" bson:"tier"`
}
type University struct {
	UniversityId   string         `json:"universityId" bson:"universityId"`
	UniversityName string         `json:"universityName" bson:"universityName"`
	NickName       string         `json:"nickName" bson:"nickName"`
	PreferredImage PreferredImage `json:"preferredImage" bson:"preferred_image"`
	Img            string        `bson:"img"`
}

type Universities []University

type PipeComponents struct {
	c []string
	p []map[string]interface{}
}

func (parts *PipeComponents) a(cypher string, params map[string]interface{}) {
	parts.c = append(parts.c, cypher)
	parts.p = append(parts.p, params)
}

func (parts PipeComponents) execute(conn bolt.Conn) {
	pipeline, err := conn.PreparePipeline(parts.c...)
	if err != nil {
		panic(err)
	}
	logrus.Info(pipeline)
	logrus.Info(parts.p)

	pipelineResults, err := pipeline.ExecPipeline(parts.p...)
	if err != nil {
		panic(err)
	}

	for _, result := range pipelineResults {
		numResult, _ := result.RowsAffected()
		fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 2 (per each iteration)
	}

	err = pipeline.Close()
	if err != nil {
		panic(err)
	}
}

func (ct CollegeTeams) SaveCollegeTeams() {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)
	p := &PipeComponents{}
	for _, val := range ct {
		cypher_query := `
			MATCH (s:Sport {gracenote_sport_id:{sports_id}})
			MATCH (u:University {gracenote_university_id: {uni_id}})
			MERGE (t:CollegeTeam {team_brand_id:{brand_id}, sports_id:{sports_id}, name:{team_name}, nickname:{nick}, propername:{propername}, abbreviation:{abbr}, img:{uri}})
			MERGE (t)-[:PLAYS]->(s)
			MERGE (t)-[:BELONGS_TO]->(u)
		`
		params := map[string]interface{}{
			"brand_id":fmt.Sprint(val.TeamBrandId),
			"sports_id":fmt.Sprint(val.SportsId),
			"team_name":val.TeamBrandName,
			"nick": val.Nickname,
			"propername": val.ProperName,
			"abbr": val.Abbreviation,
			"uni_id": fmt.Sprintf(val.University.UniversityId),
			"uri": fmt.Sprint(val.PreferredImage.Uri),
		}
		p.a(cypher_query, params)

		//ProcessCypher(conn, cypher_query, params)

		//cypher_query = `
		//	MATCH (t:Team {team_brand_id:{brand_id}})
		//
		//	MATCH (o:Org {gracenote_organization_id:{org_id}})
		//
		//	MERGE (t)-[:MEMBER_OF]->(o)
		//`
		//params = map[string]interface{}{
		//	"brand_id":fmt.Sprint(val.TeamBrandId),
		//	"sports_id":fmt.Sprint(val.SportsId),
		//}
		//ProcessCypher(conn, cypher_query, params)
	}
	p.execute(conn)

}

func (teams Teams) SaveTeams(org_id float64) {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)
	p := &PipeComponents{}
	//var a = append
	for _, val := range teams {
		cypher_query := `
			MERGE (t:Team {team_brand_id:{brand_id}, sports_id:{sports_id}, name:{team_name}, nickname:{nick}, propername:{propername}, abbreviation:{abbr}, img:{img}})
		`
		params := map[string]interface{}{
			"brand_id":fmt.Sprint(val.TeamBrandId),
			"sports_id":fmt.Sprint(val.SportsId),
			"team_name":val.TeamBrandName,
			"nick": val.Nickname,
			"propername": val.ProperName,
			"abbr": val.Abbreviation,
			"org_id": org_id,
			"img":val.PreferredImage.Uri,
		}
		//ProcessCypher(conn, cypher_query, params)
		p.a(cypher_query, params)
		if org_id != 1.001 {

			cypher_query = `
			MATCH (t:Team {team_brand_id:{brand_id}})
			MATCH (s:Sport {gracenote_sport_id:{sports_id}})
			MATCH (o:Org {gracenote_organization_id:{org_id}})
			MERGE (t)-[:PLAYS]->(s)
			MERGE (t)-[:MEMBER_OF]->(o)
		`
			params = map[string]interface{}{
				"brand_id":fmt.Sprint(val.TeamBrandId),
				"sports_id":fmt.Sprint(val.SportsId),
				"org_id": fmt.Sprint(org_id),
			}
			//ProcessCypher(conn, cypher_query, params)
			p.a(cypher_query, params)
		}
	}
	p.execute(conn)

}

func (unis Universities) SaveUniversities() {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	common.Check(err)

	p := &PipeComponents{}
	for _, val := range unis {
		cypher_query := `
			MERGE (u:University {gracenote_university_id:{uni_id}, name:{uni_name}, nickname:{uni_nick}, img:{uri}})
		`
		params := map[string]interface{}{"uni_id": val.UniversityId, "uni_name": val.UniversityName, "uni_nick": val.NickName, "uri": fmt.Sprintf("%v", val.PreferredImage.Uri)}
		//ProcessCypher(conn, cypher_query, params)
		//logrus.Info("before", p.c)
		//logrus.Info(p.p)
		p.a(cypher_query, params)
		//logrus.Info(p.c)
		//logrus.Info(p.p)
	}

	p.execute(conn)

}

func (sl SportsList) SaveSportsList() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	driver := bolt.NewDriver()
	logrus.Info(os.Getenv("NEO4JBOLT"))
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	defer conn.Close()
	//conn.ExecPipeline()
	common.Check(err)

	p := &PipeComponents{}
	for _, val := range sl {
		cypher_query := `MERGE (s:Sport {sport_name:{sport_name}, gracenote_sport_id:{sports_id}})`
		params := map[string]interface{}{"sport_name": val.SportsName, "sports_id": val.SportsId}
		//ProcessCypher(conn, cypher_query, params)
		p.a(cypher_query, params)
		for _, org := range val.Organizations {
			cypher_query := `MERGE (o:Org {organization:{org_name}, gracenote_organization_id:{org_id}, img:{uri}})`
			params := map[string]interface{}{"org_name": org.OrganizationName, "org_id": org.OrganizationId, "uri": fmt.Sprintf("%v", org.PreferredImage.Uri)}
			//ProcessCypher(conn, cypher_query, params)
			p.a(cypher_query, params)
			cypher_query = `MATCH (s:Sport {gracenote_sport_id: {s}})
					MATCH (o:Org {gracenote_organization_id: {o}})
					MERGE (o)-[:REPRESENTS]->(s)
					`
			params = map[string]interface{}{"s": val.SportsId, "o": org.OrganizationId}
			//ProcessCypher(conn, cypher_query, params)

			p.a(cypher_query, params)
		}
	}
	p.execute(conn)
}

func ProcessCypher(conn bolt.Conn, cypher_template string, params map[string]interface{}) {
	stmt, err := conn.PrepareNeo(cypher_template)
	common.Check(err)
	result, err := stmt.ExecNeo(params)
	common.Check(err)
	logrus.WithFields(logrus.Fields{
		"Cypher QueryResult": result,
		"params": params,
	}).Info()
	if val, ok := result.RowsAffected(); ok != nil {
		//numResult, err := result.RowsAffected()
		common.Check(err)
		logrus.Info("CREATED ROWS", val)

	}
	//Closing the statment will also close the rows
	stmt.Close()
}

func QueryNeo(cypher_query string, params map[string]interface{}) ([][]interface{}, map[string]interface{}) {
	driver := bolt.NewDriver()
	logrus.Info(os.Getenv("NEO4JBOLT"))
	conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	common.Check(err)
	defer conn.Close()

	data, rowMetaData, _, _ := conn.QueryNeoAll(cypher_query, params)

	return data, rowMetaData
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
		"api_key":   ApiKey,
		"imageSize": "md",
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

func GetTeamsInAnOrganization() {
	//driver := bolt.NewDriver()
	//logrus.Info(os.Getenv("NEO4JBOLT"))
	//conn, err := driver.OpenNeo(os.Getenv("NEO4JBOLT"))
	//common.Check(err)
	//defer conn.Close()
	//
	//data, rowMetaData, _, _ := conn.QueryNeoAll(`MATCH (n:Org) RETURN n`, nil)

	data, rowMetaData := QueryNeo(`MATCH (n:Org) RETURN n`, nil)

	orgs := []Organization{}

	logrus.Info(rowMetaData)
	for _, val := range data {
		org := &Organization{}
		the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
		err = bson.Unmarshal(the_bson, &org)
		common.Check(err)
		orgs = append(orgs, *org)
	}

	for _, val := range orgs {
		sClient := http.Client{}
		url := fmt.Sprintf("%v/%v/%v", SportsUri, "organizations", val.OrganizationId)
		req, err := http.NewRequest("GET", url, nil)
		common.Check(err)
		params := map[string]string{
			"api_key":   ApiKey,
			"imageSize": "Md",
		}
		common.BuildQuery(req, params)
		res, err := sClient.Do(req)

		common.Check(err)
		if res.StatusCode == 200 {

			decoder := json.NewDecoder(res.Body)
			logrus.Info(res.Body)

			orgs := []Organization{}

			if err := decoder.Decode(&orgs); err == nil {
				//logrus.Info(orgs)
				common.Check(err)
				org := orgs[0]
				org.Teams.SaveTeams(org.OrganizationId.(float64))
			}
		}
	}
}
func GetTeamsAtAUniversity(uni_id string) {
	sClient := &http.Client{}
	url := fmt.Sprintf("%v/%v/%v", SportsUri, "universities", uni_id)
	req, err := http.NewRequest("GET", url, nil)
	common.Check(err)
	params := map[string]string{
		"includeTeam": "true",
		"api_key":   ApiKey,
	}
	common.BuildQuery(req, params)
	res, err := sClient.Do(req)
	common.Check(err)
	decoder := json.NewDecoder(res.Body)
	college_teams := &CollegeTeams{}
	err = decoder.Decode(&college_teams)
	common.Check(err)
	college_teams.SaveCollegeTeams()
}

func GetTeamDetails(teamBrandId string) {

	sClient := &http.Client{}
	url := fmt.Sprintf("%v/%v/%v", SportsUri, "teams", teamBrandId)
	req, err := http.NewRequest("GET", url, nil)
	common.Check(err)
	params := map[string]string{
		"includeTeam": "true",
		"api_key":   ApiKey,
	}
	common.BuildQuery(req, params)
	res, err := sClient.Do(req)
	common.Check(err)
	decoder := json.NewDecoder(res.Body)

	teams := &Teams{}
	err = decoder.Decode(&teams)
	common.Check(err)
	teams.SaveTeams(1.001)

}

func GetAllTeamsDetails() {

	data, rowMetaData := QueryNeo(`MATCH (n:Team) RETURN n`, nil)

	logrus.Info(rowMetaData)
	//for indx := range data {
	for i := 0; i * 20 < len(data); i += 1 {
		//start := indx * 20
		//end := start + 20
		//chunk := data[start: end]
		var chunk [][]interface{}
		chunk, data = data[:12], data[12:]
		//raw := []interface{}{}
		teams := &Teams{}
		for _, val := range chunk {
			t := &Team{}
			the_bson, err := bson.Marshal(val[0].(graph.Node).Properties)
			err = bson.Unmarshal(the_bson, &t)
			common.Check(err)
			*teams = append(*teams, *t)
		}


		ids := []string{}

		for _, val := range *teams {
			ids = append(ids, val.TeamBrandId)
		}

		GetTeamDetails(strings.Join(ids, ","))
	}

}