// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
)

func Test_Perform_Authentification_With_Not_Found_Credentials(t *testing.T) {

	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	a := &MockedAuthenticator{}
	a.On("GetCredentials", r).Return(nil, nil)

	cw := &web.C{Env: make(map[interface{}]interface{})}

	performed, output := performAuthentication(cw, a, w, r)

	assert.False(t, performed)
	assert.False(t, output)
}

func Test_Perform_Authentification_With_Not_User_Found(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	c := map[string]string{
		"login":    "thomas",
		"password": "password",
	}

	a := &MockedAuthenticator{}
	a.On("GetCredentials", r).Return(c, nil)
	a.On("GetUser", c).Return(nil, nil)
	a.On("OnAuthenticationFailure", r, w, nil).Return(true)

	cw := &web.C{Env: make(map[interface{}]interface{})}

	performed, output := performAuthentication(cw, a, w, r)

	assert.True(t, performed)
	assert.True(t, output)
}

func Test_Perform_Authentification_With_Invalid_Credentials(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	c := map[string]string{
		"login":    "thomas",
		"password": "password",
	}

	u := &DefaultGuardUser{
		Username: "thomas",
		Password: "password",
	}

	a := &MockedAuthenticator{}
	a.On("GetCredentials", r).Return(c, nil)
	a.On("GetUser", c).Return(u, nil)
	a.On("CheckCredentials", c, u).Return(InvalidCredentials)
	a.On("OnAuthenticationFailure", r, w, InvalidCredentials).Return(true)

	cw := &web.C{Env: make(map[interface{}]interface{})}

	performed, output := performAuthentication(cw, a, w, r)

	assert.True(t, performed)
	assert.True(t, output)
}

func Test_Perform_Authentification_With_Valid_User(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	c := map[string]string{
		"login":    "thomas",
		"password": "password",
	}

	u := &DefaultGuardUser{
		Username: "thomas",
		Password: "password",
	}

	token := &DefaultGuardToken{
		Username: "thomas",
	}

	a := &MockedAuthenticator{}
	a.On("GetCredentials", r).Return(c, nil)
	a.On("GetUser", c).Return(u, nil)
	a.On("CheckCredentials", c, u).Return(nil)
	a.On("CreateAuthenticatedToken", u).Return(token, nil)
	a.On("OnAuthenticationSuccess", r, w, token).Return(true)

	cw := &web.C{Env: make(map[interface{}]interface{})}

	performed, output := performAuthentication(cw, a, w, r)

	assert.True(t, performed)
	assert.True(t, output)
	assert.Equal(t, token, cw.Env["guard_token"])
}
