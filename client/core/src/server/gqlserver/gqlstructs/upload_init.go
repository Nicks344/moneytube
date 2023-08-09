package gqlstructs

import (
	"github.com/graphql-go/graphql"
)

var uploadOutput = map[*graphql.Object]FieldsList{
	TableUploadTaskOutput:        TableUploadTaskFields,
	UploadTaskDataOutput:         UploadTaskDataFields,
	UploadOptionsOutput:          UploadVideoOptionsFields,
	UploadTitlesOptionsOutput:    UploadTitlesOptionsFields,
	UploadVideoOptionsOutput:     UploadVideoOptionsFields,
	UploadEnvelopesOptionsOutput: UploadEnvelopesOptionsFields,
	UploadCommentsOptionsOutput:  UploadCommentsOptionsFields,
	HintsOptionsOutput:           HintsOptionsFields,
	HintOptionsOutput:            HintOptionsFields,
	UploadDataTemplatesOutput:    UploadDataTemplateFields,
	UploadDataSchedulerOutput:    UploadDataSchedulerFields,
	UploadTaskSchedulerOutput:    UploadTaskSchedulerFields,
}

var uploadInput = map[*graphql.InputObject]FieldsList{
	UploadTaskDataInput:         UploadTaskDataFields,
	UploadDataTemplateInput:     UploadDataTemplateFields,
	UploadOptionsInput:          UploadVideoOptionsFields,
	UploadTitlesOptionsInput:    UploadTitlesOptionsFields,
	UploadVideoOptionsInput:     UploadVideoOptionsFields,
	UploadEnvelopesOptionsInput: UploadEnvelopesOptionsFields,
	UploadCommentsOptionsInput:  UploadCommentsOptionsFields,
	HintsOptionsInput:           HintsOptionsFields,
	HintOptionsInput:            HintOptionsFields,
	UploadDataSchedulerInput:    UploadDataSchedulerFields,
}

func init() {
	for obj, fields := range uploadOutput {
		initOutputSchema(obj, fields, false)
	}

	for obj, fields := range uploadInput {
		initInputSchema(obj, fields)
	}
}
