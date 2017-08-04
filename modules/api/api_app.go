// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"container/list"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		app.Set("gonode.api", func(app *goapp.App) interface{} {
			return &Api{
				Manager:    app.Get("gonode.manager").(*base.PgNodeManager),
				Version:    "1.0.0",
				Logger:     app.Get("logger").(*log.Logger),
				Authorizer: app.Get("security.authorizer").(security.AuthorizationChecker),
			}
		})

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

		mux := app.Get("goji.mux").(*web.Mux)

		mux.Get(conf.Api.Prefix+"/:version/nodes/stream", Api_GET_Stream(app))
		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid", Api_GET_Node(app))
		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid/revisions", Api_GET_Node_Revisions(app))
		mux.Get(conf.Api.Prefix+"/:version/nodes/:uuid/revisions/:rev", Api_GET_Node_Revision(app))
		mux.Post(conf.Api.Prefix+"/:version/nodes", Api_POST_Nodes(app))
		mux.Put(conf.Api.Prefix+"/:version/nodes/:uuid", Api_PUT_Nodes(app))
		mux.Put(conf.Api.Prefix+"/:version/nodes/move/:uuid/:parentUuid", Api_PUT_Nodes_Move(app))
		mux.Delete(conf.Api.Prefix+"/:version/nodes/:uuid", Api_DELETE_Nodes(app))
		mux.Get(conf.Api.Prefix+"/:version/nodes", Api_GET_Nodes(app))
		mux.Get(conf.Api.Prefix+"/:version/hello", Api_GET_Hello(app))
		mux.Put(conf.Api.Prefix+"/:version/notify/:name", Api_PUT_Notify(app))
		mux.Get(conf.Api.Prefix+"/:version/handlers/node", Api_GET_Handlers_Node(app))
		mux.Get(conf.Api.Prefix+"/:version/handlers/view", Api_GET_Handlers_View(app))
		mux.Get(conf.Api.Prefix+"/:version/services", Api_GET_Services(app))

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
}
