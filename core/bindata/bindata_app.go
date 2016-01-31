// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/assets"
	"github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Config(func(app *goapp.App) error {
		assets.UpdateRootDir(conf.BinData.BasePath)

		return nil
	})

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.pongo", func(app *goapp.App) interface{} {

			return pongo2.NewSet("gonode.bindata", &PongoTemplateLoader{
				Asset: assets.Asset,
				Paths: conf.BinData.Templates,
			})
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {

		if !app.Has("goji.mux") {
			return nil
		}

		mux := app.Get("goji.mux").(*web.Mux)
		logger := app.Get("logger").(*log.Logger)

		for _, bindata := range conf.BinData.Assets {
			ConfigureBinDataMux(mux, bindata.Public, bindata.Private, bindata.Index, logger)
		}

		return nil
	})
}
