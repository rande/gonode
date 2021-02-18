// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
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
	"github.com/rande/gonode/modules/user"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Access_FindOne_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		// create dummy user
		u := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.user")
		u.Name = "User Dummy"

		dataUser := u.Data.(*user.User)
		dataUser.Email = "test-dummy@example.org"
		dataUser.Enabled = true
		dataUser.NewPassword = "dummy"
		dataUser.Username = "dummy"
		dataUser.Roles = []string{"ROLE_API"}

		metaUser := u.Meta.(*user.UserMeta)
		metaUser.PasswordCost = 1 // save test time

		manager.Save(u, false)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"no.role"}

		manager.Save(node, false)

		auth := test.GetAuthHeaderFromCredentials("dummy", "dummy", ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		fmt.Print(res.GetBodyAsString())

		assert.Equal(t, 403, res.StatusCode, "Should not find a node as roles no not match")
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

		auth := test.GetDefaultAuthHeader(ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 200, res.StatusCode, "Should find a node as roles match")
	})
}

func Test_Access_RemoveOne_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		// create dummy user
		u := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.user")
		u.Name = "User Dummy"

		dataUser := u.Data.(*user.User)
		dataUser.Email = "test-dummy@example.org"
		dataUser.Enabled = true
		dataUser.NewPassword = "dummy"
		dataUser.Username = "dummy"
		dataUser.Roles = []string{"ROLE_API"}

		metaUser := u.Meta.(*user.UserMeta)
		metaUser.PasswordCost = 1 // save test time

		manager.Save(u, false)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"no.role"}

		manager.Save(node, false)

		auth := test.GetAuthHeaderFromCredentials("dummy", "dummy", ts)

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

		auth := test.GetDefaultAuthHeader(ts)

		res, _ := test.RunRequest("DELETE", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 200, res.StatusCode)
	})
}

func Test_Access_Find_NoResult(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		// create dummy user
		u := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.user")
		u.Name = "User Dummy"

		dataUser := u.Data.(*user.User)
		dataUser.Email = "test-dummy@example.org"
		dataUser.Enabled = true
		dataUser.NewPassword = "dummy"
		dataUser.Username = "dummy"
		dataUser.Roles = []string{"ROLE_API", "node:api:list"}

		metaUser := u.Meta.(*user.UserMeta)
		metaUser.PasswordCost = 1 // save test time

		manager.Save(u, false)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"

		manager.Save(node, false)

		auth := test.GetAuthHeaderFromCredentials("dummy", "dummy", ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes?type=blog.post", ts.URL), nil, auth)

		assert.Equal(t, 200, res.StatusCode)

		p := test.GetPager(app, res)

		assert.Equal(t, 0, len(p.Elements))
	})
}
