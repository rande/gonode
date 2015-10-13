// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package commands

import (
	"flag"
	"github.com/mitchellh/cli"
)

type ClientCommand struct {
	Ui         cli.Ui
	configFile string
}

func (c *ClientCommand) Help() string {
	return `Start a secure local client with embedded http server, this is usefull if you like to
use the vault feature to encrypt and decrypt data from remote storage.
	`
}

func (c *ClientCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("client", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.configFile, "client.toml.dist", "packages", "")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	return 0
}

func (c *ClientCommand) Synopsis() string {
	return "start a local proxy"
}
