// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/modules/base"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {
	l.Prepare(func(app *goapp.App) error {
		c := app.Get("gonode.handler_collection").(base.HandlerCollection)
		c.Add("core.user", &UserHandler{})

		s := app.Get("gonode.node.serializer").(*base.Serializer)
		s.AddSerializer("core.user", UserSerializer)
		s.AddDeserializer("core.user", UserDeserializer)

		return nil
	})
}
