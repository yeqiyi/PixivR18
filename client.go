package pixivR18

import (
	"crypto/tls"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func NewClient() *http.Client{
	cjar,_:=cookiejar.New(&cookiejar.Options{PublicSuffixList:publicsuffix.List})
	client:=&http.Client{Jar: cjar}
	return client
}

func NewClientWithPorxy(proxyUrl string) *http.Client {
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := NewClient()
	client.Transport = tr
	return client
}

