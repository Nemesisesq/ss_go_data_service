package graphqlApi

import (
	"github.com/graphql-go/graphql"
	"github.com/nemesisesq/ss_data_service/common"
	"fmt"

	"github.com/nemesisesq/ss_data_service/dao"
)
var orgType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:"Organizations",
		Fields: graphql.Fields{
			"img": &graphql.Field{
				Type: graphql.String,
			},
			"gracenote_organization_id": &graphql.Field{
				Type:graphql.String,
			},
			"organization": &graphql.Field{
				Type:graphql.String,
			},
		},
	},
)

var orgsType = graphql.NewList(orgType)

var sportType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:"Sport",
		Fields: graphql.Fields{
			"sportsId": &graphql.Field{
				Type: graphql.String,
			},
			"sportsName": &graphql.Field{
				Type:graphql.String,
			},
		},
	},
)

var sportsType = graphql.NewList(sportType)

var teamType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:"Team",
		Fields: graphql.Fields{
			"propername":&graphql.Field{
				Type:graphql.String,
			},
			"img": &graphql.Field{
				Type:graphql.String,
			},
			"team_brand_id": &graphql.Field{
			Type:graphql.String,
			},
			"name": &graphql.Field{
			Type:graphql.String,
			},
			"nickname": &graphql.Field{
			Type:graphql.String,
			},
			"abbreviation": &graphql.Field{
			Type:graphql.String,
			},
			"sports_id": &graphql.Field{
			Type:graphql.String,
			},
		},
	},
)

var teamsType = graphql.NewList(teamType)
func Schema() *graphql.Schema {

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "World", nil
			},
		},
		"orgs" : &graphql.Field{
			Type:orgsType,
			Args:graphql.FieldConfigArgument{
				"sportId": &graphql.ArgumentConfig{
					Type:graphql.String,
				},
			},
			Resolve:func(p graphql.ResolveParams) (interface{}, error) {
				sportId := p.Args["sportId"].(string)
				orgs := dao.GetOrgs(sportId)
				return orgs, nil
			},
		},

		"teams" :&graphql.Field{
			Type: teamsType,
			Args: graphql.FieldConfigArgument{
				"sportId": &graphql.ArgumentConfig{
					Type:graphql.String,
				},
				"orgId": &graphql.ArgumentConfig{
					Type:graphql.String,
				},
			},
			Resolve:func(p graphql.ResolveParams) (interface{}, error) {

				var teams []map[string]interface{}

				if sportId, ok := p.Args["sportId"].(string); ok {

		       			teams = dao.GetTeams(sportId)
				}

				if orgId, ok := p.Args["orgId"].(string); ok {
					teams = dao.GetTeamsByOrganization(orgId)
				}

				return teams, nil
			},
		},
		"sports": &graphql.Field{
			Type:sportsType,
			Resolve:func(p graphql.ResolveParams) (interface{}, error) {
				params := map[string]interface{}{}
				sports := dao.GetSport(params)

				return sports, nil
			},
		},

		"sport": &graphql.Field{
			Type: sportType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve:func(p graphql.ResolveParams) (interface{}, error) {
				params := map[string]interface{}{}

				var sports []dao.Sport
				if id, ok := p.Args["id"]; ok {
					params["id"] = id
					sports = dao.GetSport(params)
				} else {
					sports = dao.GetSport(params)
				}

				return sports[0], nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	common.Check(err)

	return &schema

}

func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
