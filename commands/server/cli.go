// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rande/goapp"

	"net/http"

	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/plugins/api"
	"github.com/rande/gonode/plugins/bindata"
	"github.com/rande/gonode/plugins/guard"
	"github.com/rande/gonode/plugins/search"
	"github.com/rande/gonode/plugins/security"
	"github.com/rande/gonode/plugins/setup"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"log"
)

type ServerCommand struct {
	Ui         cli.Ui
	ConfigFile string
	Test       bool
	Verbose    bool
}

func (c *ServerCommand) Help() string {
	return `Serve gonode server (better be behing a http reverse proxy)`
}

func (c *ServerCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.ConfigFile, "config", "server.toml.dist", "")
	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.BoolVar(&c.Test, "test", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	conf := config.NewServerConfig()

	config.LoadConfigurationFromFile(c.ConfigFile, conf)

	c.Ui.Info("Starting GoNode Server on: " + conf.Bind)

	l := goapp.NewLifecycle()

	ConfigureServer(l, conf)

	// add plugins
	setup.ConfigureServer(l, conf)
	security.ConfigureServer(l, conf)
	search.ConfigureServer(l, conf)
	api.ConfigureServer(l, conf)
	guard.ConfigureServer(l, conf)
	bindata.ConfigureServer(l, conf)

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.ServerConfig)

		mux.Compile()

		// Install our handler at the root of the standard net/http default mux.
		// This allows packages like expvar to continue working as expected.
		http.Handle("/", mux)

		listener := bind.Socket(conf.Bind)
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

	return l.Go(goapp.NewApp())
}

func (c *ServerCommand) Synopsis() string {
	return "server local command"
}
