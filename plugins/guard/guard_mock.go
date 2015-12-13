// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockedAuthenticator struct {
	mock.Mock
}

func (m *MockedAuthenticator) getCredentials(req *http.Request) (interface{}, error) {
	args := m.Mock.Called(req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(interface{}), args.Error(1)
}
func (m *MockedAuthenticator) getUser(credentials interface{}) (GuardUser, error) {
	args := m.Mock.Called(credentials)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(GuardUser), args.Error(1)
}
func (m *MockedAuthenticator) checkCredentials(credentials interface{}, user GuardUser) error {
	args := m.Mock.Called(credentials, user)

	return args.Error(0)
}
func (m *MockedAuthenticator) createAuthenticatedToken(u GuardUser) (GuardToken, error) {
	args := m.Mock.Called(u)

	if args.Get(0) == nil {
		return nil, args.Error(0)
	}

	return args.Get(0).(GuardToken), args.Error(1)
}
func (m *MockedAuthenticator) onAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) {
	m.Mock.Called(req, res, err)
}

func (m *MockedAuthenticator) onAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) {
	m.Mock.Called(req, res, token)
}
