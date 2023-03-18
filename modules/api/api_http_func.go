// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"container/list"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
	log "github.com/sirupsen/logrus"
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

	var logger *log.Entry

	if l, ok := c.Env["logger"]; ok {
		logger = l.(*log.Entry).WithFields(log.Fields{
			"module": "api.http",
		})
	}

	if err := security.CheckAccess(token, attrs, res, req, auth); err != nil {
		if logger != nil {
			logger.WithFields(log.Fields{
				"attrs": attrs,
				"token": security.GetTokenFromContext(c),
			}).Warn("Unable to check access")
		}

		base.HandleError(req, res, err)

		return false
	}

	if logger != nil {
		logger.WithFields(log.Fields{
			"attrs": attrs,
			"token": security.GetTokenFromContext(c),
		}).Debug("Check access granted")
	}

	if err := versionChecker(c, res); err != nil {

		if logger != nil {
			logger.WithFields(log.Fields{
				"attrs": attrs,
				"token": security.GetTokenFromContext(c),
			}).Warn("invalid version")
		}

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
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

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

			if node, err := apiHandler.FindOne(c.URLParams["uuid"], options); err != nil {
				base.HandleError(req, res, err)
			} else {
				serializer.Serialize(res, node)
			}
		}
	}
}

func Api_GET_Node_Revisions(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	searchBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)
	searchParser := app.Get("gonode.search.parser.http").(*search.HttpSearchParser)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

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

		pager, err := apiHandler.Find(searchBuilder.BuildQuery(searchForm, query), searchForm.Page, searchForm.PerPage, options)

		if err != nil {
			base.HandleError(req, res, err)
		}

		for k, v := range pager.Elements {
			b := bytes.NewBuffer([]byte{})
			serializer.Serialize(b, v)
			message := json.RawMessage(b.Bytes())

			pager.Elements[k] = &message
		}

		base.Serialize(res, pager)
	}
}

func Api_GET_Node_Revision(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

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

		if node, err := apiHandler.FindOneBy(query, options); err != nil {
			base.HandleError(req, res, err)
		} else {
			serializer.Serialize(res, node)
		}
	}
}

func Api_POST_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:create"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		node := base.NewNode()
		if err := serializer.Deserialize(req.Body, node); err != nil {
			base.HandleError(req, res, err)

			return
		}

		options := base.NewAccessOptionsFromToken(token)

		if node, errors, err := apiHandler.Save(node, options); err != nil && err != base.ValidationError {
			base.HandleError(req, res, err)
		} else if errors != nil {
			res.WriteHeader(http.StatusPreconditionFailed)
			base.Serialize(res, errors)
		} else {
			res.WriteHeader(http.StatusCreated)
			serializer.Serialize(res, node)
		}
	}
}

func Api_PUT_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	apiHandler := app.Get("gonode.api").(*Api)
	handler_collection := app.Get("gonode.handler_collection").(base.Handlers)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

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

			if err != nil {
				base.HandleError(req, res, err)
				return
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
			options := base.NewAccessOptionsFromToken(token)

			node := base.NewNode()
			if err := serializer.Deserialize(req.Body, node); err != nil {
				base.HandleError(req, res, err)

				return
			}

			if node, errors, err := apiHandler.Save(node, options); err != nil && err != base.ValidationError {
				base.HandleError(req, res, err)
			} else if errors != nil {
				res.WriteHeader(http.StatusPreconditionFailed)
				base.Serialize(res, errors)
			} else {
				res.WriteHeader(http.StatusCreated)
				serializer.Serialize(res, node)
			}
		}
	}
}

func Api_PUT_Nodes_Move(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:move"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		options := base.NewAccessOptionsFromToken(token)

		if result, err := apiHandler.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], options); err != nil {
			base.HandleError(req, res, err)
		} else {
			serializer.Serialize(res, result)
		}
	}
}

func Api_DELETE_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:delete"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		options := base.NewAccessOptionsFromToken(token)

		if node, err := apiHandler.RemoveOne(c.URLParams["uuid"], options); err != nil {
			base.HandleError(req, res, err)
		} else {
			serializer.Serialize(res, node)
		}
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
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		var logger *log.Entry

		token := security.GetTokenFromContext(c)
		attrs := security.Attributes{"node:api:master", "node:api:list"}

		if !Check(c, res, req, attrs, authorizer) {
			return
		}

		if l, ok := c.Env["logger"]; ok {
			logger = l.(*log.Entry).WithFields(log.Fields{
				"module": "api.http",
			})
		}

		res.Header().Set("Content-Type", "application/json")

		searchForm := searchParser.HandleSearch(res, req)

		if searchForm == nil {
			return
		}

		query := searchBuilder.BuildQuery(searchForm, manager.SelectBuilder(base.NewSelectOptions()))

		options := base.NewAccessOptionsFromToken(token)

		pager, err := apiHandler.Find(query, searchForm.Page, searchForm.PerPage, options)

		if err != nil {
			base.HandleError(req, res, err)
		}

		for k, v := range pager.Elements {
			if logger != nil {
				logger.WithFields(log.Fields{
					"node": v,
				}).Debug("serializing row")
			}

			b := bytes.NewBuffer([]byte{})
			serializer.Serialize(b, v)
			message := json.RawMessage(b.Bytes())

			pager.Elements[k] = &message
		}

		base.Serialize(res, pager)
	}
}
