package extra

import (
	"fmt"
	"net/http"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	nc "github.com/rande/gonode/core"
	sq "github.com/lann/squirrel"
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
				searchForm.Meta[values[1]] = value
			}
		}

		// analyse Data
		for name, value := range req.Form {
			values := rexMeta.FindStringSubmatch(name)

			if len(values) == 2 {
				searchForm.Meta[values[1]] = value
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

		if len(searchForm.Source) != 0 {
			query = query.Where(sq.Eq{"source": searchForm.Source})
		}

		if searchForm.Enabled != "" {
			query = query.Where("enabled = ?", searchForm.Enabled)
		}

		if len(searchForm.Type) != 0 {
			query = query.Where(sq.Eq{"type": searchForm.Type})
		}

		if searchForm.Current != "" {
			query = query.Where("current = ?", searchForm.Current)
		}

		if searchForm.Deleted != "" {
			query = query.Where("deleted = ?", searchForm.Deleted)
		}

		if len(searchForm.Uuid) != 0 {
			query = query.Where(sq.Eq{"uuid": searchForm.Uuid})
		}

		if len(searchForm.ParentUuid) != 0 {
			query = query.Where(sq.Eq{"parent_uuid": searchForm.ParentUuid})
		}

		if searchForm.Slug != "" {
			query = query.Where("slug = ?", searchForm.Slug)
		}

		if searchForm.Revision != "" {
			query = query.Where("revision = ?", searchForm.Revision)
		}

		if len(searchForm.Status) != 0 {
			query = query.Where(sq.Eq{"status": searchForm.Status})
		}

		// Parse Meta value
		for name, value := range searchForm.Meta {
			//-- SELECT uuid, "data" #> '{tags,1}' as tags FROM nodes WHERE  "data" @> '{"tags": ["sport"]}'
			//-- SELECT uuid, "data" #> '{tags}' AS tags FROM nodes WHERE  "data" -> 'tags' ?| array['sport'];
			if len(value) > 1 {
				query = query.Where(sq.ExprSlice(fmt.Sprintf("meta->'%s' ??| array[" + sq.Placeholders(len(value))+ "]", name), len(value), value),)
			}

			if len(value) == 1 {
				query = query.Where(sq.Expr(fmt.Sprintf("meta->>'%s' = ?", name), value))
			}
		}

		// Parse Data value
		for name, value := range searchForm.Data {
			if len(value) > 1 {
				query = query.Where(sq.ExprSlice(fmt.Sprintf("data->'%s' ??| array[" + sq.Placeholders(len(value))+ "]", name), len(value), value),)
			}

			if len(value) == 1 {
				query = query.Where(sq.Expr(fmt.Sprintf("data->>'%s' = ?", name), value))
			}
		}

		api.Find(res, query, searchForm.Page, searchForm.PerPage)
	})
}

