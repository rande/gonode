// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"github.com/rande/gonode/modules/base"
	"io"
)

type DefaultHandler struct {
}

func (h *DefaultHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	data := make(map[string]interface{})
	meta := make(map[string]interface{})

	return &data, &meta
}

func (h *DefaultHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *DefaultHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *DefaultHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {

}

func (h *DefaultHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	return base.GetDownloadData()
}

func (h *DefaultHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *DefaultHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return base.DefaultHandlerStoreStream(node, r)
}
