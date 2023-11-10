// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"fmt"
	tpl "html/template"
	"net/url"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/template"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Register(func(app *goapp.App) error {

		app.Set("gonode.router", func(app *goapp.App) interface{} {
			return NewRouter(app.Get("goji.mux").(*web.Mux))
		})

		return nil
	})

	l.Config(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		mux.Use(RequestContextMiddleware)

		router := app.Get("gonode.router").(*Router)
		loader := app.Get("gonode.template").(*template.TemplateLoader)

		loader.FuncMap["path"] = func(name string, options ...interface{}) tpl.HTML {
			params := url.Values{}

			if len(options) > 0 && options[0] != nil {
				params = options[0].(url.Values)
			}

			if path, err := router.GeneratePath(name, params); err != nil {
				panic(err)
			} else {
				return tpl.HTML(path)
			}
		}

		loader.FuncMap["url"] = func(name string, options ...interface{}) tpl.HTML {

			params := url.Values{}
			requestContext := &RequestContext{}

			if len(options) > 0 && options[0] != nil {
				params = options[0].(url.Values)
			}

			if len(options) > 1 && options[1] != nil {
				requestContext = options[1].(*RequestContext)
			}

			if path, err := router.GenerateUrl(name, params, requestContext); err != nil {
				panic(err)
			} else {
				return tpl.HTML(path)
			}
		}

		loader.FuncMap["net"] = func(name string, options ...interface{}) tpl.HTML {
			params := url.Values{}

			if len(options) > 0 && options[0] != nil {
				params = options[0].(url.Values)
			}

			path, err := router.GenerateNet(name, params)

			if err != nil {
				panic(err)
			}

			return tpl.HTML(path)
		}

		loader.FuncMap["url_values"] = func(values ...interface{}) url.Values {
			var key string

			m := url.Values{}

			for i, v := range values {

				if i == 0 || i%2 == 0 {
					key = fmt.Sprintf("%v", v)
					m[key] = nil
				} else {
					m.Add(key, fmt.Sprintf("%v", v))
				}
			}

			return m
		}

		return nil
	})
}
