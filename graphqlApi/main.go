package graphqlApi

import (
	"github.com/graphql-go/graphql"
	"github.com/nemesisesq/ss_data_service/common"
	"fmt"

	"github.com/nemesisesq/ss_data_service/dao"
)


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

func Schema() *graphql.Schema {

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "World", nil
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
