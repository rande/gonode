// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"encoding/json"
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

func SendWithHttpCode(res http.ResponseWriter, code int, message string) {
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(code)

	status := "KO"
	if code >= 200 && code < 300 {
		status = "OK"
	}

	data, _ := json.Marshal(map[string]string{
		"status":  status,
		"message": message,
	})

	res.Write(data)
}

func SendWithStatus(status string, message string, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")

	if status == "KO" {
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
	}

	data, _ := json.Marshal(map[string]string{
		"status":  status,
		"message": message,
	})

	res.Write(data)
}
