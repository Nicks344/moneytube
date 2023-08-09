package queries

import (
	"errors"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
	"github.com/spf13/viper"
)

func settings() *graphql.Field {
	return &graphql.Field{
		Type:        gqlstructs.SettingsOutput,
		Description: "Get settings",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			settings := viper.AllSettings()
			return settings, nil
		},
	}
}

func macroses() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(graphql.String),
		Description: "Get macroses",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			macroses, err := model.GetMacroses()
			if err != nil {
				logger.Error(err)
				return nil, errors.New("ошибка при получении макросов")
			}

			names := make([]string, len(macroses), len(macroses))
			for i, m := range macroses {
				names[i] = m.Name
			}
			return names, nil
		},
	}
}

func macros() *graphql.Field {
	return &graphql.Field{
		Type: gqlstructs.MacrosOutput,
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Description: "Get macros",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			name := p.Args["name"].(string)
			macros, err := model.GetMacros(name)
			if err != nil {
				logger.Error(err)
				return nil, errors.New("ошибка при получении макроса")
			}

			return macros, nil
		},
	}
}
