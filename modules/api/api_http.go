// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"bufio"
	"container/list"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			return
		}
	}
}

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		app.Set("gonode.websocket.clients", func(app *goapp.App) interface{} {
			return list.New()
		})

		sub := app.Get("gonode.postgres.subscriber").(*base.Subscriber)
		sub.ListenMessage(conf.Databases["master"].Prefix+"_manager_action", func(notification *pq.Notification) (int, error) {
			logger := app.Get("logger").(*log.Logger)
			logger.WithFields(log.Fields{
				"module":  "api.websocket",
				"payload": notification.Extra,
			}).Debug("Sending message")

			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			for e := webSocketList.Front(); e != nil; e = e.Next() {
				if err := e.Value.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(notification.Extra)); err != nil {
					logger.Warn("Error writing to websocket")
				}
			}

			logger.WithFields(log.Fields{
				"module": "api.websocket",
			}).Debug("WebSocket: End Sending message")

			return base.PubSubListenContinue, nil
		})

		graceful.PreHook(func() {
			logger := app.Get("logger").(*log.Logger)
			webSocketList := app.Get("gonode.websocket.clients").(*list.List)

			logger.WithFields(log.Fields{
				"module": "api.websocket",
			}).Info("Closing websocket connections")

			for e := webSocketList.Front(); e != nil; e = e.Next() {
				e.Value.(*websocket.Conn).Close()
			}
		})

		return nil
	})

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		logger := app.Get("logger").(*log.Logger)

		logger.WithFields(log.Fields{
			"module": "api.websocket",
		}).Info("Starting PostgreSQL subcriber")

		app.Get("gonode.postgres.subscriber").(*base.Subscriber).Register()

		return nil
	})

	l.Exit(func(app *goapp.App) error {
		logger := app.Get("logger").(*log.Logger)
		logger.WithFields(log.Fields{
			"module": "api.websocket",
		}).Info("Closing PostgreSQL subcriber")

		app.Get("gonode.postgres.subscriber").(*base.Subscriber).Stop()

		logger.WithFields(log.Fields{
			"module": "api.websocket",
		}).Info("End closing PostgreSQL subcriber")

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		apiHandler := app.Get("gonode.api").(*Api)
		handler_collection := app.Get("gonode.handler_collection").(base.Handlers)
		searchBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)
		searchParser := app.Get("gonode.search.parser.http").(*search.HttpSearchParser)

		mux.Get(conf.Api.Prefix+"/:version/hello", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Hello!"))
		})

		mux.Get(conf.Api.Prefix+"/:version/nodes/stream", func(res http.ResponseWriter, req *http.Request) {
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

			go readLoop(ws)

			// ping remote client, avoid keeping open connection
			for {
				time.Sleep(2 * time.Second)
				if err := ws.WriteMessage(websocket.TextMessage, []byte("PING")); err != nil {
					return
				}
			}
		})

		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // ask for binary content
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

				if err == base.NotFoundError {
					helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
				}

				if err != nil {
					helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
				}
			}
		})

		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid/revisions", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			searchForm := searchParser.HandleSearch(res, req)

			options := base.NewSelectOptions()
			options.TableSuffix = "nodes_audit"

			query := apiHandler.SelectBuilder(options).
				Where("uuid = ?", c.URLParams["uuid"])

			apiHandler.Find(res, searchBuilder.BuildQuery(searchForm, query), searchForm.Page, searchForm.PerPage)
		})

		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid/revisions/:rev", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			options := base.NewSelectOptions()
			options.TableSuffix = "nodes_audit"

			query := apiHandler.SelectBuilder(options).
				Where("uuid = ?", c.URLParams["uuid"]).
				Where("revision = ?", c.URLParams["rev"])

			err := apiHandler.FindOneBy(query, res)

			if err == base.NotFoundError {
				helper.SendWithHttpCode(res, http.StatusNotFound, err.Error())
			}

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Post(conf.Api.Prefix+"/:version/nodes", func(res http.ResponseWriter, req *http.Request) {
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
		})

		mux.Put(conf.Api.Prefix+"/:version/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			values := req.URL.Query()

			if _, raw := values["raw"]; raw { // send binary data
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
		})

		mux.Put(conf.Api.Prefix+"/:version/nodes/move/:uuid/:parentUuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			err := apiHandler.Move(c.URLParams["uuid"], c.URLParams["parentUuid"], res)

			if err != nil {
				helper.SendWithHttpCode(res, http.StatusInternalServerError, err.Error())
			}
		})

		mux.Delete(conf.Api.Prefix+"/:version/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
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
		})

		mux.Put(conf.Api.Prefix+"/:version/notify/:name", func(c web.C, res http.ResponseWriter, req *http.Request) {
			body, _ := ioutil.ReadAll(req.Body)

			manager.Notify(c.URLParams["name"], string(body[:]))
		})

		mux.Get(conf.Api.Prefix+"/:version/nodes", func(c web.C, res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")

			searchForm := searchParser.HandleSearch(res, req)

			if searchForm == nil {
				return
			}

			query := searchBuilder.BuildQuery(searchForm, manager.SelectBuilder(base.NewSelectOptions()))

			apiHandler.Find(res, query, searchForm.Page, searchForm.PerPage)
		})

		return nil
	})
}
