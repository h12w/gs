package gs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"h12.me/html-query"
	"h12.me/socks"
)

type HTTP struct {
	Proxy  string
	client *http.Client
}

func (h HTTP) Get(uri string) WebPage {
	if h.client == nil {
		h.client = &http.Client{}
		u, err := url.Parse(h.Proxy)
		c(err)
		switch u.Scheme {
		case "socks5":
			h.client.Transport = &http.Transport{Dial: socks.DialSocksProxy(socks.SOCKS5, u.Host)}
		}
	}
	resp, err := h.client.Get(uri)
	for i := 0; i < 3; i++ {
		resp, err = h.client.Get(uri)
		if err == nil {
			break
		}
	}
	if err != nil {
		return WebPage{err: err}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return WebPage{body, err}
}

type WebPage struct {
	body []byte
	err  error
}

func (p WebPage) Save(file string) error {
	if p.err != nil {
		return p.err
	}
	f, err := os.Create(file)
	c(err)
	defer f.Close()
	_, err = f.Write(p.body)
	c(err)
	return nil
}

func (p WebPage) Parse() *query.Node {
	n, err := query.Parse(bytes.NewBuffer(p.body))
	c(err)
	return n
}
