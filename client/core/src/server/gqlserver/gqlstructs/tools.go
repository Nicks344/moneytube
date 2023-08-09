package gqlstructs

import "github.com/graphql-go/graphql"

var ToolsProgressResult = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ToolsProgressResult",
		Fields: graphql.Fields{
			"Status": &graphql.Field{
				Type: StatusEnum,
			},
			"MaxProgress": &graphql.Field{
				Type: graphql.Int,
			},
			"Progress": &graphql.Field{
				Type: graphql.Int,
			},
			"Error": &graphql.Field{
				Type: graphql.String,
			},
			"JsonData": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
