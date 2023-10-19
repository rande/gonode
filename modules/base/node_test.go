// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_NewNode(t *testing.T) {
	node := NewNode()

	assert.Equal(t, node.Weight, 1)
	assert.Equal(t, node.Revision, 1)
	assert.Equal(t, node.Enabled, true)
	assert.Equal(t, node.Status, StatusNew)

	assert.Equal(t, node.Id, 0)
}
