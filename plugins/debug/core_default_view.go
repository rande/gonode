// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"fmt"
	"github.com/rande/gonode/core"
)

type DefaultViewHandler struct {
}

func (v *DefaultViewHandler) Execute(node *core.Node, request *core.ViewRequest, response *core.ViewResponse) error {

	response.
		Set(200, fmt.Sprintf("nodes/%s.tpl", node.Type)).
		Add("node", node)

	return nil
}
