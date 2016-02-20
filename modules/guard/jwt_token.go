// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"github.com/rande/gonode/modules/base"
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
