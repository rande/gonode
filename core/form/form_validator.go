// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"net/mail"
	"net/url"
	"reflect"
)

var (
	ErrValidationError   = errors.New("the form is not valid")
	ErrRequiredValidator = errors.New("the value is required")
	ErrEmailValidator    = errors.New("the value is not a valid email")
	ErrUrlValidator      = errors.New("the value is not a valid url")
	ErrMaxValidator      = errors.New("the value is too big")
	ErrMinValidator      = errors.New("the value is too small")
	ErrOptionsValidator  = errors.New("the value is not in the list of options")
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func validateForm(fields []*FormField, form *Form) error {
	for _, field := range fields {
		if field.HasErrors {
			form.HasErrors = true
		}

		for _, validator := range field.Validators {
			if v, ok := field.InitialValue.(*Form); ok {
				if err := validateForm(v.Fields, v); err != nil {
					v.HasErrors = true
					field.HasErrors = true
					form.HasErrors = true
				}

				continue
			}

			if err := validator.Validate(field, form); err != nil {
				field.HasErrors = true
				form.HasErrors = true
				field.Errors = append(field.Errors, err.Error())
			}
		}

		if field.HasErrors {
			field.Error.Class = "text-red-500 text-xs italic"
		}
	}

	if !form.HasErrors {
		return nil
	}

	return ErrValidationError
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
				return nil
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
			if field.SubmittedValue == nil {
				return nil
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
				return nil
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
				return nil
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
				return nil
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
				return nil
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

func OptionsValidator() Validator {
	return &ValidatorFunc{
		code: "select",
		validator: func(field *FormField, form *Form) error {
			if field.SubmittedValue == nil {
				return nil
			}

			if !field.Input.Multiple {
				submittedValue := reflect.ValueOf(field.SubmittedValue)
				for _, option := range field.Options.(FieldOptions) {
					optionValue := reflect.ValueOf(option.Value)

					if optionValue.Equal(submittedValue) {
						return nil
					}
				}

				return ErrOptionsValidator
			} else {
				slice := reflect.ValueOf(field.SubmittedValue)
				for i := 0; i < slice.Len(); i++ {
					v := slice.Index(i)

					found := false
					for _, option := range field.Options.(FieldOptions) {
						optionValue := reflect.ValueOf(option.Value)

						if optionValue.Equal(v) {
							found = true
						}
					}

					if !found {
						return ErrOptionsValidator
					}
				}
			}

			return nil
		},
	}
}
