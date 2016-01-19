// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package node

import (
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Valid_UpdatedAt(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		node := app.Get("gonode.handler_collection").(core.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"

		manager.Save(node, false)

		updatedAt := node.UpdatedAt

		manager.Save(node, false)

		assert.NotEqual(t, updatedAt, node.UpdatedAt)
	})
}

func Test_New_Revision(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		node := app.Get("gonode.handler_collection").(core.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"

		manager.Save(node, false)
		assert.Equal(t, 1, node.Revision)

		manager.Save(node, false)
		assert.Equal(t, 1, node.Revision)

		manager.Save(node, true)
		assert.Equal(t, 2, node.Revision)

		manager.Save(node, true)
		assert.Equal(t, 3, node.Revision)

		manager.Save(node, false)
		assert.Equal(t, 3, node.Revision)

	})
}