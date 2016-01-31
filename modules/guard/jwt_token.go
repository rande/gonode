// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"github.com/rande/gonode/modules/base"
	"io"
	"time"
)

type JwtTokenMeta struct {
	Expiration time.Time `json:"expiration"`
}

type JwtToken struct {
	User  *base.Reference `json:"user"`
	Key   []byte          `json:"key"`
	Roles []string        `json:"roles"`
}

type JwtTokentHandler struct {
}

func (h *JwtTokentHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &JwtToken{}, &JwtTokenMeta{}
}

func (h *JwtTokentHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {

}

func (h *JwtTokentHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	return base.GetDownloadData()
}

func (h *JwtTokentHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *JwtTokentHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return base.DefaultHandlerStoreStream(node, r)
}
