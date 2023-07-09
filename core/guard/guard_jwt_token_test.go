// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
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
	"regexp"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func GetToken() *jwt.Token {
	jwtToken := jwt.New(jwt.SigningMethodHS256)

	// @todo: add support for referenced token on database
	// token.Header["kid"] = "the sha1"

	claims := jwtToken.Claims.(jwt.MapClaims)

	// Set reserved claims
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// Set shared claims
	claims["rls"] = []string{"ADMIN"}
	claims["usr"] = "thomas"

	return jwtToken
}

func Test_JwtTokenGuardAuthenticator_getCredentials_NoMatch(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/api/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)

	c, err := a.GetCredentials(req)

	assert.Nil(t, c)
	assert.Nil(t, err)
}

func Test_JwtTokenGuardAuthenticator_getCredentials_NoHeader_Request(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)

	c, err := a.GetCredentials(req)

	assert.Nil(t, c)
	assert.Nil(t, err)
}

func Test_JwtTokenGuardAuthenticator_getCredentials_Invalid_Token(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	req.Header.Set("Authorization", "Bearer XXXX")

	c, err := a.GetCredentials(req)

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func Test_JwtTokenGuardAuthenticator_getCredentials_Valid_Token_Header(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	jwtToken := GetToken()
	tokenString, _ := jwtToken.SignedString(a.Key)

	req, _ := http.NewRequest("GET", "/ressource", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	c, err := a.GetCredentials(req)

	assert.NotNil(t, c)
	assert.Nil(t, err)
	assert.Equal(t, "thomas", c.(*jwt.Token).Claims.(jwt.MapClaims)["usr"].(string))
}

func Test_JwtTokenGuardAuthenticator_getCredentials_Valid_Token_QueryString(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	jwtToken := GetToken()
	tokenString, _ := jwtToken.SignedString(a.Key)

	req, _ := http.NewRequest("GET", "/ressource?access_token="+tokenString, nil)

	c, err := a.GetCredentials(req)

	assert.NotNil(t, c)
	assert.Nil(t, err)
	assert.Equal(t, "thomas", c.(*jwt.Token).Claims.(jwt.MapClaims)["usr"].(string))
}

func Test_JwtTokenGuardAuthenticator_checkCredentials(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	// not used as the getCredentials check token validity
	c := GetToken()
	u := &DefaultGuardUser{Username: "thomas", Password: "dontcareaboutpassword"}

	err := a.CheckCredentials(c, u)

	assert.Nil(t, err)
}

func Test_JwtTokenGuardAuthenticator_createAuthenticatedToken(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	u := &DefaultGuardUser{
		Username: "Thomas",
		Roles:    []string{"ADMIN"},
	}

	token, err := a.CreateAuthenticatedToken(u)

	assert.NotNil(t, token)
	assert.Nil(t, err)
	assert.Equal(t, token.GetUsername(), "Thomas")
	assert.Equal(t, token.GetRoles(), []string{"ADMIN"})
}

func Test_JwtTokenGuardAuthenticator_onAuthenticationSuccess(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	res := httptest.NewRecorder()
	token := &DefaultGuardToken{
		Username: "thomas",
		Roles:    []string{"ADMIN"},
	}

	a.OnAuthenticationSuccess(req, res, token)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	assert.Equal(t, b.Len(), 0)
}

func Test_JwtTokenGuardAuthenticator_onAuthenticationFailure(t *testing.T) {
	a := &JwtTokenGuardAuthenticator{
		Path:     regexp.MustCompile("/*"),
		Manager:  &MockedManager{},
		Validity: 12,
		Key:      []byte("ZeKey"),
	}

	req, _ := http.NewRequest("GET", "/ressource", nil)
	res := httptest.NewRecorder()

	err := ErrInvalidCredentials

	a.OnAuthenticationFailure(req, res, err)

	b := bytes.NewBuffer([]byte(""))
	io.Copy(b, res.Body)

	v := &struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{}

	json.Unmarshal(b.Bytes(), v)

	assert.Equal(t, "KO", v.Status)
	assert.Equal(t, "Unable to validate the token", v.Message)
}
