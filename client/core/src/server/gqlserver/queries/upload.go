package queries

import (
	"sort"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/upload"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
)

func uploadTasks() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(gqlstructs.TableUploadTaskOutput),
		Description: "Get list of upload tasks",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return model.GetUploadTasks()
		},
	}
}

func getUploadTaskData() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        gqlstructs.UploadTaskDataOutput,
		Description: "Get upload tasks data",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return model.GetUploadData(p.Args["id"].(int))
		},
	}
}

func uploadDataLists() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "UploadDataList",
				Fields: graphql.Fields{
					"langs": &graphql.Field{
						Type: graphql.NewList(graphql.String),
					},
					"categories": &graphql.Field{
						Type: graphql.NewList(graphql.String),
					},
				},
			},
		),
		Description: "Get upload data lists",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			langs := []string{}
			for lang := range upload.LangCodes {
				if lang == "Нет" {
					continue
				}
				langs = append(langs, lang)
			}
			sort.Strings(langs)
			langs = append([]string{"Нет"}, langs...)

			categories := []string{}
			for category := range upload.CategoriesCodes {
				categories = append(categories, category)
			}

			return map[string]interface{}{
				"langs": langs, "categories": categories,
			}, nil
		},
	}
}

func uploadDataTemplates() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(gqlstructs.UploadDataTemplatesOutput),
		Description: "Get list of upload data templates",
		Resolve: func(p graphql.ResolveParams) (res interface{}, err error) {
			res, err = model.GetUploadDataTemplates()
			return
		},
	}
}

func getUploadDataTemplate() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"label": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        gqlstructs.UploadDataTemplatesOutput,
		Description: "Get upload data template",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return model.GetUploadDataTemplate(p.Args["label"].(string))
		},
	}
}
