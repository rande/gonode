// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrNoValue = errors.New("unable to find the value")
)

var replacers = strings.NewReplacer(".", "_", "[", "_", "]", "")

type FieldCollectionValue struct {
	Value interface{}
	Key   string
}

type FieldCollectionOptions struct {
	Items     []*FieldCollectionValue
	Configure func(value interface{}) *Form
}

type FieldOption struct {
	Label   string
	Checked bool
}

type FieldOptions map[string]FieldOption

type Attributes map[string]string

type Input struct {
	Name         string
	Template     string
	Class        string
	Style        string
	Value        string
	Placeholder  string
	Type         string
	Id           string
	Pattern      string
	List         string
	Autocomplete string
	Readonly     bool
	Checked      bool
	Multiple     bool
	Required     bool
	Autofocus    bool
	Novalidate   bool
	Size         int
	MinLength    int
	MaxLength    int
	Min          int
	Max          int
	Step         int
	Height       int
	Width        int
}

type Label struct {
	Template string
	Class    string
	Style    string
	Value    string
}

type FormField struct {
	Prefix        string // used for nested forms
	Name          string
	Module        string
	Attributes    Attributes
	Label         Label
	Input         Input
	Mandatory     bool
	InitialValue  interface{}
	SubmitedValue interface{}
	Children      []*FormField
	Touched       bool
	Submitted     bool
	Errors        []string
	HasErrors     bool
	// from go to serialized
	Marshal func(field *FormField, form *Form) error
	// from serialized to go
	Unmarshal func(field *FormField, form *Form, values url.Values) error
	// Validator
}

func (f *FormField) Get(name string) *FormField {
	for _, field := range f.Children {
		if field.Name == name {
			return field
		}
	}

	return nil
}

func (f *FormField) Value(name string) interface{} {
	for _, field := range f.Children {
		if field.Name == name {
			return field.SubmitedValue
		}
	}

	return nil
}

func (f *FormField) Add(name string, fieldType string, value interface{}) *FormField {
	field := create(name, fieldType, value)

	f.Children = append(f.Children, field)

	return field
}

type Form struct {
	Fields    []*FormField
	State     string
	HasErrors bool
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

func formMarshal(field *FormField, form *Form) error {
	subForm := field.InitialValue.(*Form)

	field.Children = subForm.Fields

	for _, subField := range field.Children {
		subField.Input.Name = fmt.Sprintf("%s.%s", field.Name, subField.Name)
		subField.Input.Id = replacers.Replace(subField.Input.Name)
		subField.Prefix = field.Input.Name + "."

		subField.Marshal(subField, form)
	}

	return nil
}

func formUnmarshal(field *FormField, form *Form, values url.Values) error {
	fmt.Printf("formUnmarshal: %s\n", field.Name)

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

		fmt.Printf("collectionMarshal: %s\n", subField.Input.Name)

		field.Children = append(field.Children, subField)

		subField.Marshal(subField, form)
	}

	return nil
}

func collectionUnmarshal(field *FormField, form *Form, values url.Values) error {
	options := field.InitialValue.(*FieldCollectionOptions)

	for _, value := range options.Items {
		subField := field.Get(value.Key)

		subField.Unmarshal(subField, form, values)
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

		submitedValue[name] = FieldOption{
			Label:   option.Label,
			Checked: value == "checked" || value == "true" || value == "1" || value == "on" || value == "yes",
		}
	}

	field.SubmitedValue = submitedValue
	field.Touched = true

	return nil
}

func (f *Form) Get(name string) *FormField {
	for _, field := range f.Fields {
		if field.Name == name {
			return field
		}
	}

	return nil
}

type FormTypes interface {
	string | bool | int
	// bool | string | int | float | []string | []int | []float | []bool | map[string]string | map[string]int | map[string]float | map[string]bool
}

func Val[T FormTypes](form *Form, name string) T {
	var value T

	for _, field := range form.Fields {
		if field.Name == name {
			value = field.SubmitedValue.(T)
			break
		}
	}

	return value
}

func (f *Form) Value(name string) interface{} {
	for _, field := range f.Fields {
		if field.Name == name {
			return field.SubmitedValue
		}
	}

	return nil
}

func create(name string, fieldType string, value interface{}) *FormField {
	field := CreateFormField()
	field.Name = name
	field.Input.Type = fieldType
	field.Label.Value = name
	field.InitialValue = value

	if fieldType == "checkbox" {
		field.Marshal = checkboxMarshal
		field.Unmarshal = checkboxUnmarshal
	}

	if fieldType == "form" {
		field.Marshal = formMarshal
		field.Unmarshal = formUnmarshal
	}

	if fieldType == "collection" {
		field.Marshal = collectionMarshal
		field.Unmarshal = collectionUnmarshal
	}

	return field
}

func CreateForm() *Form {
	return &Form{}
}

func CreateFormField() *FormField {
	return &FormField{
		Prefix: "",
		Name:   "",
		Label: Label{
			Class:    "",
			Value:    "",
			Template: "form:label.tpl",
		},
		Input: Input{
			Class:       "",
			Style:       "",
			Value:       "",
			Placeholder: "",
			Type:        "text",
			Template:    "form:fields/input.text.tpl",
			Readonly:    false,
		},
		Module:        "form",
		InitialValue:  nil,
		Mandatory:     true,
		SubmitedValue: nil,
		Children:      []*FormField{},
		Touched:       false,
		Submitted:     false,
		Errors:        []string{},
		Marshal:       defaultMarshal,
		Unmarshal:     defaultUnmarshal,
		Attributes:    Attributes{},
	}
}

func (f *Form) Add(name string, fieldType string, value interface{}) *Form {
	field := create(name, fieldType, value)

	f.Fields = append(f.Fields, field)

	return f
}

func PrepareForm(form *Form) error {
	iterateFields(form, form.Fields)

	return nil
}

func iterateFields(form *Form, fields []*FormField) {
	for _, field := range fields {
		field.Input.Name = field.Name
		field.Marshal(field, form)
		field.Input.Id = replacers.Replace(field.Input.Name)
		fmt.Printf("iterateField: %s\n", field.Input.Id)
	}
}

func BindUrlValues(form *Form, values url.Values) error {
	for _, field := range form.Fields {
		field.Unmarshal(field, form, values)
	}

	return nil
}

func getValue(field *FormField, values url.Values) (string, error) {

	if !values.Has(field.Input.Name) {
		return "", ErrNoValue
	}

	return values.Get(field.Input.Name), nil
}
