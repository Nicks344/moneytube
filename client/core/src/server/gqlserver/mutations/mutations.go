package mutations

import (
	"github.com/graphql-go/graphql"
)

func GetMutations() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addOrEditAccount":         addOrEditAccount(),
			"deleteAccount":            deleteAccount(),
			"updateAccount":            updateAccount(),
			"updateAllAccounts":        updateAllAccounts(),
			"openAccountBrowser":       openAccountBrowser(),
			"addUploadTasks":           addUploadTasks(),
			"deleteUploadTask":         deleteUploadTask(),
			"deleteAllUploadTasks":     deleteAllUploadTasks(),
			"startUploadTask":          startUploadTask(),
			"stopUploadTask":           stopUploadTask(),
			"startAllUploadTasks":      startAllUploadTasks(),
			"stopAllUploadTasks":       stopAllUploadTasks(),
			"startGetLinks":            startGetLinks(),
			"startChangeDescriptions":  startChangeDescriptions(),
			"startCreatePlaylist":      startCreatePlaylist(),
			"startDeleteVideo":         startDeleteVideo(),
			"startComment":             startComment(),
			"startGenerateVideo":       startGenerateVideo(),
			"startGenerateAudio":       startGenerateAudio(),
			"startGenerateImages":      startGenerateImages(),
			"saveSettings":             saveSettings(),
			"addOrEditMacros":          addOrEditMacros(),
			"deleteMacros":             deleteMacros(),
			"cancelTool":               cancelTool(),
			"startGenerateCopies":      startGenerateCopies(),
			"startGenerateVideoFFmpeg": startGenerateVideoFFmpeg(),
			"sendBugReport":            sendBugReport(),
			"clearTemp":                clearTemp(),
			"saveTemplate":             saveTemplate(),
			"deleteTemplate":           deleteTemplate(),
			"deleteGroup":              deleteGroup(),
			"importCookies":            importCookies(),
		},
	})
}
