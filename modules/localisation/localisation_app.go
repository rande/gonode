// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package localisation

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/template"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {

		return nil
	})

	l.Config(func(app *goapp.App) error {
		loader := app.Get("gonode.template").(*template.TemplateLoader)
		defaultLocale := "en_GB"

		for name, Func := range CreateTemplateFuncMap(defaultLocale) {
			loader.FuncMap[name] = Func
		}

		return nil
	})
}
