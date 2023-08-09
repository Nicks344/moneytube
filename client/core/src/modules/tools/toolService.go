package tools

import (
	"context"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/client/core/src/utils/videoeditor"
)

type ToolService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	toolName string

	progress    int
	maxProgress int
	isErr       bool

	toolResultListener func(tool string, input events.ToolsResultInput)
}

func (ts *ToolService) isCancelled(stop bool) bool {
	if utils.IsContextCancelled(ts.ctx) {
		if stop {
			ts.handleStopped()
		}
		return true
	}
	return false
}

func (ts *ToolService) handleStop() {
	ts.cancel()
}

func (ts *ToolService) handleStopped() {
	ts.handleToolResult(ts.toolName, events.ToolsResultInput{
		Status: events.TRStopped,
	})
}

func (ts *ToolService) handleError(err error) {
	if ts.isCancelled(false) {
		ts.handleStopped()
		return
	}

	ts.handleToolResult(ts.toolName, events.ToolsResultInput{
		Status: events.TRError,
		Error:  err.Error(),
	})

	ts.cancel()

	logger.Error(err)
	if genErr, ok := err.(*videoeditor.GenerateError); ok {
		logger.Warning(genErr.Output())
	}

	ts.isErr = true
}

func (ts *ToolService) handleProgressChanged(jsonData string) {
	if ts.isErr {
		return
	}
	status := events.TRWorking
	if ts.progress == ts.maxProgress {
		status = events.TRReady
	}

	ts.handleToolResult(ts.toolName, events.ToolsResultInput{
		Status:      status,
		Progress:    ts.progress,
		MaxProgress: ts.maxProgress,
		JsonData:    jsonData,
	})
}

func (ts *ToolService) handleToolResult(tool string, input events.ToolsResultInput) {
	if ts.toolResultListener != nil {
		ts.toolResultListener(tool, input)
	}

	events.OnToolResult(tool, input)
}

func (ts *ToolService) OnToolResult(f func(tool string, input events.ToolsResultInput)) {
	ts.toolResultListener = f
}
