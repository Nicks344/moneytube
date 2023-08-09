package gqlstructs

import (
	"github.com/graphql-go/graphql"
)

var TableUploadTaskOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "TableUploadTask",
		Fields: graphql.Fields{},
	},
)

var UploadTaskDataOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadTaskData",
		Fields: graphql.Fields{},
	},
)

var UploadDataTemplatesOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadDataTemplates",
		Fields: graphql.Fields{},
	},
)

var UploadOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadOptionsOutput",
		Fields: graphql.Fields{},
	})

var UploadVideoOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadVideoOptionsOutput",
		Fields: graphql.Fields{},
	})

var UploadTitlesOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadTitlesOptionsOutput",
		Fields: graphql.Fields{},
	})

var UploadEnvelopesOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadEnvelopesOptionsOutput",
		Fields: graphql.Fields{},
	})

var UploadCommentsOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadCommentsOptionsOutput",
		Fields: graphql.Fields{},
	})

var HintsOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "HintsOptionsOutput",
		Fields: graphql.Fields{},
	})

var HintOptionsOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "HintOptions",
		Fields: graphql.Fields{},
	})

var UploadDataSchedulerOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadDataSchedulerOutput",
		Fields: graphql.Fields{},
	})

var UploadTaskSchedulerOutput = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "UploadTaskSchedulerOutput",
		Fields: graphql.Fields{},
	})
