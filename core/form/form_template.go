// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"errors"
	"fmt"
	tpl "html/template"

	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/template"
)

var (
	ErrNoTemplate = errors.New("unable to find the template to render")
)

func createTemplateField(loader *template.TemplateLoader) func(name string, form *Form) tpl.HTML {

	return func(name string, form *Form) tpl.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/field", template.Context{
			"form":  form,
			"field": form.Get(name),
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return tpl.HTML(string(data))
	}
}

func createTemplateLabel(loader *template.TemplateLoader) func(name string, form *Form) tpl.HTML {

	return func(name string, form *Form) tpl.HTML {
		field := form.Get(name)

		data, err := loader.Execute(field.Label.Template, template.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return tpl.HTML(string(data))
	}
}

func createTemplateInput(loader *template.TemplateLoader) func(name string, form *Form) tpl.HTML {

	return func(name string, form *Form) tpl.HTML {
		field := form.Get(name)

		templates := []string{
			fmt.Sprintf("%s:form/input.%s", field.Module, field.Input.Type),
			"form:form/input.base",
		}

		for _, path := range templates {
			data, err := loader.Execute(path, template.Context{
				"form":  form,
				"field": field,
				"input": field.Input,
				"label": field.Label,
			})

			if err == nil {
				return tpl.HTML(string(data))
			}
		}

		helper.PanicOnError(ErrNoTemplate)

		return tpl.HTML("")
	}
}

func createTemplateErrors(loader *template.TemplateLoader) func(name string, form *Form) tpl.HTML {

	return func(name string, form *Form) tpl.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/errors", template.Context{
			"form":  form,
			"field": field,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return tpl.HTML(string(data))
	}
}

func createTemplateHelp(loader *template.TemplateLoader) func(name string, form *Form) tpl.HTML {

	return func(name string, form *Form) tpl.HTML {
		field := form.Get(name)

		data, err := loader.Execute("form:form/help", template.Context{
			"form":  form,
			"field": field,
			"help":  field.Help,
			"input": field.Input,
			"label": field.Label,
		})

		helper.PanicOnError(err)

		return tpl.HTML(string(data))
	}
}
