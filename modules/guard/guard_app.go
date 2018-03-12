// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/guard"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/user"
	"github.com/zenazn/goji/web"
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

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		mux := app.Get("goji.mux").(*web.Mux)
		conf := app.Get("gonode.configuration").(*config.Config)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		logger := app.Get("logger").(*log.Logger)

		auths := []guard.GuardAuthenticator{
			&guard.JwtLoginGuardAuthenticator{
				LoginPath: regexp.MustCompile(conf.Guard.Jwt.Login.Path),
				Key:       []byte(conf.Guard.Key),
				Validity:  conf.Guard.Jwt.Validity,
				Manager:   &GuardManager{manager},
				Logger:    logger,
			},
			&guard.JwtTokenGuardAuthenticator{
				Path:     regexp.MustCompile(conf.Guard.Jwt.Token.Path),
				Key:      []byte(conf.Guard.Key),
				Validity: conf.Guard.Jwt.Validity,
				Manager:  &GuardManager{manager},
				Logger:   logger,
			},
			&guard.AnonymousAuthenticator{
				DefaultRoles: conf.Guard.Anonymous.Roles,
			},
		}

		mux.Use(guard.GetGuardMiddleware(auths))

		return nil
	})
}
