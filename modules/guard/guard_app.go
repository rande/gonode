// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/core/config"
	"github.com/zenazn/goji/web"
	"regexp"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.ServerConfig)
		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		auths := []GuardAuthenticator{
			&JwtTokenGuardAuthenticator{
				Path:        regexp.MustCompile(conf.Guard.Jwt.Token.Path),
				Key:         []byte(conf.Guard.Key),
				Validity:    conf.Guard.Jwt.Validity,
				NodeManager: manager,
			},
			&JwtLoginGuardAuthenticator{
				LoginPath:   conf.Guard.Jwt.Login.Path,
				Key:         []byte(conf.Guard.Key),
				Validity:    conf.Guard.Jwt.Validity,
				NodeManager: manager,
			},
		}

		mux.Use(GetGuardMiddleware(auths))

		return nil
	})
}
