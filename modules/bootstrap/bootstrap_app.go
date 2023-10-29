// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bootstrap

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("bootstrap", GetEmbedFS())

		return nil
	})
}
