// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"io"
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
)

type MockedHttpClient struct {
	mock.Mock
}

func (c *MockedHttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	args := c.Mock.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *MockedHttpClient) Get(url string) (resp *http.Response, err error) {
	args := c.Mock.Called(url)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *MockedHttpClient) Head(url string) (resp *http.Response, err error) {
	args := c.Mock.Called(url)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *MockedHttpClient) Post(url string, bodyType string, body io.Reader) (resp *http.Response, err error) {
	args := c.Mock.Called(url, bodyType, body)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (c *MockedHttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := c.Mock.Called(url, data)

	return args.Get(0).(*http.Response), args.Error(1)
}
