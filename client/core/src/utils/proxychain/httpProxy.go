package proxychain

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func init() {
	proxy.RegisterDialerType("http", newHTTPProxy)
	proxy.RegisterDialerType("https", newHTTPProxy)
}

type httpProxy struct {
	host     string
	haveAuth bool
	username string
	password string
	forward  proxy.Dialer
}

func newHTTPProxy(uri *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	s := new(httpProxy)
	s.host = uri.Host
	s.forward = forward
	if uri.User != nil {
		s.haveAuth = true
		s.username = uri.User.Username()
		s.password, _ = uri.User.Password()
	}

	return s, nil
}

func (s *httpProxy) Dial(network, addr string) (net.Conn, error) {
	// Dial and create the https client connection.
	c, err := s.forward.Dial("tcp", s.host)
	if err != nil {
		return nil, err
	}

	// HACK. http.ReadRequest also does this.
	reqURL, err := url.Parse("http://" + addr)
	if err != nil {
		c.Close()
		return nil, err
	}
	reqURL.Scheme = ""

	req, err := http.NewRequest("CONNECT", reqURL.String(), nil)
	if err != nil {
		c.Close()
		return nil, err
	}
	req.Close = false
	if s.haveAuth {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(s.username+":"+s.password))
		req.Header.Add("Proxy-Authorization", basicAuth)
	}

	err = req.Write(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		c.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		c.Close()
		err = fmt.Errorf("Connect server using proxy error, StatusCode [%d]", resp.StatusCode)
		return nil, err
	}

	return c, nil
}
