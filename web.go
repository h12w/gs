package gs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	query "h12.io/html-query"
	"h12.io/socks"
)

type HTTP struct {
	client http.Client
	retry  int

	cacheDir string
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

func (h HTTP) Timeout(du time.Duration) HTTP {
	h.client.Timeout = du
	return h
}

func (h HTTP) Retry(n int) HTTP {
	h.retry = n
	return h
}

func (h HTTP) CacheDir(dir string) HTTP {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	h.cacheDir = dir
	return h
}

var (
	ErrNotFound = errors.New("not found")
)

func (h HTTP) cacheFilename(uri string) string {
	if h.cacheDir == "" {
		return ""
	}
	return path.Join(h.cacheDir, trimScheme(uri))
}

func (h HTTP) Get(uri string) WebPage {
	cacheFilename := h.cacheFilename(uri)
	if cacheFilename != "" {
		body, err := ioutil.ReadFile(cacheFilename)
		if err == nil {
			return WebPage{body: body}
		}
	}
	resp, err := h.client.Get(uri)
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		for i := 0; i < h.retry; i++ {
			resp, err = h.client.Get(uri)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return WebPage{err: err}
	}
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return WebPage{err: ErrNotFound}
		}
		return WebPage{err: fmt.Errorf("Status Code %d", resp.StatusCode)}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil && cacheFilename != "" {
		if err := os.MkdirAll(path.Dir(cacheFilename), 0755); err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(cacheFilename, body, 0644); err != nil {
			panic(err)
		}
	}
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
	c(p.TrySave(file))
}

func (p WebPage) TrySave(file string) error {
	if p.err != nil {
		return p.err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(p.body)
	return err
}

func (p WebPage) ParseHTML() *query.Node {
	n, err := query.Parse(bytes.NewBuffer(p.body))
	c(err)
	return n
}

func (p WebPage) ParseJSON(v interface{}) error {
	return json.Unmarshal(p.body, v)
}

func (p WebPage) Body() string {
	return string(p.body)
}

func trimScheme(uri string) string {
	uri = strings.TrimPrefix(uri, "https://")
	uri = strings.TrimPrefix(uri, "http://")
	return uri
}
