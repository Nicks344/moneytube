package drivenbrowser

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-rod/stealth"
	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/utils/proxychain"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

var netErrors = map[string]struct{}{
	"net::ERR_CONNECTION_CLOSED":        {},
	"net::ERR_TUNNEL_CONNECTION_FAILED": {},
	"net::ERR_CONNECTION_TIMED_OUT":     {},
	"net::ERR_CERT_AUTHORITY_INVALID":   {},
	"net::ERR_EMPTY_RESPONSE":           {},
	"net::ERR_TIMED_OUT":                {},
	"net::ERR_FAILED":                   {},
	"net::ERR_PROXY_CONNECTION_FAILED":  {},
	"net::ERR_SOCKS_CONNECTION_FAILED":  {},
}

type ReqHandler func(req *http.Request) (*http.Request, *http.Response)
type ErrHandler func(err error, host string)

func Start(ctx context.Context, visible bool, session, proxy, useragent string, onErr ErrHandler, onReq ReqHandler) (browser *DrivenBrowser, err error) {
	browser = &DrivenBrowser{
		CmdTimeout:   360 * time.Second,
		proxychainID: uuid.New().String(),
	}

	visible = true

	localProxy := proxychain.New(browser.proxychainID, proxychain.Config{proxy, onErr, onReq})

	extPath, err := filepath.Abs(filepath.Join(paths.Bin, "extensions", "webrtc"))
	if err != nil {
		return
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// TODO: Add UAs
		chromedp.UserAgent(getUserAgent(useragent)),
		chromedp.ProxyServer(localProxy),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("force-webrtc-ip-handling-policy", "disable_non_proxied_udp"),
		chromedp.Flag("enforce-webrtc-ip-permission-check", "1"),
		chromedp.Flag("disable-extensions-except", extPath),
		chromedp.Flag("load-extension", extPath),
		//chromedp.Headless,
		//chromedp.NoSandbox,
		//chromedp.Flag("ignore-certificate-errors", "1"),
		//chromedp.Flag("enable-parallel-downloading", "1"),
		//chromedp.Flag("enable-quic", "1"),
		//chromedp.Flag("enable-gpu-rasterization", "1"),
		//chromedp.ExecPath(paths.ChromeExe),
	)

	if visible {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	if session != "" {
		sessionDir := filepath.Join(paths.Sessions, session)
		sessionDir, err = filepath.Abs(sessionDir)
		if err != nil {
			return
		}
		opts = append(opts, chromedp.UserDataDir(sessionDir))
	}

	ctx = context.WithValue(ctx, "session", session)
	allocatorContext, cancel := chromedp.NewExecAllocator(ctx, opts...)
	browser.cancel = cancel

	browser.Ctx, _ = chromedp.NewContext(allocatorContext) //chromedp.WithDebugf(log.Printf),

	//chromedp.WithErrorf(log.Printf),
	//chromedp.WithLogf(log.Printf),

	netEnable := network.Enable()
	netEnable.MaxResourceBufferSize = 1 * 1024 * 1024 * 1024
	netEnable.MaxTotalBufferSize = 1 * 1024 * 1024 * 1024
	err = chromedp.Run(browser.Ctx, netEnable, chromedp.ActionFunc(func(ctx context.Context) (err error) {
		_, err = page.AddScriptToEvaluateOnNewDocument(stealth.JS).Do(ctx)
		if err != nil {
			return
		}
		_, err = page.AddScriptToEvaluateOnNewDocument("delete navigator.__proto__.webdriver;").Do(ctx)
		return
	}), chromedp.Navigate("about:blank"))

	if onErr != nil {
		willBeSentHosts := map[string]string{}
		chromedp.ListenTarget(browser.Ctx, func(ev interface{}) {
			switch ev := ev.(type) {
			case *network.EventLoadingFailed:
				if _, ok := netErrors[ev.ErrorText]; ok {
					host := willBeSentHosts[ev.RequestID.String()]
					onErr(errors.New(ev.ErrorText), host)
				}

			case *network.EventRequestWillBeSentExtraInfo:
				if host, ok := ev.Headers["Host"]; ok {
					willBeSentHosts[ev.RequestID.String()] = host.(string)
				}

			case *network.EventLoadingFinished:
				delete(willBeSentHosts, ev.RequestID.String())

			}
		})
	}

	return
}

func getUserAgent(ua string) string {
	if ua == "" {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
	}

	return ua
}

type browserInfo struct {
	Browser              string `json:"Browser"`
	ProtocolVersion      string `json:"Protocol-Version"`
	UserAgent            string `json:"User-Agent"`
	V8Version            string `json:"V8-Version"`
	WebKitVersion        string `json:"WebKit-Version"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}
