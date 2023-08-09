package mutations

import (
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/upload"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

func addUploadTasks() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"data": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadTaskDataInput,
			},
		},
		Type:        graphql.NewList(gqlstructs.TableUploadTaskOutput),
		Description: "Add new upload tasks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var data moneytubemodel.UploadData
			mapstructure.Decode(params.Args["data"], &data)
			tasks, err := model.SaveUploadData(&data)
			if err != nil {
				return nil, err
			}
			return tasks, nil
		},
	}
}

func deleteUploadTask() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete upload task",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			err := model.DeleteUploadTask(params.Args["id"].(int))
			return err == nil, err
		},
	}
}

func deleteAllUploadTasks() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Delete all upload tasks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			err := model.DeleteAllUploadTasks()
			return err == nil, err
		},
	}
}

func startUploadTask() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start upload task",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := upload.StartUploadTaskByID(params.Args["id"].(int), 0)
			return err == nil, err
		},
	}
}

func stopUploadTask() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Stop upload task",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			err := upload.StopUploadTask(params.Args["id"].(int))
			return err == nil, err
		},
	}
}

func startAllUploadTasks() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Start all upload tasks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			tasks, err := model.GetUploadTasksByStatuses([]int{moneytubemodel.UTSStopped, moneytubemodel.UTSError})
			if err != nil {
				return nil, err
			}
			for _, t := range tasks {
				_, err = upload.StartUploadTaskByID(int(t.ID), 0)
				if err != nil {
					return nil, err
				}
			}
			return true, nil
		},
	}
}

func stopAllUploadTasks() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Stop all upload tasks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			tasks, err := model.GetUploadTasksByStatuses([]int{moneytubemodel.UTSInProcess})
			if err != nil {
				return nil, err
			}
			for _, t := range tasks {
				err = upload.StopUploadTask(int(t.ID))
				if err != nil {
					return nil, err
				}
			}
			return true, nil
		},
	}
}

func addUploadDataTemplate() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"data": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadTaskDataInput,
			},
		},
		Type:        graphql.NewList(gqlstructs.TableUploadTaskOutput),
		Description: "Add new upload tasks",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var data moneytubemodel.UploadData
			mapstructure.Decode(params.Args["data"], &data)
			tasks, err := model.SaveUploadData(&data)
			if err != nil {
				return nil, err
			}
			return tasks, nil
		},
	}
}

func saveTemplate() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"data": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadDataTemplateInput,
			},
			"label": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Save task params as template",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var data moneytubemodel.UploadDataTemplate
			mapstructure.Decode(params.Args["data"], &data)
			data.Label = params.Args["label"].(string)
			err := model.SaveUploadDataTemplate(&data)
			if err != nil {
				return nil, err
			}
			return true, nil
		},
	}
}

func deleteTemplate() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"label": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete upload data template",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			err := model.DeleteUploadDataTemplate(params.Args["label"].(string))
			return err == nil, err
		},
	}
}
