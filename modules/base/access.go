// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/rande/gonode/core/security"
)

type AccessOptions struct {
	Token security.SecurityToken
	Roles security.Attributes
}

func NewAccessOptions(token security.SecurityToken, roles security.Attributes) *AccessOptions {
	return &AccessOptions{
		Token: token,
		Roles: roles,
	}
}

func NewAccessOptionsFromToken(token security.SecurityToken) *AccessOptions {
	return NewAccessOptions(token, GetSecurityAttributes(token.GetRoles()))
}

func GetSecurityAttributes(access []string) security.Attributes {
	attrs := security.Attributes{}
	for _, r := range access {
		attrs = append(attrs, r)
	}

	return attrs
}

type AccessVoter struct {
}

func (a *AccessVoter) Support(v interface{}) bool {
	switch v.(type) {
	case *Node:
		return true
	default:
		return false
	}
}

func (a *AccessVoter) Vote(t security.SecurityToken, o interface{}, attrs security.Attributes) (result security.VoterResult, err error) {
	result = security.ACCESS_ABSTAIN

	node := o.(*Node)

	for _, value := range node.Access {
		result = security.ACCESS_DENIED
		for _, role := range t.GetRoles() {
			if role == value {
				return security.ACCESS_GRANTED, nil
			}
		}
	}

	return result, nil
}
