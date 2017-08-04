// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"net/http"
)

type AnonymousAuthenticator struct {
	DefaultRoles []string
}

func (a *AnonymousAuthenticator) GetCredentials(req *http.Request) (interface{}, error) {
	return &struct{ Username string }{"anonymous"}, nil
}

func (a *AnonymousAuthenticator) GetUser(credentials interface{}) (GuardUser, error) {
	c := credentials.(*struct{ Username string })

	u := &DefaultGuardUser{
		Username: c.Username,
		Roles:    a.DefaultRoles,
	}

	return u, nil
}

func (a *AnonymousAuthenticator) CheckCredentials(credentials interface{}, user GuardUser) error {
	return nil
}

func (a *AnonymousAuthenticator) CreateAuthenticatedToken(user GuardUser) (GuardToken, error) {
	return &DefaultGuardToken{
		Username: user.GetUsername(),
		Roles:    user.GetRoles(),
	}, nil
}

func (a *AnonymousAuthenticator) OnAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) bool {
	return false
}

func (a *AnonymousAuthenticator) OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool {
	return false
}
