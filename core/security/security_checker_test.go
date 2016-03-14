// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

	c := &DefaultAuthorizationChecker{}
	attrs := make(Attributes, 0)
	s := &DummyStruct{}

	tk := &DefaultSecurityToken{
		Roles: []string{"ROLE_ADMIN"},
	}

	b, err := c.IsGranted(tk, attrs, s)

	assert.NoError(t, err)
	assert.False(t, b)
}
