// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Date_NumberToStr(t *testing.T) {
	var v string
	var ok bool

	v, ok = NumberToStr(1.2)
	assert.True(t, ok)
	assert.Equal(t, "1.2", v)

	v, ok = NumberToStr(1231221312323223)
	assert.True(t, ok)
	assert.Equal(t, "1231221312323223", v)
}
