package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"time"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/tidwall/gjson"

	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"

	"github.com/chromedp/chromedp"
)

const (
	CPMAllOnChannel  = 10
	CPMByLinks       = 20
	CPMByChannelLink = 30
)

type CreatePlaylistWorker struct {
	*ToolService

	AccountIDs        []int
	Mode              int
	Description       moneytubemodel.UploadOptions
	VideoLinks        moneytubemodel.UploadOptions
	ChannelLink       string
	Name              moneytubemodel.UploadOptions
	PlaylistCountFrom int
	PlaylistCountTo   int
	VideoCountFrom    int
	VideoCountTo      int

	result string
}

func (w *CreatePlaylistWorker) Start() {
	w.maxProgress = len(w.AccountIDs)
	w.handleProgressChanged("")

	switch w.Mode {
	case CPMByChannelLink:
		channelID, err := parseLinks.GetChannelIDByLink(w.ChannelLink)
		if err != nil {
			w.handleError(errors.New("Неверно указана ссылка на канал"))
			return
		}
		w.VideoLinks.List, err = serverAPI.GetVideoIDsByChannelID(channelID)
		if err != nil {
			w.handleError(err)
			return
		}

	case CPMByLinks:
		w.VideoLinks.List = utils.ClearSliceFromEmpty(w.VideoLinks.List)
		w.VideoLinks.List = parseLinks.GetVideoIDsByLinks(w.VideoLinks.List)
	}

	w.Description.List = utils.ClearSliceFromEmpty(w.Description.List)
	w.Name.List = utils.ClearSliceFromEmpty(w.Name.List)

	for _, accountID := range w.AccountIDs {
		if err := w.processAccount(accountID); err != nil {
			w.handleError(err)
			return
		}

		w.progress++
		w.handleProgressChanged(w.result)
	}
}

func (w *CreatePlaylistWorker) processAccount(accountID int) error {
	account, err := model.GetAccount(accountID)
	if err != nil {
		return err
	}

	w.result += fmt.Sprintf("Аккаунт %s:\r\n", account.Login)

	count := utils.RandRange(w.PlaylistCountFrom, w.PlaylistCountTo)

	for i := 0; i < count; i++ {
		count, err := w.newPlaylist(account)
		if err != nil {
			w.result += fmt.Sprintf("Ошибка: %s\r\n", err.Error())
			return err
		}
		w.result += fmt.Sprintf("Создан плейлист с %d видео\r\n", count)
	}

	return nil
}

func (w *CreatePlaylistWorker) newPlaylist(account moneytubemodel.Account) (count int, err error) {
	macrosData := macros.StaticMacroses{
		ChannelTitle: account.ChannelName,
		ChannelLink:  account.GetChannelLink(),
	}

	description, err := w.Description.GetOne()
	if err != nil {
		err = errors.New("закончились описания")
		return
	}
	description = macros.Execute(description, macrosData)

	name, err := w.Name.GetOne()
	if err != nil {
		err = errors.New("закончились названия")
		return
	}
	name = macros.Execute(name, macrosData)

	count = utils.RandRange(w.VideoCountFrom, w.VideoCountTo)

	videoIDs := make([]string, count)
	for i := range videoIDs {
		videoIDs[i], err = w.VideoLinks.GetOne()
		if err != nil {
			err = errors.New("закончились видео")
			return
		}
	}

	err = w.createPlaylist(account, videoIDs, name, description)
	return
}

func (w *CreatePlaylistWorker) createPlaylist(account moneytubemodel.Account, videoIDs []string, name, description string) (err error) {
	ybrowser, err := ybrowser.Start(w.ctx, &account)
	if err != nil {
		return
	}
	defer ybrowser.Stop()

	ctx, cancel := context.WithTimeout(w.ctx, time.Minute*20)
	defer cancel()
	ybrowser.Browser.BlockRequest(ctx, "playlist/create")

	reqCh := ybrowser.ListenForRequests(ctx, "playlist/create")

	var added bool
	for _, videoID := range videoIDs {

		err = ybrowser.GoTo("https://www.youtube.com/watch?v=" + videoID)
		if err != nil {
			return
		}

		var disabled bool
		err = ybrowser.Browser.Run(
			drivenbrowser.NewLogAction(chromedp.WaitEnabled("//yt-formatted-string[text()='Сохранить']/../..", chromedp.BySearch), "wait for menu"),
		)
		if err != nil {
			logger.Error(fmt.Errorf("cannot found button 'Сохранить' '%s'", videoID))
			return
		}

		if disabled {
			logger.Error(fmt.Errorf("cannot add to playlist video with id '%s'", videoID))
			continue
		}

		time.Sleep(5 * time.Second)

		err = ybrowser.Browser.Run(
			drivenbrowser.NewLogAction(chromedp.Click("#button[aria-label='Добавить в плейлист']", chromedp.ByQuery), "click add to playlist"),
			drivenbrowser.NewLogAction(chromedp.Click(".ytd-add-to-playlist-create-renderer", chromedp.ByQuery), "click to create playlist"),
			drivenbrowser.NewLogAction(chromedp.SendKeys("#labelAndInputContainer input", name, chromedp.ByQuery), "set title"),
			drivenbrowser.NewLogAction(chromedp.Click("#actions.ytd-add-to-playlist-create-renderer #button", chromedp.ByQuery), "click to submit playlist"),
		)
		if err != nil {
			return
		}

		added = true
		break
	}

	if !added {
		err = errors.New("unknown error while create playlist")
	}

	reqData := <-reqCh

	var data CreatePlaylistData
	err = json.Unmarshal([]byte(reqData.Request.PostData), &data)
	if err != nil {
		return
	}

	data.Description = description
	data.PrivacyStatus = "PUBLIC"
	data.VideoIds = videoIDs

	var headers req.Header
	err = mapstructure.Decode(reqData.Request.Headers, &headers)
	if err != nil {
		return
	}

	cookiesArr, err := ybrowser.Browser.GetCookies()
	if err != nil {
		return
	}
	cookies := ""
	for _, cookie := range cookiesArr {
		if strings.Contains(cookie.Domain, "youtube") {
			cookies += cookie.Name + "=" + cookie.Value + "; "
		}
	}
	headers["Cookie"] = cookies

	request := req.New()
	if ybrowser.Account.Proxy != "" {
		request.SetProxyUrl(ybrowser.Account.Proxy)
	}
	resp, err := request.Post(reqData.Request.URL, headers, req.BodyJSON(data))
	if err != nil {
		return
	}

	body := resp.String()

	playlistID := gjson.Get(body, "playlistId").String()
	if playlistID == "" {
		errText := gjson.Get(body, "error.message").String()
		if errText != "" {
			err = errors.New(errText)
			return
		}

		data.Context = nil
		d, _ := json.Marshal(data)
		logger.Warning("unknown error while create playlist\r\nbody:\r\n" + string(d) + "\r\nanswer:\r\n" + body)
		err = errors.New("unknown error")
		return
	}

	return
}

type CreatePlaylistData struct {
	Context       interface{} `json:"context"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	PrivacyStatus string      `json:"privacyStatus"`
	VideoIds      []string    `json:"videoIds"`
}
