// Copyright © 2014-2017 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"strings"

	sq "github.com/lann/squirrel"
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
}

type SearchPGSQL struct {
}

func (s *SearchPGSQL) BuildQuery(searchForm *SearchForm, query sq.SelectBuilder) sq.SelectBuilder {

	for _, order := range searchForm.OrderBy {
		helper.PanicIf(len(order.SubField) == 0, "OrderBy field name is empty")

		query = query.OrderBy(GetJsonQuery(order.SubField, "->") + " " + order.Operation)
	}

	if searchForm.Uuid != nil {
		query = query.Where(sq.Eq{"uuid": searchForm.Uuid.Value})
	}

	if searchForm.Type != nil {
		query = query.Where(sq.Eq{"type": searchForm.Type.Value})
	}

	if searchForm.Name != nil {
		query = query.Where("name = ?", searchForm.Name.Value)
	}

	if searchForm.Slug != nil {
		query = query.Where("slug = ?", searchForm.Slug.Value)
	}

	query = GetJsonSearchQuery(query, searchForm.Data, "data")
	query = GetJsonSearchQuery(query, searchForm.Meta, "meta")

	if searchForm.Status != nil {
		query = query.Where(sq.Eq{"status": searchForm.Status.Value})
	}

	if searchForm.Weight != nil {
		query = query.Where("weight = ?", searchForm.Weight.Value)
	}

	if searchForm.Revision != nil {
		query = query.Where("revision = ?", searchForm.Revision.Value)
	}

	if searchForm.Enabled != nil {
		query = query.Where("enabled = ?", searchForm.Enabled.Value)
	}

	if searchForm.Deleted != nil {
		query = query.Where("deleted = ?", searchForm.Deleted.Value)
	}

	if searchForm.Current != nil {
		query = query.Where("current = ?", searchForm.Current.Value)
	}

	if searchForm.UpdatedBy != nil {
		query = query.Where(sq.Eq{"updated_by": searchForm.UpdatedBy.Value})
	}

	if searchForm.CreatedBy != nil {
		query = query.Where(sq.Eq{"created_by": searchForm.CreatedBy.Value})
	}

	if searchForm.ParentUuid != nil {
		query = query.Where(sq.Eq{"parent_uuid": searchForm.ParentUuid.Value})
	}

	if searchForm.SetUuid != nil {
		query = query.Where(sq.Eq{"set_uuid": searchForm.SetUuid.Value})
	}

	if searchForm.Source != nil {
		query = query.Where(sq.Eq{"source": searchForm.Source.Value})
	}

	return query
}
