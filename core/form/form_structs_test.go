// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestUser struct {
	Name     string
	Enabled  bool
	Hidden   bool
	Email    string
	Position int8
	Ratio    float32
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
	form.Add("options", "checkbox", FieldOptions{
		"enabled": {Label: "Enabled", Checked: true},
		"hidden":  {Label: "Hidden", Checked: false},
	})

	subForm := CreateForm(nil)
	subForm.Add("title", "text", "The title")
	subForm.Add("Body", "text", "The body")
	subForm.Add("options", "checkbox", FieldOptions{
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

	// add a field not linked an entity
	subForm := CreateForm(blog)
	subForm.Add("Title", "text")
	subForm.Add("Body", "text")
	subForm.Add("IsValidated", "boolean")
	subForm.Add("options", "checkbox", FieldOptions{
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
		"post.Title":            []string{"New title"},
		"post.Body":             []string{"New Body"},
		"post.options[enabled]": []string{"false"},
		"post.options[hidden]":  []string{"true"},
	}

	BindUrlValues(form, v)

	assert.Equal(t, "John Doe", form.Get("Name").InitialValue)
	assert.Equal(t, "Thomas", form.Get("Name").SubmitedValue)

	assert.Equal(t, "New title", form.Get("post").Get("Title").SubmitedValue)

	if v, ok := form.Get("post").Get("options").SubmitedValue.(FieldOptions); ok {
		assert.Equal(t, false, v["enabled"].Checked)
		assert.Equal(t, true, v["hidden"].Checked)
	} else {
		t.Error("options is not a FieldOptions")
	}

	AttachValues(form)

	assert.Equal(t, "Thomas", user.Name)
	assert.Equal(t, false, user.Enabled)
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
			tagForm.Add("options", "checkbox", FieldOptions{
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
