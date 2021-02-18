// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Panic_If_True(t *testing.T) {
	assert.Panics(t, func() {
		PanicIf(true, "Panic !!!")
	})
}

func Test_Panic_If_False(t *testing.T) {
	assert.NotPanics(t, func() {
		PanicIf(false, "Should not panic !!!")
	})
}

func Test_Panic_Unless_False(t *testing.T) {
	assert.Panics(t, func() {
		PanicUnless(false, "Panic !!!")
	})
}

func Test_Panic_Unless_True(t *testing.T) {
	assert.NotPanics(t, func() {
		PanicUnless(true, "Should not panic !!!")
	})
}

func Test_Panic_Callback(t *testing.T) {
	called := false

	assert.Panics(t, func() {
		PanicIf(true, "Panic !!!", func() {
			called = true
		})
	})

	assert.True(t, called)
}
