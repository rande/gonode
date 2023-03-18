// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Errors(t *testing.T) {
	errors := NewErrors()

	assert.False(t, errors.HasErrors())

	errors.AddError("field", "myerror")

	assert.True(t, errors.HasErrors())
	assert.True(t, errors.HasError("field"))
	assert.False(t, errors.HasError("foobar"))
}
