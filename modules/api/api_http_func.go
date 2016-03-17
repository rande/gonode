// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bufio"
	"container/list"
	"errors"
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

var (
	InvalidVersion  = errors.New("Invalid version")
	AccessForbidden = errors.New("Access Forbidden")
)

func getToken(c web.C) security.SecurityToken {
	if _, ok := c.Env["guard_token"]; !ok { // no token
		return nil
	}

	return c.Env["guard_token"].(security.SecurityToken)
}

func checkAccess(options *ApiOptions, res http.ResponseWriter, req *http.Request, auth security.AuthorizationChecker) error {

	if options.Token == nil { // no token
		helper.SendWithHttpCode(res, http.StatusForbidden, "Access Forbidden: missing token")

		return AccessForbidden
	}

	r, err := auth.IsGranted(options.Token, options.Roles, req)

	if err != nil {
		helper.SendWithHttpCode(res, http.StatusForbidden, "Access Forbidden: error while checking credentials")

		return err
	}

	if r == false {
		helper.SendWithHttpCode(res, http.StatusForbidden, "Access Forbidden: not enough credentials")

		return AccessForbidden
	}

	return nil
}

func versionChecker(c web.C, res http.ResponseWriter) error {
	if c.URLParams["version"] == "v1.0" { // for now there is only one version
		return nil
	}

	helper.SendWithHttpCode(res, http.StatusBadRequest, "Invalid version")

	return InvalidVersion
}

func Api_GET_Hello(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Write([]byte("Hello!"))
	}
}

func Api_GET_Stream(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:stream"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
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

		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:get"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		values := req.URL.Query()

		if _, raw := values["raw"]; raw {
			// ask for binary content
			reference, err := base.GetReferenceFromString(c.URLParams["uuid"])

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

				return
			}

			node := manager.Find(reference)

			if node == nil {
				helper.SendWithHttpCode(res, http.StatusNotFound, "Element not found")

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
			err := apiHandler.FindOne(c.URLParams["uuid"], res)

			if err == nil {
				return
			}

			statusCode := http.StatusInternalServerError

			if err == base.NotFoundError {
				statusCode = http.StatusNotFound
			}

			helper.SendWithHttpCode(res, statusCode, err.Error())
		}
	}
}

func Api_GET_Node_Revisions(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	searchBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)
	searchParser := app.Get("gonode.search.parser.http").(*search.HttpSearchParser)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:revisions"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		searchForm := searchParser.HandleSearch(res, req)

		selectOptions := base.NewSelectOptions()
		selectOptions.TableSuffix = "nodes_audit"

		query := apiHandler.SelectBuilder(selectOptions).
			Where("uuid = ?", c.URLParams["uuid"])

		apiHandler.Find(res, searchBuilder.BuildQuery(searchForm, query), searchForm.Page, searchForm.PerPage)
	}
}

func Api_GET_Node_Revision(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:revision"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		selectOptions := base.NewSelectOptions()
		selectOptions.TableSuffix = "nodes_audit"

		query := apiHandler.SelectBuilder(selectOptions).
			Where("uuid = ?", c.URLParams["uuid"]).
			Where("revision = ?", c.URLParams["rev"])

		err := apiHandler.FindOneBy(query, res)

		if err == base.NotFoundError {
			helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
		}

		if err != nil {
			helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
		}
	}
}

func Api_POST_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:create"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		w := bufio.NewWriter(res)

		err := apiHandler.Save(req.Body, w)

		if err == base.RevisionError {
			res.WriteHeader(http.StatusConflict)
		}

		if err == base.ValidationError {
			res.WriteHeader(http.StatusPreconditionFailed)
		}

		res.WriteHeader(http.StatusCreated)

		w.Flush()
	}
}

func Api_PUT_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	apiHandler := app.Get("gonode.api").(*Api)
	handler_collection := app.Get("gonode.handler_collection").(base.Handlers)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:update"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		values := req.URL.Query()

		if _, raw := values["raw"]; raw {
			// send binary data
			reference, err := base.GetReferenceFromString(c.URLParams["uuid"])

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, "Unable to parse the reference")

				return
			}

			node := manager.Find(reference)

			if node == nil {
				helper.SendWithHttpCode(res, http.StatusNotFound, "Element not found")
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
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			} else {
				manager.Save(node, false)

				helper.SendWithHttpCode(res, http.StatusOK, "binary stored")
			}

		} else {
			w := bufio.NewWriter(res)

			err := apiHandler.Save(req.Body, w)

			if err == base.RevisionError {
				res.WriteHeader(http.StatusConflict)
			}

			if err == base.ValidationError {
				res.WriteHeader(http.StatusPreconditionFailed)
			}

			w.Flush()
		}
	}
}

func Api_PUT_Nodes_Move(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:move"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		err := apiHandler.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], res)

		if err != nil {
			helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
		}
	}
}

func Api_DELETE_Nodes(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	apiHandler := app.Get("gonode.api").(*Api)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:delete"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		err := apiHandler.RemoveOne(c.URLParams["uuid"], res)

		if err == base.NotFoundError {
			helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
			return
		}

		if err == base.AlreadyDeletedError {
			helper.SendWithHttpCode(res, http.StatusGone, err.Error())
			return
		}

		if err != nil {
			helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
		}
	}
}

func Api_PUT_Notify(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	authorizer := app.Get("security.authorizer").(security.AuthorizationChecker)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:notify"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
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
		options := &ApiOptions{
			Token: getToken(c),
			Roles: security.Attributes{"node:api:master", "node:api:list"},
		}

		if err := checkAccess(options, res, req, authorizer); err != nil {
			return
		}

		if err := versionChecker(c, res); err != nil {
			return
		}

		res.Header().Set("Content-Type", "application/json")

		searchForm := searchParser.HandleSearch(res, req)

		if searchForm == nil {
			return
		}

		query := searchBuilder.BuildQuery(searchForm, manager.SelectBuilder(base.NewSelectOptions()))

		apiHandler.Find(res, query, searchForm.Page, searchForm.PerPage)
	}
}
