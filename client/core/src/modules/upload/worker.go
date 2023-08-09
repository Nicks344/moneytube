package upload

import (
	"errors"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var uploaders = map[int]*VideoUploader{}

func StartUploadTaskByID(id int, uploads int) (done chan UploadResult, err error) {
	task, err := model.GetUploadTask(id)
	if err != nil {
		return
	}

	return StartUploadTask(&task, uploads), nil
}

type UploadResult struct {
	Uploads int
	Error   error
}

func StartUploadTask(task *moneytubemodel.UploadTask, count int) (done chan UploadResult) {
	done = make(chan UploadResult)

	uploader := NewUploader(task)
	uploaders[task.ID] = &uploader
	go func() {
		uploads, err := uploader.Run(count)
		delete(uploaders, task.ID)
		if err != nil {
			uploader.handleError(err)
		}

		done <- UploadResult{
			Uploads: uploads,
			Error:   err,
		}
	}()

	return
}

func StopUploadTask(id int) error {
	uploader, ok := uploaders[id]
	if !ok {
		return errors.New("cannot find uploader")
	}
	uploader.Stop()
	return nil
}
