package subscriptions

import (
	"github.com/graphql-go/graphql"
)

func GetSubscriptions() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Subscriptions",
		Fields: graphql.Fields{
			"onAccountUpdated":    onAccountUpdated(),
			"onUploadTaskUpdated": onUploadTaskUpdated(),
			"onToolResult":        onToolResult(),
		},
	})
}
