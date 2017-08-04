// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"errors"
	"net/http"

	"github.com/zenazn/goji/web"
)

var (
	AccessForbiddenError = errors.New("Access Forbidden")
)

// Bare interface to used inside a request lifecycle
type SecurityToken interface {
	// return the current username for the current token
	GetUsername() string

	// return the related roles linked to the current token
	GetRoles() []string
}

// Default implementation to the GuardToken
type DefaultSecurityToken struct {
	Username string
	Roles    []string
}

func (t *DefaultSecurityToken) GetUsername() string {
	return t.Username
}

func (t *DefaultSecurityToken) GetRoles() []string {
	return t.Roles
}

func GetTokenFromContext(c web.C) SecurityToken {
	if _, ok := c.Env["guard_token"]; !ok { // no token
		return nil
	}

	return c.Env["guard_token"].(SecurityToken)
}

func CheckAccess(token SecurityToken, attrs Attributes, res http.ResponseWriter, req *http.Request, auth AuthorizationChecker) error {
	if token == nil { // no token
		return AccessForbiddenError
	}

	r, err := auth.IsGranted(token, attrs, req)

	if err != nil {
		return err
	}

	if r == false {
		return AccessForbiddenError
	}

	return nil
}
