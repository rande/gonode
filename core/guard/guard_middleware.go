// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

func GetGuardMiddleware(auths []GuardAuthenticator) func(c *web.C, h http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		var logger *log.Entry

		fn := func(w http.ResponseWriter, r *http.Request) {

			if _, ok := c.Env["logger"]; ok {
				logger = c.Env["logger"].(*log.Entry)
			}

			// handle security here
			for _, authenticator := range auths {
				if logger != nil {
					logger.WithFields(log.Fields{
						"module": "core.guard.middleware",
						"type":   fmt.Sprintf("%T", authenticator),
					}).Debug("Starting authentificator process")
				}

				performed, output := performAuthentication(c, authenticator, w, r)

				if performed && output {
					if logger != nil {
						logger.WithFields(log.Fields{
							"module": "core.guard.middleware",
							"type":   fmt.Sprintf("%T", authenticator),
						}).Debug("Authentification performed and output sent, skipping next middleware")
					}

					return
				} else if performed {
					if logger != nil {
						logger.WithFields(log.Fields{
							"module": "core.guard.middleware",
							"type":   fmt.Sprintf("%T", authenticator),
						}).Debug("Authentification performed, start next middleware")
					}

					break
				}

				if logger != nil {
					logger.WithFields(log.Fields{
						"module":    "core.guard.middleware",
						"type":      fmt.Sprintf("%T", authenticator),
						"performed": performed,
						"output":    output,
					}).Debug("Ignoring authenticator")
				}
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// false means, no authentification has been done
func performAuthentication(c *web.C, a GuardAuthenticator, w http.ResponseWriter, r *http.Request) (bool, bool) {
	var o bool

	// get credentials from request
	credentials, err := a.GetCredentials(r)

	if err == InvalidCredentialsFormat {
		o = a.OnAuthenticationFailure(r, w, err)

		return true, o
	}

	// no credentials, return
	if credentials == nil { // nothing to do, next one
		return false, false
	}

	// ok get the current user for the current credentials
	user, err := a.GetUser(credentials)

	if err != nil || user == nil {
		o = a.OnAuthenticationFailure(r, w, err)

		return true, o
	}

	// check if the request's credentials match user credentials
	if err = a.CheckCredentials(credentials, user); err != nil {
		o = a.OnAuthenticationFailure(r, w, err)

		return true, o
	}

	// create a valid security token for the user
	token, err := a.CreateAuthenticatedToken(user)

	c.Env["guard_token"] = token

	// complete the process
	o = a.OnAuthenticationSuccess(r, w, token)

	return true, o
}
