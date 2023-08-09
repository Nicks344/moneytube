package mutations

import (
	"github.com/Nicks344/moneytube/client/core/src/modules/bugreport"

	"github.com/graphql-go/graphql"
)

func sendBugReport() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"type": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"description": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"dataJSON": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.String,
		Description: "Send bug report",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return bugreport.Send(params.Args["type"].(string), params.Args["description"].(string), params.Args["dataJSON"].(string))
		},
	}
}
