// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prism

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	sq "github.com/lann/squirrel"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/modules/base"
	"github.com/zenazn/goji/web"
)

func RenderPrism(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)
	handlers := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		var node *base.Node
		var logger *log.Entry

		if _, ok := c.Env["logger"]; ok {
			logger = c.Env["logger"].(*log.Entry)
		}

		format := "html"
		if uuid, ok := c.URLParams["uuid"]; ok {
			reference, _ := base.GetReferenceFromString(uuid)

			node = manager.Find(reference)

			if _, ok := c.URLParams["format"]; ok {
				format = c.URLParams["format"]
			}

		} else {
			// get the path
			lookupPaths := []string{req.URL.Path}

			path := req.URL.Path
			s := strings.Split(req.URL.Path, ".")

			if len(s) > 1 {
				format = s[len(s)-1]
				path = strings.Join(s[0:len(s)-1], "/")
				lookupPaths = append(lookupPaths, path)
			}

			if logger != nil {
				logger.WithFields(log.Fields{
					"module":           "prism.view",
					"node_lookup_path": lookupPaths,
				}).Debug("Search valid node")
			}

			query := manager.SelectBuilder(base.NewSelectOptions()).Where(sq.Eq{"path": lookupPaths})

			nodes := manager.FindBy(query, 0, 2)

			switch nodes.Len() {
			case 0:
				node = nil
			case 1:
				node = nodes.Front().Value.(*base.Node)
			case 2:
				for e := nodes.Front(); e != nil; e = e.Next() {
					n := e.Value.(*base.Node)
					if n.Path == req.URL.Path {
						node = n
						break
					}
				}

			}
		}

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
			res.Header().Set("X-Content-Uuid", node.Uuid.String())

			response.Add("node", node)

			handler := handlers.Get(node)

			if !handler.Support(node, request, response) {
				// the execute method already take care of the rendering, nothing to do
				response.Template = "pages/bad_request.tpl"
				response.StatusCode = http.StatusBadRequest

				if logger != nil {
					logger.WithFields(log.Fields{
						"module":         "prism.view",
						"node_uuid":      node.Uuid.String(),
						"request_format": request.Format,
					}).Debug("ViewHandler does not support current request")
				}
			} else {
				err := handler.Execute(node, request, response)

				if err != nil {
					if logger != nil {
						logger.WithFields(log.Fields{
							"module":          "prism.view",
							"node_uuid":       node.Uuid.String(),
							"node_type":       node.Type,
							"request_format":  request.Format,
							"view_template":   response.Template,
							"response_status": response.StatusCode,
							"error":           err.Error(),
							"view_handler":    fmt.Sprintf("%T", handler),
						}).Warn("Error while executing ViewHandler")
					}

					response.Template = "pages/internal_error.tpl"
					response.StatusCode = http.StatusInternalServerError
				}

				if response.Template == "" {
					return
				}
				// pongo does not support template
				// context["base_template"] = "layouts/base.tpl"
			}

		} else {
			response.Template = "pages/not_found.tpl"
			response.StatusCode = http.StatusNotFound
		}

		if logger != nil {
			logger.WithFields(log.Fields{
				"module":            "prism.view",
				"request_format":    request.Format,
				"response_template": response.Template,
				"response_status":   response.StatusCode,
			}).Debug("Render node from ViewHandler")
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

func PrismPath(router *router.Router) func(nv *pongo2.Value, vparams ...*pongo2.Value) *pongo2.Value {

	return func(nv *pongo2.Value, vparams ...*pongo2.Value) *pongo2.Value {
		var route string

		if nv.Interface() == nil {
			return pongo2.AsSafeValue("no-node")
		}

		node := nv.Interface().(*base.Node)

		params := url.Values{}
		if len(vparams) > 0 {
			params = vparams[0].Interface().(url.Values)
		}

		if len(node.Path) > 0 {
			params.Set("path", node.Path[1:])

			if len(params.Get("format")) > 0 {
				route = "prism_path_format"
			} else {
				route = "prism_path"
			}
		} else {
			params.Set("uuid", node.Uuid.String())

			if len(params.Get("format")) > 0 {
				route = "prism_format"
			} else {
				route = "prism"
			}
		}

		path, err := router.GeneratePath(route, params)

		if err != nil {
			panic(err)
		}

		return pongo2.AsSafeValue(path)
	}
}

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		r := app.Get("gonode.router").(*router.Router)
		prefix := ""

		r.Handle("prism_format", prefix+"/prism/:uuid.:format", RenderPrism(app))
		r.Handle("prism", prefix+"/prism/:uuid", RenderPrism(app))
		r.Handle("prism_path_catch_all", prefix+"/*", RenderPrism(app))

		// this should be never call, only there for route generation
		r.Handle("prism_path_format", prefix+"/:path.:format", func(c web.C, res http.ResponseWriter, req *http.Request) {})
		r.Handle("prism_path", prefix+"/:path", func(c web.C, res http.ResponseWriter, req *http.Request) {})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {

		router := app.Get("gonode.router").(*router.Router)
		pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)
		pongo.Globals["prism_path"] = PrismPath(router)

		return nil
	})
}
