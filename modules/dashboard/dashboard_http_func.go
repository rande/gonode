// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/search"
	"github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
)

func Dashboard_GET_Index(app *goapp.App) ViewHandlerInterface {
	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		return HtmlResponse(200, "dashboard:pages/index")
	}
}

func Dashboard_GET_Login(app *goapp.App) ViewHandlerInterface {
	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		return HtmlResponse(200, "dashboard:pages/login")
	}
}

func Dashboard_GET_ToDo(app *goapp.App) ViewHandlerInterface {
	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		return HtmlResponse(200, "dashboard:dash/todo")
	}
}

func Dashboard_GET_Node_List(app *goapp.App) ViewHandlerInterface {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	logger := app.Get("logger").(*logrus.Logger)

	queryBuilder := app.Get("gonode.search.pgsql").(*search.SearchPGSQL)

	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		form := search.NewSearchForm()
		form.PerPage = 32
		form.Page = 1

		options := base.NewAccessOptionsFromToken(security.GetTokenFromContext(c))
		pager := search.GetPager(form, manager, queryBuilder, options)

		return HtmlResponse(200, "dashboard:node/list.tpl").
			Add("pager", pager)
	}
}

type NodeForm struct {
	Name string `schema:"Name,required"`
	Slug string `schema:"Slug,required"`
	// Data       interface{} `schema:"data"`
	// Meta       interface{} `schema:"meta"`
	Status     int            `schema:"Status,required"`
	Weight     int            `schema:"Weight"`
	Enabled    bool           `schema:"Enabled,required"`
	ParentUuid base.Reference `schema:"ParentUuid"`
}

func Dashboard_GET_Node_Edit(app *goapp.App) ViewHandlerInterface {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	logger := app.Get("logger").(*logrus.Logger)

	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		var node *base.Node

		if uuid, ok := c.URLParams["uuid"]; ok {
			reference, err := base.GetReferenceFromString(uuid)

			if err != nil {
				panic(err)
			}

			node = manager.Find(reference)
		}

		if node == nil {
			return HtmlResponse(404, "dashboard:error/404.tpl")
		}

		nodeForm := &NodeForm{}
		nodeForm.Name = node.Name
		nodeForm.Slug = node.Slug
		nodeForm.Status = node.Status
		nodeForm.Weight = node.Weight
		nodeForm.Enabled = node.Enabled
		nodeForm.ParentUuid = node.ParentUuid

		logger.Info(fmt.Sprintf("%#v", nodeForm))

		return HtmlResponse(200, "dashboard:node/edit.tpl").
			Add("node", node).
			Add("err", nil).
			Add("form", nodeForm)
	}
}

func Dashboard_GET_Node_Update(app *goapp.App) ViewHandlerInterface {
	manager := app.Get("gonode.manager").(*base.PgNodeManager)
	// logger := app.Get("logger").(*logrus.Logger)

	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {
		var node *base.Node

		// load the node
		if uuid, ok := c.URLParams["uuid"]; ok {
			reference, err := base.GetReferenceFromString(uuid)

			if err != nil {
				panic(err)
			}

			node = manager.Find(reference)
		}

		// parse the form
		req.ParseForm()

		nodeForm := &NodeForm{}

		decoder := schema.NewDecoder()
		errors := decoder.Decode(nodeForm, req.Form)

		// if error, then render the form with error
		if errors != nil {
			return HtmlResponse(200, "dashboard:node/edit.tpl").
				Add("node", node).
				Add("form", nodeForm).
				Add("errors", errors)
		}

		// update the node
		node.Name = nodeForm.Name
		node.Slug = nodeForm.Slug
		node.Status = nodeForm.Status
		node.Weight = nodeForm.Weight
		node.Enabled = nodeForm.Enabled
		node.ParentUuid = nodeForm.ParentUuid

		node, err := manager.Save(node, true)

		if err != nil {
			panic(err)
		}

		return RedirectResponse(302, "/dashboard/node/"+node.Uuid.String())
	}
}

func Dashboard_HANDLE_Node_Create(app *goapp.App) ViewHandlerInterface {
	handlers := app.Get("gonode.handler_collection").(base.Handlers)

	return func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse {

		handlerType := req.URL.Query().Get("type")

		if len(handlerType) == 0 {
			return HtmlResponse(200, "dashboard:node/create.tpl").
				Add("keys", handlers.GetTypes())
		}

		if len(handlerType) > 0 && !handlers.HasType(handlerType) {
			return HtmlResponse(404, "dashboard:node/create.tpl").
				Add("keys", handlers.GetTypes()).
				Add("error", fmt.Sprintf("The type `%s` does not exists", handlerType))
		}

		handler := handlers.GetByType(handlerType)

		nodeForm := &NodeForm{}
		var errors error

		if req.Method == "POST" {
			decoder := schema.NewDecoder()
			errors = decoder.Decode(nodeForm, req.Form)
		}

		return HtmlResponse(200, "dashboard:node/create.tpl").
			Add("keys", handlers.GetTypes()).
			Add("handler", handler).
			Add("form", nodeForm).
			Add("errors", errors)
	}
}
