package ybrowser

import (
	"context"
	"errors"
	"fmt"

	"github.com/meandrewdev/logger"

	"strings"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	webStartPage = "https://accounts.google.com/signin/v2/identifier?hl=ru&service=youtube&continue=https%3A%2F%2Fwww.youtube.com%2Fsignin%3Ffeature%3Dsign_in_button%26hl%3Den%26app%3Ddesktop%26next%3D%252F%26action_handle_signin%3Dtrue&passive=true&uilel=3&flowName=GlifWebSignIn&flowEntry=ServiceLogin"
)

func (ybrowser *YoutubeBrowser) GetIDFromLocation() (id string, err error) {
	channelHref := ""
	err = ybrowser.Browser.Run(chromedp.Location(&channelHref))
	if err != nil {
		return
	}

	split1 := strings.Split(channelHref, "/")
	id = strings.Split(split1[len(split1)-1], "?")[0]
	return
}

func (ybrowser *YoutubeBrowser) OpenStartPage() error {
	return ybrowser.Browser.Run(chromedp.Navigate(webStartPage))
}

func (ybrowser *YoutubeBrowser) Auth() (err error) {
	err = ybrowser.OpenStartPage()
	if err != nil {
		return
	}

	isLogined, err := ybrowser.isLogined()
	if err != nil {
		return
	}

	if !isLogined {
		err = ybrowser.doLogin()
	}

	return
}

func (ybrowser *YoutubeBrowser) GetID() (id string, err error) {
	err = ybrowser.GoTo("https://studio.youtube.com")
	if err != nil {
		return
	}

	id, err = ybrowser.GetIDFromLocation()
	if err != nil {
		return
	}

	if id == "create_channel" {
		err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Click(".create-channel-submit", chromedp.ByQuery), "click to create channel"))
		if err != nil {
			return
		}
	}

	for id == "create_channel" || id == "" {
		id, err = ybrowser.GetIDFromLocation()
		if err != nil {
			return
		}
	}

	return id, nil
}

func (ybrowser *YoutubeBrowser) doLogin() error {
	var val bool
	err := ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.body.innerHTML.includes('Выберите аккаунт')", &val))
	if err != nil {
		return err
	}

	if val {
		time.Sleep(2 * time.Second)
		err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.Click("div[data-item-index='0']", chromedp.ByQuery), "click to existing account"))
		if err != nil {
			return err
		}
	} else {
		err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.WaitVisible("input[type=email]", chromedp.ByQuery), "wait for login input"),
			chromedp.Sleep(2*time.Second),
			drivenbrowser.NewLogAction(chromedp.SendKeys("input[type=email]", ybrowser.Account.Login, chromedp.ByQuery), "input login"),
			drivenbrowser.NewLogAction(chromedp.Click("#identifierNext", chromedp.ByQuery), "click to next"))

		if err != nil {
			return err
		}

		for {
			var exist bool
			err := ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.querySelector('input[type=password]') != null", &exist))
			if err != nil {
				return err
			}
			if exist {
				break
			}

			if ybrowser.isInvalid("input[type=email]") {
				// TODO: remove hardcode
				return errors.New("invalid login")
			}
			/*
				err = ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.querySelector('#captchaimg').width > 0", &exist))
				if err != nil {
					return err
				}
			*/
			if exist {
				// TODO: remove hardcode
				return errors.New("captcha required")
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	err = ybrowser.Browser.Run(drivenbrowser.NewLogAction(chromedp.SendKeys("input[type=password]", ybrowser.Account.Password), "input password"),
		chromedp.WaitVisible("#passwordNext"),
		chromedp.Sleep(time.Second*1),
		drivenbrowser.NewLogAction(chromedp.Click("#passwordNext"), "click to next"))

	if err != nil {
		return err
	}

	showed := false

	for {
		login, err := ybrowser.isLogined()
		if err != nil {
			return err
		}
		if login {
			logger.Notice("auth completed")
			break
		}
		if ybrowser.isInvalid("input[name=password]") {
			// TODO: remove hardcode
			return errors.New("invalid password")
		}

		if !showed && ybrowser.isSecure() {
			// TODO: remove hardcode
			return errors.New("identity confirmation")
		}

		if err = ybrowser.skipRestore(); err != nil {
			return err
		}

		time.Sleep(100 * time.Millisecond)
	}

	return ybrowser.skipPersonalization()
}

func (ybrowser *YoutubeBrowser) skipPersonalization() (err error) {
	err = ybrowser.Browser.RunWithTimeout(5*time.Second, chromedp.Click("//span[text()='Не сейчас']", chromedp.BySearch))
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}

func (ybrowser *YoutubeBrowser) skipRestore() (err error) {
	var submitParams bool
	err = ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.querySelector('title').innerText && document.querySelector('title').innerText.includes('Параметры восстановления аккаунта')", &submitParams))
	if err != nil {
		return err
	}
	if submitParams {
		var buttons []*cdp.Node
		err = ybrowser.Browser.Run(chromedp.Nodes("div[role=button]", &buttons, chromedp.ByQueryAll))
		if err != nil {
			return err
		}
		err = ybrowser.Browser.Run(chromedp.MouseClickNode(buttons[1]))
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return
}

func (ybrowser *YoutubeBrowser) isSecure() bool {
	var val bool
	ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.getElementById('headingText').innerHTML.includes('Подтвердите, что это именно вы')", &val))
	if val {
		return true
	}
	ybrowser.Browser.Run(chromedp.EvaluateAsDevTools("document.body.innerHTML.includes('Сменить пароль')", &val))
	return val
}

func (ybrowser *YoutubeBrowser) isLogined() (bool, error) {
	location, err := ybrowser.GetLocation()
	if err != nil {
		return false, err
	}

	if strings.Contains(location.Host, "accounts.google.com") {
		return false, nil
	}

	cookies, err := ybrowser.Browser.GetCookies()
	if err != nil {
		return false, err
	}

	for _, cookie := range cookies {
		if strings.Contains(cookie.Domain, "youtube") && cookie.Name == "APISID" {
			expiresAt := time.Unix(int64(cookie.Expires), 0)
			expired := time.Now().After(expiresAt)
			fmt.Printf("found APISID, expires at %s, is expired: %v", expiresAt.String(), expired)
			return !expired, nil
		}
	}

	return false, nil
}

func (ybrowser *YoutubeBrowser) isInvalid(selector string) bool {
	var isInvalid string
	ybrowser.Browser.Run(chromedp.EvaluateAsDevTools(fmt.Sprintf("document.querySelector('%s').attributes['aria-invalid'].value", selector), &isInvalid))
	return isInvalid == "true"
}
