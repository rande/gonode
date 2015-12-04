// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rande/goapp"

	"github.com/rande/gonode/core"

	"net/http"

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

	config := NewServerConfig()

	core.LoadConfiguration(c.ConfigFile, config)

	c.Ui.Info("Starting GoNode Server on: " + config.Bind)

	l := goapp.NewLifecycle()

	ConfigureServer(l, config)
	ConfigureHttpApi(l)

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		mux := app.Get("goji.mux").(*web.Mux)
		config := app.Get("gonode.configuration").(*ServerConfig)

		mux.Compile()

		// Install our handler at the root of the standard net/http default mux.
		// This allows packages like expvar to continue working as expected.
		http.Handle("/", mux)

		listener := bind.Socket(config.Bind)
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
