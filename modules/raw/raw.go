// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"
	"github.com/rande/gonode/core"
	"io"
)

type Raw struct {
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
	Name        string `json:"name"`
}

type RawMeta struct {
}

type RawHandler struct {
}

func (h *RawHandler) GetStruct() (core.NodeData, core.NodeMeta) {
	return &Raw{}, &RawMeta{}
}

func (h *RawHandler) PreInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *RawHandler) PreUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *RawHandler) PostInsert(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *RawHandler) PostUpdate(node *core.Node, m core.NodeManager) error {
	return nil
}

func (h *RawHandler) Validate(node *core.Node, m core.NodeManager, errors core.Errors) {
}

func (h *RawHandler) GetDownloadData(node *core.Node) *core.DownloadData {
	return core.GetDownloadData()
}

func (h *RawHandler) Load(data []byte, meta []byte, node *core.Node) error {
	return core.HandlerLoad(h, data, meta, node)
}

func (h *RawHandler) StoreStream(node *core.Node, r io.Reader) (int64, error) {
	return core.DefaultHandlerStoreStream(node, r)
}

type RawViewHandler struct {
}

func (v *RawViewHandler) Execute(node *core.Node, request *core.ViewRequest, response *core.ViewResponse) error {
	raw := node.Data.(*Raw)

	values := request.HttpRequest.URL.Query()

	if _, ok := values["dl"]; ok { // ask for binary content
		response.HttpResponse.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", raw.Name))
	}

	response.HttpResponse.Header().Set("Content-Type", raw.ContentType)
	response.HttpResponse.Write(raw.Content)

	return nil
}
