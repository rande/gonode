// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"net/http"

	"github.com/mitchellh/cli"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

type ServerCommand struct {
	Ui         cli.Ui
	ConfigFile string
	Test       bool
	Verbose    bool
	Configure  func(configFile string) *goapp.Lifecycle
}

func (c *ServerCommand) Help() string {
	return `Serve gonode server (better be behing a http reverse proxy)`
}

func (c *ServerCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.Ui.Output(c.Help())
	}

	cmdFlags.StringVar(&c.ConfigFile, "config", "server.toml.dist", "")
	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.BoolVar(&c.Test, "test", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	l := c.Configure(c.ConfigFile)

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.Config)
		logger := app.Get("logger").(*log.Logger)

		mux.Compile()

		// Install our handler at the root of the standard net/http default mux.
		// This allows packages like expvar to continue working as expected.
		http.Handle("/", mux)

		listener := bind.Socket(conf.Bind)
		logger.WithFields(log.Fields{
			"module": "command.cli",
		}).Debugf("Starting Goji on %s", listener.Addr())

		graceful.HandleSignals()
		bind.Ready()

		graceful.PreHook(func() {
			logger.WithFields(log.Fields{
				"module": "command.cli",
			}).Debug("Goji received signal, gracefully stopping")
		})

		graceful.PostHook(func() {
			logger.WithFields(log.Fields{
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
