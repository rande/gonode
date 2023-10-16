// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

const (
	ENC_TYPE_URLENCODED = "application/x-www-form-urlencoded"
	ENC_TYPE_MULTIPART  = "multipart/form-data"
)

var (
	ErrInvalidState = errors.New("invalid state")
)

const (
	Initialized int = 0
	Prepared        = 1
	Submitted       = 2
	Validated       = 3
	Processed       = 4
)

func Process(form *Form, request *http.Request) error {

	if err := PrepareForm(form); err != nil {
		return err
	}

	if request.Method == "GET" {
		return nil
	}

	if form.EncType == ENC_TYPE_URLENCODED {
		if err := request.ParseForm(); err != nil {
			return err
		}
	}

	if form.EncType == ENC_TYPE_MULTIPART {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			return err
		}
	}

	if err := BindUrlValues(form, request.Form); err != nil {
		return err
	}

	if err := ValidateForm(form); err != nil {
		return err
	}

	if form.Data == nil {
		return nil
	}

	if err := AttachValues(form); err != nil {
		return err
	}

	return nil
}

// Iterate over all fields and call the marshaller function to transform the Go value into
// a serialized value used in the HTML form. This also setup the different attribures required
// by the HTML form.
func PrepareForm(form *Form) error {
	fmt.Printf("PrepareForm > State: %d, expected: %d\n", form.State, Initialized)

	if form.State != Initialized {
		return ErrInvalidState
	}

	form.State = Prepared

	iterateFields(form, form.Fields)

	return nil
}

// Iterate over all fields and bind the submitted values to the SubmittedValue field
// in the form. This will not attach the value to the underlying data structure, use AttachValues
// for this.
//
// This will call all unmarshallers defined on each form field.
func BindUrlValues(form *Form, values url.Values) error {
	if form.State != Prepared {
		return ErrInvalidState
	}

	form.State = Submitted

	for _, field := range form.Fields {
		unmarshal(field, form, values)
	}

	if form.HasErrors {
		return ErrInvalidSubmittedData
	}

	return nil
}

func ValidateForm(form *Form) error {
	if form.State != Submitted {
		return ErrInvalidState
	}

	form.State = Validated

	return validateForm(form.Fields, form)
}

// Iterate over all submitted values and assign them to the related data structure
// if possible. This will not work work if the form has some errors.
func AttachValues(form *Form) error {
	if form.State != Validated {
		return ErrInvalidState
	}

	// cannot attach value if no entity is linked
	if form.Data == nil {
		return nil
	}

	if form.HasErrors {
		return ErrInvalidSubmittedData
	}

	attachValues(form.Fields)

	return nil
}

func attachValues(fields []*FormField) {
	// fmt.Println("attachValues > Attach Value")

	for _, field := range fields {
		// fmt.Printf("attachValues > Field name: %s\n", field.Name)
		if v, ok := field.InitialValue.(*Form); ok {
			// fmt.Printf("attachValues > Sub form: %s\n", field.Name)
			if v.Data == nil {
				// fmt.Printf("attachValues > skipping no data attached: %s\n", field.Name)
				continue
			}

			attachValues(v.Fields)
			continue
		}

		if !field.Touched {
			// fmt.Printf("attachValues > Field not touched: %s, skipping\n", field.Name)
			continue
		} else {
			// fmt.Printf("attachValues > Field touched: %s, updating\n", field.Name)
		}

		if field.reflect.Kind() == reflect.Invalid {
			// fmt.Printf("attachValues > Invalid type: %s\n", field.Name)
			continue
		}

		newValue := reflect.ValueOf(field.SubmittedValue)
		// fmt.Printf("attachValues > Field name: %s, value: %s\n", field.Name, newValue)
		if newValue.CanConvert(field.reflect.Type()) {
			field.reflect.Set(newValue.Convert(field.reflect.Type()))
		} else {
			// fmt.Printf("attachValues > Unable to convert type: Type, field: %s (kind: %s) submitted: %s, value: %s\n", field.Name, field.reflect.Kind(), newValue.Kind(), newValue.Interface())
		}
	}
}
