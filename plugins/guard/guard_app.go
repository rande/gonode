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
				Path:        regexp.MustCompile(conf.Auth.Jwt.Token.Path),
				Key:         []byte(conf.Auth.Key),
				Validity:    conf.Auth.Jwt.Validity,
				NodeManager: manager,
			},
			&JwtLoginGuardAuthenticator{
				LoginPath:   conf.Auth.Jwt.Login.Path,
				Key:         []byte(conf.Auth.Key),
				Validity:    conf.Auth.Jwt.Validity,
				NodeManager: manager,
			},
		}

		mux.Use(GetGuardMiddleware(auths))

		return nil
	})
}
