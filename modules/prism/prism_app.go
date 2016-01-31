// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prism

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/modules/base"
	"github.com/zenazn/goji/web"
	"net/http"
)

func RenderPrism(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)
	handlers := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		reference, _ := base.GetReferenceFromString(c.URLParams["uuid"])

		format := "html"
		if _, ok := c.URLParams["format"]; ok {
			format = c.URLParams["format"]
		}

		node := manager.Find(reference)

		res.Header().Set("Content-Type", "text/html; charset=UTF-8")

		request := &base.ViewRequest{
			Context:     c,
			HttpRequest: req,
			Format:      format,
		}

		response := base.NewViewResponse(res)

		if _, ok := c.Env["request_context"]; ok {
			response.Add("request_context", c.Env["request_context"])
		} else {
			response.Add("request_context", nil)
		}

		response.Add("request", req)

		if node != nil {
			response.Add("node", node)

			err := handlers.Get(node).Execute(node, request, response)

			helper.PanicOnError(err)

			// the execute method already take care of the rendering, nothing to do
			if response.Template == "" {
				return
			}

			// pongo does not support template
			// context["base_template"] = "layouts/base.tpl"

		} else {
			response.Template = "pages/not_found.tpl"
			response.StatusCode = 404
		}

		tpl, err := pongo.FromFile(response.Template)

		helper.PanicOnError(err)

		data, err := tpl.ExecuteBytes(response.Context)

		if err != nil {
			res.Header().Set("Content-Type", "text/html; charset=UTF-8")
			res.WriteHeader(500)
			res.Write([]byte("<html><head><title>Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Sorry, an unexpected error occurs on the server...</p></body></html>"))

			panic(err)
		} else {
			res.WriteHeader(response.StatusCode)
			res.Write(data)
		}
	}
}

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Prepare(func(app *goapp.App) error {
		r := app.Get("gonode.router").(*router.Router)
		prefix := ""

		r.Handle("prism_format", prefix+"/prism/:uuid.:format", RenderPrism(app))
		r.Handle("prism", prefix+"/prism/:uuid", RenderPrism(app))

		return nil
	})
}
