// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/modules/template"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Config(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("form", GetEmbedFS())

		loader := app.Get("gonode.template").(*template.TemplateLoader)

		loader.FuncMap["form_field"] = createTemplateField(loader)
		loader.FuncMap["form_label"] = createTemplateLabel(loader)
		loader.FuncMap["form_input"] = createTemplateInput(loader)
		loader.FuncMap["form_help"] = createTemplateHelp(loader)
		loader.FuncMap["form_errors"] = createTemplateErrors(loader)

		return nil
	})
}
