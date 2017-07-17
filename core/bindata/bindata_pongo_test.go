// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	"errors"
	"testing"

	"github.com/flosch/pongo2"
	"github.com/stretchr/testify/assert"
)

func Test_PongoTemplateLoader_Valid_Template(t *testing.T) {
	loader := &PongoTemplateLoader{
		Asset: func(name string) ([]byte, error) {

			assert.Equal(t, "vfs/file.tpl", name)

			return []byte("Content"), nil
		},
		Paths: []string{
			"vfs",
		},
	}

	path := loader.Abs("/", "file.tpl")

	assert.Equal(t, "vfs/file.tpl", path)

	// pongo add "/"
	r, err := loader.Get("vfs/file.tpl")

	assert.NoError(t, err)

	data := make([]byte, len("Content"))

	r.Read(data)

	assert.Equal(t, []byte("Content"), data)
}

func Test_PongoTemplateLoader_Integration(t *testing.T) {

	loader := &PongoTemplateLoader{
		Asset: func(name string) ([]byte, error) {

			assert.Equal(t, "vfs/file.tpl", name)

			return []byte("Value: {{ k.v }}"), nil
		},
		Paths: []string{
			"vfs",
		},
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	tpl, err := pongo.FromFile("file.tpl")

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
	loader := &PongoTemplateLoader{
		Asset: func(name string) ([]byte, error) {
			return nil, errors.New("template not found")
		},
		Paths: []string{},
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	_, err := pongo.FromFile("file.tpl")

	assert.Error(t, err)
}

func Test_PongoTemplateLoader_Integration_InvalidSyntax(t *testing.T) {

	loader := &PongoTemplateLoader{
		Asset: func(name string) ([]byte, error) {
			return []byte("Value: {% if true %}hellor{%else%}"), nil
		},
		Paths: []string{},
	}

	pongo := pongo2.NewSet("gonode_test", loader)
	_, err := pongo.FromFile("file.tpl")

	assert.Error(t, err)
}
