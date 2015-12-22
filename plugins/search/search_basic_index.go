// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"github.com/rande/gonode/core"
	"io"
)

type IndexMeta struct {
	Format string `json:"format"`
}

type Index struct {
	Page       int64               `jsonjson:"page"`
	PerPage    int64               `json:"per_page"`
	OrderBy    []string            `json:"order_by"`
	Uuid       string              `json:"uuid"`
	Type       []string            `json:"type"`
	Name       string              `json:"name"`
	Slug       string              `json:"slug"`
	Data       map[string][]string `json:"data"`
	Meta       map[string][]string `json:"meta"`
	Status     []string            `json:"status"`
	Weight     []string            `json:"weight"`
	Revision   string              `json:"revision"`
	Enabled    string              `json:"enabled"`
	Deleted    bool                `json:"deleted"`
	Current    string              `json:"current"`
	UpdatedBy  []string            `json:"updated_by"`
	CreatedBy  []string            `json:"created_by"`
	ParentUuid []string            `json:"parent_uuid"`
	SetUuid    []string            `json:"set_uuid"`
	Source     []string            `json:"source"`
}

type IndexHandler struct {
}

func (h *IndexHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &Index{
		Deleted: false,
	}, &IndexMeta{}
}

func (h *IndexHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *IndexHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *IndexHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *IndexHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *IndexHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *IndexHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *IndexHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *IndexHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}
