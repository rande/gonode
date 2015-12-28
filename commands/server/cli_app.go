// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"database/sql"
	"github.com/hypebeast/gojistatic"
	sq "github.com/lann/squirrel"
	pq "github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/plugins/api"
	"github.com/rande/gonode/plugins/blog"
	"github.com/rande/gonode/plugins/debug"
	"github.com/rande/gonode/plugins/feed"
	"github.com/rande/gonode/plugins/media"
	"github.com/rande/gonode/plugins/search"
	"github.com/rande/gonode/plugins/user"
	"github.com/rande/gonode/plugins/vault"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
			return conf
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

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.vault.fs", func(app *goapp.App) interface{} {
			return &vault.Vault{
				BaseKey: []byte(""),
				Algo:    "no_op",
				Driver: &vault.DriverFs{
					Root: conf.Filesystem.Path,
				},
			}
		})

		app.Set("gonode.http_client", func(app *goapp.App) interface{} {
			return &http.Client{}
		})

		app.Set("gonode.handler_collection", func(app *goapp.App) interface{} {
			return core.HandlerCollection{
				"default": &debug.DefaultHandler{},
				"media.image": &media.ImageHandler{
					Vault: app.Get("gonode.vault.fs").(*vault.Vault),
				},
				"media.youtube": &media.YoutubeHandler{},
				"blog.post":     &blog.PostHandler{},
				"core.user":     &user.UserHandler{},
				"core.index":    &search.IndexHandler{},
			}
		})

		app.Set("gonode.view_handler_collection", func(app *goapp.App) interface{} {
			return core.ViewHandlerCollection{
				"default": &debug.DefaultViewHandler{},
				"core.index": &search.IndexViewHandler{
					Search:    app.Get("gonode.search.pgsql").(*search.SearchPGSQL),
					Manager:   app.Get("gonode.manager").(*core.PgNodeManager),
					MaxResult: 128,
				},
			}
		})

		app.Set("gonode.manager", func(app *goapp.App) interface{} {
			return &core.PgNodeManager{
				Logger:   app.Get("logger").(*log.Logger),
				Db:       app.Get("gonode.postgres.connection").(*sql.DB),
				ReadOnly: false,
				Handlers: app.Get("gonode.handler_collection").(core.Handlers),
				Prefix:   conf.Databases["master"].Prefix,
			}
		})

		app.Set("gonode.postgres.connection", func(app *goapp.App) interface{} {
			sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
			db, err := sql.Open("postgres", conf.Databases["master"].DSN)

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
			return &api.Api{
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
			return core.NewSubscriber(
				conf.Databases["master"].DSN,
				app.Get("logger").(*log.Logger),
			)
		})

		app.Set("gonode.listener.youtube", func(app *goapp.App) interface{} {
			client := app.Get("gonode.http_client").(*http.Client)

			return &media.YoutubeListener{
				HttpClient: client,
			}
		})

		app.Set("gonode.listener.file_downloader", func(app *goapp.App) interface{} {
			return &media.ImageDownloadListener{
				Vault:      app.Get("gonode.vault.fs").(*vault.Vault),
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
			listener := app.Get("gonode.listener.youtube").(*media.YoutubeListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("media_file_download", func(app *goapp.App) core.SubscriberHander {
			manager := app.Get("gonode.manager").(*core.PgNodeManager)
			listener := app.Get("gonode.listener.file_downloader").(*media.ImageDownloadListener)

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
