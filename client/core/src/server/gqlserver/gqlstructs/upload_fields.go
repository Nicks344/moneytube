package gqlstructs

import (
	"github.com/graphql-go/graphql"
)

var TableUploadTaskFields = FieldsList{
	"ID":           Uint,
	"DetailsID":    Uint,
	"Account":      AccountOutput,
	"Progress":     graphql.Int,
	"Count":        graphql.Int,
	"Status":       graphql.Int,
	"ErrorMessage": graphql.String,
	"IsScheduled":  graphql.Boolean,
	"ScheduleTime": graphql.String,
	"Scheduler": FieldTypes{
		Out: UploadTaskSchedulerOutput,
	},
}

var UploadTaskSchedulerFields = FieldsList{
	"Enabled":         graphql.Boolean,
	"SecondStartTime": graphql.String,
	"Progress":        graphql.Int,
}

var UploadTaskDataFields = FieldsList{
	"ID":         Uint,
	"AccountIDs": graphql.NewList(Uint),
	"template":   UploadDataFields,
}

var UploadDataTemplateFields = FieldsList{
	"Label":    graphql.String,
	"template": UploadDataFields,
}

var UploadDataFields = FieldsList{
	"WithProcessing":    graphql.Boolean,
	"DisableComments":   graphql.Boolean,
	"AgeRestrictions":   graphql.Boolean,
	"OrderComments":     graphql.Boolean,
	"ShowRating":        graphql.Boolean,
	"NotifySubscribers": graphql.Boolean,
	"ShowStatistic":     graphql.Boolean,
	"Language":          graphql.String,
	"Category":          graphql.String,
	"PauseFrom":         graphql.Int,
	"PauseTo":           graphql.Int,
	"UploadCountFrom":   graphql.Int,
	"UploadCountTo":     graphql.Int,
	"CommentMode":       graphql.Int,
	"Videos": FieldTypes{
		In:  UploadVideoOptionsInput,
		Out: UploadVideoOptionsOutput,
	},
	"Descriptions": FieldTypes{
		In:  UploadOptionsInput,
		Out: UploadOptionsOutput,
	},
	"Envelopes": FieldTypes{
		In:  UploadEnvelopesOptionsInput,
		Out: UploadEnvelopesOptionsOutput,
	},
	"Titles": FieldTypes{
		In:  UploadTitlesOptionsInput,
		Out: UploadTitlesOptionsOutput,
	},
	"Tags": FieldTypes{
		In:  UploadOptionsInput,
		Out: UploadOptionsOutput,
	},
	"Comments": FieldTypes{
		In:  UploadCommentsOptionsInput,
		Out: UploadCommentsOptionsOutput,
	},
	"Hints": FieldTypes{
		In:  HintsOptionsInput,
		Out: HintsOptionsOutput,
	},
	"IsScheduled":                  graphql.Boolean,
	"ScheduleTime":                 graphql.String,
	"ScheduleStep":                 graphql.Int,
	"SkipErrors":                   graphql.Boolean,
	"WaitVideoInFolder":            graphql.Boolean,
	"ClearFilesAfterSuccessUpload": graphql.Boolean,
	"IsDeferred":                   graphql.Boolean,
	"DeferTime":                    graphql.String,
	"DeferStep":                    graphql.Int,
	"Scheduler": FieldTypes{
		In:  UploadDataSchedulerInput,
		Out: UploadDataSchedulerOutput,
	},
}

var UploadDataSchedulerFields = FieldsList{
	"Enabled":     graphql.Boolean,
	"StartTime":   graphql.String,
	"Interval":    graphql.Int,
	"UploadCount": graphql.Int,
	"StopVariant": graphql.Int,
	"StopCount":   graphql.Int,
}

var UploadOptionsFields = FieldsList{
	"List":     graphql.NewList(graphql.String),
	"Cycle":    graphql.Boolean,
	"IsRandom": graphql.Boolean,
}

var UploadVideoOptionsFields = FieldsList{
	"RenameToTitle": graphql.Boolean,
	"options":       UploadOptionsFields,
}

var UploadTitlesOptionsFields = FieldsList{
	"IsGetFilename": graphql.Boolean,
	"options":       UploadOptionsFields,
}

var UploadEnvelopesOptionsFields = FieldsList{
	"IsRandomFromPropose": graphql.Boolean,
	"options":             UploadOptionsFields,
}

var UploadCommentsOptionsFields = FieldsList{
	"AddComment": graphql.Boolean,
	"FixComment": graphql.Boolean,
	"options":    UploadOptionsFields,
}

var HintsOptionsFields = FieldsList{
	"HintsList": FieldTypes{
		Out: graphql.NewList(HintOptionsOutput),
		In:  graphql.NewList(HintOptionsInput),
	},
	"AddHints": graphql.Boolean,
}

var HintOptionsFields = FieldsList{
	"Type": graphql.Int,
	"Time": graphql.String,
	"Data": FieldTypes{
		Out: UploadOptionsOutput,
		In:  UploadOptionsInput,
	},
	"Message": FieldTypes{
		Out: UploadOptionsOutput,
		In:  UploadOptionsInput,
	},
	"Teaser": FieldTypes{
		Out: UploadOptionsOutput,
		In:  UploadOptionsInput,
	},
}
