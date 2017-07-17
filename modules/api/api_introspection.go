// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"

	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/zenazn/goji/web"
)

type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func Api_GET_Handlers_Node(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	collections := app.Get("gonode.handler_collection").(base.HandlerCollection)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		ch := make([]*base.HandlerMetadata, 0)

		for code, h := range collections {
			var m *base.HandlerMetadata

			if cm, ok := h.(base.MetadataHandler); ok {
				m = cm.GetMetadata()
			} else {
				m = base.NewHandlerMetadata()
				m.Name = code
			}

			m.Code = code

			ch = append(ch, m)
		}

		serializer.Serialize(res, ch)
	}
}

func Api_GET_Handlers_View(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	collections := app.Get("gonode.view_handler_collection").(base.ViewHandlerCollection)
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		ch := make([]*base.HandlerViewMetadata, 0)

		for code, h := range collections {
			var m *base.HandlerViewMetadata

			if cm, ok := h.(base.ViewMetadataHandler); ok {
				m = cm.GetViewMetadata()
			} else {
				m = base.NewViewHandlerMetadata()
				m.Name = code
			}

			m.Code = code

			ch = append(ch, m)
		}

		serializer.Serialize(res, ch)
	}
}

func Api_GET_Services(app *goapp.App) func(c web.C, res http.ResponseWriter, req *http.Request) {
	serializer := app.Get("gonode.node.serializer").(*base.Serializer)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		ch := make([]*Service, 0)

		for _, key := range app.GetKeys() {
			ch = append(ch, &Service{
				Name: key,
				Type: fmt.Sprintf("%T", app.Get(key)),
			})
		}

		serializer.Serialize(res, ch)
	}
}
