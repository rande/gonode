// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
	"net/url"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.Config) {
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

		pongo.Globals["path"] = func(name, params *pongo2.Value) *pongo2.Value {
			path, err := router.GeneratePath(name.String(), params.Interface().(url.Values))

			if err != nil {
				panic(err)
			}

			return pongo2.AsSafeValue(path)
		}

		pongo.Globals["url"] = func(name, params, requestContext *pongo2.Value) *pongo2.Value {
			path, err := router.GenerateUrl(name.String(), params.Interface().(url.Values), requestContext.Interface().(*RequestContext))

			if err != nil {
				panic(err)
			}

			return pongo2.AsSafeValue(path)
		}

		pongo.Globals["net"] = func(name, params *pongo2.Value) *pongo2.Value {
			path, err := router.GenerateNet(name.String(), params.Interface().(url.Values))

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
