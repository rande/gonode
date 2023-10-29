// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"testing"

	"github.com/flosch/pongo2"
	"github.com/rande/gonode/core/embed/fixtures"
	"github.com/stretchr/testify/assert"
)

// func Test_PongoTemplateLoader_Valid_Template(t *testing.T) {
// 	assets := NewAssets()
// 	assets.Add("testmodule", GetTestEmbedFS())

// 	loader := &PongoTemplateLoader{
// 		Assets: assets,
// 		BasePath: "fixtures/",
// 	}

// 	path := loader.Abs("testmodule:/", "testmodule:file.tpl")

// 	assert.Equal(t, "vfs/file.tpl", path)

// 	// pongo add "/"
// 	r, err := loader.Get("vfs/file.tpl")

// 	assert.NoError(t, err)

// 	data := make([]byte, len("Content"))

// 	r.Read(data)

// 	assert.Equal(t, []byte("Content"), data)
// }

func Test_PongoTemplateLoader_Integration(t *testing.T) {
	embeds := NewEmbeds()
	embeds.Add("testmodule", fixtures.GetTestEmbedFS())

	loader := &PongoTemplateLoader{
		Embeds: embeds,
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	tpl, err := pongo.FromFile("testmodule:file.tpl")

	assert.NoError(t, err)

	c := pongo2.Context{
		"k": struct {
			v string
		}{
			v: "foobar",
		},
	}

	result, err := tpl.Execute(c)

	assert.NoError(t, err)

	assert.Equal(t, "Value: foobar", result)
}

func Test_PongoTemplateLoader_Integration_NotFound(t *testing.T) {
	embeds := NewEmbeds()
	embeds.Add("testmodule", fixtures.GetTestEmbedFS())

	loader := &PongoTemplateLoader{
		Embeds: embeds,
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	_, err := pongo.FromFile("testmodule:file_does_not_exist.tpl")

	assert.Error(t, err)
}

func Test_PongoTemplateLoader_Integration_InvalidSyntax(t *testing.T) {

	embeds := NewEmbeds()
	embeds.Add("testmodule", fixtures.GetTestEmbedFS())

	loader := &PongoTemplateLoader{
		Embeds: embeds,
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	_, err := pongo.FromFile("testmodule:invalid_syntax.tpl")

	assert.Error(t, err)
}
