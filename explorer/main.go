// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/hypebeast/gojistatic"
	"github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/explorer/helper"
	"github.com/rande/gonode/extra"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"net/http"
	"os"
)

func init() {
	bind.WithFlag()
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
	//  graceful.DoubleKickWindow(2 * time.Second)
}

func main() {
	app := goapp.NewApp()

	l := goapp.NewLifecycle()

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
			return extra.GetConfiguration("./config.toml")
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
		mux := app.Get("goji.mux").(*web.Mux)

		prefix := ""

		mux.Put(prefix+"/data/purge", func(res http.ResponseWriter, req *http.Request) {

			manager := app.Get("gonode.manager").(*nc.PgNodeManager)
			configuration := app.Get("gonode.configuration").(*extra.Config)

			prefix := configuration.Databases["master"].Prefix

			tx, _ := manager.Db.Begin()
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes"`, prefix))
			manager.Db.Exec(fmt.Sprintf(`DELETE FROM "%s_nodes_audit"`, prefix))
			err := tx.Commit()

			if err != nil {
				helper.Send("KO", err.Error(), res)
			} else {
				helper.Send("OK", "Data purged!", res)
			}
		})

		mux.Put(prefix+"/data/load", func(res http.ResponseWriter, req *http.Request) {
			manager := app.Get("gonode.manager").(*nc.PgNodeManager)
			nodes := manager.FindBy(manager.SelectBuilder(), 0, 10)

			if nodes.Len() != 0 {
				helper.Send("KO", "Table contains data, purge the data first!", res)

				return
			}

			err := helper.LoadFixtures(manager, 100)

			if err != nil {
				helper.Send("KO", err.Error(), res)
			} else {
				helper.Send("OK", "Data loaded!", res)
			}
		})

		return nil
	})

	extra.ConfigureApp(l)
	extra.ConfigureGoji(l)

	l.Run(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)

		if !flag.Parsed() {
			flag.Parse()
		}

		mux.Compile()

		// Install our handler at the root of the standard net/http default mux.
		// This allows packages like expvar to continue working as expected.
		http.Handle("/", mux)

		listener := bind.Default()
		log.Println("Starting Goji on", listener.Addr())

		graceful.HandleSignals()
		bind.Ready()

		graceful.PreHook(func() { log.Printf("Goji received signal, gracefully stopping") })
		graceful.PostHook(func() { log.Printf("Goji stopped") })

		err := graceful.Serve(listener, http.DefaultServeMux)

		if err != nil {
			log.Fatal(err)
		}

		graceful.Wait()

		return nil
	})

	os.Exit(l.Go(app))

}
