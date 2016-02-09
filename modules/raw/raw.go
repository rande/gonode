// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"
	"github.com/rande/gonode/modules/base"
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

func (h *RawHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	return &Raw{}, &RawMeta{}
}

func (h *RawHandler) PreInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *RawHandler) PreUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *RawHandler) PostInsert(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *RawHandler) PostUpdate(node *base.Node, m base.NodeManager) error {
	return nil
}

func (h *RawHandler) Validate(node *base.Node, m base.NodeManager, errors base.Errors) {
}

func (h *RawHandler) GetDownloadData(node *base.Node) *base.DownloadData {
	return base.GetDownloadData()
}

func (h *RawHandler) Load(data []byte, meta []byte, node *base.Node) error {
	return base.HandlerLoad(h, data, meta, node)
}

func (h *RawHandler) StoreStream(node *base.Node, r io.Reader) (int64, error) {
	return base.DefaultHandlerStoreStream(node, r)
}

type RawViewHandler struct {
}

func (v *RawViewHandler) Support(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) bool {
	return true
}

func (v *RawViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {
	raw := node.Data.(*Raw)

	values := request.HttpRequest.URL.Query()

	if _, ok := values["dl"]; ok { // ask for binary content
		response.HttpResponse.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", raw.Name))
	}

	response.HttpResponse.Header().Set("Content-Type", raw.ContentType)
	response.HttpResponse.Write(raw.Content)

	return nil
}
