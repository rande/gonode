// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Form_Init(t *testing.T) {

	form := &Form{}
	form.Add("name", "text", "John Doe")
	form.Add("email", "email", "john.doe@gmail.com")

	assert.False(t, form.HasErrors)
	assert.False(t, form.Get("name").HasErrors)
	assert.False(t, form.Get("name").Touched)
	assert.False(t, form.Get("name").Submitted)
	assert.True(t, form.Get("name").Mandatory)
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

func Test_Bind_Form_Nested_Basic(t *testing.T) {
	form := &Form{}
	form.Add("name", "text", "John Doe")
	form.Add("options", "checkbox", FieldOptions{
		"enabled": {Label: "Enabled", Checked: true},
		"hidden":  {Label: "Hidden", Checked: false},
	})

	subForm := &Form{}
	subForm.Add("title", "text", "The title")
	subForm.Add("body", "text", "The body")
	subForm.Add("options", "checkbox", FieldOptions{
		"admin": {Label: "Is Admin", Checked: false},
	})

	form.Add("post", "form", subForm)

	PrepareForm(form)

	assert.Equal(t, "name", form.Get("name").Input.Id)
	assert.Equal(t, "name", form.Get("name").Name)

	assert.Equal(t, "post_title", form.Get("post").Get("title").Input.Id)
	assert.Equal(t, "post.title", form.Get("post").Get("title").Input.Name)

	assert.NotNil(t, form.Get("post"))
	assert.NotNil(t, form.Get("post").Get("options"))
	assert.NotNil(t, form.Get("post").Get("options").Get("admin"))
	assert.Equal(t, "post.options[admin]", form.Get("post").Get("options").Get("admin").Input.Name)
	assert.Equal(t, "post_options_admin", form.Get("post").Get("options").Get("admin").Input.Id)

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

	assert.True(t, form.Get("post").Get("options").SubmitedValue.(FieldOptions)["admin"].Checked)
}
