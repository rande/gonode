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
	Position int8
	Ratio    float32
	DOB      time.Time
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

	form := &Form{}
	form.Add("name", "text", "John Doe")

	PrepareForm(form)

	v := url.Values{
		"name": []string{"Thomas"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmitedValue)
}

func Test_Bind_Form_Basic_Errors(t *testing.T) {

	form := &Form{}
	form.Add("name", "text", "John Doe")

	PrepareForm(form)

	v := url.Values{
		"name": []string{"Thomas"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmitedValue)

	ValidateForm(form)
}

func Test_FormField_Validate(t *testing.T) {
	form := CreateForm(nil)
	field := form.Add("name", "email").AddValidators(RequiredValidator(), EmailValidator())

	field.Touched = true
	field.SubmitedValue = "john.doe@example.com"

	result := validateForm(form.Fields, form)

	assert.True(t, result)

	field.Touched = false
	field.SubmitedValue = nil

	result = validateForm(form.Fields, form)

	assert.False(t, result)

	assert.Equal(t, 2, len(field.Errors))
	assert.Equal(t, ErrValidatorRequired.Error(), field.Errors[0])
	assert.Equal(t, ErrValidatorEmail.Error(), field.Errors[1])
}

func Test_FormField_Validate_TypeMismatch(t *testing.T) {
	form := CreateForm(nil)
	field := form.Add("name", "int", "hello").AddValidators(RequiredValidator(), EmailValidator())

	PrepareForm(form)

	v := url.Values{
		"name": []string{"foo"},
	}

	BindUrlValues(form, v)

	result := ValidateForm(form)

	assert.False(t, result)

	assert.Equal(t, 2, len(field.Errors))
	assert.Equal(t, ErrInvalidType.Error(), field.Errors[0])
	assert.Equal(t, ErrValidatorEmail.Error(), field.Errors[1])
}

func Test_Bind_Form_Basic_Struct(t *testing.T) {
	user := &TestUser{
		Name:    "John Doe",
		Enabled: true,
		Hidden:  false,
	}

	form := CreateForm(user)
	form.Add("Name", "text")

	PrepareForm(form)

	v := url.Values{
		"Name": []string{"Thomas"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("Name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("Name").SubmitedValue)

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
	form.Add("options", "checkbox", CheckboxOptions{
		"enabled": {Label: "Enabled", Checked: true},
		"hidden":  {Label: "Hidden", Checked: false},
	})

	subForm := CreateForm(nil)
	subForm.Add("title", "text", "The title")
	subForm.Add("Body", "text", "The body")
	subForm.Add("options", "checkbox", CheckboxOptions{
		"validated": {Label: "Is Validated", Checked: true},
	})

	form.Add("post", "form", subForm)

	PrepareForm(form)

	assert.Equal(t, "name", form.Get("name").Input.Id)
	assert.Equal(t, "name", form.Get("name").Name)

	assert.Equal(t, "post_title", form.Get("post").Get("title").Input.Id)
	assert.Equal(t, "post.title", form.Get("post").Get("title").Input.Name)

	assert.NotNil(t, form.Get("post"))
	assert.NotNil(t, form.Get("post").Get("options"))
	assert.NotNil(t, form.Get("post").Get("options").Get("validated"))

	assert.Equal(t, "post.options[validated]", form.Get("post").Get("options").Get("validated").Input.Name)
	assert.Equal(t, "post_options_validated", form.Get("post").Get("options").Get("validated").Input.Id)

	assert.Equal(t, "post.options", form.Get("post").Get("options").Input.Name)
	assert.Equal(t, "post_options", form.Get("post").Get("options").Input.Id)

	v := url.Values{
		"name":                []string{"Thomas"},
		"options[enabled]":    []string{"false"},
		"options[hidden]":     []string{"true"},
		"post.title":          []string{"le titre"},
		"post.body":           []string{"le corps du texte"},
		"post.options[admin]": []string{"true"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("name").SubmitedValue)

	assert.Equal(t, "le titre", form.Get("post").Get("title").SubmitedValue)
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
	form.Add("Position", "int")
	form.Add("Ratio", "float")
	form.Add("DOB", "date")

	// add a field not linked an entity
	subForm := CreateForm(blog)
	subForm.Add("Title", "text")
	subForm.Add("Body", "text")
	subForm.Add("IsValidated", "boolean")
	subForm.Add("options", "checkbox", CheckboxOptions{
		"enabled": {Label: "Enabled", Checked: true},
		"hidden":  {Label: "Hidden", Checked: false},
	})

	form.Add("post", "form", subForm)

	PrepareForm(form)

	assert.Equal(t, "Name", form.Get("Name").Input.Id)
	assert.Equal(t, "Name", form.Get("Name").Name)

	assert.Equal(t, "post_Title", form.Get("post").Get("Title").Input.Id)
	assert.Equal(t, "post.Title", form.Get("post").Get("Title").Input.Name)

	assert.NotNil(t, form.Get("post"))
	assert.NotNil(t, form.Get("post").Get("options"))
	assert.NotNil(t, form.Get("post").Get("options").Get("enabled"))

	assert.Equal(t, "post.options[enabled]", form.Get("post").Get("options").Get("enabled").Input.Name)
	assert.Equal(t, "post_options_enabled", form.Get("post").Get("options").Get("enabled").Input.Id)

	assert.Equal(t, "post.options", form.Get("post").Get("options").Input.Name)
	assert.Equal(t, "post_options", form.Get("post").Get("options").Input.Id)

	v := url.Values{
		"Name":                  []string{"Thomas"},
		"Enabled":               []string{"no"},
		"Position":              []string{"12"},
		"Ratio":                 []string{"1.2"},
		"post.Title":            []string{"New title"},
		"post.Body":             []string{"New Body"},
		"post.options[enabled]": []string{"false"},
		"post.options[hidden]":  []string{"true"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("Name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("Name").SubmitedValue)

	assert.Equal(t, "New title", form.Get("post").Get("Title").SubmitedValue)

	if v, ok := form.Get("post").Get("options").SubmitedValue.(CheckboxOptions); ok {
		assert.Equal(t, false, v["enabled"].Checked)
		assert.Equal(t, true, v["hidden"].Checked)
	} else {
		t.Error("options is not a CheckboxOptions")
	}

	AttachValues(form)

	assert.Equal(t, "Thomas", user.Name)
	assert.Equal(t, false, user.Enabled)
	assert.Equal(t, float32(1.2), user.Ratio)
	assert.Equal(t, int8(12), user.Position)
	assert.Equal(t, "New title", blog.Title)
	assert.Equal(t, true, blog.IsValidated) // not submitted
}

func Test_Bind_Form_Collection(t *testing.T) {
	values := []*FieldCollectionValue{
		{Key: "0", Value: &TestTag{Id: 1, Name: "tag1", Enabled: true}},
		{Key: "1", Value: &TestTag{Id: 1, Name: "tag2", Enabled: true}},
	}

	form := &Form{}
	form.Add("tags", "collection", &FieldCollectionOptions{
		Items: values,
		Configure: func(value interface{}) *Form {
			tag := value.(*TestTag)

			tagForm := &Form{}
			tagForm.Add("name", "text", tag.Name)
			tagForm.Add("options", "checkbox", CheckboxOptions{
				"enabled": {Label: "Is Enabled", Checked: tag.Enabled},
			})

			return tagForm
		},
	})

	PrepareForm(form)

	assert.Equal(t, form.Get("tags").Get("0").Get("name").Input.Name, "tags[0].name")
	assert.Equal(t, form.Get("tags").Get("0").Get("name").InitialValue, "tag1")

	v := url.Values{
		"tags[0].name":             []string{"TAG 1"},
		"tags[0].options[enabled]": []string{"false"},
		"tags[1].name":             []string{"TAG 2"},
		"tags[1].options[enabled]": []string{"false"},
	}

	BindUrlValues(form, v)

	assert.NotNil(t, form.Get("tags").Get("0").Get("name").SubmitedValue, "TAG 1")
}