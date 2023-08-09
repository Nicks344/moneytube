package queries

import (
	"github.com/graphql-go/graphql"
)

func GetQueries() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Queries",
		Fields: graphql.Fields{
			"accounts":      getAccounts(),
			"exportCookies": exportCookies(),
			//"uploadTask":  uploadTask(),
			"uploadTasks":           uploadTasks(),
			"settings":              settings(),
			"macroses":              macroses(),
			"macros":                macros(),
			"getUploadTaskData":     getUploadTaskData(),
			"uploadDataLists":       uploadDataLists(),
			"uploadDataTemplates":   uploadDataTemplates(),
			"getUploadDataTemplate": getUploadDataTemplate(),
		},
	})
}
