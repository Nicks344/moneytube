package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func OpenAndShow(acc moneytubemodel.Account) error {
	ytBrowser, err := ybrowser.StartVisible(context.Background(), &acc)
	if err != nil {
		return err
	}
	err = ytBrowser.OpenStartPage()
	if err != nil {
		return err
	}

	return nil
}

func WebAuth(acc moneytubemodel.Account) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	ytBrowser, err := ybrowser.Start(ctx, &acc)
	if err != nil {
		return "", err
	}
	defer ytBrowser.Stop()

	if err := ytBrowser.Auth(); err != nil {
		return "", err
	}

	return ytBrowser.GetID()
}

func GetInfo(acc moneytubemodel.Account) (moneytubemodel.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	ytBrowser, err := ybrowser.Start(ctx, &acc)
	if err != nil {
		return acc, err
	}
	defer ytBrowser.Stop()

	acc.ChannelID, err = ytBrowser.GetID()
	if err != nil {
		return acc, err
	}

	var welcome bool
	err = ytBrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.querySelector('ytcp-warm-welcome-dialog #welcome-dialog') != null", &welcome))
	if err != nil {
		return acc, err
	}
	if welcome {
		err = ytBrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Click("ytcp-warm-welcome-dialog #welcome-dialog", chromedp.ByQuery), "click to hide welcome dialog"))
		if err != nil {
			return acc, err
		}
	}

	channelName := ""
	err = ytBrowser.Browser.Run(
		drivenbrowser.NewLogAction(chromedp.WaitVisible("entity-name", chromedp.ByID), "wait for channel name"),
		drivenbrowser.NewLogAction(
			chromedp.EvaluateAsDevTools("document.getElementById('entity-name').innerText", &channelName),
			"get channel name"),
	)
	if err != nil {
		return acc, err
	}

	acc.ChannelName = strings.Trim(channelName, " \r\n")

	stats, err := ytBrowser.GetChannelStats()
	if err != nil {
		return acc, err
	}

	acc.SubscribersCount = stats.SubscriberCount
	acc.VideoCount = stats.VideoCount
	acc.ViewsCount = stats.TotalVideoViewsCount

	playlists, err := ytBrowser.GetPlaylistsCount()
	if err != nil {
		return acc, err
	}

	acc.PlaylistsCount = playlists

	return acc, nil
}

func ExportCookies(acc moneytubemodel.Account, file string) error {
	if !strings.HasSuffix(file, ".json") && !strings.HasSuffix(file, ".txt") {
		return errors.New("Расширение файла должно быть .txt или .json")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	logger.Warning("start browser")
	ytBrowser, err := ybrowser.StartVisible(ctx, &acc)
	if err != nil {
		logger.Warning("error start browser: " + err.Error())
		return err
	}
	defer ytBrowser.Stop()

	logger.Warning("get cookies")
	cookies, err := ytBrowser.Browser.GetCookies()
	if err != nil {
		logger.Warning("error get cookies: " + err.Error())
		return err
	}

	if strings.HasSuffix(file, ".txt") {
		cookiesLines := make([]string, len(cookies))
		for i, cookie := range cookies {
			cookiesLines[i] = fmt.Sprintf("%s=%s;domain=%s;path=%s;expires=%.2f;httpOnly=%t;secure=%t;sameSite=%s;priority=%s;sameParty=%t",
				cookie.Name, cookie.Value, cookie.Domain, cookie.Path, cookie.Expires, cookie.HTTPOnly, cookie.Secure, cookie.SameSite, cookie.Priority, cookie.SameParty)
		}

		return os.WriteFile(file, []byte(strings.Join(cookiesLines, "; ")), 0666)
	} else if strings.HasSuffix(file, ".json") {
		data, err := json.MarshalIndent(cookies, "", "\t")
		if err != nil {
			logger.Warning("error MarshalIndent: " + err.Error())
			return err
		}

		return os.WriteFile(file, data, 0777)
	}

	return errors.New("Расширение файла должно быть .txt или .json")
}

func ImportCookies(acc moneytubemodel.Account, file string) error {
	if !strings.HasSuffix(file, ".json") && !strings.HasSuffix(file, ".txt") {
		return errors.New("Расширение файла должно быть .txt или .json")
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var cookies []network.Cookie
	switch {
	case strings.HasSuffix(file, ".json"):
		if err := json.Unmarshal(data, &cookies); err != nil {
			return err
		}

	case strings.HasSuffix(file, ".txt"):
		cookiesLines := strings.Split(string(data), "; ")
		cookies = make([]network.Cookie, len(cookiesLines))
		for i, cookieLine := range cookiesLines {
			cookie := network.Cookie{
				SourceScheme: network.CookieSourceSchemeSecure,
				SourcePort:   443,
			}

			cookieParts := strings.Split(cookieLine, ";")
			for _, cookiePart := range cookieParts {
				cookiePartParts := strings.Split(cookiePart, "=")
				switch cookiePartParts[0] {
				case "domain":
					cookie.Domain = cookiePartParts[1]
				case "path":
					cookie.Path = cookiePartParts[1]
				case "expires":
					cookie.Expires, _ = strconv.ParseFloat(cookiePartParts[1], 64)
				case "httpOnly":
					cookie.HTTPOnly, _ = strconv.ParseBool(cookiePartParts[1])
				case "secure":
					cookie.Secure, _ = strconv.ParseBool(cookiePartParts[1])
				case "sameSite":
					cookie.SameSite = network.CookieSameSite(cookiePartParts[1])
				case "priority":
					cookie.Priority = network.CookiePriority(cookiePartParts[1])
				case "sameParty":
					cookie.SameParty, _ = strconv.ParseBool(cookiePartParts[1])
				default:
					cookie.Name = cookiePartParts[0]
					cookie.Value = strings.Join(cookiePartParts[1:], "=")
				}
			}

			cookies[i] = cookie
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	ytBrowser, err := ybrowser.StartVisible(ctx, &acc)
	if err != nil {
		return err
	}
	defer ytBrowser.Stop()

	return ytBrowser.Browser.SetCookies(cookies)
}
