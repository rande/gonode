// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"database/sql"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	sq "github.com/lann/squirrel"
	pq "github.com/lib/pq"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/vault"
	"github.com/rande/gonode/modules/base"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

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

		app.Set("gonode.http_client", func(app *goapp.App) interface{} {
			return &http.Client{}
		})

		app.Set("gonode.vault.fs", func(app *goapp.App) interface{} {
			return &vault.Vault{
				BaseKey: []byte(""),
				Algo:    "no_op",
				Driver: &vault.DriverFs{
					Root: conf.Filesystem.Path,
				},
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

		app.Set("gonode.postgres.subscriber", func(app *goapp.App) interface{} {
			return base.NewSubscriber(
				conf.Databases["master"].DSN,
				app.Get("logger").(*log.Logger),
			)
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		// need to find a way to trigger the handler registration
		sub := app.Get("gonode.postgres.subscriber").(*base.Subscriber)

		sub.ListenMessage("core_sleep", func(app *goapp.App) base.SubscriberHander {
			return func(notification *pq.Notification) (int, error) {
				duration, _ := time.ParseDuration(notification.Extra)
				time.Sleep(duration)

				return base.PubSubListenContinue, nil
			}
		}(app))

		return nil
	})

	l.Exit(func(app *goapp.App) error {
		logger := app.Get("logger").(*log.Logger)
		logger.WithFields(log.Fields{
			"module": "commands.server",
		}).Info("Closing PostgreSQL connection")

		db := app.Get("gonode.postgres.connection").(*sql.DB)
		err := db.Close()

		if err != nil {
			logger.WithFields(log.Fields{
				"module": "commands.server",
				"error":  err,
			}).Warn("Error while closing the connection")
		}

		logger.WithFields(log.Fields{
			"module": "commands.server",
		}).Info("End closing PostgreSQL connection")

		return err
	})
}
