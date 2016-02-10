// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/rande/goapp"
	"github.com/rande/gonode/commands/server"
	"github.com/rande/gonode/core/bindata"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/api"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/guard"
	"github.com/rande/gonode/modules/prism"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/setup"
	"github.com/rande/gonode/modules/user"
	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"testing"
)

func GetPager(app *goapp.App, res *Response) *api.ApiPager {
	p := &api.ApiPager{}

	serializer := app.Get("gonode.node.serializer").(*base.Serializer)
	serializer.Deserialize(res.Body, p)

	// the Element is a [string]interface so we need to convert it back to []byte
	// and then unmarshal again with the correct structure
	for k, v := range p.Elements {
		raw, _ := json.Marshal(v)

		n := base.NewNode()
		json.Unmarshal(raw, n)

		p.Elements[k] = n
	}

	return p
}

func GetNode(app *goapp.App, res *Response) *base.Node {
	n := base.NewNode()

	serializer := app.Get("gonode.node.serializer").(*base.Serializer)
	serializer.Deserialize(res.Body, n)

	return n
}

func GetLifecycle(file string) *goapp.Lifecycle {

	l := goapp.NewLifecycle()

	conf := config.NewServerConfig()
	conf.Test = true

	config.LoadConfigurationFromFile(file, conf)

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
			return conf
		})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("logger", func(app *goapp.App) interface{} {
			logger := log.New()
			logger.Out = os.Stdout
			logger.Level = log.WarnLevel

			return logger
		})

		app.Set("goji.mux", func(app *goapp.App) interface{} {
			mux := web.New()

			mux.Use(middleware.Logger)
			mux.Use(middleware.Recoverer)

			return mux
		})

		return nil
	})

	server.ConfigureServer(l, conf)
	security.ConfigureServer(l, conf)
	search.ConfigureServer(l, conf)
	api.ConfigureServer(l, conf)
	setup.ConfigureServer(l, conf)
	node_guard.ConfigureServer(l, conf)
	bindata.ConfigureServer(l, conf)
	prism.ConfigureServer(l, conf)
	router.ConfigureServer(l, conf)
	base.ConfigureServer(l, conf)

	return l
}

type Response struct {
	*http.Response
	RawBody  []byte
	bodyRead bool
}

func (r Response) GetBody() []byte {
	var err error

	if !r.bodyRead {
		r.RawBody, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		r.bodyRead = true
	}

	return r.RawBody
}

func GetAuthHeader(t *testing.T, ts *httptest.Server) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", GetAuthToken(t, ts)),
	}
}

func GetAuthToken(t *testing.T, ts *httptest.Server) string {
	res, _ := RunRequest("POST", fmt.Sprintf("%s/login", ts.URL), url.Values{
		"username": {"test-admin"},
		"password": {"admin"},
	})

	assert.Equal(t, 200, res.StatusCode)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	v := &struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{}

	json.Unmarshal(b.Bytes(), v)

	return v.Token
}

func RunRequest(method string, path string, options ...interface{}) (*Response, error) {
	var body interface{}
	var headers map[string]string

	if len(options) > 0 {
		body = options[0]
	}

	if len(options) > 1 {
		headers = options[1].(map[string]string)
	}

	client := &http.Client{}
	var req *http.Request
	var err error

	switch v := body.(type) {
	case nil:
		req, err = http.NewRequest(method, path, nil)
	case *strings.Reader:
		req, err = http.NewRequest(method, path, v)
	case io.Reader:
		req, err = http.NewRequest(method, path, v)

	case url.Values:
		req, err = http.NewRequest(method, path, strings.NewReader(v.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	default:
		panic(fmt.Sprintf("please add a new test case for %T", body))
	}

	if headers != nil {
		for name, value := range headers {
			req.Header.Set(name, value)
		}
	}

	helper.PanicOnError(err)

	resp, err := client.Do(req)

	return &Response{Response: resp}, err
}

func RunHttpTest(t *testing.T, f func(t *testing.T, ts *httptest.Server, app *goapp.App)) {

	l := GetLifecycle("../config_test.toml")

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		var err error
		var res *Response

		mux := app.Get("goji.mux").(*web.Mux)

		ts := httptest.NewServer(mux)

		defer func() {
			state.Out <- goapp.Control_Stop

			ts.Close()

			if r := recover(); r != nil {
				assert.Equal(t, false, true, fmt.Sprintf("RunHttpTest: Panic recovered, message=%s\n\n%s", r, string(debug.Stack()[:])))
			}
		}()

		res, err = RunRequest("POST", ts.URL+"/setup/uninstall", nil)
		helper.PanicIf(res.StatusCode != http.StatusOK, fmt.Sprintf("Expected code 200, get %d\n%s", res.StatusCode, string(res.GetBody()[:])))
		helper.PanicOnError(err)

		res, err = RunRequest("POST", ts.URL+"/setup/install", nil)
		helper.PanicIf(res.StatusCode != http.StatusOK, fmt.Sprintf("Expected code 200, get %d\n%s", res.StatusCode, string(res.GetBody()[:])))
		helper.PanicOnError(err)

		// create a valid user
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		u := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.user")
		u.Name = "User ZZ"
		data := u.Data.(*user.User)
		data.Email = "test-admin@example.org"
		data.Enabled = true
		data.NewPassword = "admin"
		data.Username = "test-admin"
		data.Roles = []string{"ADMIN"}

		meta := u.Meta.(*user.UserMeta)
		meta.PasswordCost = 1 // save test time

		_, err = manager.Save(u, false)
		helper.PanicOnError(err)

		f(t, ts, app)

		return nil
	})

	l.Go(goapp.NewApp())
}
