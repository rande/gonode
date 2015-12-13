// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Perform_Authentification_With_Not_Found_Credentials(t *testing.T) {

	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	a := &MockedAuthenticator{}
	a.On("getCredentials", r).Return(nil, nil)

	performed, err := performAuthentication(a, w, r)

	assert.False(t, performed)
	assert.Nil(t, err)
}

func Test_Perform_Authentification_With_Not_User_Found(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foobar", nil)
	w := httptest.NewRecorder()

	c := map[string]string{
		"login":    "thomas",
		"password": "password",
	}

	a := &MockedAuthenticator{}
	a.On("getCredentials", r).Return(c, nil)
	a.On("getUser", c).Return(nil, nil)

	performed, err := performAuthentication(a, w, r)

	assert.True(t, performed)
	assert.Nil(t, err)
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
	a.On("getCredentials", r).Return(c, nil)
	a.On("getUser", c).Return(u, nil)
	a.On("checkCredentials", c, u).Return(InvalidCredentials)
	a.On("onAuthenticationFailure", r, w, InvalidCredentials)

	performed, err := performAuthentication(a, w, r)

	assert.True(t, performed)
	assert.Equal(t, InvalidCredentials, err)
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
	a.On("getCredentials", r).Return(c, nil)
	a.On("getUser", c).Return(u, nil)
	a.On("checkCredentials", c, u).Return(nil)
	a.On("createAuthenticatedToken", u).Return(token, nil)
	a.On("onAuthenticationSuccess", r, w, token)

	performed, err := performAuthentication(a, w, r)

	assert.True(t, performed)
	assert.Nil(t, err)
}
