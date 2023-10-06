// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"fmt"
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
	}

	form := &Form{}

	field.SubmittedValue = int32(8)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = int32(12)
	assert.NotNil(t, validator.Validate(field, form))
}

func Test_Min_Validator(t *testing.T) {
	validator := MinValidator(int64(10))

	assert.Equal(t, "min", validator.Code())

	field := &FormField{
		Name:    "Name",
		Touched: true,
	}

	form := &Form{}

	field.SubmittedValue = int64(10)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = int64(10000)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = int64(9)
	assert.NotNil(t, validator.Validate(field, form))
}

func Test_Min_Validator_Float(t *testing.T) {
	validator := MinValidator(float32(10.12312))

	assert.Equal(t, "min", validator.Code())

	field := &FormField{
		Name:    "Name",
		Touched: true,
	}

	form := &Form{}

	field.SubmittedValue = float32(10.12312)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = float32(145.12312)
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = float32(9.12312)
	assert.NotNil(t, validator.Validate(field, form))
}

func Test_MaxLength(t *testing.T) {
	msg := "abc"
	fmt.Printf("size: %d\n", len(msg))
	fmt.Printf("size: %d\n", len([]rune(msg)))
	fmt.Printf("size: %d\n", uint32(len([]rune(msg))))

	field := &FormField{
		Name:    "Name",
		Touched: true,
	}

	form := &Form{}

	field.SubmittedValue = "🇫🇷"
	validator := MaxLengthValidator(1, "bytes")
	assert.Equal(t, "max_length", validator.Code())
	assert.NotNil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MaxLengthValidator(10, "bytes")
	assert.NotNil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MaxLengthValidator(16, "bytes")
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MaxLengthValidator(10, "rune")
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷"
	validator = MaxLengthValidator(1, "rune")
	assert.NotNil(t, validator.Validate(field, form))

	field.SubmittedValue = "abc"
	validator = MaxLengthValidator(10, "bytes")
	assert.Nil(t, validator.Validate(field, form))
}

func Test_MinLength(t *testing.T) {
	msg := "🇫🇷🇫🇷"
	fmt.Printf("size: %d\n", len(msg))
	fmt.Printf("size: %d\n", len([]rune(msg)))
	fmt.Printf("size: %d\n", uint32(len([]rune(msg))))

	field := &FormField{
		Name:    "Name",
		Touched: true,
	}

	form := &Form{}

	field.SubmittedValue = "🇫🇷"
	validator := MinLengthValidator(1, "bytes")
	assert.Equal(t, "min_length", validator.Code())
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MinLengthValidator(10, "bytes")
	assert.Nil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MinLengthValidator(16, "bytes")
	assert.NotNil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷🇫🇷"
	validator = MinLengthValidator(10, "rune")
	assert.NotNil(t, validator.Validate(field, form))

	field.SubmittedValue = "🇫🇷"
	validator = MinLengthValidator(1, "rune")
	assert.Nil(t, validator.Validate(field, form))
}
