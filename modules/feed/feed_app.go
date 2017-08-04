// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package feed

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Prepare(func(app *goapp.App) error {
		c := app.Get("gonode.handler_collection").(base.HandlerCollection)
		c.Add("feed.index", &FeedHandler{})

		cv := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
		cv.Add("feed.index", &FeedViewHandler{
			Search:  app.Get("gonode.search.pgsql").(*search.SearchPGSQL),
			Manager: app.Get("gonode.manager").(*base.PgNodeManager),
		})

		return nil
	})
}
