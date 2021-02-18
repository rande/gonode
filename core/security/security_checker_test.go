// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyStruct struct {
}

func Test_DefaultAuthorizationChecker(t *testing.T) {
	d := &MockedDecisionVoter{}
	c := &DefaultAuthorizationChecker{
		DecisionVoter: d,
	}

	attrs := make(Attributes, 0)
	s := &DummyStruct{}

	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_ADMIN"},
	}

	d.On("Decide", tk, attrs, s).Return(false)

	b, err := c.IsGranted(tk, attrs, s)

	assert.NoError(t, err)
	assert.False(t, b)
}
