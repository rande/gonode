// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/rande/gonode/core/security"
	log "github.com/sirupsen/logrus"
)

// Extract token from cookie request.
type CookieExtractor struct {
	CookieName string
}

func (e *CookieExtractor) ExtractToken(req *http.Request) (string, error) {
	// Make sure form is parsed
	token, err := req.Cookie(e.CookieName)

	if err == http.ErrNoCookie {
		return "", request.ErrNoTokenInRequest
	}

	return token.Value, nil
}

// Extractor for OAuth2 access tokens.  Looks in 'Authorization'
// header then 'access_token' argument for a token.
var OAuth2Extractor = &request.MultiExtractor{
	&CookieExtractor{
		CookieName: "access_token",
	},
	request.AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

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
		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module": "core.guard.jwt_token",
				"path":   req.URL.Path,
				"regexp": a.Path,
			}).Info("Invalid path, skipping")
		}

		return nil, nil
	}

	credentials, err := request.ParseFromRequest(req, OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
					"algo":   token.Header["alg"],
				}).Info("Invalid signing method")
			}

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Key), nil
	})

	// no error, no credentials, nothing to do.
	if credentials == nil && err == nil {
		return nil, nil
	}

	if credentials != nil && err == nil {
		claims := credentials.Claims.(jwt.MapClaims)
		if _, ok := claims["usr"]; !ok {

			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
				}).Info("invalid credentials, missing usr field")
			}

			return nil, ErrInvalidCredentialsFormat
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

	// we have an error
	if err == request.ErrNoTokenInRequest {
		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module": "core.guard.jwt_token",
				"error":  err.Error(),
			}).Debug("No token in request, skipping")
		}

		return nil, nil
	}

	if e, ok := err.(*jwt.ValidationError); ok {
		if (e.Errors & jwt.ValidationErrorExpired) != 0 {
			if a.Logger != nil {
				a.Logger.WithFields(log.Fields{
					"module": "core.guard.jwt_token",
					"error":  err.Error(),
					"code":   e.Errors,
				}).Warn("Token expired, please refresh it")
			}

			return nil, ErrTokenExpired
		}
	}

	if a.Logger != nil {
		a.Logger.WithFields(log.Fields{
			"module": "core.guard.jwt_token",
			"error":  err.Error(),
		}).Warn("Invalid credentials format")
	}

	return nil, ErrInvalidCredentialsFormat
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

	return nil, ErrUnableRetrieveUser
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

	if req.Header.Get("Accept") == "application/json" {
		res.Header().Set("Content-Type", "application/json")

		res.WriteHeader(http.StatusForbidden)

		data, _ := json.Marshal(map[string]string{
			"status":  "KO",
			"message": "Unable to validate the token",
		})

		res.Write(data)

		return true
	}

	if err == ErrTokenExpired {
		security.RenderExpiredToken(res)
	} else {
		security.RenderForbidden(res)
	}

	return true
}

func (a *JwtTokenGuardAuthenticator) OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool {
	// nothing to do

	return false
}
