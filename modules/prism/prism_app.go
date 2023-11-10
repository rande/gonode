// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prism

import (
	"errors"
	"fmt"
	tpl "html/template"
	"net/http"
	"net/url"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/template"
	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

func RenderPrism(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	loader := app.Get("gonode.template").(*template.TemplateLoader)
	handlers := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		var node *base.Node
		var logger *log.Entry

		token := security.GetTokenFromContext(c)

		if _, ok := c.Env["logger"]; ok {
			logger = c.Env["logger"].(*log.Entry)
		}

		// the user must have roles,
		// this test should be useless as the first check is present.
		// however this might need to be usefull as the authorizer might have different implementations.
		if len(token.GetRoles()) == 0 {
			if logger != nil {
				logger.WithFields(log.Fields{
					"module": "prism.view",
				}).Debug("No roles associate with current token")
			}

			base.HandleError(req, res, base.ErrAccessForbidden)
			return
		}

		format := "html"
		if uuid, ok := c.URLParams["uuid"]; ok {
			reference, err := base.GetReferenceFromString(uuid)

			if err != nil {
				base.HandleError(req, res, err)
				return
			}

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

			query := manager.
				SelectBuilder(base.NewSelectOptions()).
				Where(sq.Eq{"path": lookupPaths})

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
			if granted, err := authorizer.IsGranted(token, nil, node); err != nil {
				if logger != nil {
					logger.WithFields(log.Fields{
						"module": "prism.view",
						"error":  err,
					}).Debug("Authorization generates an error")
				}

				base.HandleError(req, res, err)
				return
			} else if !granted {

				if logger != nil {
					logger.WithFields(log.Fields{
						"module": "prism.view",
					}).Debug("Authorization not granded to access the node")
				}

				base.HandleError(req, res, base.ErrAccessForbidden)
				return
			}

			response.Add("node", node)

			handler := handlers.Get(node)

			if !handler.Support(node, request, response) {
				// the execute method already take care of the rendering, nothing to do
				response.Template = "prism:pages/bad_request"
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

					response.Template = "prism:pages/internal_error"
					response.StatusCode = http.StatusInternalServerError
				}

				if response.Template == "" {
					return
				}
			}

		} else {
			response.Template = "prism:pages/not_found"
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

		data, err := loader.Execute(response.Template, response.Context)

		if err != nil {
			res.Header().Set("Content-Type", "text/html; charset=UTF-8")
			res.WriteHeader(500)
			res.Write([]byte("<html><head><title>Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Sorry, an unexpected error occurs on the server...</p></body></html>"))

			fmt.Printf("Template: %s\n", response.Template)
			fmt.Printf("Error: %s %s\n", err, errors.Unwrap(err))
			panic(err)
		} else {
			res.WriteHeader(response.StatusCode)
			res.Write(data)
		}
	}
}

func PrismPath(router *router.Router) func(node *base.Node, params ...interface{}) tpl.HTML {

	return func(node *base.Node, options ...interface{}) tpl.HTML {
		var route string

		params := url.Values{}
		if len(options) > 0 {
			params = options[0].(url.Values)
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

		return tpl.HTML(path)
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

	l.Register(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("prism", GetEmbedFS())

		return nil
	})

	l.Config(func(app *goapp.App) error {
		loader := app.Get("gonode.template").(*template.TemplateLoader)
		router := app.Get("gonode.router").(*router.Router)

		loader.FuncMap["prism_path"] = PrismPath(router)

		return nil
	})
}
