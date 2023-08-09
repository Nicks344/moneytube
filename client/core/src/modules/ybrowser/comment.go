package ybrowser

import (
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"

	"time"

	"github.com/chromedp/chromedp"
)

func (ybrowser *YoutubeBrowser) Comment(link, comment string) error {
	location, err := ybrowser.GetLocation()
	if err != nil {
		return err
	}

	ids := parseLinks.GetVideoIDsByLinks([]string{link, location.String()})

	if len(ids) != 2 || ids[0] != ids[1] {
		err = ybrowser.GoTo(link)
		if err != nil {
			return err
		}
	}

	return ybrowser.Browser.Run(
		drivenbrowser.NewLogAction(chromedp.WaitVisible("comments", chromedp.ByID), "wait for comments"),
		drivenbrowser.NewLogAction(chromedp.EvaluateAsDevTools("document.getElementById('comments').scrollIntoView()", &[]byte{}), "scroll to comments"),
		drivenbrowser.NewLogAction(chromedp.WaitVisible("ytd-comment-simplebox-renderer", chromedp.ByQuery), "wait for comment input"),
		drivenbrowser.NewLogAction(chromedp.Click("ytd-comment-simplebox-renderer", chromedp.ByQuery), "click to comment input"),
		chromedp.WaitVisible(".ytd-commentbox #contenteditable-root", chromedp.ByQuery),
		drivenbrowser.NewLogAction(chromedp.SendKeys(".ytd-commentbox #contenteditable-root", comment, chromedp.ByQuery), "input comment"),
		chromedp.WaitEnabled("submit-button", chromedp.ByID),
		drivenbrowser.NewLogAction(chromedp.Click("submit-button", chromedp.ByID), "submit comment"),
		drivenbrowser.NewLogAction(chromedp.WaitNotPresent("tp-yt-paper-spinner-lite.ytd-commentbox[active]", chromedp.ByQuery), "wait confirm"),
		drivenbrowser.NewLogAction(chromedp.Sleep(time.Second*time.Duration(5)), "sleep"),
	)
}
