// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"errors"
	"net/http"
)

var (
	InvalidCredentialsFormat        = errors.New("Invalid credentials format")
	InvalidCredentials              = errors.New("Invalid credentials")
	UnableRetrieveUser              = errors.New("Unable to retrieve the user")
	CredentialMismatch              = errors.New("Credential mismatch")
	AuthenticatedTokenCreationError = errors.New("Unable to create authentication token")
)

// Bare interface with the default requirement to check username and password
type GuardUser interface {
	GetUsername() string
	GetPassword() string
	GetRoles() []string
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
	getCredentials(req *http.Request) (interface{}, error)

	// Return the user from the credentials
	getUser(credentials interface{}) (GuardUser, error)

	// Check if the provided credentials are valid for the current user
	checkCredentials(credentials interface{}, user GuardUser) error

	// Return a security token related to the user
	createAuthenticatedToken(u GuardUser) (GuardToken, error)

	// Action when the authentication fail.
	// On a default form login, it can be used to redirect the user to login page
	onAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error)

	// Action when the authentication success
	// On a default form login, it can be used to redirect the user to protected page
	// or the homepage
	onAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken)
}
