package gqlstructs

import "github.com/graphql-go/graphql"

var AccountInput = graphql.FieldConfigArgument{
	"account": &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name: "AccountFields",
				Fields: graphql.InputObjectConfigFieldMap{
					"ID": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"Login": &graphql.InputObjectFieldConfig{
						Type: graphql.String,
					},
					"Password": &graphql.InputObjectFieldConfig{
						Type: graphql.String,
					},
					"Proxy": &graphql.InputObjectFieldConfig{
						Type: graphql.String,
					},
					"Group": &graphql.InputObjectFieldConfig{
						Type: graphql.String,
					},
				},
			}),
	},
	"cookieFile": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var AccountOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Account",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: Uint,
			},
			"Login": &graphql.Field{
				Type: graphql.String,
			},
			"Password": &graphql.Field{
				Type: graphql.String,
			},
			"Proxy": &graphql.Field{
				Type: graphql.String,
			},
			"Group": &graphql.Field{
				Type: graphql.String,
			},
			"ChannelName": &graphql.Field{
				Type: graphql.String,
			},
			"VideoCount": &graphql.Field{
				Type: graphql.Int,
			},
			"PlaylistsCount": &graphql.Field{
				Type: graphql.Int,
			},
			"SubscribersCount": &graphql.Field{
				Type: graphql.Int,
			},
			"ViewsCount": &graphql.Field{
				Type: graphql.Int,
			},
			"Status": &graphql.Field{
				Type: graphql.Int,
			},
			"ErrorMessage": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
