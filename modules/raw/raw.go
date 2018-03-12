// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package raw

import (
	"fmt"

	"github.com/rande/gonode/modules/base"
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
