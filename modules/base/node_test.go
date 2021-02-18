// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

func Test_Reference(t *testing.T) {

	input := []byte("\"64200eae-2539-4d92-a371-f906757f314d\"")

	ref := Reference{}
	err := ref.UnmarshalJSON(input)

	assert.Nil(t, err)
	assert.NotNil(t, ref.UUID)
}

func Test_Marshal_Empty_Reference(t *testing.T) {
	ref := Reference{}
	data, err := ref.MarshalJSON()

	assert.Len(t, data, 2) // => ""
	assert.Nil(t, err)
}

func Test_UnMarshal_NoData_Reference(t *testing.T) {
	input := []byte("")

	ref := Reference{}
	err := ref.UnmarshalJSON(input)

	assert.Error(t, err)
	assert.Equal(t, err, InvalidUuidLengthError)
	assert.Nil(t, ref.UUID)
}

func Test_UnMarshal_Empty_String_Reference(t *testing.T) {
	input := []byte("\"\"")

	ref := Reference{}
	err := ref.UnmarshalJSON(input)

	assert.NoError(t, err)
}
