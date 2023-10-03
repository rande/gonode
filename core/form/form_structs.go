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

type CheckboxOption struct {
	Label   string
	Checked bool
}

type FieldOption struct {
	Label string
}

type CheckboxOptions map[string]*CheckboxOption

type Attributes map[string]string

type Input struct {
	Name         string
	Template     string
	Class        string
	Style        string
	Value        string
	Placeholder  string // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#placeholder
	Type         string
	Id           string
	Pattern      string // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#pattern
	List         string
	Autocomplete string // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input#autocomplete
	Readonly     bool
	Checked      bool
	Multiple     bool
	Required     bool
	Autofocus    bool
	Novalidate   bool
	Size         int
	MinLength    int
	MaxLength    int
	Min          interface{}
	Max          interface{}
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
	Help          string
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
	Validators    []Validator
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

	field.Input.Class = "shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
	field.Label.Class = "block text-gray-700 text-sm font-bold mb-2"

	if fieldType == "color" {
		field.Input.Class = ""
	}

	if fieldType == "range" {
		field.Input.Class = ""
		field.Marshal = numberMarshal
	}

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
		field.Input.Type = "checkbox"
		field.Input.Class = ""
	}

	if fieldType == "number" {
		field.Marshal = numberMarshal
		field.Unmarshal = numberUnmarshal
	}

	if fieldType == "date" {
		field.Marshal = dateMarshal
		field.Unmarshal = dateUnmarshal
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
		Help:          "",
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
		Validators:    []Validator{},
	}
}

func (f *FormField) AddValidator(validator Validator) *FormField {
	f.Validators = append(f.Validators, validator)
	return f
}

func (f *FormField) AddValidators(validators ...Validator) *FormField {
	f.Validators = append(f.Validators, validators...)
	return f
}

func (f *FormField) ResetValidators() *FormField {
	f.Validators = []Validator{}
	return f
}

func (f *FormField) RemoveValidator(code string) *FormField {
	validators := []Validator{}

	for _, validator := range f.Validators {
		if validator.Code() != code {
			validators = append(validators, validator)
		}
	}

	f.Validators = validators

	return f
}

func (f *FormField) SetModule(value string) *FormField {
	f.Module = value
	return f
}

func (f *FormField) SetHelp(value string) *FormField {
	f.Help = value
	return f
}

func (f *FormField) SetClass(value string) *FormField {
	f.Input.Class = value
	return f
}

func (f *FormField) SetMin(min interface{}) *FormField {
	f.Input.Min = min

	kind := reflect.ValueOf(min).Kind()
	var v Validator

	if kind == reflect.Int {
		v = MinValidator(min.(int))
	} else if kind == reflect.Int8 {
		v = MinValidator(min.(int8))
	} else if kind == reflect.Int16 {
		v = MinValidator(min.(int16))
	} else if kind == reflect.Int32 {
		v = MinValidator(min.(int32))
	} else if kind == reflect.Int64 {
		v = MinValidator(min.(int64))
	} else if kind == reflect.Uint {
		v = MinValidator(min.(uint))
	} else if kind == reflect.Uint8 {
		v = MinValidator(min.(uint8))
	} else if kind == reflect.Uint16 {
		v = MinValidator(min.(int16))
	} else if kind == reflect.Uint32 {
		v = MinValidator(min.(int32))
	} else if kind == reflect.Uint64 {
		v = MinValidator(min.(int64))
	} else if kind == reflect.Float32 {
		v = MinValidator(min.(float32))
	} else if kind == reflect.Float64 {
		v = MinValidator(min.(float64))
	} else {
		panic(fmt.Sprintf("Unable to handle type: %s", kind))
	}

	f.RemoveValidator("min").AddValidator(v)

	return f
}

