package proxychain

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/meandrewdev/logger"

	"golang.org/x/net/proxy"
	"gopkg.in/elazarl/goproxy.v1"
)

var proxyChains = map[string]*proxyChain{}
var lc = sync.Mutex{}

type Config struct {
	Proxy string
	OnErr func(error, string)
	OnReq func(req *http.Request) (*http.Request, *http.Response)
}

func New(id string, conf Config) string {
	pc := &proxyChain{
		Config: conf,
	}
	lc.Lock()
	proxyChains[id] = pc
	lc.Unlock()
	addr := pc.create()
	return "http://" + addr
}

func Close(id string) {
	if chain, ok := proxyChains[id]; ok {
		chain.destroy()
		delete(proxyChains, id)
	}
}

type proxyChain struct {
	Config

	addr   string
	server *http.Server
}

func (pc *proxyChain) create() string {
	gproxy := goproxy.NewProxyHttpServer()
	//gproxy.Logger = log.New(io.Discard, "", 0)

	if pc.Proxy != "" {

		logger.Warning("set proxy: " + pc.Proxy)
		proxyURL, err := url.Parse(pc.Proxy)
		if err != nil {
			logger.Warning("set proxy error: " + err.Error())
		}
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			logger.Warning("set proxy error: " + err.Error())
		}

		gproxy.Tr = &http.Transport{
			Dial:            dialer.Dial,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	pc.addr = fmt.Sprintf("127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)
	pc.server = &http.Server{
		Addr:    pc.addr,
		Handler: gproxy,
	}
	go pc.server.Serve(listener)
	return pc.addr
}

func (pc *proxyChain) destroy() {
	pc.server.Close()
}
