// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Role_Voter_Default_Behavior(t *testing.T) {
	v := &RoleVoter{
		Prefix: "ROLE_",
	}
	attrs := make(Attributes, 0)
	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{},
	}

	assert.True(t, v.Support(s))

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, ACCESS_ABSTAIN)
}

func Test_Role_Voter_Default_Access_Behavior(t *testing.T) {
	v := &RoleVoter{
		Prefix: "ROLE_",
	}
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_ADMIN"},
	}

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, ACCESS_GRANTED)
}

func Test_Role_Voter_Default_Denied_Behavior(t *testing.T) {
	v := &RoleVoter{
		Prefix: "ROLE_",
	}
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{},
	}

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, ACCESS_DENIED)
}

func Test_Role_Voter_Invalid_Attributes(t *testing.T) {
	v := &RoleVoter{
		Prefix: "ROLE_",
	}
	attrs := make(Attributes, 0)
	attrs = append(attrs, &DummyStruct{})
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_ADMIN"},
	}

	r, err := v.Vote(tk, s, attrs)

	assert.NoError(t, err)
	assert.Equal(t, r, ACCESS_GRANTED)
}
