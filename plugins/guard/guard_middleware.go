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
				if performed, err := performAuthentication(authenticator, w, r); performed && err != nil {
					// continue ???
					// an issue occured while the authenticator can handle the authentication
					// move to the next one for now ...

					continue
				} else if performed && err == nil {
					// nothing to do, auth end (login, or content sent ...)
					return
				} else {
					// nothing to do move to the next authenticator
					continue
				}
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// false means, no authentification has been done
func performAuthentication(a GuardAuthenticator, w http.ResponseWriter, r *http.Request) (bool, error) {
	// get credentials from request
	credentials, _ := a.getCredentials(r)

	// no credentials, return
	if credentials == nil { // nothing to do, next one
		return false, nil
	}

	// ok get the current user for the current credentials
	user, err := a.getUser(credentials)

	if err != nil || user == nil {
		return true, err
	}

	// check if the request's credentials match user credentials
	if err = a.checkCredentials(credentials, user); err != nil {
		a.onAuthenticationFailure(r, w, err)

		return true, err
	}

	// create a valid security token for the user
	token, err := a.createAuthenticatedToken(user)

	// complete the process
	a.onAuthenticationSuccess(r, w, token)

	return true, nil
}
