package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"github.com/twinj/uuid"
)

func Test_Errors(t *testing.T) {
	errors := NewErrors()

	assert.False(t, errors.HasErrors())

	errors.AddError("field", "myerror")

	assert.True(t, errors.HasErrors())
	assert.True(t, errors.HasError("field"))
	assert.False(t, errors.HasError("foobar"))
}
