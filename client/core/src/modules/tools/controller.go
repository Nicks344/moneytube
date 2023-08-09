package tools

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"

	"github.com/mitchellh/mapstructure"
)

const (
	GenerateCopies      = "GenerateCopies"
	GenerateVideo       = "GenerateVideo"
	GenerateVideoFFmpeg = "GenerateVideoFFmpeg"
	GenerateImages      = "GenerateImages"
	GenerateAudio       = "GenerateAudio"
	CommentAndLike      = "CommentAndLike"
	ChangeDescription   = "ChangeDescription"
	CreatePlaylist      = "CreatePlaylist"
	DeleteVideo         = "DeleteVideo"
	GetLinks            = "GetLinks"
)

var toolServices = map[string]*ToolService{}

func Start(tool string, args map[string]interface{}) (ts *ToolService, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	if ts, ok := toolServices[tool]; ok {
		ts.handleStop()
	}

	ts = &ToolService{
		ctx:      ctx,
		cancel:   cancel,
		toolName: tool,
	}

	toolServices[tool] = ts

	switch tool {
	case GenerateVideo:
		worker := VideoGeneratorWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		if dataJson, ok := args["dataJson"]; ok {
			err = json.Unmarshal([]byte(dataJson.(string)), &worker.Data)
			if err != nil {
				return
			}
		}
		go worker.Start()
		return

	case GenerateImages:
		worker := ImagesGeneratorWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		if dataJson, ok := args["dataJson"]; ok {
			err = json.Unmarshal([]byte(dataJson.(string)), &worker.Data)
			if err != nil {
				return
			}
		}
		go worker.Start()
		return

	case GenerateAudio:
		worker := AudioGeneratorWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case GenerateCopies:
		worker := CopiesGeneratorWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case CommentAndLike:
		worker := CommentAndLikeWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case ChangeDescription:
		worker := ChangeDescriptionWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case CreatePlaylist:
		worker := CreatePlaylistWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case DeleteVideo:
		worker := DeleteVideoWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case GetLinks:
		worker := GetLinksWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	case GenerateVideoFFmpeg:
		worker := VideoFFmpegGeneratorWorker{
			ToolService: ts,
		}
		err = mapstructure.Decode(args, &worker)
		if err != nil {
			return
		}
		go worker.Start()
		return

	default:
		err = errors.New("tool not found")
		return
	}

}

func Cancel(tool string) error {
	if ts, ok := toolServices[tool]; ok {
		ts.handleStop()
		delete(toolServices, tool)
		events.OnToolResult(tool, events.ToolsResultInput{
			Status: events.TRStopping,
		})
		if tool == GenerateVideo {
			stopAEProcess()
		}
		return nil
	}

	return errors.New("tool not started")
}
