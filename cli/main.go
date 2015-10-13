// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mitchellh/cli"
	"github.com/rande/gonode/commands"
	"log"
	"os"
)

func main() {

	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gonode", "0.0.1-DEV")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &commands.ServerCommand{
				Ui: ui,
			}, nil
		},
		"client": func() (cli.Command, error) {
			return &commands.ClientCommand{
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
