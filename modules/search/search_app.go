// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/modules/base"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.search.pgsql", func(app *goapp.App) interface{} {
			return &SearchPGSQL{}
		})

		app.Set("gonode.search.parser.http", func(app *goapp.App) interface{} {
			return &HttpSearchParser{
				MaxResult: conf.Search.MaxResult,
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		c := app.Get("gonode.handler_collection").(base.HandlerCollection)
		c.Add("search.index", &IndexHandler{})

		cv := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
		cv.Add("search.index", &IndexViewHandler{
			Search:    app.Get("gonode.search.pgsql").(*SearchPGSQL),
			Manager:   app.Get("gonode.manager").(*base.PgNodeManager),
			MaxResult: 128,
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("search", GetEmbedFS())

		return nil
	})
}
