// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"database/sql"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/security"
)

func GetValue(source interface{}, field string) interface{} {
	v := reflect.ValueOf(source).Elem().FieldByName(field)

	// field does not exist ... return nil
	// need to add a way to configure alias
	if !v.IsValid() {
		return nil
	}

	return v.Interface()
}

func Configure(l *goapp.Lifecycle, conf *config.Config) {

	l.Register(func(app *goapp.App) error {
		app.Set("gonode.handler_collection", func(app *goapp.App) interface{} {
			return HandlerCollection{}
		})

		app.Set("gonode.view_handler_collection", func(app *goapp.App) interface{} {
			return ViewHandlerCollection{}
		})

		app.Set("gonode.security.voter.access", func(app *goapp.App) interface{} {
			return &AccessVoter{}
		})

		app.Set("gonode.security.voter.role", func(app *goapp.App) interface{} {
			return &security.RoleVoter{
				Prefix: "node:",
			}
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		app.Set("gonode.manager", func(app *goapp.App) interface{} {
			return &PgNodeManager{
				Logger:   app.Get("logger").(*log.Logger),
				Db:       app.Get("gonode.postgres.connection").(*sql.DB),
				ReadOnly: false,
				Handlers: app.Get("gonode.handler_collection").(Handlers),
				Prefix:   conf.Databases["master"].Prefix,
			}
		})

		app.Set("gonode.node.serializer", func(app *goapp.App) interface{} {
			s := NewSerializer()
			s.Handlers = app.Get("gonode.handler_collection").(Handlers)

			return s
		})

		return nil
	})

	l.Prepare(func(app *goapp.App) error {
		pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)

		pongo.Globals["node_data"] = func(vnode, vname *pongo2.Value) *pongo2.Value {
			node := vnode.Interface().(*Node)

			return pongo2.AsValue(GetValue(node.Data, vname.String()))
		}

		pongo.Globals["node_meta"] = func(vnode, vname *pongo2.Value) *pongo2.Value {
			node := vnode.Interface().(*Node)

			return pongo2.AsValue(GetValue(node.Meta, vname.String()))
		}

		return nil
	})
}
