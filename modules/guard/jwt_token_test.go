// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node_guard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JwtToken_Handler(t *testing.T) {
	h := &JwtTokenHandler{}

	node, meta := h.GetStruct()

	assert.NotNil(t, node)
	assert.NotNil(t, meta)
}
