package mock

import (
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/url"
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
