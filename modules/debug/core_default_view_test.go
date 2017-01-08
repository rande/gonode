// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package debug

import (
	"testing"

	"github.com/flosch/pongo2"
	"github.com/rande/gonode/modules/base"
	"github.com/stretchr/testify/assert"
)

func Test_Default_View(t *testing.T) {
	node := &base.Node{
		Type: "foo.bar",
	}

	request := &base.ViewRequest{}

	response := &base.ViewResponse{
		StatusCode: 200,
		Context:    pongo2.Context{},
	}

	v := &DefaultViewHandler{}

	err := v.Execute(node, request, response)

	assert.NoError(t, err)
	assert.Equal(t, node, response.Context["node"])
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "nodes/foo.bar.tpl", response.Template)
}
