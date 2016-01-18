// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package guard

import (
	"github.com/rande/gonode/core"
	"io"
	"time"
)

type JwtTokenMeta struct {
	Expiration time.Time `json:"expiration"`
}

type JwtToken struct {
	User  *core.Reference `json:"user"`
	Key   []byte          `json:"key"`
	Roles []string        `json:"roles"`
}

type JwtTokentHandler struct {
}

func (h *JwtTokentHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &JwtToken{}, &JwtTokenMeta{}
}

func (h *JwtTokentHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *JwtTokentHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {

}

func (h *JwtTokentHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *JwtTokentHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *JwtTokentHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}
