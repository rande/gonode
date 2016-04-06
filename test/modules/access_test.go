// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Access_FindOne_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"no.role"}

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 403, res.StatusCode, "Should not find node as roles no not match")
	})
}

func Test_Access_FindOne_Result(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"node:api:master"}

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 200, res.StatusCode, "Should not find node as roles no not match")
	})
}

func Test_Access_RemoveOne_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"no.role"}

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("DELETE", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 403, res.StatusCode)
	})
}

func Test_Access_RemoveOne_Result(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"node:api:master"}

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("DELETE", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 200, res.StatusCode)
	})
}

func Test_Access_Find_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes?type=blog.post", ts.URL), nil, auth)

		assert.Equal(t, 200, res.StatusCode)

		p := test.GetPager(app, res)

		assert.Equal(t, 0, len(p.Elements))
	})
}
