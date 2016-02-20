// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"github.com/rande/gonode/modules/base"
)

type DefaultHandler struct {
}

func (h *DefaultHandler) GetStruct() (base.NodeData, base.NodeMeta) {
	data := make(map[string]interface{})
	meta := make(map[string]interface{})

	return &data, &meta
}
