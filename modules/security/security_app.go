// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"log"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {
	l.Prepare(func(app *goapp.App) error {
		conf := app.Get("gonode.configuration").(*config.ServerConfig)

		if conf.Security == nil {
			return nil // nothing setup
		}

		mux := app.Get("goji.mux").(*web.Mux)

		c := cors.New(cors.Options{
			AllowedOrigins:     conf.Security.Cors.AllowedOrigins,
			AllowedHeaders:     conf.Security.Cors.AllowedHeaders,
			AllowedMethods:     conf.Security.Cors.AllowedMethods,
			ExposedHeaders:     conf.Security.Cors.ExposedHeaders,
			AllowCredentials:   conf.Security.Cors.AllowCredentials,
			MaxAge:             conf.Security.Cors.MaxAge,
			OptionsPassthrough: conf.Security.Cors.OptionsPassthrough,
		})

		if app.Has("logger") {
			log := app.Get("logger").(*log.Logger)

			c.Log = log
		}

		mux.Use(c.Handler)

		return nil
	})
}
