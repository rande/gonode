// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"errors"
	"fmt"
	"strings"
)

type VoterResult int

var (
	ACCESS_GRANTED = VoterResult(1)
	ACCESS_ABSTAIN = VoterResult(2)
	ACCESS_DENIED  = VoterResult(-1)

	NotStringableAttributeError = errors.New("Attribute is not stringable")
)

type Attributes []interface{}

func (attrs Attributes) ToStringSlice() ([]string, error) {
	roles := []string{}

	var err error

	for _, v := range attrs {

		switch s := v.(type) {
		case string:
			roles = append(roles, s)
		case fmt.Stringer:
			roles = append(roles, s.String())
		default:
			err = NotStringableAttributeError
		}
	}

	return roles, err
}

func AttributesFromString(roles []string) Attributes {
	a := make(Attributes, 0)

	for _, s := range roles {
		a = append(a, s)
	}

	return a
}

type Voter interface {
	Support(v interface{}) bool
	Vote(t SecurityToken, o interface{}, attrs Attributes) (VoterResult, error)
}

type RoleVoter struct {
	Prefix string
}

func (v *RoleVoter) supportAttribute(value interface{}) bool {
	switch role := value.(type) {
	case string:
		return strings.HasPrefix(role, v.Prefix)

	default: // invalid attrs ...
		return false
	}
}

func (v *RoleVoter) Support(o interface{}) bool {
	return true
}

func (v *RoleVoter) Vote(t SecurityToken, o interface{}, attrs Attributes) (result VoterResult, err error) {
	result = ACCESS_ABSTAIN

	for _, value := range attrs {
		if !v.supportAttribute(value) {
			continue
		}

		result = ACCESS_DENIED
		for _, role := range t.GetRoles() {
			if role == value.(string) {
				return ACCESS_GRANTED, nil
			}
		}
	}

	return result, nil
}
