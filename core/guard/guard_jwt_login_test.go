// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func Test_JwtLoginGuardAuthenticator_getCredentials_Valid_Request(t *testing.T) {
	a := &JwtLoginGuardAuthenticator{
		LoginPath: regexp.MustCompile("/login"),
		Manager:   &MockedManager{},
		Validity:  12,
		Key:       []byte("ZeKey"),
	}

	v := url.Values{
		"username": {"thomas"},
		"password": {"ZePassword"},
	}

	req, _ := http.NewRequest("POST", "/login", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c, err := a.GetCredentials(req)

	assert.NotNil(t, c)
	assert.Nil(t, err)

	cs := c.(*struct{ Username, Password string })

	assert.Equal(t, cs.Username, "thomas")
	assert.Equal(t, cs.Password, "ZePassword")
}

func Test_JwtLoginGuardAuthenticator_checkCredentials_Valid_Password(t *testing.T) {
	a := &JwtLoginGuardAuthenticator{
		LoginPath: regexp.MustCompile("/login"),
		Manager:   &MockedManager{},
		Validity:  12,
		Key:       []byte("ZeKey"),
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("ZePassword"), 1)

	c := &struct{ Username, Password string }{Username: "thomas", Password: "ZePassword"}
	u := &DefaultGuardUser{Username: "thomas", Password: string(password[:])}

	err := a.CheckCredentials(c, u)

	assert.Nil(t, err)
}

func Test_JwtLoginGuardAuthenticator_createAuthenticatedToken(t *testing.T) {
	a := &JwtLoginGuardAuthenticator{
		LoginPath: regexp.MustCompile("/login"),
		Manager:   &MockedManager{},
		Validity:  12,
		Key:       []byte("ZeKey"),
	}

	u := &DefaultGuardUser{
		Username: "Thomas",
		Password: "EncryptedPassword",
		Roles:    []string{"ADMIN"},
	}

	token, err := a.CreateAuthenticatedToken(u)

	assert.NotNil(t, token)
	assert.Nil(t, err)
	assert.Equal(t, token.GetUsername(), "Thomas")
	assert.Equal(t, token.GetRoles(), []string{"ADMIN"})
}

func Test_JwtLoginGuardAuthenticator_onAuthenticationSuccess(t *testing.T) {
	a := &JwtLoginGuardAuthenticator{
		LoginPath: regexp.MustCompile("/login"),
		Manager:   &MockedManager{},
		Validity:  12,
		Key:       []byte("ZeKey"),
	}

	req, _ := http.NewRequest("POST", "/login", nil)
	res := httptest.NewRecorder()
	token := &DefaultGuardToken{
		Username: "thomas",
		Roles:    []string{"ADMIN"},
	}

	a.OnAuthenticationSuccess(req, res, token)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	v := &struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{}

	json.Unmarshal(b.Bytes(), v)

	jwtToken, err := jwt.Parse(v.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.Key, nil
	})

	assert.Nil(t, err)
	claims := jwtToken.Claims.(jwt.MapClaims)

	assert.Equal(t, token.Username, claims["usr"])
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	assert.Equal(t, v.Token, res.Header().Get("X-Token"))
	assert.Equal(t, "Request is authenticated", v.Message)
	assert.Equal(t, "OK", v.Status)

	// @todo: I fail on basic golang conversion here ... from []interface{} to []string
	//assert.Equal(t, token.Roles, jwtToken.Claims["rls"].([]string))
}

func Test_JwtLoginGuardAuthenticator_onAuthenticationFailure(t *testing.T) {
	a := &JwtLoginGuardAuthenticator{
		LoginPath: regexp.MustCompile("/login"),
		Manager:   &MockedManager{},
		Validity:  12,
		Key:       []byte("ZeKey"),
	}

	req, _ := http.NewRequest("POST", "/login", nil)
	res := httptest.NewRecorder()

	err := InvalidCredentials

	a.OnAuthenticationFailure(req, res, err)
	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	v := &struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{}

	json.Unmarshal(b.Bytes(), v)

	assert.Equal(t, "KO", v.Status)
	assert.Equal(t, "Unable to authenticate request", v.Message)
	assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
}
