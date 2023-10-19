// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dashboard

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/embed"
	"github.com/rande/gonode/core/router"
)

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Prepare(func(app *goapp.App) error {
		app.Get("gonode.embeds").(*embed.Embeds).Add("dashboard", GetEmbedFS())

		return nil
	})

	l.Prepare(func(app *goapp.App) error {

		r := app.Get("gonode.router").(*router.Router)

		r.Get("dashboard", conf.Dashboard.Prefix, InitView(app, Dashboard_GET_Index))
		r.Get("dashboard_login", conf.Dashboard.Prefix+"/login", InitView(app, Dashboard_GET_Login))

		r.Get("dashboard_node_list", conf.Dashboard.Prefix+"/node/list", InitView(app, Dashboard_GET_Node_List))
		r.Handle("dashboard_node_create", conf.Dashboard.Prefix+"/node/create", InitView(app, Dashboard_HANDLE_Node_Create))
		r.Post("dashboard_node_update", conf.Dashboard.Prefix+"/node/:nid/update", InitView(app, Dashboard_GET_Node_Update))
		r.Handle("dashboard_node_delete", conf.Dashboard.Prefix+"/node/delete", InitView(app, Dashboard_GET_ToDo))
		r.Get("dashboard_node_edit", conf.Dashboard.Prefix+"/node/:nid", InitView(app, Dashboard_GET_Node_Edit))

		r.Handle("dashboard_settings", conf.Dashboard.Prefix+"/settings", InitView(app, Dashboard_GET_ToDo))

		return nil
	})
}
