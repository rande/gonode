// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
)

func ConfigureServer(l *goapp.Lifecycle, conf *config.ServerConfig) {

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
}
