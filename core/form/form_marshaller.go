// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrInvalidType = errors.New("value does not match the expected type")
)

func iterateFields(form *Form, fields []*FormField) {
	for _, field := range fields {
		field.Input.Name = field.Name
		marshal(field, form)
		field.Input.Id = replacers.Replace(field.Input.Name)
	}
}

func marshal(field *FormField, form *Form) {
	// we do not try to get the value, if the value is already set
	// and if form.Data is not set.
	if form.Data != nil || field.InitialValue == nil {
		field.reflect = form.reflect.FieldByName(field.Name)

		if field.reflect.Kind() != reflect.Invalid {
			field.InitialValue = field.reflect.Interface()
		}
	}

	field.Marshal(field, form)
}

func unmarshal(field *FormField, form *Form, values url.Values) {
	field.Errors = []string{}

	if values.Has(field.Input.Name) {
		field.Input.Value = values.Get(field.Input.Name)
	}

	field.Unmarshal(field, form, values)

	field.HasErrors = len(field.Errors) > 0
}

func defaultMarshal(field *FormField, form *Form) error {
	field.Input.Value = fmt.Sprintf("%s", field.InitialValue)
	field.Input.Name = fmt.Sprintf("%s%s", field.Prefix, field.Name)
	field.Input.Id = replacers.Replace(field.Input.Name)

	return nil
}

func defaultUnmarshal(field *FormField, form *Form, values url.Values) error {
	value, err := getValue(field, values)

	if err != nil {
		field.Errors = append(field.Errors, err.Error())
		field.HasErrors = true

		return err
	}

	// to do, add a validator call
	field.SubmitedValue = value
	field.Touched = true

	return nil
}

func booleanMarshal(field *FormField, form *Form) error {
	defaultMarshal(field, form)

	if v, ok := field.InitialValue.(bool); ok {
		if v {
			field.Input.Value = "yes"
		} else {
			field.Input.Value = "no"
		}
	}

	return nil
}

func booleanUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if field.HasErrors {
		return nil
	}

	if v, ok := field.SubmitedValue.(string); ok {
		if yes(v) {
			field.SubmitedValue = true
		} else {
			field.SubmitedValue = false
		}
	}

	return nil
}

func numberMarshal(field *FormField, form *Form) error {
	defaultMarshal(field, form)

	if field.InitialValue != nil {
		field.Input.Value = fmt.Sprintf("%d", field.InitialValue)
	} else {
		field.Input.Value = "0"
	}

	return nil
}

func numberUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if field.HasErrors {
		return nil
	}

	var v string
	var ok bool
	hasError := true

	if v, ok = field.SubmitedValue.(string); !ok {
		field.Errors = append(field.Errors, ErrInvalidType.Error())
		return nil
	}

	if field.reflect.Kind() == reflect.Int {
		if i, err := strconv.ParseInt(v, 10, 0); err == nil {
			field.SubmitedValue = int(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Int8 {
		if i, err := strconv.ParseInt(v, 10, 8); err == nil {
			field.SubmitedValue = int8(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Int16 {
		if i, err := strconv.ParseInt(v, 10, 16); err == nil {
			field.SubmitedValue = int16(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Int32 {
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			field.SubmitedValue = int32(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Int64 {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			field.SubmitedValue = int64(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Uint {
		if i, err := strconv.ParseUint(v, 10, 0); err == nil {
			field.SubmitedValue = uint(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Uint8 {
		if i, err := strconv.ParseUint(v, 10, 8); err == nil {
			field.SubmitedValue = uint8(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Uint16 {
		if i, err := strconv.ParseUint(v, 10, 16); err == nil {
			field.SubmitedValue = uint16(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Uint32 {
		if i, err := strconv.ParseUint(v, 10, 32); err == nil {
			field.SubmitedValue = uint32(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Uint64 {
		if i, err := strconv.ParseUint(v, 10, 64); err == nil {
			field.SubmitedValue = uint64(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Float32 {
		if i, err := strconv.ParseFloat(v, 32); err == nil {
			field.SubmitedValue = float32(i)
			hasError = false
		}
	} else if field.reflect.Kind() == reflect.Float64 {
		if i, err := strconv.ParseFloat(v, 64); err == nil {
			field.SubmitedValue = float64(i)
			hasError = false
		}
	}

	if hasError {
		field.Errors = append(field.Errors, ErrInvalidType.Error())
	}

	return nil
}

// The displayed date format will differ from the actual value —
// the displayed date is formatted based on the locale of the user's browser,
// but the parsed value is always formatted yyyy-mm-dd.
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/date
func dateMarshal(field *FormField, form *Form) error {
	defaultMarshal(field, form)

	if v, ok := field.InitialValue.(time.Time); ok {
		field.Input.Value = v.Format("2006-01-02")
	} else {
		fmt.Printf("Invalid date: %s", field.InitialValue)
	}

	return nil
}

func dateUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if v, ok := field.SubmitedValue.(string); ok {
		if t, err := time.ParseInLocation("2006-01-02", v, time.UTC); err != nil {
			field.Errors = append(field.Errors, err.Error())
		} else {
			field.SubmitedValue = t
		}
	}

	return nil
}

func formMarshal(field *FormField, form *Form) error {
	subForm := field.InitialValue.(*Form)

	subForm.Locale = form.Locale

	field.Children = subForm.Fields

	for _, subField := range field.Children {
		subField.Input.Name = fmt.Sprintf("%s.%s", field.Name, subField.Name)
		subField.Input.Id = replacers.Replace(subField.Input.Name)
		subField.Prefix = field.Input.Name + "."

		marshal(subField, subForm)
	}

	return nil
}

func formUnmarshal(field *FormField, form *Form, values url.Values) error {
	for _, child := range field.Children {
		unmarshal(child, form, values)
	}

	return nil
}

func collectionMarshal(field *FormField, form *Form) error {
	options := field.InitialValue.(*FieldCollectionOptions)

	field.Input.Name = fmt.Sprintf("%s%s", field.Prefix, field.Name)
	field.Input.Id = replacers.Replace(field.Input.Name)

	for _, value := range options.Items {
		subForm := options.Configure(value.Value)

		subField := create(value.Key, "form", subForm)
		subField.Input.Name = fmt.Sprintf("%s[%s]", field.Input.Name, value.Key)
		subField.Input.Id = replacers.Replace(subField.Input.Name)
		subField.Prefix = field.Input.Name + "."

		field.Children = append(field.Children, subField)

		marshal(subField, form)
	}

	return nil
}

func collectionUnmarshal(field *FormField, form *Form, values url.Values) error {
	options := field.InitialValue.(*FieldCollectionOptions)

	for _, value := range options.Items {
		subField := field.Get(value.Key)

		unmarshal(subField, form, values)
	}

	return nil
}

func checkboxMarshal(field *FormField, form *Form) error {
	field.Input.Name = fmt.Sprintf("%s%s", field.Prefix, field.Name)
	field.Input.Id = replacers.Replace(field.Input.Name)

	for name, option := range field.InitialValue.(CheckboxOptions) {
		// find a nice way to generate the name
		subField := CreateFormField()
		subField.Name = name
		subField.Input.Name = fmt.Sprintf("%s[%s]", field.Input.Name, name)
		subField.Input.Id = replacers.Replace(subField.Input.Name)
		subField.Label.Value = option.Label
		subField.Input.Type = "checkbox"
		subField.InitialValue = option.Checked

		field.Children = append(field.Children, subField)
	}

	return nil
}

func checkboxUnmarshal(field *FormField, form *Form, values url.Values) error {
	// we need to check for extra values!
	submitedValue := CheckboxOptions{}
	for name, option := range field.InitialValue.(CheckboxOptions) {
		value, err := getValue(field.Get(name), values)

		if err != nil {
			field.Errors = append(field.Errors, err.Error())
			field.HasErrors = true

			return err
		}

		submitedValue[name] = &CheckboxOption{
			Label:   option.Label,
			Checked: yes(value),
		}
	}

	field.SubmitedValue = submitedValue
	field.Touched = true

	return nil
}
