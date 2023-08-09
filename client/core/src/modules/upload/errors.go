package upload

import (
	"errors"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var limitError = errors.New("Загрузка недоступна, временная блокировка от YouTube")

func (vu *VideoUploader) handleError(err error) {
	if utils.IsContextCancelled(vu.ctx) {
		vu.task.Status = moneytubemodel.UTSStopped
		vu.saveTask()
		return
	}

	vu.cancel()
	logger.Error(err)
	vu.task.ErrorMessage = err.Error()
	vu.task.Status = moneytubemodel.UTSError
	vu.saveTask()
}
