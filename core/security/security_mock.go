// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"github.com/stretchr/testify/mock"
)

type MockedVoter struct {
	mock.Mock
}

func (m *MockedVoter) Support(v interface{}) bool {
	args := m.Mock.Called(v)

	return args.Bool(0)
}

func (m *MockedVoter) Vote(t SecurityToken, o interface{}, attrs Attributes) (VoterResult, error) {
	args := m.Mock.Called(t, o, attrs)

	return args.Get(0).(VoterResult), args.Error(1)
}
