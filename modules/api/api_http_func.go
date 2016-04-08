// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bufio"
	"container/list"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
	"github.com/zenazn/goji/web"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func versionChecker(c web.C, res http.ResponseWriter) error {
	if c.URLParams["version"] == "v1.0" {
		// for now there is only one version
		return nil
	}

	return base.InvalidVersionError
}

func Check(c web.C, res http.ResponseWriter, req *http.Request, attrs security.Attributes, auth security.AuthorizationChecker) bool {
	token := security.GetTokenFromContext(c)

	if err := security.CheckAccess(token, attrs, res, req, auth); err != nil {
		base.HandleError(req, res, err)

		return false
	}

	if err := versionChecker(c, res); err != nil {
		base.HandleError(req, res, err)

		return false
	}

	return true
}

func Api_GET_Hello(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		if err := versionChecker(c, res); err != nil {
			base.HandleError(req, res, err)

			return
		}

		res.Write([]byte("Hello!"))
	}
}

func Api_GET_Stream(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		attrs := security.Attributes{"node:api:master", "node:api:stream"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		webSocketList := app.Get("gonode.websocket.clients").(*list.List)

		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		ws, err := upgrader.Upgrade(res, req, nil)

		helper.PanicOnError(err)

		element := webSocketList.PushBack(ws)

		var closeDefer = func() {
			ws.Close()
			webSocketList.Remove(element)
		}

		defer closeDefer()

		go func(c *websocket.Conn) {
			for {
				if _, _, err := c.NextReader(); err != nil {
					return
				}
			}
		}(ws)

		// ping remote client, avoid keeping open connection
		for {
			time.Sleep(2 * time.Second)
			if err := ws.WriteMessage(websocket.TextMessage, []byte("PING")); err != nil {
				return
			}
		}
	}
}

func Api_GET_Node(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	apiHandler := app.Get("gonode.api").(*Api)
	handler_collection := app.Get("gonode.handler_collection").(base.Handlers)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:read"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		values := req.URL.Query()

		options := base.NewAccessOptionsFromToken(token)

		if _, raw := values["raw"]; raw {
			// ask for binary content
			reference, err := base.GetReferenceFromString(c.URLParams["uuid"])

			if err != nil {
				base.HandleError(req, res, err)

				return
			}

			node := manager.Find(reference)

			if node == nil {
				base.HandleError(req, res, base.NotFoundError)

				return
			}

			if granted, err := authorizer.IsGranted(token, options.Roles, node); err != nil {
				base.HandleError(req, res, err)
				return
			} else if !granted {
				base.HandleError(req, res, base.AccessForbiddenError)
				return
			}

			handler := handler_collection.Get(node)

			var data *base.DownloadData

			if h, ok := handler.(base.DownloadNodeHandler); ok {
				data = h.GetDownloadData(node)
			} else {
				data = base.GetDownloadData()
			}

			res.Header().Set("Content-Type", data.ContentType)

			data.Stream(node, res)
		} else {
			// send the json value
			res.Header().Set("Content-Type", "application/json")

			options := base.NewAccessOptionsFromToken(token)

			err := apiHandler.FindOne(c.URLParams["uuid"], res, options)

			base.HandleError(req, res, err)
		}
	}
}

func Api_GET_Node_Revisions(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	searchBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)
	searchParser := app.Get("gonode.search.parser.http").(*search.HttpSearchParser)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:revisions"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		searchForm := searchParser.HandleSearch(res, req)

		selectOptions := base.NewSelectOptions()
		selectOptions.TableSuffix = "nodes_audit"

		query := apiHandler.SelectBuilder(selectOptions).
			Where("uuid = ?", c.URLParams["uuid"])

		options := base.NewAccessOptionsFromToken(token)

		apiHandler.Find(res, searchBuilder.BuildQuery(searchForm, query), searchForm.Page, searchForm.PerPage, options)
	}
}

func Api_GET_Node_Revision(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:revision"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		selectOptions := base.NewSelectOptions()
		selectOptions.TableSuffix = "nodes_audit"

		query := apiHandler.SelectBuilder(selectOptions).
			Where("uuid = ?", c.URLParams["uuid"]).
			Where("revision = ?", c.URLParams["rev"])

		options := base.NewAccessOptionsFromToken(token)

		err := apiHandler.FindOneBy(query, res, options)

		base.HandleError(req, res, err)
	}
}

func Api_POST_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:create"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		w := bufio.NewWriter(res)

		options := base.NewAccessOptionsFromToken(token)

		err := apiHandler.Save(req.Body, w, options)

		if err == nil {
			res.WriteHeader(http.StatusCreated)
			w.Flush()
		} else {
			base.HandleError(req, res, err)
		}
	}
}

func Api_PUT_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	apiHandler := app.Get("gonode.api").(*Api)
	handler_collection := app.Get("gonode.handler_collection").(base.Handlers)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:update"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		values := req.URL.Query()

		if _, raw := values["raw"]; raw {
			// send binary data
			reference, err := base.GetReferenceFromString(c.URLParams["uuid"])

			if err != nil {
				base.HandleError(req, res, err)
				return
			}

			node := manager.Find(reference)

			if node == nil {
				base.HandleError(req, res, base.NotFoundError)
				return
			}

			handler := handler_collection.Get(node)

			if h, ok := handler.(base.StoreStreamNodeHandler); ok {
				_, err = h.StoreStream(node, req.Body)
			} else {
				_, err = base.DefaultHandlerStoreStream(node, req.Body)
			}

			// we don't save a new revision as we just need to attach binary to current node
			manager.Save(node, false)

			if err != nil {
				base.HandleError(req, res, err)
			} else {
				manager.Save(node, false)

				helper.SendWithHttpCode(res, http.StatusOK, "binary stored")
			}

		} else {
			w := bufio.NewWriter(res)
			options := base.NewAccessOptionsFromToken(token)

			err := apiHandler.Save(req.Body, w, options)

			if err != nil {
				base.HandleError(req, res, err)
			} else {
				w.Flush()
			}
		}
	}
}

func Api_PUT_Nodes_Move(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:move"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		options := base.NewAccessOptionsFromToken(token)

		err := apiHandler.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], res, options)

		base.HandleError(req, res, err)
	}
}

func Api_DELETE_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:delete"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		options := base.NewAccessOptionsFromToken(token)

		err := apiHandler.RemoveOne(c.URLParams["uuid"], res, options)

		base.HandleError(req, res, err)
	}
}

func Api_PUT_Notify(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		attrs := security.Attributes{"node:api:master", "node:api:notify"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		body, _ := ioutil.ReadAll(req.Body)

		manager.Notify(c.URLParams["name"], string(body[:]))
	}
}

func Api_GET_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	apiHandler := app.Get("gonode.api").(*Api)
	searchBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)
	searchParser := app.Get("gonode.search.parser.http").(*search.HttpSearchParser)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:list"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		searchForm := searchParser.HandleSearch(res, req)

		if searchForm == nil {
			return
		}

		query := searchBuilder.BuildQuery(searchForm, manager.SelectBuilder(base.NewSelectOptions()))

		options := base.NewAccessOptionsFromToken(token)

		apiHandler.Find(res, query, searchForm.Page, searchForm.PerPage, options)
	}
}
