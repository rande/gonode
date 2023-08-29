// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"fmt"
	"os"
	"testing"

	"github.com/flosch/pongo2"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"

	"github.com/rande/gonode/core/embed"
)

func TestMain(t *testing.M) {

	fmt.Println("Setup")
	v := t.Run()

	// After all tests have run `go-snaps` can check for not used snapshots
	snaps.Clean(t)

	os.Exit(v)
}

func GetPongo() *pongo2.TemplateSet {
	embeds := embed.NewEmbeds()
	embeds.Add("form", GetEmbedFS())

	pongo := pongo2.NewSet("gonode.embeds", &embed.PongoTemplateLoader{
		Embeds:   embeds,
		BasePath: "",
	})
	pongo.Options = &pongo2.Options{
		TrimBlocks:   true,
		LStripBlocks: true,
	}

	pongo.Globals["form_field"] = createPongoField(pongo)
	pongo.Globals["form_label"] = createPongoLabel(pongo)
	pongo.Globals["form_input"] = createPongoInput(pongo)

	return pongo
}
func Test_Form_Rendering(t *testing.T) {

	form := &Form{}
	form.Add("name", "text", "John Doe")
	form.Add("email", "email", "john.doe@gmail.com")

	PrepareForm(form)

	assert.False(t, form.HasErrors)

	pongo := GetPongo()
	template, err := pongo.FromFile("form:form.tpl")

	assert.Equal(t, "John Doe", form.Get("name").InitialValue)
	assert.Equal(t, "John Doe", form.Get("name").Input.Value)

	form.Get("name").Input.Pattern = "^[a-z]+$"
	form.Get("name").Input.Placeholder = "Enter the name"
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

	assert.Nil(t, err)
	assert.NotNil(t, template)

	html, err := template.Execute(pongo2.Context{
		"form": form,
	})

	// fmt.Printf("%s", err.Error())
	assert.Nil(t, err)

	snaps.MatchSnapshot(t, html)
}
