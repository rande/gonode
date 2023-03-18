// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/rande/goapp"
	"github.com/rande/gonode/app/assets"
	"github.com/rande/gonode/core/bindata"
	"github.com/rande/gonode/core/commands"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/logger"
	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/api"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/modules/debug"
	"github.com/rande/gonode/modules/feed"
	node_guard "github.com/rande/gonode/modules/guard"
	"github.com/rande/gonode/modules/media"
	"github.com/rande/gonode/modules/prism"
	"github.com/rande/gonode/modules/raw"
	"github.com/rande/gonode/modules/search"
	"github.com/rande/gonode/modules/setup"
	"github.com/rande/gonode/modules/user"
	log "github.com/sirupsen/logrus"
)

func Configure(configFile string) *goapp.Lifecycle {
	l := goapp.NewLifecycle()

	conf := config.NewConfig()
	config.LoadConfigurationFromFile(configFile, conf)

	l.Config(func(app *goapp.App) error {
		assets.UpdateRootDir(conf.BinData.BasePath)

		app.Set("gonode.asset", func(app *goapp.App) interface{} {
			return assets.Asset
		})

		return nil
	})

	base.Configure(l, conf)
	debug.Configure(l, conf)
	user.Configure(l, conf)
	raw.Configure(l, conf)
	blog.Configure(l, conf)
	media.Configure(l, conf)
	search.Configure(l, conf)
	feed.Configure(l, conf)

	logger.Configure(l, conf)
	commands.Configure(l, conf)
	security.ConfigureCors(l, conf)
	node_guard.Configure(l, conf)
	security.ConfigureSecurity(l, conf)
	api.Configure(l, conf)
	setup.Configure(l, conf)
	bindata.Configure(l, conf)
	prism.Configure(l, conf)
	router.Configure(l, conf)

	return l
}

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	c := cli.NewCLI("gonode", "0.0.2-DEV")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &commands.ServerCommand{
				Ui:        ui,
				Configure: Configure,
			}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
