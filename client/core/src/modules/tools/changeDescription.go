package tools

import (
	"errors"
	"fmt"

	"github.com/meandrewdev/logger"

	"strings"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/chromedp"
)

const (
	CDMAllOnChannel   = 10
	CDMByLinks        = 20
	CDMByPlaylistLink = 30
)

type ChangeDescriptionWorker struct {
	*ToolService

	AccountID    int
	Mode         int
	Description  moneytubemodel.UploadOptions
	VideoLinks   []string
	PlaylistLink string

	account  moneytubemodel.Account
	ybrowser *ybrowser.YoutubeBrowser
}

func (w *ChangeDescriptionWorker) Start() {
	var err error
	w.account, err = model.GetAccount(w.AccountID)
	if err != nil {
		w.handleError(err)
		return
	}

	switch w.Mode {
	case CDMAllOnChannel, CDMByPlaylistLink:
		w.maxProgress = 1

	case CDMByLinks:
		w.maxProgress = len(w.VideoLinks)
	}

	w.handleProgressChanged("")

	w.Description.ClearFromEmpty()

	w.ybrowser, err = ybrowser.Start(w.ctx, &w.account)
	if err != nil {
		w.handleError(err)
		return
	}
	defer w.ybrowser.Stop()

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

	var videoIDs []string

	logger.Notice("get video ids to change descriptions")
	switch w.Mode {
	case CDMAllOnChannel:
		if stats.VideoCount == 0 {
			w.handleError(errors.New("video list is empty"))
			return
		}
		videoIDs, err = w.ybrowser.GetAllVideoIDs()
		if err != nil {
			w.handleError(err)
			return
		}

	case CDMByLinks:
		videoIDs = parseLinks.GetVideoIDsByLinks(w.VideoLinks)

	case CDMByPlaylistLink:
		playlistID, err := parseLinks.GetPlaylistIDByLink(w.PlaylistLink)
		if err != nil {
			w.handleError(err)
			return
		}
		videoIDs, err = serverAPI.GetVideoIDsByPlaylistID(playlistID)
		if err != nil {
			w.handleError(err)
			return
		}
	}

	w.maxProgress = len(videoIDs)
	w.handleProgressChanged("")

	err = w.changeByIDs(videoIDs)
	if err != nil {
		w.handleError(err)
	}
}

func (w *ChangeDescriptionWorker) changeByIDs(ids []string) error {
	for _, id := range ids {
		description, err := w.Description.GetOne()
		if err != nil {
			return errors.New("Закончились описания")
		}

		macrosData := macros.StaticMacroses{
			ChannelTitle: w.account.ChannelName,
			ChannelLink:  w.account.GetChannelLink(),
		}

		macrosData.VideoTitle = strings.Trim(strings.ReplaceAll(macrosData.VideoTitle, "\n", ""), " ")
		macrosData.VideoLink = strings.Trim(strings.ReplaceAll(macrosData.VideoLink, "\n", ""), " ")
		macrosData.VideoTags = strings.Trim(strings.ReplaceAll(macrosData.VideoTags, "\n", ""), " ")

		description = macros.Execute(description, macrosData)

		w.ybrowser.GoTo(fmt.Sprintf("https://studio.youtube.com/video/%s/edit", id))

		logger.Notice("get video data: title, link and tags")
		oldDesc := ""
		err = w.ybrowser.Browser.Run(chromedp.WaitVisible(".title #textbox", chromedp.ByQuery),
			chromedp.TextContent(".title #textbox", &macrosData.VideoTitle, chromedp.ByQuery),
			chromedp.TextContent("a.ytcp-video-info", &macrosData.VideoLink, chromedp.ByQuery),
			chromedp.EvaluateAsDevTools("Array.from(document.querySelectorAll('#chip-bar #chip-text')).map(t => t.innerText).join(', ')", &macrosData.VideoTags),
			chromedp.TextContent(".description #textbox", &oldDesc, chromedp.ByQuery))
		if err != nil {
			return err
		}

		if description != oldDesc {
			logger.Notice("change description")
			err = w.ybrowser.Browser.Run(chromedp.SetJavascriptAttribute(".description #textbox", "innerHTML", "", chromedp.ByQuery),
				chromedp.SendKeys(".description #textbox", macros.Execute(description, macrosData), chromedp.ByQuery),
				chromedp.Click("ytcp-button#save", chromedp.ByQuery),
				chromedp.WaitVisible("ytcp-button#save[disabled]", chromedp.ByQuery))
			if err != nil {
				return err
			}
		}

		w.progress++
		w.handleProgressChanged("")
	}
	return nil
}
