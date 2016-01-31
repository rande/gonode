// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	sq "github.com/lann/squirrel"
	pq "github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/api"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/modules/debug"
	"github.com/rande/gonode/modules/feed"
	"github.com/rande/gonode/modules/media"
	"github.com/rande/gonode/modules/raw"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/user"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"net/http"
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
		app.Set("goji.mux", func(app *goapp.App) interface{} {
			mux := web.New()

			mux.Use(middleware.RequestID)
			mux.Use(middleware.Recoverer)
			mux.Use(middleware.AutomaticOptions)

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
			return base.HandlerCollection{
				"default": &debug.DefaultHandler{},
				"media.image": &media.ImageHandler{
					Vault: app.Get("gonode.vault.fs").(*vault.Vault),
				},
				"media.youtube": &media.YoutubeHandler{},
				"blog.post":     &blog.PostHandler{},
				"core.user":     &user.UserHandler{},
				"core.index":    &search.IndexHandler{},
				"feed.index":    &feed.FeedHandler{},
				"core.raw":      &raw.RawHandler{},
			}
		})

		app.Set("gonode.view_handler_collection", func(app *goapp.App) interface{} {
			return base.ViewHandlerCollection{
				"default": &debug.DefaultViewHandler{},
				"core.index": &search.IndexViewHandler{
					Search:    app.Get("gonode.search.pgsql").(*search.SearchPGSQL),
					Manager:   app.Get("gonode.manager").(*base.PgNodeManager),
					MaxResult: 128,
				},
				"feed.index": &feed.FeedViewHandler{
					Search:  app.Get("gonode.search.pgsql").(*search.SearchPGSQL),
					Manager: app.Get("gonode.manager").(*base.PgNodeManager),
				},
				"core.raw": &raw.RawViewHandler{},
				"media.image": &media.MediaViewHandler{
					Vault:         app.Get("gonode.vault.fs").(*vault.Vault),
					MaxWidth:      conf.Media.Image.MaxWidth,
					AllowedWidths: conf.Media.Image.AllowedWidths,
				},
			}
		})

		app.Set("gonode.manager", func(app *goapp.App) interface{} {
			return &base.PgNodeManager{
				Logger:   app.Get("logger").(*log.Logger),
				Db:       app.Get("gonode.postgres.connection").(*sql.DB),
				ReadOnly: false,
				Handlers: app.Get("gonode.handler_collection").(base.Handlers),
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
				Manager:    app.Get("gonode.manager").(*base.PgNodeManager),
				Version:    "1.0.0",
				Serializer: app.Get("gonode.node.serializer").(*base.Serializer),
				Logger:     app.Get("logger").(*log.Logger),
			}
		})

		app.Set("gonode.node.serializer", func(app *goapp.App) interface{} {
			s := base.NewSerializer()
			s.Handlers = app.Get("gonode.handler_collection").(base.Handlers)

			return s
		})

		app.Set("gonode.postgres.subscriber", func(app *goapp.App) interface{} {
			return base.NewSubscriber(
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
		sub := app.Get("gonode.postgres.subscriber").(*base.Subscriber)

		sub.ListenMessage("media_youtube_update", func(app *goapp.App) base.SubscriberHander {
			manager := app.Get("gonode.manager").(*base.PgNodeManager)
			listener := app.Get("gonode.listener.youtube").(*media.YoutubeListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("media_file_download", func(app *goapp.App) base.SubscriberHander {
			manager := app.Get("gonode.manager").(*base.PgNodeManager)
			listener := app.Get("gonode.listener.file_downloader").(*media.ImageDownloadListener)

			return func(notification *pq.Notification) (int, error) {
				return listener.Handle(notification, manager)
			}
		}(app))

		sub.ListenMessage("core_sleep", func(app *goapp.App) base.SubscriberHander {
			return func(notification *pq.Notification) (int, error) {

				logger := app.Get("logger").(*log.Logger)

				duration, _ := time.ParseDuration(notification.Extra)

				logger.Printf("[core_sleep] sleep ...")

				time.Sleep(duration)

				logger.Printf("[core_sleep] wake up ...")

				return base.PubSubListenContinue, nil
			}
		}(app))

		return nil
	})
}
