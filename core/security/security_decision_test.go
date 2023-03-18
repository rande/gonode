// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AffirmativeDecision_Support(t *testing.T) {
	d := &AffirmativeDecision{
		Voters: []Voter{
			&RoleVoter{Prefix: "ROLE_"},
		},
	}

	s := &DummyStruct{}

	assert.True(t, d.Support(s))
}

func Test_AffirmativeDecision_Valid(t *testing.T) {
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_USER"},
	}

	v := &MockedVoter{}
	v.On("Vote", tk, s, attrs).Return(ACCESS_GRANTED, nil)
	v.On("Support", s).Return(true)

	d := &AffirmativeDecision{
		Voters: []Voter{v},
	}

	assert.True(t, d.Decide(tk, attrs, s))
}

func Test_AffirmativeDecision_Invalid(t *testing.T) {
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_USER"},
	}

	v := &MockedVoter{}
	v.On("Vote", tk, s, attrs).Return(ACCESS_DENIED, nil)
	v.On("Support", s).Return(true)

	d := &AffirmativeDecision{
		Voters: []Voter{v},
	}

	assert.False(t, d.Decide(tk, attrs, s))
}

func Test_AffirmativeDecision_Abstain_Default(t *testing.T) {
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_USER"},
	}

	v := &MockedVoter{}
	v.On("Vote", tk, s, attrs).Return(ACCESS_ABSTAIN, nil)
	v.On("Support", s).Return(true)

	d := &AffirmativeDecision{
		Voters: []Voter{v},
	}

	assert.False(t, d.Decide(tk, attrs, s))
}
func Test_AffirmativeDecision_Abstain_ForceTrue(t *testing.T) {
	attrs := make(Attributes, 0)
	attrs = append(attrs, "ROLE_ADMIN")

	s := &DummyStruct{}
	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_USER"},
	}

	v := &MockedVoter{}
	v.On("Vote", tk, s, attrs).Return(ACCESS_ABSTAIN, nil)
	v.On("Support", s).Return(true)

	d := &AffirmativeDecision{
		Voters:                     []Voter{v},
		AllowIfAllAbstainDecisions: true, // force
	}

	assert.True(t, d.Decide(tk, attrs, s))
}
