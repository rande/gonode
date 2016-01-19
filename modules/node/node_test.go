// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Foo struct {
	Bar  string
	Whoo string
}

func Test_Node_Value(t *testing.T) {
	f := &Foo{
		Bar:  "bar",
		Whoo: "whoo",
	}

	assert.Equal(t, "bar", GetValue(f, "Bar"))
	assert.Equal(t, nil, GetValue(f, "fake"))
	assert.Equal(t, nil, GetValue(f, "whoo"))
}
