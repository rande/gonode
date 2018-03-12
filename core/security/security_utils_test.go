// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EnsureRoles(t *testing.T) {
	roles := EnsureRoles([]string{"admin"}, "admin", "editor", "reviewer")

	assert.Equal(t, roles, []string{"admin", "editor", "reviewer"})
}

func Test_EnsureRoles_Empty(t *testing.T) {
	roles := EnsureRoles([]string{})

	assert.Equal(t, roles, []string{})
}
