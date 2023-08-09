package queries

import (
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/accounts"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
)

func getAccounts() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(gqlstructs.AccountOutput),
		Description: "Get list of accounts",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return model.GetAccounts()
		},
	}
}

func exportCookies() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"file": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Export account cookies",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			acc, err := model.GetAccount(p.Args["id"].(int))
			if err != nil {
				return nil, err
			}

			file := p.Args["file"].(string)

			err = accounts.ExportCookies(acc, file)
			return err == nil, err
		},
	}
}
