// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"net/url"

	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Register(func(app *goapp.App) error {

		app.Set("gonode.router", func(app *goapp.App) interface{} {
			return NewRouter(app.Get("goji.mux").(*web.Mux))
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {

		mux := app.Get("goji.mux").(*web.Mux)
		mux.Use(RequestContextMiddleware)

		router := app.Get("gonode.router").(*Router)
		pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)

		pongo.Globals["path"] = func(args ...*pongo2.Value) *pongo2.Value {
			if len(args) == 0 {
				panic("url: missing arguments, at least one is required (name, params))")
			}

			name := args[0].String()
			params := url.Values{}

			if len(args) > 1 && args[1] != nil && args[1].Interface() != nil {
				params = args[1].Interface().(url.Values)
			}

			path, err := router.GeneratePath(name, params)

			if err != nil {
				panic(err)
			}

			return pongo2.AsSafeValue(path)
		}

		pongo.Globals["url"] = func(args ...*pongo2.Value) *pongo2.Value {
			if len(args) == 0 {
				panic("url: missing arguments, at least one is required (name string, params url.Values, request_context *RequestContext))")
			}

			name := args[0].String()
			params := url.Values{}
			requestContext := &RequestContext{}

			if len(args) > 1 && args[1] != nil && args[1].Interface() != nil {
				params = args[1].Interface().(url.Values)
			}

			if len(args) > 2 && args[2] != nil && args[2].Interface() != nil {
				requestContext = args[2].Interface().(*RequestContext)
			}

			path, err := router.GenerateUrl(name, params, requestContext)

			if err != nil {
				panic(err)
			}

			return pongo2.AsSafeValue(path)
		}

		pongo.Globals["net"] = func(args ...*pongo2.Value) *pongo2.Value {
			if len(args) == 0 {
				panic("url: missing arguments, at least one is required (name, params))")
			}

			name := args[0].String()
			params := url.Values{}

			if len(args) > 1 && args[1] != nil && args[1].Interface() != nil {
				params = args[1].Interface().(url.Values)
			}

			path, err := router.GenerateNet(name, params)

			if err != nil {
				panic(err)
			}

			return pongo2.AsSafeValue(path)
		}

		pongo.Globals["url_values"] = func(values ...*pongo2.Value) *pongo2.Value {
			var key string

			m := url.Values{}

			for i, v := range values {
				if i == 0 || i%2 == 0 {
					key = v.String()
					m[key] = nil
				} else {
					m.Add(key, v.String())
				}
			}

			return pongo2.AsValue(m)
		}

		return nil
	})
}
