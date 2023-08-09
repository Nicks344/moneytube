package ybrowser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/tidwall/gjson"

	"github.com/chromedp/chromedp"
)

func Start(ctx context.Context, acc *moneytubemodel.Account) (ytBrowser *YoutubeBrowser, err error) {
	return start(ctx, config.GetShowBrowser(), acc)
}

func StartVisible(ctx context.Context, acc *moneytubemodel.Account) (ytBrowser *YoutubeBrowser, err error) {
	return start(ctx, true, acc)
}

func StartHidden(ctx context.Context, acc *moneytubemodel.Account) (ytBrowser *YoutubeBrowser, err error) {
	return start(ctx, false, acc)
}

var hostRgx = regexp.MustCompile(`(?m).*(\.google\.|youtube).*`)

func isYoutubeHost(host string) bool {
	return hostRgx.MatchString(host)
}

func start(ctx context.Context, visible bool, acc *moneytubemodel.Account) (ytBrowser *YoutubeBrowser, err error) {
	ytBrowser = &YoutubeBrowser{
		Account: acc,
	}
	proxy, err := acc.GetFormattedProxy()
	if err != nil {
		return
	}

	onErr := func(err error, host string) {
		if !isYoutubeHost(host) {
			return
		}

		// TODO: протестировать кейс с нестабильными прокси
		logger.WarningF("[%s]: %v", host, err)
		//ytBrowser.Browser.StopWithError(errors.New("ошибка соединения"))
	}

	onReq := func(req *http.Request) (*http.Request, *http.Response) {
		if cookie, err := req.Cookie("PREF"); err == nil {
			query, err := url.ParseQuery(cookie.Value)
			if err != nil {
				return req, nil
			}

			for _, key := range []string{"hl", "al"} {
				if query.Has(key) {
					query.Set(key, "ru")
				}
			}

			req.Header["Cookie"] = []string{strings.ReplaceAll(req.Header["Cookie"][0], cookie.Value, query.Encode())}
		}

		return req, nil
	}

	ytBrowser.Browser, err = drivenbrowser.Start(ctx, visible, acc.Login, proxy, acc.UserAgent, onErr, onReq)

	if err != nil {
		return
	}

	return
}

type YoutubeBrowser struct {
	Browser *drivenbrowser.DrivenBrowser
	Account *moneytubemodel.Account
}

func (ybrowser *YoutubeBrowser) Stop() {
	if ybrowser != nil && ybrowser.Browser != nil {
		ybrowser.Browser.Stop()
	}
}

func (ybrowser *YoutubeBrowser) GoTo(url string) (err error) {
	err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Navigate(url), "go to "+url))
	if err != nil {
		return
	}

	isLogined, err := ybrowser.isLogined()
	if err != nil {
		return
	}

	//ybrowser.Browser.FullScreenshot("1.jpg")

	if !isLogined {
		logger.Warning("account is not logined, start login")
		err = ybrowser.Auth()
		if err != nil {
			return
		}
		err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Navigate(url), "go to "+url))
	}

	return
}

func (ybrowser *YoutubeBrowser) GetLocation() (location *url.URL, err error) {
	var loc string
	err = ybrowser.Browser.Run(chromedp.Location(&loc))
	if err != nil {
		return
	}
	location, err = url.ParseRequestURI(loc)
	return
}

func (ybrowser *YoutubeBrowser) GetPlaylistsCount() (int, error) {
	ctx, cancel := context.WithTimeout(ybrowser.Browser.Ctx, time.Minute*5)
	defer cancel()
	playlistBody := ybrowser.WaitForRequestBody(ctx, "list_creator_playlists")

	err := ybrowser.GoTo("https://www.youtube.com/view_all_playlists?nv=1")
	if err != nil {
		return 0, err
	}

	var playlistsAnswer string
	var replaced bool
loop:
	for {
		select {
		case playlistsAnswer = <-playlistBody:
			break loop

		default:
			if utils.IsContextCancelled(ctx) {
				return 0, errors.New("таймаут ожидания плейлистов")
			}

			if replaced {
				continue
			}

			var errText string
			if err := ybrowser.Browser.Run(chromedp.InnerHTML("error-message", &errText, chromedp.ByID)); err != nil {
				return 0, err
			}

			if errText != "" {
				err := ybrowser.GoTo(fmt.Sprintf("https://studio.youtube.com/channel/%s/content/playlists", ybrowser.Account.ChannelID))
				if err != nil {
					return 0, err
				}
			}
		}
	}

	if playlistsAnswer == "" {
		return 0, errors.New("не могу получить данные о плейлистах")
	}

	size := gjson.Get(playlistsAnswer, "playlistsTotalSize.size").String()
	return strconv.Atoi(size)
}
