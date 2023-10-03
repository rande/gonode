// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"

	"github.com/flosch/pongo2"
	"github.com/rande/gonode/core/helper"
)

var (
	ErrNoTemplate = errors.New("unable to find the template to render")
)

type AttributOption struct {
	Name  string
	Value interface{}
}

func createPongoField(pongo *pongo2.TemplateSet) func(name string, form *Form) *pongo2.Value {

	return func(name string, form *Form) *pongo2.Value {
		tpl, err := pongo.FromFile("form:field.tpl")

		field := form.Get(name)

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

func createPongoLabel(pongo *pongo2.TemplateSet) func(name string, form *Form) *pongo2.Value {

	return func(name string, form *Form) *pongo2.Value {

		field := form.Get(name)

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

func createPongoInput(pongo *pongo2.TemplateSet) func(name string, form *Form) *pongo2.Value {

	return func(name string, form *Form) *pongo2.Value {

		field := form.Get(name)

		templates := []string{
			fmt.Sprintf("%s:fields/input.%s.tpl", field.Module, field.Input.Type),
			"form:fields/input.base.tpl",
		}

		for _, path := range templates {
			tpl, err := pongo.FromFile(path)

			if err == nil {
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

		helper.PanicOnError(ErrNoTemplate)

		return nil
	}
}

func createPongoErrors(pongo *pongo2.TemplateSet) func(name string, form *Form) *pongo2.Value {

	return func(name string, form *Form) *pongo2.Value {

		field := form.Get(name)

		tpl, err := pongo.FromFile("form:errors.tpl")

		helper.PanicOnError(err)

		data, err := tpl.Execute(pongo2.Context{
			"form":      form,
			"field":     field,
			"input":     field.Input,
			"label":     field.Label,
			"error":     field.Errors,
			"hasErrors": len(field.Errors) > 0,
		})

		helper.PanicOnError(err)

		return pongo2.AsSafeValue(data)
	}
}

func createPongoHelp(pongo *pongo2.TemplateSet) func(name string, form *Form) *pongo2.Value {

	return func(name string, form *Form) *pongo2.Value {

		field := form.Get(name)

		tpl, err := pongo.FromFile("form:help.tpl")

		helper.PanicOnError(err)

		data, err := tpl.Execute(pongo2.Context{
			"form":  form,
			"field": field,
			"help":  field.Help,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return pongo2.AsSafeValue(data)
	}
}
