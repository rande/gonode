// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"database/sql"
	"fmt"
	tpl "html/template"
	"reflect"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/template"

	log "github.com/sirupsen/logrus"
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

	l.Config(func(app *goapp.App) error {
		loader := app.Get("gonode.template").(*template.TemplateLoader)

		loader.FuncMap["safe"] = func(v interface{}) tpl.HTML {
			return tpl.HTML(fmt.Sprintf("%v", v))
		}

		loader.FuncMap["node_data"] = func(node *Node, name string) interface{} {
			return GetValue(node.Data, name)
		}

		loader.FuncMap["node_meta"] = func(node *Node, name string) interface{} {
			return GetValue(node.Meta, name)
		}

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
}
