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