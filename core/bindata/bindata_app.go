// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Register(func(app *goapp.App) error {
		app.Set("gonode.pongo", func(app *goapp.App) interface{} {

			return pongo2.NewSet("gonode.bindata", &PongoTemplateLoader{
				Asset: app.Get("gonode.asset").(func(name string) ([]byte, error)),
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
		asset := app.Get("gonode.asset").(func(name string) ([]byte, error))

		for _, bindata := range conf.BinData.Assets {
			ConfigureBinDataMux(mux, asset, bindata.Public, bindata.Private, bindata.Index, logger)
		}

		return nil
	})
}
