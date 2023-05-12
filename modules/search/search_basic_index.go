// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/core/security"
	"github.com/rande/gonode/modules/base"
)

type IndexMeta struct {
	Format string `json:"format"`
}

type Index struct {
	Page       int64    `json:"page"`
	PerPage    int64    `json:"per_page"`
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

type IndexHandler struct {
}

func (h *IndexHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Index{
		Deleted: NewParam(false, "="),
	}, &IndexMeta{}
}

func (h *IndexHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *IndexHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *IndexHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *IndexHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *IndexHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {

}

func (h *IndexHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	return base.GetDownloadData()
}

func (h *IndexHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *IndexHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return base.DefaultHandlerStoreStream(node, r)
}

type IndexViewHandler struct {
	Search    *SearchPGSQL
	Manager   base.NodeManager
	MaxResult uint64
}

func (v *IndexViewHandler) Support(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) bool {
	return request.Format == "html"
}

func (v *IndexViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	var err error

	index := node.Data.(*Index)

	// we just copy over node to create search form
	search := NewSearchFormFromIndex(index)

	if v := request.HttpRequest.URL.Query().Get("per_page"); len(v) > 0 {
		if search.PerPage, err = strconv.ParseUint(v, 10, 32); err != nil {
			return err
		}
	}

	if v := request.HttpRequest.URL.Query().Get("page"); len(v) > 0 {
		if search.Page, err = strconv.ParseUint(v, 10, 32); err != nil {
			return err
		}
	}

	// check page range
	if uint64(search.PerPage) > v.MaxResult {
		helper.SendWithHttpCode(response.HttpResponse, http.StatusPreconditionFailed, "Invalid `pagination` range")

		return nil
	}

	if search.Page == 0 {
		search.Page = uint64(1)
	}

	if search.PerPage == 0 {
		search.PerPage = uint64(32)
	}

	options := base.NewAccessOptionsFromToken(security.GetTokenFromContext(request.Context))
	pager := GetPager(search, v.Manager, v.Search, options)

	response.
		Set(200, fmt.Sprintf("search:nodes/%s.tpl", node.Type)).
		Add("pager", pager)

	return nil
}
