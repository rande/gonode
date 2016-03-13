// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Register(func(app *goapp.App) error {

		app.Set("security.authorizer", func(app *goapp.App) interface{} {
			return DefaultAuthorizationChecker{
				DecisionManager: &AffirmativeDecision{
					Voters: []Voter{
						&RoleVoter{Prefix: "ROLE_"},
					},
				},
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		conf := app.Get("gonode.configuration").(*config.Config)

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

		mux.Use(c.Handler)

		return nil
	})
}
