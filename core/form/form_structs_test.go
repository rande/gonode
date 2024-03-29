// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	Name     string
	Enabled  bool
	Hidden   bool
	Email    string
	Position int32
	Ratio    float32
	DOB      time.Time
	Items    []int32
	Tags     []string
}

type TestBlogPost struct {
	Title       string
	IsValidated bool
	Body        string
}

type TestTag struct {
	Id      int
	Name    string
	Enabled bool
}

func Test_Create_Form_Empty(t *testing.T) {
	form := CreateForm(nil)

	assert.NotNil(t, form)
}

func Test_Form_Init(t *testing.T) {
	user := &TestUser{
		Name:    "John Doe",
		Enabled: true,
		Hidden:  false,
	}

	form := CreateForm(user)
	form.Add("Name", "text")

	assert.False(t, form.HasErrors)
	assert.False(t, form.Get("Name").HasErrors)
	assert.False(t, form.Get("Name").Touched)
	assert.False(t, form.Get("Name").Submitted)
	assert.True(t, form.Get("Name").Mandatory)

	PrepareForm(form)

	assert.Equal(t, form.Get("Name").InitialValue, "John Doe")
}

func Test_Bind_Form_Basic(t *testing.T) {

	form := CreateForm(nil)
	form.Add("name", "text", "John Doe")

	err := PrepareForm(form)
	assert.Nil(t, err)

	v := url.Values{
		"name": []string{"Thomas"},
	}

	err = BindUrlValues(form, v)
	assert.Nil(t, err)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmittedValue)
}

