package events

import (
	"github.com/Nicks344/moneytube/client/core/src/server/serverutils"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

const (
	TRStopped  = 10
	TRWorking  = 20
	TRStopping = 30
	TRReady    = 40
	TRError    = 50
)

type ToolsResultInput struct {
	Status      int
	Progress    int
	MaxProgress int
	Error       string
	JsonData    string
}

func OnAccountUpdated(account moneytubemodel.Account) {
	serverutils.OnGQLEvent("onAccountUpdated", "account", account)
}

func OnUploadTaskUpdated(task moneytubemodel.UploadTask) {
	serverutils.OnGQLEvent("onUploadTaskUpdated", "tasks", task)
}

func OnToolResult(tool string, input ToolsResultInput) {
	serverutils.OnGQLEvent("onToolResult", tool, input)
}

func OnGenerateVideoResult(input ToolsResultInput) {
	serverutils.OnGQLEvent("onGenerateVideoResult", "video", input)
}

func OnGenerateAudioResult(input ToolsResultInput) {
	serverutils.OnGQLEvent("onGenerateAudioResult", "audio", input)
}
