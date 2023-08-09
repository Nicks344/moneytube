package mutations

import (
	"github.com/Nicks344/moneytube/client/core/src/modules/tools"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/gqlstructs"

	"github.com/graphql-go/graphql"
)

func startGenerateVideo() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"aepFile": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"compositionName": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"fromRenderQueue": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"resolution": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"resultPath": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"ext": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"outputExt": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"quality": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"isCompress": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"layers": &graphql.ArgumentConfig{
				Type: graphql.NewList(gqlstructs.LayerInfoInput),
			},
			"dataJson": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start generate video",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GenerateVideo", params.Args)
			return err != nil, err
		},
	}
}

func startGenerateImages() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"psdFile": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"resultPath": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"layers": &graphql.ArgumentConfig{
				Type: graphql.NewList(gqlstructs.LayerInfoInput),
			},
			"dataJson": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start generate images",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GenerateImages", params.Args)
			return err != nil, err
		},
	}
}

func startGenerateAudio() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"api": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"cropAudio": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"cropAudioVariant": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"maxTime": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"textFiles": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.String),
			},
			"resultPath": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"lang": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"voice": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"speed": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"threads": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"gspeechProps": &graphql.ArgumentConfig{
				Type: gqlstructs.GSpeechPropsInput,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start generate audio",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GenerateAudio", params.Args)
			return err != nil, err
		},
	}
}

func startGenerateCopies() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"resultPath": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"count": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"imageOverlay": &graphql.ArgumentConfig{
				Type: gqlstructs.ImageOverlayDataInput,
			},
			"textOverlay": &graphql.ArgumentConfig{
				Type: gqlstructs.TextOverlayDataInput,
			},
			"intro": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"introDurationSec": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"outro": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"outroDurationSec": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"cutMethod": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"cutSecondsFrom": &graphql.ArgumentConfig{
				Type: graphql.Float,
			},
			"cutSecondsTo": &graphql.ArgumentConfig{
				Type: graphql.Float,
			},
			"cutParts": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start generate copies",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GenerateCopies", params.Args)
			return err != nil, err
		},
	}
}

func startGenerateVideoFFmpeg() *graphql.Field {
	return &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"resultPath": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"ext": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"deleteImages": &graphql.ArgumentConfig{
				Type: graphql.Boolean,
			},
			"audioFile": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"ifVideoLonger": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"ifAudioLonger": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"intro": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"outro": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"slideDurationFrom": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"slideDurationTo": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoDurationType": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoDurationFrom": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoDurationTo": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"slideCountFrom": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"slideCountTo": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"videoCount": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"threads": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"fps": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"overlay": &graphql.ArgumentConfig{
				Type: gqlstructs.ImageOverlayDataInput,
			},
		},
		Type:        graphql.Boolean,
		Description: "Start generate video(ffmpeg)",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			_, err := tools.Start("GenerateVideoFFmpeg", params.Args)
			return err != nil, err
		},
	}
}
