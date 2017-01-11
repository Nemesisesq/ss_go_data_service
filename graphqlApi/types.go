package graphqlApi

import "github.com/graphql-go/graphql"

var orgType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Organizations",
		Fields: graphql.Fields{
			"img": &graphql.Field{
				Type: graphql.String,
			},
			"gracenote_organization_id": &graphql.Field{
				Type: graphql.String,
			},
			"organization": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var favoriteEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "Favorite",
	Description: "A boolean indicating whether or not the item is to be favorited or un-favorited",
	Values: graphql.EnumValueConfigMap{
		"true": &graphql.EnumValueConfig{
			Value:       true,
			Description: "This value indicates that the object for the Id submitted will be favorited",
		},
		"false": &graphql.EnumValueConfig{
			Value:       false,
			Description: "This inidcated that the object for the Id submitted will be un-favorited",
		},
	},
})

var teamLabelEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "TeamType",
	Description: "A Value indicated whether the team being queried or modified is a College Team or a professional team",
	Values: graphql.EnumValueConfigMap{
		"college": &graphql.EnumValueConfig{
			Value:       "CollegeTeam",
			Description: "This value indicates that the object for the Id submitted is a CollegeTeam",
		},
		"pro": &graphql.EnumValueConfig{
			Value:       "Team",
			Description: "This inidcated that the object for the Id submitted is a ProTeam",
		},
	},
})
var orgsType = graphql.NewList(orgType)

var sportType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Sport",
		Fields: graphql.Fields{
			"sportsId": &graphql.Field{
				Type: graphql.String,
			},
			"sportsName": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var favoriteStatusType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FavoriteStatus",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

var sportsType = graphql.NewList(sportType)

var teamType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Team",
		Fields: graphql.Fields{
			"propername": &graphql.Field{
				Type: graphql.String,
			},
			"img": &graphql.Field{
				Type: graphql.String,
			},
			"team_brand_id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"nickname": &graphql.Field{
				Type: graphql.String,
			},
			"abbreviation": &graphql.Field{
				Type: graphql.String,
			},
			"sports_id": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var teamsType = graphql.NewList(teamType)
