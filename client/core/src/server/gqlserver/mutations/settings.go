package mutations

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func saveSettings() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"configJSON": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Save settings",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			configStr := params.Args["configJSON"].(string)
			var configMap map[string]interface{}
			err := json.Unmarshal([]byte(configStr), &configMap)
			if err != nil {
				return false, err
			}
			viper.MergeConfigMap(configMap)
			err = viper.WriteConfig()
			if err != nil {
				return false, err
			}
			return true, nil
		},
	}
}

func addOrEditMacros() *graphql.Field {
	return &graphql.Field{
		Args:        gqlstructs.MacrosInput,
		Type:        gqlstructs.MacrosOutput,
		Description: "Add new or edit macros",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var macros moneytubemodel.Macros
			mapstructure.Decode(params.Args["macros"], &macros)
			err := model.SaveMacros(macros)
			if err != nil {
				return nil, err
			}
			return macros, nil
		},
	}
}

func deleteMacros() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Delete macros",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			err := model.DeleteMacros(p.Args["name"].(string))
			return err == nil, err
		},
	}
}

func clearTemp() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.Boolean,
		Description: "Clear temp dir",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			if err := os.RemoveAll(paths.Temp); err != nil {
				return false, err
			}

			time.Sleep(500 * time.Millisecond)

			if err := os.MkdirAll(paths.Temp, 0666); err != nil {
				return false, err
			}
			return true, nil
		},
	}
}
