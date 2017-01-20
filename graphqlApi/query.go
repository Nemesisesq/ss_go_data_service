package graphqlApi

import (
	"github.com/graphql-go/graphql"
	"github.com/nemesisesq/ss_data_service/dao"
)

var queryFields = graphql.Fields{

	"hello": &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "World", nil
		},
	},
	"orgs": &graphql.Field{
		Type: orgsType,
		Args: graphql.FieldConfigArgument{
			"sportId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			sportId := p.Args["sportId"].(string)
			orgs := dao.GetOrgs(sportId)
			return orgs, nil
		},
	},


	"teams": &graphql.Field{
		Type: teamsType,
		Args: graphql.FieldConfigArgument{
			"sportId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"orgId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {

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
		Type: sportsType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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

	"favorites": &graphql.Field {
		Type: favsType,
		Args: graphql.FieldConfigArgument{
			"user_id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//params := map[string]interface{}{}
			user_id := p.Args["user_id"].(string)
			favs := dao.GetFavorties(user_id)

			return favs, nil
		},

	},
}
