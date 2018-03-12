// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
)

// this authenticator will create a JWT Token from a standard form
type JwtLoginGuardAuthenticator struct {
	LoginPath *regexp.Regexp
	Manager   GuardManager
	Validity  int64
	Key       []byte
	Logger    *log.Logger
}

func (a *JwtLoginGuardAuthenticator) GetCredentials(req *http.Request) (interface{}, error) {
	if req.Method != "POST" {
		a.Logger.WithFields(log.Fields{
			"module": "core.guard.jwt_login",
			"method": req.Method,
		}).Debug("Invalid HTTP method")

		return nil, nil
	}

	if !a.LoginPath.Match([]byte(req.URL.Path)) {
		a.Logger.WithFields(log.Fields{
			"module": "core.guard.jwt_login",
			"path":   req.URL.Path,
		}).Debug("Invalid Path")

		return nil, nil
	}

	loginForm := &struct {
		Username string `schema:"username" json:"username"`
		Password string `schema:"password" json:"password"`
	}{}

	if req.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(loginForm); err != nil {
			a.Logger.WithFields(log.Fields{
				"module": "core.guard.jwt_login",
				"error":  err,
				"method": "application/json",
			}).Info("Unable to decode JSON data")

			return nil, InvalidCredentialsFormat
		}
	} else {
		req.ParseForm()

		decoder := schema.NewDecoder()
		if err := decoder.Decode(loginForm, req.Form); err != nil {
			a.Logger.WithFields(log.Fields{
				"module": "core.guard.jwt_login",
				"error":  err,
				"method": "form-data",
			}).Info("Unable to decode POST parameters")

			return nil, InvalidCredentialsFormat
		}
	}

	if a.Logger != nil {
		a.Logger.WithFields(log.Fields{
			"module":   "core.guard.jwt_login",
			"username": loginForm.Username,
		}).Info("Starting authentification process")
	}

	return &struct{ Username, Password string }{loginForm.Username, loginForm.Password}, nil
}

func (a *JwtLoginGuardAuthenticator) GetUser(credentials interface{}) (GuardUser, error) {
	c := credentials.(*struct{ Username, Password string })

	user, err := a.Manager.GetUser(c.Username)

	if err != nil {
		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module":   "core.guard.jwt_login",
				"error":    err.Error(),
				"username": c.Username,
			}).Error("An error occurs when retrieving the user")
		}

		return user, err
	}

	if user != nil {
		return user, nil
	}

	if a.Logger != nil {
		a.Logger.WithFields(log.Fields{
			"module":   "core.guard.jwt_login",
			"username": c.Username,
		}).Info("Unable to found the user")
	}

	return nil, UnableRetrieveUser
}

func (a *JwtLoginGuardAuthenticator) CheckCredentials(credentials interface{}, user GuardUser) error {
	c := credentials.(*struct{ Username, Password string })

	if err := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(c.Password)); err != nil { // equal

		if a.Logger != nil {
			a.Logger.WithFields(log.Fields{
				"module":   "core.guard.jwt_login",
				"username": c.Username,
			}).Info("Invalid credentials")
		}

		return InvalidCredentials
	}

	if a.Logger != nil {
		a.Logger.WithFields(log.Fields{
			"module":   "core.guard.jwt_login",
			"username": c.Username,
		}).Info("Valid credentials")
	}

	return nil
}

func (a *JwtLoginGuardAuthenticator) CreateAuthenticatedToken(user GuardUser) (GuardToken, error) {
	return &DefaultGuardToken{
		Username: user.GetUsername(),
		Roles:    user.GetRoles(),
	}, nil
}

func (a *JwtLoginGuardAuthenticator) OnAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) bool {
	// nothing to do
	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusForbidden)

	data, _ := json.Marshal(map[string]string{
		"status":  "KO",
		"message": "Unable to authenticate request",
	})

	res.Write(data)

	return true
}

func (a *JwtLoginGuardAuthenticator) OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	// @todo: add support for referenced token on database
	// token.Header["kid"] = "the sha1"

	// Set reserved claims
	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// Set shared claims
	claims["rls"] = token.GetRoles()
	claims["usr"] = token.GetUsername()

	// Sign and get the complete encoded token as a string
	tokenString, _ := jwtToken.SignedString([]byte(a.Key))

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("X-Token", tokenString)

	data, _ := json.Marshal(map[string]string{
		"status":  "OK",
		"message": "Request is authenticated",
		"token":   tokenString,
	})

	res.Write(data)

	return true
}
