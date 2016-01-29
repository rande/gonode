// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/guard"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/user"
	"github.com/zenazn/goji/web"
	"regexp"
)

type GuardManager struct {
	m *base.PgNodeManager
}

func (g *GuardManager) GetUser(username string) (guard.GuardUser, error) {
	query := g.m.SelectBuilder(base.NewSelectOptions()).
		Where("type = 'core.user' AND data->>'username' = ?", username)

	if node := g.m.FindOneBy(query); node != nil {
		return node.Data.(*user.User), nil
	}

	return nil, nil
}

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.ServerConfig)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		auths := []guard.GuardAuthenticator{
			&guard.JwtTokenGuardAuthenticator{
				Path:     regexp.MustCompile(conf.Guard.Jwt.Token.Path),
				Key:      []byte(conf.Guard.Key),
				Validity: conf.Guard.Jwt.Validity,
				Manager:  &GuardManager{manager},
			},
			&guard.JwtLoginGuardAuthenticator{
				LoginPath: conf.Guard.Jwt.Login.Path,
				Key:       []byte(conf.Guard.Key),
				Validity:  conf.Guard.Jwt.Validity,
				Manager:   &GuardManager{manager},
			},
		}

		mux.Use(guard.GetGuardMiddleware(auths))

		return nil
	})
}
