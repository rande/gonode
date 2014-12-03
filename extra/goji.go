package extra

import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	nc "github.com/rande/gonode/core"
	"github.com/gorilla/schema"
)

func ConfigureGoji(manager *nc.PgNodeManager, prefix string) {
	api := &nc.Api{
		Manager: manager,
		Version: "1.0.0",
	}

	goji.Get(prefix + "/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		api.FindOne(c.URLParams["uuid"], res)
	})

	goji.Post(prefix + "/nodes", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		api.Save(req.Body, res)
	})

	goji.Put(prefix + "/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		api.Save(req.Body, res)
	})

	goji.Delete(prefix + "/nodes/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		api.RemoveOne(c.URLParams["uuid"], res)
	})

	goji.Get(prefix + "/nodes", func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Type", "application/json")

		query := api.SelectBuilder()

		req.ParseForm()

		searchForm := nc.GetSearchForm()
		decoder := schema.NewDecoder()
		decoder.Decode(searchForm, req.Form)

		// analyse Meta
		for name, value := range req.Form {
			values := rexMeta.FindStringSubmatch(name)

			if len(values) == 2 {
				searchForm.Meta[values[1]] = value[0]
			}
		}

		// analyse Data
		for name, value := range req.Form {
			values := rexMeta.FindStringSubmatch(name)

			if len(values) == 2 {
				searchForm.Meta[values[1]] = value[0]
			}
		}

		if searchForm.Page < 1 {
			searchForm.Page = 1
		}

		if searchForm.PerPage > 128 {
			searchForm.PerPage = 32
		}

		if searchForm.PerPage < 1 {
			searchForm.PerPage = 32
		}

		if searchForm.Source != "" {
			query = query.Where("source = ?", searchForm.Source)
		}

		if searchForm.Enabled != "" {
			query = query.Where("enabled = ?", searchForm.Enabled)
		}

		if searchForm.Type != "" {
			query = query.Where("type = ?", searchForm.Type)
		}

		if searchForm.Current != "" {
			query = query.Where("current = ?", searchForm.Current)
		}

		if searchForm.Deleted != "" {
			query = query.Where("deleted = ?", searchForm.Deleted)
		}

		if searchForm.Uuid != "" {
			query = query.Where("uuid = ?", searchForm.Uuid)
		}

		if searchForm.ParentUuid != "" {
			query = query.Where("parent_uuid = ?", searchForm.ParentUuid)
		}

		if searchForm.Slug != "" {
			query = query.Where("slug = ?", searchForm.Slug)
		}

		if searchForm.Revision != "" {
			query = query.Where("revision = ?", searchForm.Revision)
		}

		if searchForm.Status != "" {
			query = query.Where("status = ?", searchForm.Status)
		}

		for name, value := range searchForm.Meta {
			query = query.Where(fmt.Sprintf("meta->>'%s' = ?", name), value)
		}

		for name, value := range searchForm.Data {
			query = query.Where(fmt.Sprintf("data->>'%s' = ?", name), value)
		}


//		query = query.Where("meta->>'Foo' = ?", "markdown")

		//
		//		if len(req.FormValue("set")) > 0 {
		//			query["set"] = bson.RegEx{req.FormValue("set"), "i"}
		//		}
		//
		//		if _, ok := req.Form["types"]; ok {
		//			query["type"] = bson.M{
		//				"$in": req.Form["types"],
		//			}
		//		}

		api.Find(res, query, searchForm.Page, searchForm.PerPage)
	})
}

