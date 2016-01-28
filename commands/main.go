// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"github.com/rande/gonode/commands/dev"
	"github.com/rande/gonode/commands/server"
	"os"
)

func main() {

	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gonode", "0.0.1-DEV")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &server.ServerCommand{
				Ui: ui,
			}, nil
		},
		"dev:service:list": func() (cli.Command, error) {
			return &dev.DevListServicesCommand{
				Ui: ui,
			}, nil
		},
		"dev:node:list": func() (cli.Command, error) {
			return &dev.DevListNodeTypesCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
