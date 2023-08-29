// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"fmt"

	"github.com/flosch/pongo2"
	"github.com/rande/gonode/core/helper"
)

type AttributOption struct {
	Name  string
	Value interface{}
}

func createPongoField(pongo *pongo2.TemplateSet) func(field *FormField, form *Form) *pongo2.Value {

	return func(field *FormField, form *Form) *pongo2.Value {
		tpl, err := pongo.FromFile("form:field.tpl")

		helper.PanicOnError(err)

		data, err := tpl.Execute(pongo2.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return pongo2.AsSafeValue(data)
	}
}

func createPongoLabel(pongo *pongo2.TemplateSet) func(field *FormField, form *Form) *pongo2.Value {

	return func(field *FormField, form *Form) *pongo2.Value {

		tpl, err := pongo.FromFile(field.Label.Template)

		helper.PanicOnError(err)

		data, err := tpl.Execute(pongo2.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return pongo2.AsSafeValue(data)
	}
}

func createPongoInput(pongo *pongo2.TemplateSet) func(field *FormField, form *Form) *pongo2.Value {

	return func(field *FormField, form *Form) *pongo2.Value {

		tpl, err := pongo.FromFile(fmt.Sprintf("%s:fields/input.%s.tpl", field.Module, field.Input.Type))

		helper.PanicOnError(err)

		data, err := tpl.Execute(pongo2.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return pongo2.AsSafeValue(data)
	}
}
