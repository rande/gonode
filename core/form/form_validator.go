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
	ErrValidatorRequired = errors.New("the value is required")
	ErrValidatorEmail    = errors.New("the value is not a valid email")
	ErrValidatorUrl      = errors.New("the value is not a valid url")
)

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

			if err := validator(field, form); err != nil {
				isValid = false
				field.HasErrors = true
				field.Errors = append(field.Errors, err.Error())
			}
		}
	}

	return isValid
}

func RequiredValidator() func(field *FormField, form *Form) error {
	return func(field *FormField, form *Form) error {
		if !field.Touched {
			return ErrValidatorRequired
		}

		if field.SubmitedValue == nil {
			return ErrValidatorRequired
		}

		return nil
	}
}

func EmailValidator() func(field *FormField, form *Form) error {
	return func(field *FormField, form *Form) error {
		if !field.Touched {
			return ErrValidatorEmail
		}

		if field.SubmitedValue == nil {
			return ErrValidatorEmail
		}

		if _, err := mail.ParseAddress(field.SubmitedValue.(string)); err != nil {
			return ErrValidatorEmail
		}

		return nil
	}
}

func UrlValidator() func(field *FormField, form *Form) error {
	return func(field *FormField, form *Form) error {
		if !field.Touched {
			return ErrValidatorUrl
		}

		if field.SubmitedValue == nil {
			return ErrValidatorUrl
		}

		if _, err := url.Parse(field.SubmitedValue.(string)); err != nil {
			return ErrValidatorUrl
		}

		return nil
	}
}
