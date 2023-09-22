// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"net/url"
	"reflect"
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

type FieldOptions map[string]*FieldOption

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
	reflect reflect.Value
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

func (f *FormField) Add(name string, fieldType string, options ...interface{}) *FormField {
	var value interface{} = nil

	if len(options) > 0 {
		value = options[0]
	}

	field := create(name, fieldType, value)

	f.Children = append(f.Children, field)

	return field
}

type Form struct {
	Data      interface{}
	Fields    []*FormField
	State     string
	HasErrors bool
	reflect   reflect.Value
	Locale    string
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
	// bool | string | int | float | []string | []int | []float | []bool | map[string]string | map[string]int | map
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

	if fieldType == "bool" {
		field.Marshal = booleanMarshal
		field.Unmarshal = booleanUnmarshal
	}

	if fieldType == "int" {
		field.Marshal = numberMarshal
		field.Unmarshal = intUnmarshal
	}

	if fieldType == "float" {
		field.Marshal = numberMarshal
		field.Unmarshal = floatUnmarshal
	}

	if fieldType == "uint" {
		field.Marshal = numberMarshal
		field.Unmarshal = unintUnmarshal
	}

	return field
}

func CreateForm(data interface{}) *Form {

	if data == nil {
		return &Form{}
	}

	return &Form{
		Data:    data,
		reflect: reflect.ValueOf(data).Elem(),
	}
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

func (f *Form) Add(name string, fieldType string, options ...interface{}) *Form {
	var value interface{} = nil

	if len(options) > 0 {
		value = options[0]
	}

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
		marshal(field, form)
		field.Input.Id = replacers.Replace(field.Input.Name)
	}
}

func BindUrlValues(form *Form, values url.Values) error {
	for _, field := range form.Fields {
		unmarshal(field, form, values)
	}

	return nil
}

func AttachValues(form *Form) error {
	// cannot attach value if no entity is linked
	if form.Data == nil {
		return nil
	}

	attachValues(form.Fields)

	return nil
}

func getValue(field *FormField, values url.Values) (string, error) {
	if !values.Has(field.Input.Name) {
		return "", ErrNoValue
	}

	return values.Get(field.Input.Name), nil
}

func attachValues(fields []*FormField) {
	// fmt.Println("Attach Value")

	for _, field := range fields {
		// fmt.Printf("Field name: %s\n", field.Name)

		if v, ok := field.InitialValue.(*Form); ok {
			// fmt.Printf("Sub form: %s\n", field.Name)
			if v.Data == nil {
				// fmt.Printf("skipping no data attached: %s\n", field.Name)
				continue
			}

			attachValues(v.Fields)

			continue
		}

		if field.reflect.Kind() == reflect.Invalid {
			// fmt.Printf("Invalid type: %s\n", field.Name)
			continue
		}

		newValue := reflect.ValueOf(field.SubmitedValue)

		if field.reflect.Kind() != newValue.Kind() {
			// fmt.Printf("Incompatible Type, field: %s (kind: %s) submitted: %s\n", field.Name, field.reflect.Kind(), newValue.Kind())
			continue
		}

		// fmt.Printf("Type, field: %s (kind: %s) submitted: %s, value: %s\n", field.Name, field.reflect.Kind(), newValue.Kind(), newValue.Interface())

		field.reflect.Set(newValue)
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
	field.Unmarshal(field, form, values)
}

func yes(value string) bool {
	return value == "checked" || value == "true" || value == "1" || value == "on" || value == "yes"
}