func (f *FormField) SetMax(max interface{}) *FormField {
	f.Input.Max = max

	kind := reflect.ValueOf(max).Kind()
	var v Validator

	if kind == reflect.Int {
		v = MaxValidator(max.(int))
	} else if kind == reflect.Int8 {
		v = MaxValidator(max.(int8))
	} else if kind == reflect.Int16 {
		v = MaxValidator(max.(int16))
	} else if kind == reflect.Int32 {
		v = MaxValidator(max.(int32))
	} else if kind == reflect.Int64 {
		v = MaxValidator(max.(int64))
	} else if kind == reflect.Uint {
		v = MaxValidator(max.(uint))
	} else if kind == reflect.Uint8 {
		v = MaxValidator(max.(uint8))
	} else if kind == reflect.Uint16 {
		v = MaxValidator(max.(int16))
	} else if kind == reflect.Uint32 {
		v = MaxValidator(max.(int32))
	} else if kind == reflect.Uint64 {
		v = MaxValidator(max.(int64))
	} else if kind == reflect.Float32 {
		v = MaxValidator(max.(float32))
	} else if kind == reflect.Float64 {
		v = MaxValidator(max.(float64))
	} else {
		panic(fmt.Sprintf("Unable to handle type: %s", kind))
	}

	f.RemoveValidator("max").AddValidator(v)

	return f
}

func (f *FormField) SetStep(value int) *FormField {
	f.Input.Step = value
	return f
}

func (f *FormField) SetMaxLength(value int) *FormField {
	f.Input.MaxLength = value
	return f
}

func (f *FormField) SetMinLength(value int) *FormField {
	f.Input.MinLength = value
	return f
}

func (f *FormField) SetSize(value int) *FormField {
	f.Input.Size = value
	return f
}

func (f *FormField) SetHeight(value int) *FormField {
	f.Input.Height = value
	return f
}

func (f *FormField) SetWidth(value int) *FormField {
	f.Input.Width = value
	return f
}

func (f *FormField) SetNovalidation(value bool) *FormField {
	f.Input.Novalidate = value
	return f
}

func (f *FormField) SetAutofocus(value bool) *FormField {
	f.Input.Autofocus = value
	return f
}

func (f *FormField) SetRequired(value bool) *FormField {
	f.Input.Required = value
	return f
}

func (f *FormField) SetMultiple(value bool) *FormField {
	f.Input.Multiple = value
	return f
}

func (f *FormField) SetChecked(value bool) *FormField {
	f.Input.Checked = value
	return f
}

func (f *FormField) SetReadonly(value bool) *FormField {
	f.Input.Readonly = value
	return f
}

func (f *FormField) SetPattern(value string) *FormField {
	f.Input.Pattern = value
	return f
}

func (f *FormField) SetType(value string) *FormField {
	f.Input.Type = value
	return f
}
func (f *FormField) SetPlaceholder(value string) *FormField {
	f.Input.Placeholder = value
	return f
}

func (f *FormField) SetTemplate(value string) *FormField {
	f.Input.Template = value
	return f
}

func (f *FormField) SetList(value string) *FormField {
	f.Input.List = value
	return f
}

func (f *FormField) SetAutocomplete(value string) *FormField {
	f.Input.Autocomplete = value
	return f
}

func (f *Form) Add(name string, fieldType string, options ...interface{}) *FormField {
	var value interface{} = nil

	if len(options) > 0 {
		value = options[0]
	}

	field := create(name, fieldType, value)

	f.Fields = append(f.Fields, field)

	return field
}

func PrepareForm(form *Form) error {
	iterateFields(form, form.Fields)

	return nil
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

		if !field.Touched {
			// fmt.Printf("Field not touched: %s, skipping\n", field.Name)
			continue
		}

		if field.reflect.Kind() == reflect.Invalid {
			// fmt.Printf("Invalid type: %s\n", field.Name)
			continue
		}

		newValue := reflect.ValueOf(field.SubmitedValue)

		if newValue.CanConvert(field.reflect.Type()) {
			field.reflect.Set(newValue.Convert(field.reflect.Type()))
		} else {
			fmt.Printf("Unable to convert type: Type, field: %s (kind: %s) submitted: %s, value: %s\n", field.Name, field.reflect.Kind(), newValue.Kind(), newValue.Interface())
		}
	}
}

func yes(value string) bool {
	return value == "checked" || value == "true" || value == "1" || value == "on" || value == "yes"
}
