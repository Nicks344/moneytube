package subscriptions

import (
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
)

func onAccountUpdated() *graphql.Field {
	return &graphql.Field{
		Type:        gqlstructs.AccountOutput,
		Description: "On account updated",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return params.Context.Value(params.Args["id"].(string)), nil
		},
	}
}
