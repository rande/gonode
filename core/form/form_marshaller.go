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
	"time"
)

var (
	ErrInvalidType = errors.New("value does not match the expected type")
)

func iterateFields(form *Form, fields []*FormField) {
	for _, field := range fields {
		field.Input.Name = field.Name
		configure(field, form)
		marshal(field, form)
		field.Input.Id = replacers.Replace(field.Input.Name)
	}
}

func configure(field *FormField, form *Form) {
	var rv reflect.Value

	if form.Data != nil && field.InitialValue == nil {
		field.reflectField = form.reflect.FieldByName(field.Name)
		rv = field.reflectField
	} else if field.InitialValue != nil {
		field.reflectValue = reflect.ValueOf(field.InitialValue)
		rv = field.reflectValue
	} else {
		panic(fmt.Sprintf("Unable to find the field type: %s, %v", field.Name, field.InitialValue))
	}

	if rv.Kind() != reflect.Invalid {
		field.InitialValue = rv.Interface()
	}

	var marshaller Marshaller = defaultMarshal
	var unmarshaller Unmarshaller = defaultUnmarshal
	kind := "text"

	if rv.Kind() == reflect.String {
		marshaller = defaultMarshal
		unmarshaller = defaultUnmarshal
		kind = "text"
	} else if rv.Kind() == reflect.Int ||
		rv.Kind() == reflect.Int8 ||
		rv.Kind() == reflect.Int16 ||
		rv.Kind() == reflect.Int32 ||
		rv.Kind() == reflect.Int64 ||
		rv.Kind() == reflect.Uint ||
		rv.Kind() == reflect.Uint8 ||
		rv.Kind() == reflect.Uint16 ||
		rv.Kind() == reflect.Uint32 ||
		rv.Kind() == reflect.Uint64 ||
		rv.Kind() == reflect.Float32 ||
		rv.Kind() == reflect.Float64 {
		marshaller = numberMarshal
		unmarshaller = numberUnmarshal
		kind = "number"
	} else if rv.Kind() == reflect.Bool {
		marshaller = booleanMarshal
		unmarshaller = booleanUnmarshal
		kind = "boolean"

	} else if rv.Kind() == reflect.Struct {
		if rv.Type() == reflect.TypeOf(time.Time{}) {
			marshaller = dateMarshal
			unmarshaller = dateUnmarshal
			kind = "date"
		}
	} else if rv.Kind() == reflect.Ptr {
		if rv.Type() == reflect.TypeOf(&Form{}) {
			marshaller = formMarshal
			unmarshaller = formUnmarshal
			kind = "form"
		} else if rv.Type() == reflect.TypeOf(&FieldCollectionOptions{}) {
			marshaller = collectionMarshal
			unmarshaller = collectionUnmarshal
			kind = "collection"
		}
	}

	if field.Input.Type == "" {
		field.Input.Type = kind
	}

	if field.Marshaller == nil {
		field.Marshaller = marshaller
	}

	if field.Unmarshaller == nil {
		field.Unmarshaller = unmarshaller
	}
}

func marshal(field *FormField, form *Form) {
	field.Marshaller(field, form)
}

func unmarshal(field *FormField, form *Form, values url.Values) {
	field.Errors = []string{}

	if values.Has(field.Input.Name) {
		field.Input.Value = values.Get(field.Input.Name)
	}

	field.Unmarshaller(field, form, values)

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
	field.SubmittedValue = value
	field.Touched = true

	return nil
}

func booleanMarshal(field *FormField, form *Form) error {
	defaultMarshal(field, form)

	if v, ok := BoolToStr(field.InitialValue); ok && !field.HasErrors {
		field.Input.Value = v
	}

	return nil
}

func booleanUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if v, ok := StrToBool(field.SubmittedValue); ok && !field.HasErrors {
		field.SubmittedValue = v
	}

	return nil
}

func numberMarshal(field *FormField, form *Form) error {
	defaultMarshal(field, form)

	if field.InitialValue != nil {
		v, _ := NumberToStr(field.InitialValue)
		field.Input.Value = v
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

	var rv reflect.Value
	var v string
	var ok bool

	if field.reflectValue.IsValid() {
		rv = field.reflectValue
	} else {
		rv = field.reflectField
	}

	if v, ok = field.SubmittedValue.(string); !ok {
		field.Errors = append(field.Errors, ErrInvalidType.Error())
		return nil
	}

	if value, ok := StrToNumber(v, rv.Kind()); ok {
		field.SubmittedValue = value
	} else {
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

	if v, ok := field.SubmittedValue.(string); ok {
		if t, err := time.ParseInLocation("2006-01-02", v, time.UTC); err != nil {
			field.Errors = append(field.Errors, err.Error())
		} else {
			field.SubmittedValue = t
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
		configure(subField, subForm)
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

		configure(subField, form)
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

	for i, option := range field.InitialValue.(FieldOptions) {
		// find a nice way to generate the name
		subField := CreateFormField()
		subField.Name = fmt.Sprintf("%d", i)
		subField.Input.Name = fmt.Sprintf("%s[%s]", field.Input.Name, subField.Name)
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
	submitedValues := FieldOptions{}
	for i, option := range field.InitialValue.(FieldOptions) {
		name := fmt.Sprintf("%d", i)
		value, err := getValue(field.Get(name), values)

		if err != nil {
			field.Errors = append(field.Errors, err.Error())
			field.HasErrors = true

			return err
		}

		submitedValues = append(submitedValues, &FieldOption{
			Label:   option.Label,
			Checked: yes(value),
			Id:      option.Id,
			Value:   option.Value, // should not be set
		})
	}

	field.SubmittedValue = submitedValues
	field.Touched = true

	return nil
}
