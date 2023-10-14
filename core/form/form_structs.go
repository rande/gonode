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
	ErrNoValue              = errors.New("unable to find the value")
	ErrInvalidSubmittedData = errors.New("invalid submitted data")
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
	Value   interface{}
	Id      string
	Checked bool
}

type FieldOptions []*FieldOption

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
	Size         uint16
	MinLength    uint32
	MaxLength    uint32
	Min          interface{}
	Max          interface{}
	Step         uint32
	Height       uint16
	Width        uint16
	Options      FieldOptions
}

type Label struct {
	Template string
	Class    string
	Style    string
	Value    string
}

type Error struct {
	Template string
	Class    string
	Style    string
	Value    string
}

type Marshaller func(field *FormField, form *Form) error

type Unmarshaller func(field *FormField, form *Form, values url.Values) error

type FormField struct {
	Prefix         string // used for nested forms
	Name           string
	Module         string
	Help           string
	Attributes     Attributes
	Label          Label
	Input          Input
	Error          Error
	Mandatory      bool
	InitialValue   interface{}
	SubmittedValue interface{}
	Children       []*FormField
	Touched        bool
	Submitted      bool
	Errors         []string
	HasErrors      bool
	Validators     []Validator
	// from go to serialized
	Marshaller Marshaller
	// from serialized to go
	Unmarshaller Unmarshaller
	// Validator
	reflect reflect.Value
	Options interface{}
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
			return field.SubmittedValue
		}
	}

	return nil
}

func (f *FormField) Add(name string, options ...interface{}) *FormField {
	field := create(name, options...)

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
			value = field.SubmittedValue.(T)
			break
		}
	}

	return value
}

func (f *Form) Value(name string) interface{} {
	for _, field := range f.Fields {
		if field.Name == name {
			return field.SubmittedValue
		}
	}

	return nil
}

func create(name string, options ...interface{}) *FormField {
	field := CreateFormField()

	field.Name = name
	field.Label.Value = name

	if len(options) > 0 {
		field.Input.Type = options[0].(string)
	}

	if len(options) > 1 {
		field.InitialValue = options[1]
	}

	if len(options) > 2 {
		field.Options = options[2]
	}

	field.Input.Class = "shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
	field.Label.Class = "block text-gray-700 text-sm font-bold mb-2"

	if field.Input.Type == "color" {
		field.Input.Class = ""
	}

	if field.Input.Type == "range" {
		field.Input.Class = ""
	}

	if field.Input.Type == "checkbox" {
		field.Marshaller = checkboxMarshal
		field.Unmarshaller = checkboxUnmarshal
	}

	if field.Input.Type == "select" {
		field.Marshaller = selectMarshal
		field.Unmarshaller = selectUnmarshal
	}

	// if fieldType == "form" {
	// 	field.Marshaller = formMarshal
	// 	field.Unmarshaller = formUnmarshal
	// }

	if field.Input.Type == "collection" {
		field.Marshaller = collectionMarshal
		field.Unmarshaller = collectionUnmarshal
	}

	if field.Input.Type == "boolean" {
		// field.Marshaller = booleanMarshal
		// field.Unmarshaller = booleanUnmarshal
		field.Input.Type = "checkbox"
		field.Input.Class = ""
	}

	// if fieldType == "number" {
	// 	field.Marshaller = numberMarshal
	// 	field.Unmarshaller = numberUnmarshal
	// }

	// if fieldType == "date" {
	// 	field.Marshaller = dateMarshal
	// 	field.Unmarshaller = dateUnmarshal
	// }

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
		Module:         "form",
		Help:           "",
		InitialValue:   nil,
		Mandatory:      true,
		SubmittedValue: nil,
		Children:       []*FormField{},
		Touched:        false,
		Submitted:      false,
		Errors:         []string{},
		Attributes:     Attributes{},
		Validators:     []Validator{},
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

func (f *FormField) SetStep(value uint32) *FormField {
	f.Input.Step = value
	return f
}

func (f *FormField) SetMaxLength(value uint32) *FormField {
	f.Input.MaxLength = value

	f.RemoveValidator("max_length").AddValidator(MaxLengthValidator(value, "bytes"))

	return f
}

func (f *FormField) SetMinLength(value uint32) *FormField {
	f.Input.MinLength = value

	f.RemoveValidator("min_length").AddValidator(MinLengthValidator(value, "bytes"))

	return f
}

func (f *FormField) SetSize(value uint16) *FormField {
	f.Input.Size = value
	return f
}

func (f *FormField) SetHeight(value uint16) *FormField {
	f.Input.Height = value
	return f
}

func (f *FormField) SetWidth(value uint16) *FormField {
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

func (f *FormField) SetMarshaller(marshaller Marshaller) *FormField {
	f.Marshaller = marshaller
	return f
}

func (f *FormField) SetUnmarshaller(unmarshaller Unmarshaller) *FormField {
	f.Unmarshaller = unmarshaller
	return f
}

func (f *Form) Add(name string, options ...interface{}) *FormField {
	field := create(name, options...)

	f.Fields = append(f.Fields, field)

	return field
}

// Iterate over all fields and call the marshaller function to transform the Go value into
// a serialized value used in the HTML form. This also setup the different attribures required
// by the HTML form.
func PrepareForm(form *Form) error {
	iterateFields(form, form.Fields)

	return nil
}

// Iterate over all fields and bind the submitted values to the SubmittedValue field
// in the form. This will not attach the value to the underlying data structure, use AttachValues
// for this.
//
// This will call all unmarshallers defined on each form field.
func BindUrlValues(form *Form, values url.Values) error {
	for _, field := range form.Fields {
		unmarshal(field, form, values)
	}

	if form.HasErrors {
		return ErrInvalidSubmittedData
	}

	return nil
}

// Iterate over all submitted values and assign them to the related data structure
// if possible. This will not work work if the form has some errors.
func AttachValues(form *Form) error {
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

func getValue(field *FormField, values url.Values) (string, error) {
	if !values.Has(field.Input.Name) {
		return "", ErrNoValue
	}

	return values.Get(field.Input.Name), nil
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
