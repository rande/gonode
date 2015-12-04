// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dev

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rande/goapp"

	"github.com/rande/gonode/core"
	"github.com/rande/gonode/commands/server"

	"fmt"
)

type DevListNodeTypesCommand struct {
	Ui         cli.Ui
	ConfigFile string
	Test       bool
	Verbose    bool
}

func (c *DevListNodeTypesCommand) Help() string {
	return `List registered node types`
}

func (c *DevListNodeTypesCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.ConfigFile, "config", "server.toml.dist", "")
	cmdFlags.BoolVar(&c.Verbose, "verbose", false, "")
	cmdFlags.BoolVar(&c.Test, "test", false, "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	config := server.NewServerConfig()

	core.LoadConfiguration(c.ConfigFile, config)

	l := goapp.NewLifecycle()

	server.ConfigureServer(l, config)
	server.ConfigureHttpApi(l)

	c.Ui.Info("Node types available")

	l.Run(func(app *goapp.App, state *goapp.GoroutineState) error {

		handlers := app.Get("gonode.handler_collection").(core.HandlerCollection)

		for _, k := range handlers.GetKeys() {
			c.Ui.Info(fmt.Sprintf(" > % -40s - %T", k, handlers.GetByCode(k)))
		}

		c.Ui.Info(fmt.Sprintf("Found %d node types", len(app.GetKeys())))

		state.Out <- goapp.Control_Stop

		return nil
	})

	return l.Go(goapp.NewApp())
}

func (c *DevListNodeTypesCommand) Synopsis() string {
	return "list registered node types"
}
