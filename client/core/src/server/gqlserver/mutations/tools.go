package mutations

import (
	"github.com/Nicks344/moneytube/client/core/src/modules/tools"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
)

func startGetLinks() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"channels": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
		Type:        graphql.Boolean,
		Description: "Start get links",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GetLinks", params.Args)
			return err != nil, err
		},
	}
}

func startChangeDescriptions() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"accountID": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"mode": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"description": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadOptionsInput,
			},
			"videoLinks": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
			"playlistLink": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start change descriptions",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("ChangeDescription", params.Args)
			return err != nil, err
		},
	}
}

func startCreatePlaylist() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"accountIDs": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
			"mode": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"description": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadOptionsInput,
			},
			"videoLinks": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadOptionsInput,
			},
			"channelLink": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"name": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadOptionsInput,
			},
			"playlistCountFrom": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"playlistCountTo": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoCountFrom": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoCountTo": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start create playlist",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("CreatePlaylist", params.Args)
			return err != nil, err
		},
	}
}

func startDeleteVideo() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"accountID": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"mode": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoLinks": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
		Type:        graphql.Boolean,
		Description: "Start delete video",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("DeleteVideo", params.Args)
			return err != nil, err
		},
	}
}

func startComment() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"accountIDs": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
			"maxComments": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"comment": &graphql.ArgumentConfig{
				Type: gqlstructs.UploadOptionsInput,
			},
			"like": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"videoLinks": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
		Type:        graphql.Boolean,
		Description: "Start comment video",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("CommentAndLike", params.Args)
			return err != nil, err
		},
	}
}

func cancelTool() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"tool": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start get links",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			tools.Cancel(params.Args["tool"].(string))
			return true, nil
		},
	}
}
