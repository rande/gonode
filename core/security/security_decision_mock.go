// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"github.com/stretchr/testify/mock"
)

type MockedDecisionVoter struct {
	mock.Mock
}

func (m *MockedDecisionVoter) Support(o interface{}) bool {
	args := m.Mock.Called(o)

	return args.Bool(0)
}

func (m *MockedDecisionVoter) Decide(t SecurityToken, attrs Attributes, o interface{}) bool {
	args := m.Mock.Called(t, attrs, o)

	return args.Bool(0)
}
