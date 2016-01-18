// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"github.com/rande/gonode/core"
	"io"
)

type DefaultHandler struct {
}

func (h *DefaultHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	data := make(map[string]interface{})
	meta := make(map[string]interface{})

	return &data, &meta
}

func (h *DefaultHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *DefaultHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *DefaultHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *DefaultHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *DefaultHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}
