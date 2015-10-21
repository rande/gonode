// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"github.com/rande/goapp"

	"fmt"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/handlers"
	"github.com/rande/gonode/test/fixtures"
	"github.com/rande/gonode/vault"
	"net/http"

	"database/sql"
	sq "github.com/lann/squirrel"
	pq "github.com/lib/pq"

	"log"
	"time"

	"github.com/hypebeast/gojistatic"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"os"
)

func ConfigureServer(l *goapp.Lifecycle, config *core.ServerConfig) {

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
			return config
		})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("logger", func(app *goapp.App) interface{} {
			return log.New(os.Stdout, "", log.Lshortfile)
		})

		app.Set("goji.mux", func(app *goapp.App) interface{} {
			mux := web.New()

			mux.Use(middleware.RequestID)
			mux.Use(middleware.Logger)
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.AutomaticOptions)
			mux.Use(gojistatic.Static("dist", gojistatic.StaticOptions{SkipLogging: true, Prefix: "dist"}))

			return mux
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		if !config.Test {
			return nil
		}

		mux := app.Get("goji.mux").(*web.Mux)

		prefix := ""

		mux.Put(prefix+"/data/purge", func(res http.ResponseWriter, req *http.Request) {

			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			configuration := app.Get("gonode.configuration").(*core.ServerConfig)

			prefix := configuration.Databases["master"].Prefix

			tx, _ := manager.Db.Begin()
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes_audit"`, prefix))
			err := tx.Commit()

			if err != nil {
				core.SendWithStatus("KO", err.Error(), res)
			} else {
				core.SendWithStatus("OK", "Data purged!", res)
			}
		})

		mux.Put(prefix+"/data/load", func(res http.ResponseWriter, req *http.Request) {
			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			nodes := manager.FindBy(manager.SelectBuilder(), 0, 10)

			if nodes.Len() != 0 {
				core.SendWithStatus("KO", "Table contains data, purge the data first!", res)

				return
			}

			err := fixtures.LoadFixtures(manager, 100)

			if err != nil {
				core.SendWithStatus("KO", err.Error(), res)
			} else {
				core.SendWithStatus("OK", "Data loaded!", res)
			}
		})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.vault.fs", func(app *goapp.App) interface{} {
			configuration := app.Get("gonode.configuration").(*core.ServerConfig)

			return &vault.VaultFs{
				BaseKey: []byte(""),
				Algo:    "no_op",
				Root:    configuration.Filesystem.Path,
			}
		})

		app.Set("gonode.http_client", func(app *goapp.App) interface{} {
			return &http.Client{}
		})

		app.Set("gonode.handler_collection", func(app *goapp.App) interface{} {
			return core.HandlerCollection{
				"default": &handlers.DefaultHandler{},
				"media.image": &handlers.ImageHandler{
					Vault: app.Get("gonode.vault.fs").(vault.Vault),
				},
				"media.youtube": &handlers.YoutubeHandler{},
				"blog.post":     &handlers.PostHandler{},
				"core.user":     &handlers.UserHandler{},
			}
		})

		app.Set("gonode.manager", func(app *goapp.App) interface{} {
			configuration := app.Get("gonode.configuration").(*core.ServerConfig)

			return &core.PgNodeManager{
				Logger:   app.Get("logger").(*log.Logger),
				Db:       app.Get("gonode.postgres.connection").(*sql.DB),
				ReadOnly: false,
				Handlers: app.Get("gonode.handler_collection").(core.Handlers),
				Prefix:   configuration.Databases["master"].Prefix,
			}
		})

		app.Set("gonode.postgres.connection", func(app *goapp.App) interface{} {

			configuration := app.Get("gonode.configuration").(*core.ServerConfig)

			sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
			db, err := sql.Open("postgres", configuration.Databases["master"].DSN)

			db.SetMaxIdleConns(8)
			db.SetMaxOpenConns(64)

			if err != nil {
				log.Fatal(err)
			}

			err = db.Ping()
			if err != nil {
				log.Fatal(err)
			}

			return db
		})

		app.Set("gonode.api", func(app *goapp.App) interface{} {
			return &core.Api{
				Manager:    app.Get("gonode.manager").(*core.PgNodeManager),
				Version:    "1.0.0",
				Serializer: app.Get("gonode.node.serializer").(*core.Serializer),
				Logger:     app.Get("logger").(*log.Logger),
			}
		})

		app.Set("gonode.node.serializer", func(app *goapp.App) interface{} {
			s := core.NewSerializer()
			s.Handlers = app.Get("gonode.handler_collection").(core.Handlers)

			return s
		})

		app.Set("gonode.postgres.subscriber", func(app *goapp.App) interface{} {
			configuration := app.Get("gonode.configuration").(*core.ServerConfig)

			return core.NewSubscriber(
				configuration.Databases["master"].DSN,
				app.Get("logger").(*log.Logger),
			)
		})

		app.Set("gonode.listener.youtube", func(app *goapp.App) interface{} {
			client := app.Get("gonode.http_client").(*http.Client)

			return &handlers.YoutubeListener{
				HttpClient: client,
			}
		})

		app.Set("gonode.listener.file_downloader", func(app *goapp.App) interface{} {
			return &handlers.ImageDownloadListener{
				Vault:      app.Get("gonode.vault.fs").(vault.Vault),
				HttpClient: app.Get("gonode.http_client").(*http.Client),
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		// need to find a way to trigger the handler registration
		sub := app.Get("gonode.postgres.subscriber").(*core.Subscriber)

		sub.ListenMessage("media_youtube_update", func(app *goapp.App) core.SubscriberHander {
			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			listener := app.Get("gonode.listener.youtube").(*handlers.YoutubeListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("media_file_download", func(app *goapp.App) core.SubscriberHander {
			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			listener := app.Get("gonode.listener.file_downloader").(*handlers.ImageDownloadListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("core_sleep", func(app *goapp.App) core.SubscriberHander {
			return func(notification *pq.Notification) (int, error) {

				logger := app.Get("logger").(*log.Logger)

				duration, _ := time.ParseDuration(notification.Extra)

				logger.Printf("[core_sleep] sleep ...")

				time.Sleep(duration)

				logger.Printf("[core_sleep] wake up ...")

				return core.PubSubListenContinue, nil
			}
		}(app))

		return nil
	})
}
