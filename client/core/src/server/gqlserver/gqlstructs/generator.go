package gqlstructs

import "github.com/graphql-go/graphql"

var LayerInfoInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "LayerInfo",
		Fields: graphql.InputObjectConfigFieldMap{
			"type": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

var LayerDataInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "LayerData",
		Fields: graphql.InputObjectConfigFieldMap{
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"data": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
		},
	})

var ImageOverlayDataInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "ImageOverlayData",
		Fields: graphql.InputObjectConfigFieldMap{
			"overlaySrc": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"enabled": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"x": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"y": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"from": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"to": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	})

var TextOverlayDataInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "TextOverlayData",
		Fields: graphql.InputObjectConfigFieldMap{
			"text": &graphql.InputObjectFieldConfig{
				Type: UploadOptionsInput,
			},
			"enabled": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"x": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"y": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"from": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"to": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"background": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"color": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"backgroundColor": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"font": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"size": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"bold": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"italic": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
	})

var GSpeechPropsInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GSpeechProps",
		Fields: graphql.InputObjectConfigFieldMap{
			"profiles": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
			"pitch": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	})
