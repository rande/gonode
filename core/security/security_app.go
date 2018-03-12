// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"regexp"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
)

func ConfigureSecurity(l *goapp.Lifecycle, conf *config.Config) {
	l.Register(func(app *goapp.App) error {

		app.Set("security.authorizer", func(app *goapp.App) interface{} {
			voters := []Voter{}
			for _, id := range conf.Security.Voters {
				voters = append(voters, app.Get(id).(Voter))
			}

			return &DefaultAuthorizationChecker{
				DecisionVoter: &AffirmativeDecision{
					Voters: voters,
				},
			}
		})

		app.Set("security.voter.role", func(app *goapp.App) interface{} {
			return &RoleVoter{Prefix: "ROLE_"}
		})

		app.Set("security.voter.is", func(app *goapp.App) interface{} { // to remove
			return &RoleVoter{Prefix: "IS_"}
		})

		app.Set("security.access_checker", func(app *goapp.App) interface{} {
			access := make([]*AccessRule, 0)

			for _, c := range conf.Security.Access {

				attrs := make(Attributes, 0)

				for _, r := range c.Roles {
					attrs = append(attrs, r)
				}

				r := &AccessRule{
					Path:  regexp.MustCompile(c.Path),
					Roles: attrs,
				}

				access = append(access, r)
			}

			return &AccessChecker{
				Rules: access,
				DecisionVoter: &AffirmativeDecision{
					Voters: []Voter{
						&RoleVoter{},
					},
					AllowIfAllAbstainDecisions: false,
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
			Debug:              true,
		})

		mux.Use(c.Handler)

		access := app.Get("security.access_checker").(*AccessChecker)

		mux.Use(AccessCheckerMiddleware(access))

		return nil
	})
}

func ConfigureCors(l *goapp.Lifecycle, conf *config.Config) {

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
			Debug:              true,
		})

		mux.Use(c.Handler)

		return nil
	})
}
