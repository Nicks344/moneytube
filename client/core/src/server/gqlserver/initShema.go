package gqlserver

import (
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/mutations"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/queries"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/subscriptions"
	"github.com/Nicks344/moneytube/client/core/src/server/serverutils"

	"github.com/graphql-go/graphql"
)

func Init() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        queries.GetQueries(),
		Mutation:     mutations.GetMutations(),
		Subscription: subscriptions.GetSubscriptions(),
		Types: []graphql.Type{
			gqlstructs.StatusEnum,
		},
	})
	if err != nil {
		panic("failed to create new schema, error: " + err.Error())
	}
	serverutils.SetSchema(schema)
}
