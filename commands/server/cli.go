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

	log "github.com/Sirupsen/logrus"
	"github.com/rande/gonode/core/bindata"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/logger"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/api"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/guard"
	"github.com/rande/gonode/modules/prism"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/setup"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
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

	// add modules
	logger.ConfigureServer(l, conf)
	setup.ConfigureServer(l, conf)
	security.ConfigureServer(l, conf)
	search.ConfigureServer(l, conf)
	api.ConfigureServer(l, conf)
	node_guard.ConfigureServer(l, conf)
	prism.ConfigureServer(l, conf)
	router.ConfigureServer(l, conf)
	base.ConfigureServer(l, conf)

	// must be last for now
	bindata.ConfigureServer(l, conf)

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.ServerConfig)

		mux.Compile()

		// Install our handler at the root of the standard net/http default mux.
		// This allows packages like expvar to continue working as expected.
		http.Handle("/", mux)

		listener := bind.Socket(conf.Bind)
		log.WithFields(log.Fields{
			"module": "command.cli",
		}).Debug("Starting Goji on %s", listener.Addr())

		graceful.HandleSignals()
		bind.Ready()

		graceful.PreHook(func() {
			log.WithFields(log.Fields{
				"module": "command.cli",
			}).Debug("Goji received signal, gracefully stopping")
		})
		graceful.PostHook(func() {
			log.WithFields(log.Fields{
				"module": "command.cli",
			}).Debug("Goji stopped")
		})

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
