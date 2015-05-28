package gs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"h12.me/html-query"
	"h12.me/socks"
)

type HTTP struct {
	client http.Client
}

func (h HTTP) Proxy(proxy string) HTTP {
	u, err := url.Parse(proxy)
	c(err)
	switch u.Scheme {
	case "socks5":
		h.client.Transport = &http.Transport{Dial: socks.DialSocksProxy(socks.SOCKS5, u.Host)}
	}
	return h
}

func (h HTTP) Get(uri string) WebPage {
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
	if resp.StatusCode != http.StatusOK {
		return WebPage{err: fmt.Errorf("Status Code %d", resp.StatusCode)}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return WebPage{body, err}
}

type WebPage struct {
	body []byte
	err  error
}

func (p WebPage) Load(file string) WebPage {
	f, err := os.Open(file)
	c(err)
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	c(err)
	return WebPage{body: body}
}

func (p WebPage) Save(file string) {
	c(p.err)
	err := p.TrySave(file)
	c(err)
}

func (p WebPage) TrySave(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(p.body)
	return err
}

func (p WebPage) Parse() *query.Node {
	n, err := query.Parse(bytes.NewBuffer(p.body))
	c(err)
	return n
}
