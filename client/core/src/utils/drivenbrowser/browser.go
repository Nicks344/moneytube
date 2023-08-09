package drivenbrowser

import (
	"context"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils/proxychain"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type DrivenBrowser struct {
	Ctx        context.Context
	CmdTimeout time.Duration

	cancel       context.CancelFunc
	proxychainID string
	err          error

	stopLock sync.Mutex
	stopped  bool
}

func (dbrowser *DrivenBrowser) FullScreenshot(filename string) error {
	var buf []byte
	err := dbrowser.Run(chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			buf, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, buf, 0644)
	return err
}

func (dbrowser *DrivenBrowser) Run(actions ...chromedp.Action) error {
	ctx, _ := context.WithTimeout(dbrowser.Ctx, time.Minute*20)
	err := chromedp.Run(ctx, actions...)
	if dbrowser.err != nil {
		return dbrowser.err
	}
	return err
}

func (dbrowser *DrivenBrowser) RunWithTimeout(timeout time.Duration, actions ...chromedp.Action) error {
	ctx, _ := context.WithTimeout(dbrowser.Ctx, timeout)
	return chromedp.Run(ctx, actions...)
}

func (dbrowser *DrivenBrowser) StopWithError(err error) {
	dbrowser.err = err
	dbrowser.Stop()
}

func (dbrowser *DrivenBrowser) BlockRequest(ctx context.Context, mask string) {
	dbrowser.Run(network.SetBlockedURLS([]string{"*" + mask + "*"}))
}

func (dbrowser *DrivenBrowser) ListenResponse(ctx context.Context, f func(req *network.EventResponseReceived)) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		go func() {
			switch ev := ev.(type) {
			case *network.EventResponseReceived:
				f(ev)
			}
		}()
	})
}

func (dbrowser *DrivenBrowser) ListenRequest(ctx context.Context, f func(req *network.EventRequestWillBeSent)) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		go func() {
			switch ev := ev.(type) {
			case *network.EventRequestWillBeSent:
				f(ev)
			}
		}()
	})
}

func (dbrowser *DrivenBrowser) Minimize() error {
	return dbrowser.Run(chromedp.ActionFunc(func(ctx context.Context) error {
		windowsID, _, err := browser.GetWindowForTarget().Do(ctx)
		if err != nil {
			return err
		}

		return browser.SetWindowBounds(windowsID, &browser.Bounds{
			WindowState: browser.WindowStateMinimized,
		}).Do(ctx)
	}))
}

func (dbrowser *DrivenBrowser) Stop() error {
	dbrowser.stopLock.Lock()
	defer dbrowser.stopLock.Unlock()

	if dbrowser.stopped {
		return nil
	}

	dbrowser.stopped = true
	proxychain.Close(dbrowser.proxychainID)
	return chromedp.Cancel(dbrowser.Ctx)
}
