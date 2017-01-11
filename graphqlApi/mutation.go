package graphqlApi

import (
	"github.com/graphql-go/graphql"
	"github.com/nemesisesq/ss_data_service/streamsavvy"
)

var mutationFields = graphql.Fields{
	"toggleTeam": &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			"label": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(teamLabelEnum)),
			},
			"team_brand_id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"favorite": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.NewList(favoriteEnum)),
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
			label := params.Args["label"].(string)

			streamsavvy.ProcessTeamForFavorites(userId, email, team_brand_id, label, favorite)

			return "ok", nil
		},
	},
	"toggleShow": &graphql.Field{
		Type: favoriteStatusType,
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

			status := streamsavvy.ProcessShowForFavorites(userId, email, guidebox_id, favorite)

			return map[string]interface{}{"status":status}, nil
		},
	},
	"toggleSport": &graphql.Field{
		Type: favoriteStatusType,
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

			status := streamsavvy.ProcessSportForFavorites(userId, email, sportsId, favorite)

			return map[string]interface{}{"status":status}, nil
		},
	},
}
