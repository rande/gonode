// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"

	"github.com/rande/gonode/core/embed"
)

func TestMain(t *testing.M) {
	v := t.Run()

	// After all tests have run `go-snaps` can check for not used snapshots
	snaps.Clean(t)

	os.Exit(v)
}

func GetTemplate() *embed.TemplateLoader {
	embeds := embed.NewEmbeds()
	embeds.Add("form", GetEmbedFS())

	loader := &embed.TemplateLoader{
		Embeds:   embeds,
		BasePath: "",
	}

	funcMap := map[string]interface{}{
		"form_field":  createTemplateField(loader),
		"form_label":  createTemplateLabel(loader),
		"form_input":  createTemplateInput(loader),
		"form_help":   createTemplateHelp(loader),
		"form_errors": createTemplateErrors(loader),
	}

	loader.Templates = embed.GetTemplates(embeds, funcMap)

	return loader
}

func Test_Form_Rendering(t *testing.T) {
	now := time.Date(2022, time.April, 1, 1, 1, 1, 1, time.UTC)

	form := CreateForm(nil)
	form.Action = "/update"

	form.Add("name", "text", "John Doe")
	form.Add("email", "email", "john.doe@gmail.com")
	form.Add("date", "date", now)

	PrepareForm(form)

	assert.False(t, form.HasErrors)

	loader := GetTemplate()

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "John Doe", form.Get("name").Input.Value)

	form.Get("name").Input.Pattern = "^[a-z]+$"
	form.Get("name").Input.Placeholder = "Enter the name"
	form.Get("name").Input.Readonly = true
	form.Get("name").Input.Required = true
	form.Get("name").Input.Size = 10
	form.Get("name").Input.Autofocus = true
	form.Get("name").Input.Autocomplete = "on"
	form.Get("name").Input.Min = 10
	form.Get("name").Input.Max = 100
	form.Get("name").Input.Step = 10
	form.Get("name").Input.MinLength = 10
	form.Get("name").Input.MaxLength = 100

	data, err := loader.Execute("form:form/form", embed.Context{
		"form": form,
		"foo":  "bar",
	})

	assert.Nil(t, err)

	snaps.MatchSnapshot(t, string(data))
}

func Test_Form_Rendering_Error(t *testing.T) {
	form := CreateForm(nil)
	form.Add("position", "number", 1).
		AddValidators(RequiredValidator(), EmailValidator()).
		SetHelp("The position")

	PrepareForm(form)

	loader := GetTemplate()

	// -- Render form
	data, err := loader.Execute("form:form/form", embed.Context{
		"form": form,
	})

	assert.Nil(t, err)
	snaps.MatchSnapshot(t, string(data))

	// -- Bind form with request values
	v := url.Values{
		"position": []string{"foo"},
	}

	BindUrlValues(form, v)

	// -- validate form
	result := ValidateForm(form)

	assert.NotNil(t, result)

	// -- render form with errors
	data, err = loader.Execute("form:form/form", embed.Context{
		"form": form,
	})

	assert.Nil(t, err)

	snaps.MatchSnapshot(t, string(data))
}
