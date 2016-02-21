// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/config"
	"reflect"
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

func ConfigureServer(l *goapp.Lifecycle, conf *config.Config) {

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
