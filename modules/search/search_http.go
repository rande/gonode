// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"net/http"
	"regexp"

	"github.com/gorilla/schema"
	"github.com/rande/gonode/core/helper"
)

var (
	rexOrderBy = regexp.MustCompile(`(^[a-z,_.A-Z]*),(DESC|ASC|desc|asc)$`)
	rexMeta    = regexp.MustCompile(`meta\.([a-zA-Z]*)`)
	rexData    = regexp.MustCompile(`data\.([a-zA-Z]*)`)
)

type HttpSearchForm struct {
	Page       int64               `schema:"page"`
	PerPage    int64               `schema:"per_page"`
	OrderBy    []string            `schema:"order_by"`
	Uuid       []string            `schema:"uuid"`
	Type       []string            `schema:"type"`
	Name       string              `schema:"name"`
	Slug       string              `schema:"slug"`
	Data       map[string][]string `schema:"data"`
	Meta       map[string][]string `schema:"meta"`
	Status     []string            `schema:"status"`
	Weight     []string            `schema:"weight"`
	Revision   string              `schema:"revision"`
	Enabled    string              `schema:"enabled"`
	Deleted    string              `schema:"deleted"`
	Current    string              `schema:"current"`
	UpdatedBy  []string            `schema:"updated_by"`
	CreatedBy  []string            `schema:"created_by"`
	ParentUuid []string            `schema:"parent_uuid"`
	SetUuid    []string            `schema:"set_uuid"`
	Source     []string            `schema:"source"`
}

func GetHttpSearchForm() *HttpSearchForm {
	return &HttpSearchForm{
		Data:    make(map[string][]string),
		Meta:    make(map[string][]string),
		OrderBy: []string{"created_at,DESC"},
	}
}

type HttpSearchParser struct {
	MaxResult uint64
}

func (h *HttpSearchParser) HandleSearch(res http.ResponseWriter, req *http.Request) *SearchForm {
	req.ParseForm()

	searchForm := NewSearchForm()
	httpSearchForm := GetHttpSearchForm()
	decoder := schema.NewDecoder()
	decoder.Decode(httpSearchForm, req.Form)

	// check page range
	if httpSearchForm.Page < 0 || httpSearchForm.PerPage < 0 || uint64(httpSearchForm.PerPage) > h.MaxResult {
		helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid `pagination` range")

		return nil
	}

	if httpSearchForm.Page < 1 {
		httpSearchForm.Page = 1
	}

	if httpSearchForm.PerPage < 1 || httpSearchForm.PerPage > 256 {
		httpSearchForm.PerPage = 32
	}

	searchForm.PerPage = uint64(httpSearchForm.PerPage)
	searchForm.Page = uint64(httpSearchForm.Page)

	for _, order := range httpSearchForm.OrderBy {
		r := rexOrderBy.FindAllStringSubmatch(order, -1)

		if r == nil {
			helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid `order_by` condition")

			return nil
		}

		searchForm.OrderBy = append(searchForm.OrderBy, NewParam(nil, r[0][2], r[0][1]))
	}

	for _, uuid := range httpSearchForm.Uuid {
		searchForm.Uuid = append(searchForm.Uuid, NewParam(uuid))
	}

	for _, nodeType := range httpSearchForm.Type {
		searchForm.Type = append(searchForm.Type, NewParam(nodeType))
	}

	if len(httpSearchForm.Name) > 0 {
		searchForm.Name = NewParam(httpSearchForm.Name, "=")
	}

	if len(httpSearchForm.Slug) > 0 {
		searchForm.Slug = NewParam(httpSearchForm.Slug, "=")
	}

	// analyse Data
	for name, value := range req.Form {
		values := rexData.FindStringSubmatch(name)

		if len(values) == 2 {
			searchForm.Data = append(searchForm.Data, NewParam(value, "=", values[1]))
		}
	}

	// analyse Meta
	for name, value := range req.Form {
		values := rexMeta.FindStringSubmatch(name)

		if len(values) == 2 {
			searchForm.Meta = append(searchForm.Meta, NewParam(value, "=", values[1]))
		}
	}

	for _, status := range httpSearchForm.Status {
		searchForm.Status = append(searchForm.Status, NewParam(status))
	}

	for _, weight := range httpSearchForm.Weight {
		searchForm.Weight = append(searchForm.Weight, NewParam(weight, "="))
	}

	if len(httpSearchForm.Revision) > 0 {
		searchForm.Revision = NewParam(httpSearchForm.Revision, "=")
	}

	if httpSearchForm.Enabled == "true" || httpSearchForm.Enabled == "t" || httpSearchForm.Enabled == "1" {
		searchForm.Enabled = NewParam(true, "=")
	} else if httpSearchForm.Enabled == "false" || httpSearchForm.Enabled == "f" || httpSearchForm.Enabled == "0" {
		searchForm.Enabled = NewParam(false, "=")
	} else if len(httpSearchForm.Enabled) > 0 {
		helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid `enabled` condition")

		return nil
	}

	// TODO: only admin token can view deleted node
	if httpSearchForm.Deleted == "true" || httpSearchForm.Deleted == "t" || httpSearchForm.Deleted == "1" {
		searchForm.Deleted = NewParam(true, "=")
	} else if httpSearchForm.Deleted == "false" || httpSearchForm.Deleted == "f" || httpSearchForm.Deleted == "0" {
		searchForm.Deleted = NewParam(false, "=")
	} else if len(httpSearchForm.Deleted) > 0 {
		helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid `deleted `condition")

		return nil
	}

	if httpSearchForm.Current == "true" || httpSearchForm.Current == "t" || httpSearchForm.Current == "1" {
		searchForm.Current = NewParam(true, "=")
	} else if httpSearchForm.Current == "false" || httpSearchForm.Current == "f" || httpSearchForm.Current == "0" {
		searchForm.Current = NewParam(false, "=")
	} else if len(httpSearchForm.Current) > 0 {
		helper.SendWithHttpCode(res, http.StatusPreconditionFailed, "Invalid `current` condition")

		return nil
	}

	for _, updatedBy := range httpSearchForm.UpdatedBy {
		searchForm.UpdatedBy = append(searchForm.UpdatedBy, NewParam(updatedBy, "="))
	}

	for _, source := range httpSearchForm.Source {
		searchForm.Source = append(searchForm.Source, NewParam(source, "="))
	}

	for _, createdBy := range httpSearchForm.CreatedBy {
		searchForm.CreatedBy = append(searchForm.CreatedBy, NewParam(createdBy, "="))
	}

	for _, setUuid := range httpSearchForm.SetUuid {
		searchForm.SetUuid = append(searchForm.SetUuid, NewParam(setUuid, "="))
	}

	for _, parentUuid := range httpSearchForm.ParentUuid {
		searchForm.ParentUuid = append(searchForm.ParentUuid, NewParam(parentUuid, "="))
	}

	return searchForm
}
