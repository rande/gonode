package gonode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewNode(t *testing.T) {
	node := NewNode()

	assert.Equal(t, node.Weight, 1)
	assert.Equal(t, node.Revision, 1)
	assert.Equal(t, node.Enabled, true)
	assert.Equal(t, node.Status, StatusDraft)

	assert.Equal(t, node.Id(), 0)
}
