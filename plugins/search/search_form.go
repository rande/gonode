// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

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
	Uuid       *Param   `json:"uuid"`
	Type       *Param   `json:"type"`
	Name       *Param   `json:"name"`
	Slug       *Param   `json:"slug"`
	Data       []*Param `json:"data"`
	Meta       []*Param `json:"meta"`
	Status     *Param   `json:"status"`
	Weight     *Param   `json:"weight"`
	Revision   *Param   `json:"revision"`
	Enabled    *Param   `json:"enabled"`
	Deleted    *Param   `json:"deleted"`
	Current    *Param   `json:"current"`
	UpdatedBy  *Param   `json:"updated_by"`
	CreatedBy  *Param   `json:"created_by"`
	ParentUuid *Param   `json:"parent_uuid"`
	SetUuid    *Param   `json:"set_uuid"`
	Source     *Param   `json:"source"`
}

func NewSearchForm() *SearchForm {
	return &SearchForm{
		OrderBy: make([]*Param, 0),
		Data:    make([]*Param, 0),
		Meta:    make([]*Param, 0),
		Deleted: NewParam(false, "="),
	}
}
