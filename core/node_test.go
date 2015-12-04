// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/twinj/uuid"
)

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
