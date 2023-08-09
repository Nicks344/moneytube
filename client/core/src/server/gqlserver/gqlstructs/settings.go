package gqlstructs

import "github.com/graphql-go/graphql"

var SettingsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Settings",
		Fields: graphql.Fields{
			"ae_exe": &graphql.Field{
				Type: graphql.String,
			},
			"ae_memory_persent": &graphql.Field{
				Type: graphql.Int,
			},
			"ys_api_key": &graphql.Field{
				Type: graphql.String,
			},
			"gs_api_key": &graphql.Field{
				Type: graphql.String,
			},
			"vrs_api_key": &graphql.Field{
				Type: graphql.String,
			},
			"yt_api_key": &graphql.Field{
				Type: graphql.String,
			},
			"show_browser": &graphql.Field{
				Type: graphql.Boolean,
			},
			"ae_lang": &graphql.Field{
				Type: graphql.String,
			},
			"speech_pro_login": &graphql.Field{
				Type: graphql.String,
			},
			"speech_pro_id": &graphql.Field{
				Type: graphql.String,
			},
			"speech_pro_password": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var MacrosInput = graphql.FieldConfigArgument{
	"macros": &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name: "MacrosInput",
				Fields: graphql.InputObjectConfigFieldMap{
					"ID": &graphql.InputObjectFieldConfig{
						Type: graphql.Int,
					},
					"Name": &graphql.InputObjectFieldConfig{
						Type: graphql.String,
					},
					"Data": &graphql.InputObjectFieldConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
			}),
	},
}

var MacrosOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "MacrosOutput",
		Fields: graphql.Fields{
			"ID": &graphql.Field{
				Type: Uint,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Data": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)
