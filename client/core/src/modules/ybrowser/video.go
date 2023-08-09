package ybrowser

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/meandrewdev/logger"

	"github.com/tidwall/gjson"

	"strconv"
	"strings"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"

	"github.com/chromedp/chromedp"
	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
)

const VideosOnPage = 30

func (ybrowser *YoutubeBrowser) OpenVideoPage() error {
	location, err := ybrowser.GetLocation()
	if err != nil {
		return err
	}

	videoPageURL := fmt.Sprintf("https://studio.youtube.com/channel/%s/videos/upload", ybrowser.Account.ChannelID)

	if !strings.Contains(location.String(), videoPageURL) {
		err = ybrowser.GoTo(videoPageURL)
		time.Sleep(1 * time.Second)
	}

	return err
}

func (ybrowser *YoutubeBrowser) GoToSecondPage() (isLastPage bool, err error) {
	var temp string
	err = ybrowser.Browser.Run(chromedp.AttributeValue("navigate-after", "disabled", &temp, &isLastPage, chromedp.ByID))
	if err != nil {
		return
	}
	if isLastPage {
		return
	}
	err = ybrowser.Browser.Run(chromedp.ScrollIntoView("navigate-after", chromedp.ByID))
	if err != nil {
		return
	}
	err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Click("navigate-after", chromedp.ByID), "click to next page"))
	return
}

type ChannelStats struct {
	VideoCount           int
	SubscriberCount      int
	TotalVideoViewsCount int
}

func (ybrowser *YoutubeBrowser) GetChannelStats() (stats ChannelStats, err error) {
	logger.Notice("get channel stats")
	err = ybrowser.OpenVideoPage()
	if err != nil {
		return
	}

	var pageStats struct {
		VideoCount           string `json:"videoCount"`
		SubscriberCount      string `json:"subscriberCount"`
		TotalVideoViewsCount string `json:"totalVideoViewCount"`
	}

	err = ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("window.yt.config_.CHUNKED_PREFETCH_DATA[0].data.then(data => window.channelData = data)", &[]byte{}))
	if err != nil {
		return
	}

	err = ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("window.channelData.metric", &pageStats))
	if err != nil {
		return
	}

	stats.VideoCount, err = strconv.Atoi(pageStats.VideoCount)
	if err != nil {
		return
	}

	stats.SubscriberCount, err = strconv.Atoi(pageStats.SubscriberCount)
	if err != nil {
		return
	}

	stats.TotalVideoViewsCount, err = strconv.Atoi(pageStats.TotalVideoViewsCount)
	return
}

type ListVideos struct {
	Context  interface{} `json:"context"`
	Order    string      `json:"order"`
	PageSize int         `json:"pageSize"`
	Mask     interface{} `json:"mask"`
	Filter   interface{} `json:"filter"`
}

func (ybrowser *YoutubeBrowser) GetAllVideoIDs() (ids []string, err error) {
	reqCh := ybrowser.ListenForRequests(ybrowser.Browser.Ctx, "list_creator_videos")

	if err := ybrowser.GoTo(fmt.Sprintf("https://studio.youtube.com/channel/%s/videos/upload", ybrowser.Account.ChannelID)); err != nil {
		return nil, err
	}

	reqData := <-reqCh

	var data ListVideos
	err = json.Unmarshal([]byte(reqData.Request.PostData), &data)
	if err != nil {
		return
	}

	data.PageSize = 100000

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
	jsonArr := gjson.Get(body, "videos").Array()
	for _, jstr := range jsonArr {
		ids = append(ids, jstr.Get("videoId").String())
	}

	return
}

type DeleteVideosReq struct {
	Context interface{}     `json:"context"`
	Videos  DeleteVideosIds `json:"videos"`
}

type DeleteVideosIds struct {
	VideoIds []string `json:"videoIds"`
}

func (ybrowser *YoutubeBrowser) DeleteVideo(ids []string) (err error) {
	if err := ybrowser.OpenVideoPage(); err != nil {
		return err
	}

	ybrowser.Browser.BlockRequest(ybrowser.Browser.Ctx, "enqueue_creator_bulk_delete")
	reqCh := ybrowser.ListenForRequests(ybrowser.Browser.Ctx, "enqueue_creator_bulk_delete")

	err = ybrowser.Browser.RunWithTimeout(20*time.Second, drivenbrowser.NewLogAction(chromedp.WaitVisible("page-size", chromedp.ByID), "wait for video"))
	if err != nil {
		return err
	}

	err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.WaitNotPresent(".loading-text.ytcp-bulk-actions", chromedp.ByQuery), "wait for previous delete"))
	if err != nil {
		return err
	}

	var selected bool
	for !selected {
		err = ybrowser.Browser.Run(
			drivenbrowser.NewLogAction(chromedp.Click("ytcp-video-row #checkbox-container", chromedp.ByQuery), "click on first video"),
			chromedp.EvaluateAsDevTools("document.querySelector('.row-selected') != null", &selected),
		)
		if err != nil {
			return err
		}
	}

	err = ybrowser.Browser.Run(
		drivenbrowser.NewLogAction(chromedp.Click("additional-action-options", chromedp.ByID), "click on additional options"),
		drivenbrowser.NewLogAction(chromedp.Click("paper-item[test-id=DELETE]", chromedp.ByQuery), "click on delete"),
		drivenbrowser.NewLogAction(chromedp.Click("confirm-checkbox", chromedp.ByID), "click to confirm delete checkbox"),
		drivenbrowser.NewLogAction(chromedp.Click("confirm-button", chromedp.ByID), "click to confirm button"))
	if err != nil {
		return err
	}

	reqData := <-reqCh

	var data DeleteVideosReq
	err = json.Unmarshal([]byte(reqData.Request.PostData), &data)
	if err != nil {
		return
	}

	data.Videos.VideoIds = ids

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
	jsonErr := gjson.Get(body, "error")
	if jsonErr.Exists() {
		return errors.New(jsonErr.Get("message").String())
	}

	return nil
}
