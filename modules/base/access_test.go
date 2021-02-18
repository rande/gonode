// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"testing"

	"github.com/rande/gonode/core/security"
	"github.com/stretchr/testify/assert"
)

func Test_Role_Node_Voter_Default_Behavior(t *testing.T) {
	v := &AccessVoter{}

	attrs := make(security.Attributes, 0)
	s := &Node{}
	tk := &security.DefaultSecurityToken{
		Roles: []string{},
	}

	assert.True(t, v.Support(s))

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, security.ACCESS_ABSTAIN)
}

func Test_Role_Node_Voter_Default_Access_Behavior(t *testing.T) {
	v := &AccessVoter{}

	attrs := make(security.Attributes, 0)

	s := &Node{}
	s.Access = []string{"node:api:master"}

	tk := &security.DefaultSecurityToken{
		Roles: []string{"node:api:master"},
	}

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, security.ACCESS_GRANTED)
}

func Test_Role_Node_Voter_Default_Denied_Behavior(t *testing.T) {
	v := &AccessVoter{}

	attrs := make(security.Attributes, 0)

	s := &Node{}
	s.Access = []string{"node:api:master"}

	tk := &security.DefaultSecurityToken{
		Roles: []string{},
	}

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, security.ACCESS_DENIED)
}
