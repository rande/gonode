// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidCredentialsFormat   = errors.New("invalid credentials format")
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrUnableRetrieveUser         = errors.New("unable to retrieve the user")
	ErrCredentialMismatch         = errors.New("credential mismatch")
	ErrAuthenticatedTokenCreation = errors.New("unable to create authentication token")
	ErrTokenExpired               = errors.New("token expired")
)

// Bare interface with the default requirement to check username and password
type GuardUser interface {
	GetUsername() string
	GetPassword() string
	GetRoles() []string
}

type DefaultGuardUser struct {
	Username string
	Password string
	Roles    []string
}

func (u *DefaultGuardUser) GetUsername() string {
	return u.Username
}

func (u *DefaultGuardUser) GetPassword() string {
	return u.Password
}

func (u *DefaultGuardUser) GetRoles() []string {
	return u.Roles
}

// Bare interface to used inside a request lifecycle
type GuardToken interface {
	// return the current username for the current token
	GetUsername() string

	// return the related roles linked to the current token
	GetRoles() []string
}

// Default implementation to the GuardToken
type DefaultGuardToken struct {
	Username string
	Roles    []string
}

func (t *DefaultGuardToken) GetUsername() string {
	return t.Username
}

func (t *DefaultGuardToken) GetRoles() []string {
	return t.Roles
}

type GuardAuthenticator interface {
	// This method is call on each request.
	// If the method return nil as interface{} value, it means the authenticator
	// cannot handle the request
	GetCredentials(req *http.Request) (interface{}, error)

	// Return the user from the credentials
	GetUser(credentials interface{}) (GuardUser, error)

	// Check if the provided credentials are valid for the current user
	CheckCredentials(credentials interface{}, user GuardUser) error

	// Return a security token related to the user
	CreateAuthenticatedToken(u GuardUser) (GuardToken, error)

	// Action when the authentication fail.
	// On a default form login, it can be used to redirect the user to login page
	// return true if the workflows must be stopped (ie, the authenticator was written
	// bytes on the response. false if not.
	OnAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) bool

	// Action when the authentication success
	// On a default form login, it can be used to redirect the user to protected page
	// or the homepage
	// return true if the workflows must be stopped (ie, the authenticator was written
	// bytes on the response. false if not.
	OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool
}

type GuardManager interface {
	GetUser(username string) (GuardUser, error)
}
