// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"fmt"
	"net/url"
	"strconv"
)

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
	return defaultMarshal(field, form)
}

func intUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if field.HasErrors {
		return nil
	}

	if v, ok := field.SubmitedValue.(string); ok {
		if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			field.Errors = append(field.Errors, err.Error())
		} else {
			field.SubmitedValue = i
		}
	}

	return nil
}

func unintUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if field.HasErrors {
		return nil
	}

	if v, ok := field.SubmitedValue.(string); ok {
		if i, err := strconv.ParseUint(v, 10, 64); err != nil {
			field.Errors = append(field.Errors, err.Error())
		} else {
			field.SubmitedValue = i
		}
	}

	return nil
}

func floatUnmarshal(field *FormField, form *Form, values url.Values) error {
	defaultUnmarshal(field, form, values)

	if field.HasErrors {
		return nil
	}

	if v, ok := field.SubmitedValue.(string); ok {
		if i, err := strconv.ParseFloat(v, 64); err != nil {
			field.Errors = append(field.Errors, err.Error())
		} else {
			field.SubmitedValue = i
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
		child.Unmarshal(child, form, values)
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

	for name, option := range field.InitialValue.(FieldOptions) {
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
	submitedValue := FieldOptions{}
	for name, option := range field.InitialValue.(FieldOptions) {
		value, err := getValue(field.Get(name), values)

		if err != nil {
			field.Errors = append(field.Errors, err.Error())
			field.HasErrors = true

			return err
		}

		submitedValue[name] = &FieldOption{
			Label:   option.Label,
			Checked: yes(value),
		}
	}

	field.SubmitedValue = submitedValue
	field.Touched = true

	return nil
}
