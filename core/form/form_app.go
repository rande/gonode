// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("form", GetEmbedFS())

		pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)

		pongo.Globals["form_field"] = createPongoField(pongo)
		pongo.Globals["form_label"] = createPongoLabel(pongo)
		pongo.Globals["form_input"] = createPongoInput(pongo)
		pongo.Globals["form_help"] = createPongoHelp(pongo)
		pongo.Globals["form_errors"] = createPongoErrors(pongo)

		return nil
	})
}
