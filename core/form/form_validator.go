// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"net/mail"
	"net/url"
)

var (
	ErrRequiredValidator = errors.New("the value is required")
	ErrEmailValidator    = errors.New("the value is not a valid email")
	ErrUrlValidator      = errors.New("the value is not a valid url")
	ErrMaxValidator      = errors.New("the value is too big")
	ErrMinValidator      = errors.New("the value is too small")
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func ValidateForm(form *Form) bool {
	return validateForm(form.Fields, form)
}

func validateForm(fields []*FormField, form *Form) bool {
	isValid := true

	for _, field := range fields {
		if field.HasErrors {
			isValid = false
		}

		for _, validator := range field.Validators {
			if v, ok := field.InitialValue.(*Form); ok {
				if !validateForm(v.Fields, v) {
					isValid = false
					field.HasErrors = true
				}

				continue
			}

			if err := validator.Validate(field, form); err != nil {
				isValid = false
				field.HasErrors = true
				field.Errors = append(field.Errors, err.Error())
			}
		}

		if field.HasErrors {
			field.Error.Class = "text-red-500 text-xs italic"
		}
	}

	return isValid
}

type Validator interface {
	Validate(field *FormField, form *Form) error
	Code() string
}

type ValidatorFunc struct {
	code      string
	validator func(field *FormField, form *Form) error
}

func (v *ValidatorFunc) Validate(field *FormField, form *Form) error {
	return v.validator(field, form)
}

func (v *ValidatorFunc) Code() string {
	return v.code
}

func RequiredValidator() Validator {
	return &ValidatorFunc{
		code: "required",
		validator: func(field *FormField, form *Form) error {
			if !field.Touched {
				return ErrRequiredValidator
			}

			if field.SubmittedValue == nil {
				return ErrRequiredValidator
			}

			return nil
		},
	}
}

func EmailValidator() Validator {
	return &ValidatorFunc{
		code: "email",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return ErrEmailValidator
			}

			if _, err := mail.ParseAddress(field.SubmittedValue.(string)); err != nil {
				return ErrEmailValidator
			}

			return nil
		},
	}
}

func UrlValidator() Validator {
	return &ValidatorFunc{
		code: "url",
		validator: func(field *FormField, form *Form) error {
			if !field.Touched {
				return ErrUrlValidator
			}

			if field.SubmittedValue == nil {
				return ErrUrlValidator
			}

			if _, err := url.Parse(field.SubmittedValue.(string)); err != nil {
				return ErrUrlValidator
			}

			return nil
		},
	}
}

func MinValidator[T number](min T) Validator {
	return &ValidatorFunc{
		code: "min",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return ErrMinValidator
			}

			if field.SubmittedValue.(T) < min {
				return ErrMinValidator
			}

			return nil
		},
	}
}

func MaxValidator[T number](max T) Validator {
	return &ValidatorFunc{
		code: "max",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return ErrMaxValidator
			}

			if field.SubmittedValue.(T) > max {
				return ErrMaxValidator
			}

			return nil
		},
	}
}

func MaxLengthValidator(max uint32, mode string) Validator {
	return &ValidatorFunc{
		code: "max_length",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return ErrMaxValidator
			}

			if mode == "bytes" {
				if uint32(len(field.SubmittedValue.(string))) > max {
					return ErrMaxValidator
				}
			} else if mode == "rune" {
				if uint32(len([]rune(field.SubmittedValue.(string)))) > max {
					return ErrMaxValidator
				}
			} else {
				panic("Invalid mode")
			}

			return nil
		},
	}
}

func MinLengthValidator(min uint32, mode string) Validator {
	return &ValidatorFunc{
		code: "min_length",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return ErrMaxValidator
			}

			if mode == "bytes" {
				if uint32(len(field.SubmittedValue.(string))) <= min {
					return ErrMaxValidator
				}
			} else if mode == "rune" {
				if uint32(len([]rune(field.SubmittedValue.(string)))) <= min {
					return ErrMaxValidator
				}
			} else {
				panic("Invalid mode")
			}

			return nil
		},
	}
}
