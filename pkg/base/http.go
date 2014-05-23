package base

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/Unknwon/com"
)

var httpClient *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{Transport: tr}
}

func HttpGetJSON(url string, values url.Values, data interface{}) error {
	if values != nil {
		url = url + "&" + values.Encode()
	}
	return com.HttpGetJSON(httpClient, url, data)
}
