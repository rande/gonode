// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dev

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rande/goapp"

	"github.com/rande/gonode/commands/server"
	"github.com/rande/gonode/modules/api"

	"fmt"
	"github.com/rande/gonode/modules/config"
)

type DevListServicesCommand struct {
	Ui         cli.Ui
	ConfigFile string
	Test       bool
	Verbose    bool
}

func (c *DevListServicesCommand) Help() string {
	return `List registered services`
}

func (c *DevListServicesCommand) Run(args []string) int {

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

	l := goapp.NewLifecycle()

	server.ConfigureServer(l, conf)
	api.ConfigureServer(l, conf)

	c.Ui.Info("Listing services available for the server configuration")

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {
		for _, k := range app.GetKeys() {
			c.Ui.Info(fmt.Sprintf(" > % -40s - %T - v := app.Get(\"%s\").(%T)", k, app.Get(k), k, app.Get(k)))
		}

		c.Ui.Info(fmt.Sprintf("Found %d services", len(app.GetKeys())))

		state.Out <- goapp.Control_Stop

		return nil
	})

	return l.Go(goapp.NewApp())
}

func (c *DevListServicesCommand) Synopsis() string {
	return "list registered services"
}
