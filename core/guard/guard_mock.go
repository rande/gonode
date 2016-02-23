// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

func (m *MockedAuthenticator) GetCredentials(req *http.Request) (interface{}, error) {
	args := m.Mock.Called(req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(interface{}), args.Error(1)
}

func (m *MockedAuthenticator) GetUser(credentials interface{}) (GuardUser, error) {
	args := m.Mock.Called(credentials)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(GuardUser), args.Error(1)
}

func (m *MockedAuthenticator) CheckCredentials(credentials interface{}, user GuardUser) error {
	args := m.Mock.Called(credentials, user)

	return args.Error(0)
}

func (m *MockedAuthenticator) CreateAuthenticatedToken(u GuardUser) (GuardToken, error) {
	args := m.Mock.Called(u)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(GuardToken), args.Error(1)
}

func (m *MockedAuthenticator) OnAuthenticationFailure(req *http.Request, res http.ResponseWriter, err error) bool {
	args := m.Mock.Called(req, res, err)

	return args.Bool(0)
}

func (m *MockedAuthenticator) OnAuthenticationSuccess(req *http.Request, res http.ResponseWriter, token GuardToken) bool {
	args := m.Mock.Called(req, res, token)

	return args.Bool(0)
}

type MockedManager struct {
	mock.Mock
}

func (m *MockedManager) GetUser(username string) (GuardUser, error) {
	args := m.Mock.Called(username)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(GuardUser), args.Error(1)
}
