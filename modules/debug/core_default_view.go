// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"fmt"
	"github.com/rande/gonode/modules/base"
)

type DefaultViewHandler struct {
}

func (v *DefaultViewHandler) Execute(node *base.Node, request *base.ViewRequest, response *base.ViewResponse) error {

	response.
		Set(200, fmt.Sprintf("nodes/%s.tpl", node.Type)).
		Add("node", node)

	return nil
}
