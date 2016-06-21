// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// this authenticator will create a JWT Token from a standard form
type JwtTokenGuardAuthenticator struct {
	Path     *regexp.Regexp
	Manager  GuardManager
	Validity int64
	Key      []byte
	Logger   *log.Logger
}

func (a *JwtTokenGuardAuthenticator) GetCredentials(req *http.Request) (interface{}, error) {
	if !a.Path.Match([]byte(req.URL.Path)) {
		return nil, nil
	}

	if credentials, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
					"algo":   token.Header["alg"],
				}).Info("Invalid signing method")
			}

			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Key), nil
	}); err != nil {

		if err == request.ErrNoTokenInRequest {
			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
					"error":  err.Error(),
				}).Debug("No token in request, skipping")
			}

			return nil, nil
		} else {
			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
					"error":  err.Error(),
				}).Warn("Invalid credentials format")
			}

			return nil, InvalidCredentialsFormat
		}

	} else {

		claims := credentials.Claims.(jwt.MapClaims)
		if _, ok := claims["usr"]; !ok {

			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
				}).Info("invalid credentials, missing usr field")
			}

			return nil, InvalidCredentialsFormat
		}

		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module":   "core.guard.jwt_token",
				"username": claims["usr"],
				"roles":    claims["rls"],
			}).Info("valid credentials")
		}

		return credentials, nil
	}
}

func (a *JwtTokenGuardAuthenticator) GetUser(credentials interface{}) (GuardUser, error) {
	jwtToken := credentials.(*jwt.Token)

	claims := jwtToken.Claims.(jwt.MapClaims)

	user, err := a.Manager.GetUser(claims["usr"].(string))

	if err != nil {
		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module":   "core.guard.jwt_token",
				"error":    err.Error(),
				"username": claims["usr"].(string),
			}).Error("An error occurs when retrieving the user")
		}

		return user, err
	}

	if user != nil {
		return user, nil
	}

	if a.Logger != nil {
		a.Logger.WithFields(log.Fields{
			"module":   "core.guard.jwt_token",
			"username": claims["usr"].(string),
		}).Info("Unable to found the user")
	}

	return nil, UnableRetrieveUser
}

func (a *JwtTokenGuardAuthenticator) CheckCredentials(credentials interface{}, user GuardUser) error {
	// nothing to do ...

	return nil
}

func (a *JwtTokenGuardAuthenticator) CreateAuthenticatedToken(user GuardUser) (GuardToken, error) {
	return &DefaultGuardToken{
		Username: user.GetUsername(),
		Roles:    user.GetRoles(),
	}, nil
}

func (a *JwtTokenGuardAuthenticator) OnAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) bool {
	// nothing to do
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusForbidden)

	data, _ := json.Marshal(map[string]string{
		"status":  "KO",
		"message": "Unable to validate the token",
	})

	res.Write(data)

	return true
}

func (a *JwtTokenGuardAuthenticator) OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool {
	// nothing to do

	return false
}
