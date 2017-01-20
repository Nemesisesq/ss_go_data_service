package graphqlApi

import (
	"github.com/graphql-go/graphql"
	"github.com/nemesisesq/ss_data_service/streamsavvy"
	"github.com/nemesisesq/ss_data_service/dao"
)

var mutationFields = graphql.Fields{
	"toggleTeam": &graphql.Field{
		Type: favsType,
		Args: graphql.FieldConfigArgument{
			"label": &graphql.ArgumentConfig{
				Type: graphql.NewList(teamLabelEnum),
			},
			"team_brand_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"favorite": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"userId": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},

			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			userId := params.Args["userId"].(string)
			email := params.Args["email"].(string)
			team_brand_id := params.Args["team_brand_id"].(string)
			favorite := params.Args["favorite"].(bool)
			var label string
			if label, ok := params.Args["label"]; ok {
				label = label.(string)
			}

			_ = streamsavvy.ProcessTeamForFavorites(userId, email, team_brand_id, label, favorite)

			/* TODO Cleand this up this was a boolean returned from the function.*/
			favs := dao.GetFavorties(userId)

			return favs, nil
		},
	},
	"toggleShow": &graphql.Field{
		Type: favsType,
		Args: graphql.FieldConfigArgument{
			"favorite": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"guidebox_id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"userId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},

			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			userId := params.Args["userId"].(string)
			email := params.Args["email"].(string)
			guidebox_id := params.Args["guidebox_id"].(int)
			favorite := params.Args["favorite"].(bool)


			_ = streamsavvy.ProcessShowForFavorites(userId, email, guidebox_id, favorite)
			/* TODO Cleand this up this was a boolean returned from the function.*/
			favs := dao.GetFavorties(userId)

			return favs, nil
		},
	},
	"toggleSport": &graphql.Field{
		Type: favsType,
		Args: graphql.FieldConfigArgument{
			"sportsId": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},

			"userId": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},

			"email": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"favorite": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			userId := params.Args["userId"].(string)
			email := params.Args["email"].(string)
			sportsId := params.Args["sportsId"].(int)
			favorite := params.Args["favorite"].(bool)

			_ = streamsavvy.ProcessSportForFavorites(userId, email, sportsId, favorite)


			/* TODO Cleand this up this was a boolean returned from the function.*/
			favs := dao.GetFavorties(userId)

			return favs, nil
		},
	},
}
