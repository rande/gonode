// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/plugins/user"
	"net/http"
)

// this authenticator will create a JWT Token from a standard form
type JwtGuardTokenAuthenticator struct {
	Path        string
	NodeManager core.NodeManager
	Validity    int64
	Key         []byte
}

func (a *JwtGuardTokenAuthenticator) getCredentials(req *http.Request) (interface{}, error) {
	// Authorization: Bearer <token>
	auth := req.Header.Get("Authorization")

	if auth == "" { // no header, no error
		return nil, nil
	}

	if len(auth) < 8 {
		return nil, InvalidCredentialsFormat
	}

	if credentials, err := jwt.Parse(auth[7:], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Key), nil
	}); err != nil {
		return nil, InvalidCredentialsFormat
	} else {
		return credentials, nil
	}
}

func (a *JwtGuardTokenAuthenticator) getUser(credentials interface{}) (GuardUser, error) {
	jwtToken := credentials.(*jwt.Token)

	query := a.NodeManager.
		SelectBuilder().
		Where("type = 'core.user' AND data->>'username' = ?", jwtToken.Claims["usr"].(string))

	if node := a.NodeManager.FindOneBy(query); node != nil {
		return node.Data.(*user.User), nil
	}

	return nil, UnableRetrieveUser
}

func (a *JwtGuardTokenAuthenticator) checkCredentials(credentials interface{}, user GuardUser) error {
	// nothing to do ...

	return nil
}

func (a *JwtGuardTokenAuthenticator) createAuthenticatedToken(user GuardUser) (GuardToken, error) {
	return &DefaultGuardToken{
		Username: user.GetUsername(),
		Roles:    user.GetRoles(),
	}, nil
}

func (a *JwtGuardTokenAuthenticator) onAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) {
	// nothing to do
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusForbidden)

	data, _ := json.Marshal(map[string]string{
		"status":  "KO",
		"message": "Unable to validate token",
	})

	res.Write(data)
}

func (a *JwtGuardTokenAuthenticator) onAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) {
	// nothing to do
}
