// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package raw

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/base"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Prepare(func(app *goapp.App) error {

		c := app.Get("gonode.handler_collection").(base.HandlerCollection)
		c.Add("core.raw", &RawHandler{})

		cv := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
		cv.Add("core.raw", &RawViewHandler{})

		return nil
	})
}
