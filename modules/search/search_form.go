// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"net/url"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/rande/gonode/core/squirrel"
	"github.com/rande/gonode/modules/base"
)

type Param struct {
	SubField  string      `json:"sub_field"`
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}

func NewParam(v interface{}, options ...string) *Param {
	params := &Param{
		Value: v,
	}

	if len(options) > 0 {
		params.Operation = options[0]
	} else {
		params.Operation = "="
	}

	if len(options) > 1 {
		params.SubField = options[1]
	}

	return params
}

type SearchForm struct {
	Page       uint64   `json:"page"`
	PerPage    uint64   `json:"per_page"`
	OrderBy    []*Param `json:"order_by"`
	Uuid       []*Param `json:"uuid"`
	Type       []*Param `json:"type"`
	Name       *Param   `json:"name"`
	Slug       *Param   `json:"slug"`
	Data       []*Param `json:"data"`
	Meta       []*Param `json:"meta"`
	Status     []*Param `json:"status"`
	Weight     []*Param `json:"weight"`
	Revision   *Param   `json:"revision"`
	Enabled    *Param   `json:"enabled"`
	Deleted    *Param   `json:"deleted"`
	Current    *Param   `json:"current"`
	UpdatedBy  []*Param `json:"updated_by"`
	CreatedBy  []*Param `json:"created_by"`
	ParentUuid []*Param `json:"parent_uuid"`
	SetUuid    []*Param `json:"set_uuid"`
	Source     []*Param `json:"source"`
}

func addUrlValue(values url.Values, name string, param *Param) {
	if param == nil || param.Value == nil {
		return
	}

	values.Add(name, fmt.Sprintf("%v", param.Value))
}

func addUrlValues(values url.Values, name string, param []*Param) {
	if param == nil {
		return
	}

	for _, param := range param {
		addUrlValue(values, name, param)
	}
}

func (s *SearchForm) UrlValues() url.Values {
	values := url.Values{}
	values.Add("page", strconv.FormatUint(s.Page, 10))
	values.Add("per_page", strconv.FormatUint(s.PerPage, 10))

	addUrlValues(values, "uuid", s.Uuid)
	addUrlValues(values, "type", s.Type)
	addUrlValues(values, "status", s.Status)
	addUrlValues(values, "weight", s.Weight)
	addUrlValues(values, "updated_by", s.UpdatedBy)
	addUrlValues(values, "order_by", s.OrderBy)
	addUrlValues(values, "created_by", s.CreatedBy)
	addUrlValues(values, "parent_uuid", s.ParentUuid)
	addUrlValues(values, "set_uuid", s.SetUuid)
	addUrlValues(values, "source", s.Source)

	addUrlValue(values, "name", s.Name)
	addUrlValue(values, "slug", s.Slug)
	addUrlValue(values, "revision", s.Revision)
	addUrlValue(values, "enabled", s.Enabled)
	addUrlValue(values, "deleted", s.Deleted)
	addUrlValue(values, "current", s.Current)

	for _, param := range s.Meta {
		if param.Value == nil {
			continue
		}
		values.Add(fmt.Sprintf("meta.%s", param.SubField), fmt.Sprintf("%v", param.Value))
	}

	for _, param := range s.Data {
		if param.Value == nil {
			continue
		}
		values.Add(fmt.Sprintf("data.%s", param.SubField), fmt.Sprintf("%v", param.Value))
	}

	return values
}

func NewSearchForm() *SearchForm {
	return &SearchForm{
		OrderBy: make([]*Param, 0),
		Data:    make([]*Param, 0),
		Meta:    make([]*Param, 0),
		Deleted: NewParam(false, "="),
	}
}

func NewSearchFormFromIndex(index *Index) *SearchForm {
	// we just copy over node to create search form
	search := NewSearchForm()
	search.OrderBy = index.OrderBy
	search.Uuid = index.Uuid
	search.Type = index.Type
	search.Name = index.Name
	search.Slug = index.Slug
	search.Data = index.Data
	search.Meta = index.Meta
	search.Status = index.Status
	search.Weight = index.Weight
	search.Revision = index.Revision
	search.Enabled = index.Enabled
	search.Deleted = index.Deleted
	search.Current = index.Current
	search.UpdatedBy = index.UpdatedBy
	search.CreatedBy = index.CreatedBy
	search.ParentUuid = index.ParentUuid
	search.SetUuid = index.SetUuid
	search.Source = index.Source

	return search
}

func GetPager(form *SearchForm, manager base.NodeManager, engine *SearchPGSQL, options *base.AccessOptions) *SearchPager {
	query := engine.BuildQuery(form, manager.SelectBuilder(base.NewSelectOptions()))

	// apply security access
	if options != nil && len(options.Roles) > 0 {
		value, _ := options.Roles.ToStringSlice()

		query = query.Where(squirrel.NewExprSlice("\"access\" && ARRAY["+sq.Placeholders(len(options.Roles))+"]", value))
	}

	list := manager.FindBy(query, (form.Page-1)*form.PerPage, form.PerPage+1)

	pager := &SearchPager{
		Page:     form.Page,
		PerPage:  form.PerPage,
		Elements: make([]*base.Node, 0),
		Previous: uint64(0),
		Next:     uint64(0),
		Form:     form,
	}

	if form.Page > 1 {
		pager.Previous = form.Page - 1
	}

	counter := uint64(0)
	for e := list.Front(); e != nil; e = e.Next() {
		if counter == form.PerPage {
			pager.Next = form.Page + 1
			break
		}
		pager.Elements = append(pager.Elements, e.Value.(*base.Node))

		counter++
	}

	return pager
}
