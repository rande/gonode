// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

func GetGuardMiddleware(auths []GuardAuthenticator) func(c *web.C, h http.Handler) http.Handler {
	return func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// handle security here
			for _, authenticator := range auths {
				performed, output := performAuthentication(c, authenticator, w, r)

				if performed && output {
					return
				} else if performed {
					break
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
	credentials, err := a.getCredentials(r)

	if err == InvalidCredentialsFormat {
		o = a.onAuthenticationFailure(r, w, err)

		return true, o
	}

	// no credentials, return
	if credentials == nil { // nothing to do, next one
		return false, false
	}

	// ok get the current user for the current credentials
	user, err := a.getUser(credentials)

	if err != nil || user == nil {
		o = a.onAuthenticationFailure(r, w, err)

		return true, o
	}

	// check if the request's credentials match user credentials
	if err = a.checkCredentials(credentials, user); err != nil {
		o = a.onAuthenticationFailure(r, w, err)

		return true, o
	}

	// create a valid security token for the user
	token, err := a.createAuthenticatedToken(user)

	c.Env["guard_token"] = token

	// complete the process
	o = a.onAuthenticationSuccess(r, w, token)

	return true, o
}
