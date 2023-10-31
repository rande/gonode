// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package blog

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/modules/base"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Prepare(func(app *goapp.App) error {
		app.Get("gonode.handler_collection").(base.HandlerCollection).Add("blog.post", &PostHandler{})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("blog", GetEmbedFS())

		return nil
	})
}