func Test_Bind_Form_Basic_Errors(t *testing.T) {

	form := CreateForm(nil)
	form.Add("name", "text", "John Doe")

	PrepareForm(form)

	v := url.Values{
		"name": []string{"Thomas"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmittedValue)

	ValidateForm(form)
}

func Test_FormField_Validate(t *testing.T) {
	form := CreateForm(nil)
	field := form.Add("name", "email").AddValidators(RequiredValidator(), EmailValidator())

	field.Touched = true
	field.SubmittedValue = "john.doe@example.com"

	result := validateForm(form.Fields, form)

	assert.Nil(t, result)

	field.Touched = false
	field.SubmittedValue = "john.doe"

	result = validateForm(form.Fields, form)

	assert.NotNil(t, result)

	assert.Equal(t, 2, len(field.Errors))
	assert.Equal(t, ErrRequiredValidator.Error(), field.Errors[0])
	assert.Equal(t, ErrEmailValidator.Error(), field.Errors[1])
}

func Test_FormField_Validate_MinMax(t *testing.T) {
	form := CreateForm(nil)
	field := form.Add("age", "number", 30).SetMin(18).SetMax(48).SetMin(20)

	field.Touched = true
	field.SubmittedValue = 22

	result := validateForm(form.Fields, form)
	assert.Nil(t, result)

	field.SubmittedValue = 19
	result = validateForm(form.Fields, form)
	assert.NotNil(t, result)
}

func Test_FormField_Validate_TypeMismatch(t *testing.T) {
	form := CreateForm(nil)
	field := form.Add("name", "number", 2).AddValidators(RequiredValidator(), EmailValidator())

	PrepareForm(form)

	v := url.Values{
		"name": []string{"foo"},
	}

	BindUrlValues(form, v)

	result := ValidateForm(form)

	assert.NotNil(t, result)

	assert.Equal(t, 2, len(field.Errors))
	assert.Equal(t, ErrInvalidType.Error(), field.Errors[0])
	assert.Equal(t, ErrEmailValidator.Error(), field.Errors[1])
}

func Test_Bind_Form_Basic_Struct(t *testing.T) {
	user := &TestUser{
		Name:    "John Doe",
		Enabled: true,
		Hidden:  false,
	}

	form := CreateForm(user)
	form.Add("Name")

	PrepareForm(form)

	v := url.Values{
		"Name": []string{"Thomas"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("Name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("Name").SubmittedValue)

	ValidateForm(form)
	AttachValues(form)

	assert.Equal(t, "Thomas", user.Name)
}

func Test_Reflect(t *testing.T) {

	user := &TestUser{
		Name:    "Old Name",
		Enabled: true,
		Hidden:  false,
	}

	type Field struct {
		Name      string
		Value     interface{}
		Submitted interface{}
		reflect   reflect.Value
	}

	name := Field{
		Name:      "Name",
		Value:     nil,
		Submitted: nil,
	}

	enabled := Field{
		Name:      "Enabled",
		Value:     nil,
		Submitted: nil,
	}

	v := reflect.ValueOf(user).Elem()

	// Name
	name.reflect = v.FieldByName("Name")
	name.Value = name.reflect.Interface()
	name.Submitted = "New Name"
	name.reflect.Set(reflect.ValueOf(name.Submitted))

	// Enabled
	enabled.reflect = v.FieldByName("Enabled")
	enabled.Value = enabled.reflect.Interface()
	enabled.Submitted = true
	enabled.reflect.Set(reflect.ValueOf(enabled.Submitted))

	assert.Equal(t, user.Name, "New Name")
	assert.Equal(t, name.Value, "Old Name")
	assert.Equal(t, user.Enabled, true)
	assert.Equal(t, enabled.Value, true)
}

func Test_Bind_Form_Nested_Basic(t *testing.T) {
	form := CreateForm(nil)
	form.Add("name", "text", "John Doe")
	form.Add("options", "checkbox", nil, FieldOptions{
		{Label: "Enabled", Checked: true, Value: "enabled"},
		{Label: "Hidden", Checked: false, Value: "hidden"},
	})

	subForm := CreateForm(nil)
	subForm.Add("title", "text", "The title")
	subForm.Add("Body", "text", "The body")
	subForm.Add("options", "checkbox", nil, FieldOptions{
		{Label: "Is Validated", Checked: true, Value: "validated"},
	})

	form.Add("post", "form", subForm)

	PrepareForm(form)

	assert.Equal(t, "name", form.Get("name").Input.Id)
	assert.Equal(t, "name", form.Get("name").Name)

	assert.Equal(t, "post_title", form.Get("post").Get("title").Input.Id)
	assert.Equal(t, "post.title", form.Get("post").Get("title").Input.Name)

	assert.NotNil(t, form.Get("post"))
	assert.NotNil(t, form.Get("post").Get("options"))
	assert.NotNil(t, form.Get("post").Get("options").Get("0"))

	assert.Equal(t, "post.options[0]", form.Get("post").Get("options").Get("0").Input.Name)
	assert.Equal(t, "post_options_0", form.Get("post").Get("options").Get("0").Input.Id)

	assert.Equal(t, "post.options", form.Get("post").Get("options").Input.Name)
	assert.Equal(t, "post_options", form.Get("post").Get("options").Input.Id)

	v := url.Values{
		"name":                []string{"Thomas"},
		"options[0]":          []string{"false"},
		"options[1]":          []string{"true"},
		"post.title":          []string{"le titre"},
		"post.body":           []string{"le corps du texte"},
		"post.options[admin]": []string{"true"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmittedValue)

	assert.Equal(t, "le titre", form.Get("post").Get("title").SubmittedValue)
}

func Test_Bind_Form_Nested_Basic_Struct(t *testing.T) {
	user := &TestUser{
		Name:     "John Doe",
		Enabled:  true,
		Hidden:   false,
		Position: 1,
		Ratio:    0.2,
		DOB:      time.Date(1981, time.November, 30, 0, 0, 0, 0, time.UTC),
	}

	blog := &TestBlogPost{
		Title:       "Old title",
		IsValidated: true,
		Body:        "Old body",
	}

	form := CreateForm(user)
	form.Add("Name", "text")
	form.Add("Enabled", "bool")
	form.Add("Position", "number")
	form.Add("Ratio", "number")
	form.Add("DOB", "date")

	// add a field not linked an entity
	subForm := CreateForm(blog)
	subForm.Add("Title", "text")
	subForm.Add("Body", "text")
	subForm.Add("IsValidated", "boolean")
	subForm.Add("options", "checkbox", nil, FieldOptions{
		{Label: "Enabled", Checked: true, Value: "enabled"},
		{Label: "Hidden", Checked: false, Value: "hidden"},
	})

	form.Add("post", "form", subForm)

	PrepareForm(form)

	assert.Equal(t, "Name", form.Get("Name").Input.Id)
	assert.Equal(t, "Name", form.Get("Name").Name)

	assert.Equal(t, "post_Title", form.Get("post").Get("Title").Input.Id)
	assert.Equal(t, "post.Title", form.Get("post").Get("Title").Input.Name)

	assert.NotNil(t, form.Get("post"))
	assert.NotNil(t, form.Get("post").Get("options"))
	assert.NotNil(t, form.Get("post").Get("options").Get("0"))

	assert.Equal(t, "post.options[0]", form.Get("post").Get("options").Get("0").Input.Name)
	assert.Equal(t, "post_options_0", form.Get("post").Get("options").Get("0").Input.Id)

	assert.Equal(t, "post.options", form.Get("post").Get("options").Input.Name)
	assert.Equal(t, "post_options", form.Get("post").Get("options").Input.Id)

	v := url.Values{
		"Name":            []string{"Thomas"},
		"Enabled":         []string{"no"},
		"Position":        []string{"12"},
		"Ratio":           []string{"1.2"},
		"post.Title":      []string{"New title"},
		"post.Body":       []string{"New Body"},
		"post.options[0]": []string{"false"},
		"post.options[1]": []string{"true"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("Name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("Name").SubmittedValue)

	assert.Equal(t, "New title", form.Get("post").Get("Title").SubmittedValue)

	if v, ok := form.Get("post").Get("options").SubmittedValue.(FieldOptions); ok {
		assert.Equal(t, false, v[0].Checked)
		assert.Equal(t, true, v[1].Checked)
	} else {
		t.Error("options is not a FieldOptions")
	}

	ValidateForm(form)
	err := AttachValues(form)

	assert.Nil(t, err)

	assert.Equal(t, "Thomas", user.Name)
	assert.Equal(t, false, user.Enabled)
	assert.Equal(t, float32(1.2), user.Ratio)
	assert.Equal(t, int32(12), user.Position)
	assert.Equal(t, "New title", blog.Title)
	assert.Equal(t, true, blog.IsValidated) // not submitted
}

func Test_Bind_Form_Collection(t *testing.T) {
	values := []*FieldCollectionValue{
		{Key: "0", Value: &TestTag{Id: 1, Name: "tag1", Enabled: true}},
		{Key: "1", Value: &TestTag{Id: 1, Name: "tag2", Enabled: true}},
	}

	form := CreateForm(nil)
	form.Add("tags", "collection", nil, &FieldCollectionOptions{
		Items: values,
		Configure: func(value interface{}) *Form {
			tag := value.(*TestTag)

			tagForm := CreateForm(nil)
			tagForm.Add("name", "text", tag.Name)
			tagForm.Add("options", "checkbox", nil, FieldOptions{
				{Label: "Is Enabled", Checked: tag.Enabled, Value: "enabled"},
			})

			return tagForm
		},
	})

	PrepareForm(form)

	assert.Equal(t, form.Get("tags").Get("0").Get("name").Input.Name, "tags[0].name")
	assert.Equal(t, form.Get("tags").Get("0").Get("name").InitialValue, "tag1")

	v := url.Values{
		"tags[0].name":       []string{"TAG 1"},
		"tags[0].options[0]": []string{"false"},
		"tags[1].name":       []string{"TAG 2"},
		"tags[1].options[0]": []string{"false"},
	}

	BindUrlValues(form, v)

	assert.NotNil(t, form.Get("tags").Get("0").Get("name").SubmittedValue, "TAG 1")
}

func Test_Bind_Form_Select(t *testing.T) {
	user := &TestUser{
		Name:     "John Doe",
		Enabled:  true,
		Hidden:   false,
		Position: int32(1),
		Items:    []int32{1, 2},
	}

	form := CreateForm(user)
	form.Add("Enabled", "select", nil, FieldOptions{
		{Label: "No", Value: false},
		{Label: "Yes", Value: true},
	})

	form.Add("Position", "select", nil, FieldOptions{
		{Label: "1", Value: int32(1)},
		{Label: "2", Value: int32(2)},
		{Label: "3", Value: int32(3)},
		{Label: "4", Value: int32(4)},
	})

	form.Add("Items", "select", nil, FieldOptions{
		{Label: "Food", Value: int32(1)},
		{Label: "Car", Value: int32(2)},
		{Label: "Travel", Value: int32(3)},
		{Label: "Games", Value: int32(4)},
	})

	PrepareForm(form)

	assert.True(t, form.Get("Items").Input.Multiple)
	assert.True(t, form.Get("Items").Get("0").Input.Checked)
	assert.True(t, form.Get("Items").Get("1").Input.Checked)
	assert.False(t, form.Get("Items").Get("2").Input.Checked)
	assert.False(t, form.Get("Items").Get("3").Input.Checked)

	assert.Equal(t, form.Get("Items").Get("3").Input.Name, "Items[3]")

	v := url.Values{
		"Enabled":  []string{"0"},
		"Position": []string{"3"},
		"Items":    []string{"1", "3"},
	}

	err := BindUrlValues(form, v)
	assert.Nil(t, err)

	ValidateForm(form)

	err = AttachValues(form)
	assert.Nil(t, err)

	assert.Equal(t, false, user.Enabled)
	assert.Equal(t, int32(3), user.Position)
	assert.Equal(t, []int32{1, 3}, user.Items)
}

func Test_Bind_Form_Select_Invalid_Multiple_Type(t *testing.T) {
	user := &TestUser{
		Name:     "John Doe",
		Enabled:  true,
		Hidden:   false,
		Position: int32(1),
		Items:    []int32{1, 2},
	}

	form := CreateForm(user)
	form.Add("Items", "select", nil, FieldOptions{
		{Label: "Food", Value: int32(1)},
		{Label: "Car", Value: int32(2)},
		{Label: "Travel", Value: int32(3)},
		{Label: "Games", Value: int32(4)},
	})

	PrepareForm(form)

	v := url.Values{
		"Items": []string{"foo", "3"},
	}

	err := BindUrlValues(form, v)
	assert.NotNil(t, err)

	assert.Equal(t, 1, len(form.Get("Items").Errors))
	assert.Equal(t, "unable to convert value to the correct type", form.Get("Items").Errors[0])

	ValidateForm(form)

	err = AttachValues(form)

	assert.NotNil(t, err)
}

func Test_Bind_Form_Select_Invalid_Multiple_Type_Scope(t *testing.T) {
	user := &TestUser{
		Name:     "John Doe",
		Enabled:  true,
		Hidden:   false,
		Position: int32(1),
		Items:    []int32{1, 2},
	}

	form := CreateForm(user)
	form.Add("Items", "select", nil, FieldOptions{
		{Label: "Food", Value: int32(1)},
		{Label: "Car", Value: int32(2)},
		{Label: "Travel", Value: int32(3)},
		{Label: "Games", Value: int32(4)},
	})

	PrepareForm(form)

	v := url.Values{
		"Items": []string{"5", "3"}, // 5 is not a valid value
	}

	err := BindUrlValues(form, v)
	assert.Nil(t, err)

	err = ValidateForm(form)

	assert.NotNil(t, err)
	assert.Equal(t, 1, len(form.Get("Items").Errors))
	assert.Equal(t, "the value is not in the list of options", form.Get("Items").Errors[0])

	err = AttachValues(form)

	assert.NotNil(t, err)
}

func Test_Bind_Form_Select_Invalid_Type(t *testing.T) {
	user := &TestUser{
		Name:     "John Doe",
		Enabled:  true,
		Hidden:   false,
		Position: int32(1),
		Items:    []int32{1, 2},
	}

	form := CreateForm(user)
	form.Add("Position", "select", nil, FieldOptions{
		{Label: "1", Value: int32(1)},
		{Label: "2", Value: int32(2)},
	})

	PrepareForm(form)

	v := url.Values{
		"Position": []string{"foo"},
	}

	BindUrlValues(form, v)

	assert.True(t, form.HasErrors)
	assert.Equal(t, 1, len(form.Get("Position").Errors))
	assert.Equal(t, int32(1), user.Position)
	assert.Equal(t, "value does not match the expected type", form.Get("Position").Errors[0])

	ValidateForm(form)

	err := AttachValues(form)

	assert.NotNil(t, err)
}

func Test_Form_Tags_As_String(t *testing.T) {
	user := &TestUser{
		Tags: []string{"foo", "bar"},
	}

	form := CreateForm(user)
	form.Add("Tags")

	PrepareForm(form)

	assert.Equal(t, "foo, bar", form.Get("Tags").Input.Value)

	v := url.Values{
		"Tags": []string{"boo, far"},
	}

	BindUrlValues(form, v)

	ValidateForm(form)

	AttachValues(form)

	assert.False(t, form.HasErrors)
	assert.Equal(t, []string{"boo", "far"}, user.Tags)
}
