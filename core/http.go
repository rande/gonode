package core

import (
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
	Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

func GetHttpClient() HttpClient {
	return &http.Client{}
}
