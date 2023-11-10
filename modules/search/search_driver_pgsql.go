// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"net/url"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/squirrel"
	"github.com/rande/gonode/modules/base"
)

func GetJsonQuery(left string, sep string) string {
	fields := strings.Split(left, ".")

	c := ""
	for p, f := range fields {
		if p == 0 {
			c += f
		} else {
			c += fmt.Sprintf(sep+"'%s'", f)
		}
	}

	return c
}

func GetJsonSearchQuery(query sq.SelectBuilder, params []*Param, field string) sq.SelectBuilder {
	//-- SELECT uuid, "data" #> '{tags,1}' as tags FROM nodes WHERE  "data" @> '{"tags": ["sport"]}'
	//-- SELECT uuid, "data" #> '{tags}' AS tags FROM nodes WHERE  "data" -> 'tags' ?| array['sport'];
	for _, param := range params {
		value := param.Value.([]string)

		if len(value) > 1 {
			name := GetJsonQuery(field+"."+param.SubField, "->")
			query = query.Where(squirrel.NewExprSlice(fmt.Sprintf("%s ??| ARRAY["+sq.Placeholders(len(value))+"]", name), value))
		}

		if len(value) == 1 {
			name := GetJsonQuery(field+"."+param.SubField, "->>")
			query = query.Where(sq.Expr(fmt.Sprintf("%s = ?", name), value[0]))
		}
	}

	return query
}

type SearchPager struct {
	Elements []*base.Node
	Page     uint64
	PerPage  uint64
	Next     uint64
	Previous uint64
	Form     *SearchForm
}

func (s *SearchPager) PageQuery(page uint64) url.Values {
	params := s.Form.UrlValues()

	params.Set("page", fmt.Sprintf("%v", page))

	return params
}

type SearchPGSQL struct {
}

func AddEqClause(column string, query sq.SelectBuilder, params []*Param) sq.SelectBuilder {
	if len(params) == 0 {
		return query
	}

	values := []interface{}{}
	for _, param := range params {
		values = append(values, param.Value)
	}

	query = query.Where(sq.Eq{column: values})

	return query
}

func AddClause(column string, query sq.SelectBuilder, params *Param) sq.SelectBuilder {
	if params == nil {
		return query
	}

	query = query.Where(fmt.Sprintf("%s = ?", column), params.Value)

	return query
}

func (s *SearchPGSQL) BuildQuery(searchForm *SearchForm, query sq.SelectBuilder) sq.SelectBuilder {
	for _, order := range searchForm.OrderBy {
		helper.PanicIf(len(order.SubField) == 0, "OrderBy field name is empty")

		query = query.OrderBy(GetJsonQuery(order.SubField, "->") + " " + order.Operation)
	}

	query = AddEqClause("uuid", query, searchForm.Uuid)
	query = AddEqClause("type", query, searchForm.Type)
	query = AddEqClause("updated_by", query, searchForm.UpdatedBy)
	query = AddEqClause("created_by", query, searchForm.CreatedBy)
	query = AddEqClause("parent_uuid", query, searchForm.ParentUuid)
	query = AddEqClause("set_uuid", query, searchForm.SetUuid)
	query = AddEqClause("source", query, searchForm.Source)
	query = AddEqClause("status", query, searchForm.Status)
	query = AddEqClause("weight", query, searchForm.Weight)

	query = AddClause("name", query, searchForm.Name)
	query = AddClause("slug", query, searchForm.Slug)
	query = AddClause("revision", query, searchForm.Revision)
	query = AddClause("enabled", query, searchForm.Enabled)
	query = AddClause("deleted", query, searchForm.Deleted)
	query = AddClause("current", query, searchForm.Current)

	query = GetJsonSearchQuery(query, searchForm.Data, "data")
	query = GetJsonSearchQuery(query, searchForm.Meta, "meta")

	return query
}
