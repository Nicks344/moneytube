package tools

import (
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

const (
	DVMAllOnChannel = 10
	DVMByLinks      = 20
)

type DeleteVideoWorker struct {
	*ToolService

	AccountID  int
	Mode       int
	VideoLinks []string

	account  moneytubemodel.Account
	ybrowser *ybrowser.YoutubeBrowser
}

func (w *DeleteVideoWorker) Start() {
	var err error
	w.account, err = model.GetAccount(w.AccountID)
	if err != nil {
		w.handleError(err)
		return
	}

	if w.Mode == DVMByLinks {
		w.maxProgress = len(w.VideoLinks)
	}
	w.maxProgress = 1

	w.handleProgressChanged("")

	w.ybrowser, err = ybrowser.Start(w.ctx, &w.account)
	if err != nil {
		w.handleError(err)
		return
	}
	defer w.ybrowser.Stop()

	if w.isCancelled(true) {
		return
	}

	var ids []string

	switch w.Mode {
	case DVMAllOnChannel:
		ids, err = w.ybrowser.GetAllVideoIDs()
		if err != nil {
			w.handleError(err)
			return
		}

	case DVMByLinks:
		err = w.ybrowser.OpenVideoPage()
		if err != nil {
			w.handleError(err)
			return
		}

		stats, err := w.ybrowser.GetChannelStats()
		if err != nil {
			w.handleError(err)
			return
		}

		if stats.VideoCount == 0 {
			w.progress = 1
			w.handleProgressChanged("")
			return
		}
		ids = parseLinks.GetVideoIDsByLinks(w.VideoLinks)
	}

	if err := w.ybrowser.DeleteVideo(ids); err != nil {
		w.handleError(err)
	}

	w.progress = w.maxProgress
	w.handleProgressChanged("")
}
