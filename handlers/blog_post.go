// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	nc "github.com/rande/gonode/core"
	"github.com/spf13/afero"
	"io"
)

type PostMeta struct {
	Format string `json:"format"`
}

type Post struct {
	Title    string   `json:"title"`
	SubTitle string   `json:"sub_title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
}

type PostHandler struct {
}

func (h *PostHandler) GetStruct() (nc.NodeData, nc.NodeMeta) {
	return &Post{}, &PostMeta{}
}

func (h *PostHandler) PreInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PreUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PostInsert(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) PostUpdate(node *nc.Node, m nc.NodeManager) error {
	return nil
}

func (h *PostHandler) Validate(node *nc.Node, m nc.NodeManager, errors nc.Errors) {

}

func (h *PostHandler) GetDownloadData(node *nc.Node) *nc.DownloadData {
	return nc.GetDownloadData()
}

func (h *PostHandler) Load(data []byte, meta []byte, node *nc.Node) error {
	return nc.HandlerLoad(h, data, meta, node)
}

func (h *PostHandler) StoreStream(node *nc.Node, r io.Reader) (afero.File, int64, error) {
	return nc.DefaultHandlerStoreStream(node, r)
}
