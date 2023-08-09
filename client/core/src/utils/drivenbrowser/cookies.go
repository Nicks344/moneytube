package drivenbrowser

import (
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func (dbrowser *DrivenBrowser) GetCookies() ([]*network.Cookie, error) {
	c := chromedp.FromContext(dbrowser.Ctx)
	return network.GetAllCookies().Do(cdp.WithExecutor(dbrowser.Ctx, c.Target))
}

func (dbrowser *DrivenBrowser) SetCookies(cookies []network.Cookie) error {
	c := chromedp.FromContext(dbrowser.Ctx)
	ctx := cdp.WithExecutor(dbrowser.Ctx, c.Target)
	for _, cookie := range cookies {
		expr := cdp.TimeSinceEpoch(time.Unix(int64(cookie.Expires), 0))
		err := network.SetCookie(cookie.Name, cookie.Value).
			WithExpires(&expr).
			WithDomain(cookie.Domain).
			WithPath(cookie.Path).
			WithHTTPOnly(cookie.HTTPOnly).
			WithSecure(cookie.Secure).
			WithPriority(cookie.Priority).
			WithSameParty(cookie.SameParty).
			WithSameSite(cookie.SameSite).
			WithPath(cookie.Path).
			WithSourcePort(cookie.SourcePort).
			WithSourceScheme(cookie.SourceScheme).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("could not set cookie %s: %s", cookie.Name, err.Error())
		}

	}

	return nil
}
