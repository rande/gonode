// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
)

func Test_Cors_Default_Values(t *testing.T) {

	l := goapp.NewLifecycle()
	conf := config.NewConfig()

	config.LoadConfigurationFromString(`[security]
    [security.cors]
    allowed_origins = ["*"]
    allowed_methods = ["GET", "PUT", "POST"]
    allowed_headers = ["Origin", "Accept", "Content-Type", "Authorization"]
`, conf)

	l.Prepare(func(app *goapp.App) error {
		app.Set("goji.mux", func(app *goapp.App) interface{} {
			mux := web.New()

			return mux
		})

		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
			return conf
		})

		return nil
	})

	ConfigureCors(l, conf)

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {

		mux := app.Get("goji.mux").(*web.Mux)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/login", nil)
		req.Header.Add("Origin", "http://foobar.com")

		mux.ServeHTTP(res, req)

		assert.Equal(t, "Origin", res.Header().Get("Vary"))
		assert.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))

		state.Out <- 1

		return nil
	})

	l.Go(goapp.NewApp())
}
