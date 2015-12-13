// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"bytes"
	"encoding/json"
	"github.com/rande/gonode/core"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
	"testing"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

func GetToken() *jwt.Token {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	// @todo: add support for referenced token on database
	// token.Header["kid"] = "the sha1"

	// Set reserved claims
	jwtToken.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// Set shared claims
	jwtToken.Claims["rls"] = []string{"ADMIN"}
	jwtToken.Claims["usr"] = "thomas"

	return jwtToken
}

func Test_JwtGuardTokenAuthenticator_getCredentials_NoHeader_Request(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)

	c, err := a.getCredentials(req)

	assert.Nil(t, c)
	assert.Nil(t, err)
}

func Test_JwtGuardTokenAuthenticator_getCredentials_Invalid_Token(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	req.Header.Set("Authorization", "Bearer XXXX")

	c, err := a.getCredentials(req)

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func Test_JwtGuardTokenAuthenticator_getCredentials_Valid_Token(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	jwtToken := GetToken()
	tokenString, _ := jwtToken.SignedString(a.Key)

	req, _ := http.NewRequest("GET", "/ressource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	c, err := a.getCredentials(req)

	assert.NotNil(t, c)
	assert.Nil(t, err)
	assert.Equal(t, "thomas", c.(*jwt.Token).Claims["usr"].(string))
}

func Test_JwtGuardTokenAuthenticator_checkCredentials(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	// not used as the getCredentials check token validity
	c := GetToken()
	u := &DefaultGuardUser{Username: "thomas", Password: "dontcareaboutpassword"}

	err := a.checkCredentials(c, u)

	assert.Nil(t, err)
}

func Test_JwtGuardTokenAuthenticator_createAuthenticatedToken(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	u := &DefaultGuardUser{
		Username: "Thomas",
		Roles:    []string{"ADMIN"},
	}

	token, err := a.createAuthenticatedToken(u)

	assert.NotNil(t, token)
	assert.Nil(t, err)
	assert.Equal(t, token.GetUsername(), "Thomas")
	assert.Equal(t, token.GetRoles(), []string{"ADMIN"})
}

func Test_JwtGuardTokenAuthenticator_onAuthenticationSuccess(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	res := httptest.NewRecorder()
	token := &DefaultGuardToken{
		Username: "thomas",
		Roles:    []string{"ADMIN"},
	}

	a.onAuthenticationSuccess(req, res, token)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	assert.Equal(t, b.Len(), 0)
}

func Test_JwtGuardTokenAuthenticator_onAuthenticationFailure(t *testing.T) {
	a := &JwtGuardTokenAuthenticator{
		Path:        "/",
		NodeManager: &core.MockedManager{},
		Validity:    12,
		Key:         []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	res := httptest.NewRecorder()

	err := InvalidCredentials

	a.onAuthenticationFailure(req, res, err)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	v := &struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{}

	json.Unmarshal(b.Bytes(), v)

	assert.Equal(t, "KO", v.Status)
	assert.Equal(t, "Unable to validate token", v.Message)
}
