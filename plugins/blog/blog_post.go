// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package blog

import (
	"github.com/rande/gonode/core"
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

func (h *PostHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &Post{}, &PostMeta{}
}

func (h *PostHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *PostHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *PostHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *PostHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *PostHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *PostHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *PostHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *PostHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}
