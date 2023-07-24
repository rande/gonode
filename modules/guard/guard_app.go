// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"regexp"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/guard"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/user"
	log "github.com/sirupsen/logrus"
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

		ignore := []*regexp.Regexp{}

		for _, path := range conf.Guard.Jwt.Token.Ignore {
			ignore = append(ignore, regexp.MustCompile(path))
		}

		auths := []guard.GuardAuthenticator{
			&guard.JwtLoginGuardAuthenticator{
				EndPoint: regexp.MustCompile(conf.Guard.Jwt.Login.EndPoint),
				Key:      []byte(conf.Guard.Key),
				Validity: conf.Guard.Jwt.Validity,
				Manager:  &GuardManager{manager},
				Logger:   logger,
			},
			&guard.JwtTokenGuardAuthenticator{
				Apply:     regexp.MustCompile(conf.Guard.Jwt.Token.Apply),
				Ignore:    ignore,
				Key:       []byte(conf.Guard.Key),
				Validity:  conf.Guard.Jwt.Validity,
				Manager:   &GuardManager{manager},
				Logger:    logger,
				LoginPage: conf.Guard.Jwt.Login.Page,
			},
			&guard.AnonymousAuthenticator{
				DefaultRoles: conf.Guard.Anonymous.Roles,
			},
		}

		mux.Use(guard.GetGuardMiddleware(auths))

		return nil
	})
}
