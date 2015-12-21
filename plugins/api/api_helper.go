package api

import (
	"github.com/gorilla/schema"
	sq "github.com/lann/squirrel"
	"github.com/rande/gonode/helper"
	"github.com/zenazn/goji/web"
	"net/http"
)

func HandleSearch(apiHandler *Api, query sq.SelectBuilder, c web.C, res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	searchForm := GetSearchForm()
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
		values := rexData.FindStringSubmatch(name)

		if len(values) == 2 {
			searchForm.Data[values[1]] = value
		}
	}

	if searchForm.Page < 0 || searchForm.PerPage < 0 || searchForm.PerPage > 128 {
		helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid pagination range")

		return
	}

	if searchForm.Page == 0 {
		searchForm.Page = 1
	}

	if searchForm.PerPage == 0 {
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

	if searchForm.Deleted != "" { // TODO: only admin token can view deleted node
		query = query.Where("deleted = ?", searchForm.Deleted)
	} else {
		query = query.Where("deleted = ?", "f")
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

	for _, order := range searchForm.OrderBy {
		r := rexOrderBy.FindAllStringSubmatch(order, -1)

		if r == nil {
			helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid order_by condition")

			return
		}

		query = query.OrderBy(GetJsonQuery(r[0][1], "->") + " " + r[0][2])
	}

	query = GetJsonSearchQuery(query, searchForm.Meta, "meta")
	query = GetJsonSearchQuery(query, searchForm.Data, "data")

	apiHandler.Find(res, query, uint64(searchForm.Page), uint64(searchForm.PerPage))
}
