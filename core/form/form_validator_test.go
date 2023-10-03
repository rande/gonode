// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Max_Validator(t *testing.T) {

	v := int32(10)

	validator := MaxValidator(v)

	assert.Equal(t, "max", validator.Code())

	field := &FormField{
		Name:    "Name",
		Touched: true,
		reflect: reflect.ValueOf(v),
	}

	form := &Form{}

	field.SubmitedValue = int32(8)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmitedValue = int32(12)
	assert.NotNil(t, validator.Validate(field, form))
}

func Test_Min_Validator(t *testing.T) {
	validator := MinValidator(int64(10))

	assert.Equal(t, "min", validator.Code())

	field := &FormField{
		Name:    "Name",
		Touched: true,
		reflect: reflect.ValueOf(int64(0)),
	}

	form := &Form{}

	field.SubmitedValue = int64(10)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmitedValue = int64(10000)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmitedValue = int64(9)
	assert.NotNil(t, validator.Validate(field, form))
}

func Test_Min_Validator_Float(t *testing.T) {
	validator := MinValidator(float32(10.12312))

	assert.Equal(t, "min", validator.Code())

	field := &FormField{
		Name:    "Name",
		Touched: true,
		reflect: reflect.ValueOf(int64(0)),
	}

	form := &Form{}

	field.SubmitedValue = float32(10.12312)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmitedValue = float32(145.12312)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmitedValue = float32(9.12312)
	assert.NotNil(t, validator.Validate(field, form))
}
