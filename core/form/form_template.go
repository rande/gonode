// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/core/helper"
)

var (
	ErrNoTemplate = errors.New("unable to find the template to render")
)

func createTemplateField(loader *embed.TemplateLoader) func(name string, form *Form) template.HTML {

	return func(name string, form *Form) template.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/field", embed.Context{
			"form":  form,
			"field": form.Get(name),
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return template.HTML(string(data))
	}
}

func createTemplateLabel(loader *embed.TemplateLoader) func(name string, form *Form) template.HTML {

	return func(name string, form *Form) template.HTML {
		field := form.Get(name)

		data, err := loader.Execute(field.Label.Template, embed.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return template.HTML(string(data))
	}
}

func createTemplateInput(loader *embed.TemplateLoader) func(name string, form *Form) template.HTML {

	return func(name string, form *Form) template.HTML {
		field := form.Get(name)

		templates := []string{
			fmt.Sprintf("%s:form/input.%s", field.Module, field.Input.Type),
			"form:form/input.base",
		}

		for _, path := range templates {
			data, err := loader.Execute(path, embed.Context{
				"form":  form,
				"field": field,
				"input": field.Input,
				"label": field.Label,
			})

			if err == nil {
				return template.HTML(string(data))
			}
		}

		helper.PanicOnError(ErrNoTemplate)

		return template.HTML("")
	}
}

func createTemplateErrors(loader *embed.TemplateLoader) func(name string, form *Form) template.HTML {

	return func(name string, form *Form) template.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/errors", embed.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return template.HTML(string(data))
	}
}

func createTemplateHelp(loader *embed.TemplateLoader) func(name string, form *Form) template.HTML {

	return func(name string, form *Form) template.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/help", embed.Context{
			"form":  form,
			"field": field,
			"help":  field.Help,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return template.HTML(string(data))
	}
}
