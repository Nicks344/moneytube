package clibridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	rpc "github.com/hprose/hprose-golang/rpc/websocket"
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/tools"
	"github.com/Nicks344/moneytube/client/core/src/modules/upload"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/mitchellh/mapstructure"
)

func Serve(port int) {
	service := rpc.NewWebSocketService()

	service.AddFunction("GetUploadTemplate", getUploadTemplate)
	service.AddFunction("RunTool", runTool)
	service.Publish("ToolResult", 0, 0)

	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), service); err != nil {
		log.Println(err)
	}
}

func getUploadTemplate(name string) (tmpl moneytubemodel.UploadDataTemplate, err error) {
	tmpl, err = model.GetUploadDataTemplateByName(name)
	return
}

func runTool(cmd string, input interface{}, context *rpc.WebSocketContext) (string, error) {
	switch cmd {
	case "AddUploadTask":
		return addUploadTask(input)

	case "StartUploadTask":
		err := startUploadTask(input, context)
		return "", err

	default:
		var data map[string]interface{}
		if err := mapstructure.Decode(input, &data); err != nil {
			return "", err
		}

		ts, err := tools.Start(cmd, data)
		if err != nil {
			return "", err
		}

		ready := false
		ts.OnToolResult(func(tool string, res events.ToolsResultInput) {
			if tool != cmd {
				return
			}

			context.Clients().Push("ToolResult", res)

			if res.Status == events.TRError {
				err = errors.New(res.Error)
			}

			if res.Status != events.TRWorking {
				ready = true
			}
		})

		for !ready {
			time.Sleep(1 * time.Second)
		}

		return "", err
	}
}

func startUploadTask(input interface{}, context *rpc.WebSocketContext) error {
	id, ok := input.(int)
	if !ok {
		return errors.New("cannot parse task id")
	}

	task, err := model.GetUploadTask(id)
	if err != nil {
		return err
	}

	taskResChan := upload.StartUploadTask(&task, 0)
	<-taskResChan

	errMsg := ""
	if task.Status == moneytubemodel.ASError {
		errMsg = task.ErrorMessage
	}
	context.Clients().Push("ToolResult", events.ToolsResultInput{
		Status:      task.Status,
		Progress:    task.Progress,
		MaxProgress: task.Count,
		Error:       errMsg,
	})

	return nil
}

func addUploadTask(input interface{}) (string, error) {
	var data moneytubemodel.UploadData

	if err := mapstructure.Decode(input, &data); err != nil {
		return "", err
	}

	tasks, err := model.SaveUploadData(&data)
	if err != nil {
		return "", err
	}

	for _, task := range tasks {
		events.OnUploadTaskUpdated(task)
	}

	ids := []int{}
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}

	res, err := json.Marshal(ids)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
