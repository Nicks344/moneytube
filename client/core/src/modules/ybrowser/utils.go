package ybrowser

import (
	"context"
	"strings"
	"time"

	"github.com/meandrewdev/logger"

	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func (ybrowser *YoutubeBrowser) WaitForRequestDone(ctx context.Context, urlSlice string) chan bool {
	result := make(chan bool)

	browserCtx, cancel := context.WithCancel(ybrowser.Browser.Ctx)
	ybrowser.Browser.ListenResponse(browserCtx, func(req *network.EventResponseReceived) {
		if strings.Contains(req.Response.URL, urlSlice) {
			result <- true
			cancel()
		}
	})

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			result <- false
			return

		case <-browserCtx.Done():
			return
		}
	}()

	return result
}

func (ybrowser *YoutubeBrowser) WaitForRequestBody(ctx context.Context, urlSlice string) chan string {
	result := make(chan string)

	browserCtx, cancel := context.WithCancel(ybrowser.Browser.Ctx)
	ybrowser.Browser.ListenResponse(browserCtx, func(req *network.EventResponseReceived) {
		if strings.Contains(req.Response.URL, urlSlice) {
			go func() {
				result <- GetBody(browserCtx, req)
				cancel()
			}()
		}
	})

	go func() {
		select {
		case <-ctx.Done():
			cancel()
			result <- ""
			return

		case <-browserCtx.Done():
			return
		}
	}()

	return result
}

func (ybrowser *YoutubeBrowser) ListenForRequestBody(ctx context.Context, urlSlice string) chan string {
	result := make(chan string, 10)

	browserCtx, cancel := context.WithCancel(ybrowser.Browser.Ctx)
	ybrowser.Browser.ListenResponse(browserCtx, func(req *network.EventResponseReceived) {
		if strings.Contains(req.Response.URL, urlSlice) {
			go func() {
				result <- GetBody(browserCtx, req)
			}()

		}
	})

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return result
}

func (ybrowser *YoutubeBrowser) ListenForRequests(ctx context.Context, urlSlice string) chan *network.EventRequestWillBeSent {
	result := make(chan *network.EventRequestWillBeSent, 10)

	browserCtx, cancel := context.WithCancel(ybrowser.Browser.Ctx)
	ybrowser.Browser.ListenRequest(browserCtx, func(req *network.EventRequestWillBeSent) {
		if strings.Contains(req.Request.URL, urlSlice) && req.Request.HasPostData {
			result <- req
		}
	})

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return result
}

func GetBody(browserCtx context.Context, req *network.EventResponseReceived) string {
	bodyParams := network.GetResponseBody(req.RequestID)
	c := chromedp.FromContext(browserCtx)
	b, err := bodyParams.Do(cdp.WithExecutor(browserCtx, c.Target))
	if err != nil {
		e, ok := err.(*cdproto.Error)
		if !ok {
			logger.Error(err)
			return ""
		}

		if e.Code == -32000 {
			time.Sleep(500 * time.Millisecond)
			return GetBody(browserCtx, req)
		}

		logger.Error(e)
		return ""

	}
	return string(b)
}
