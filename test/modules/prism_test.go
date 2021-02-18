// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Prism_Blog_Archive(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// GIVEN
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(node, false)

		archive := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.index")
		archive.Name = "Blog Archive"
		archive.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(archive, false)

		// WHEN
		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/prism/%s", ts.URL, archive.Uuid))

		// THEN
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func Test_Prism_Bad_Request(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		// create a valid user into the database ...
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(node, false)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/prism/%s.json", ts.URL, node.Uuid))

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func Test_Prism_Format(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		home := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.root")
		home.Name = "Homepage"
		home.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(home, false)

		raw := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.raw")
		raw.Name = "Humans.txt"
		raw.Slug = "humans.txt"
		raw.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(raw, false)
		manager.Move(raw.Uuid, home.Uuid)

		raw2 := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.raw")
		raw2.Name = "Humans"
		raw2.Slug = "humans"
		raw2.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(raw2, false)
		manager.Move(raw2.Uuid, home.Uuid)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/humans", ts.URL))

		assert.Equal(t, http.StatusOK, res.StatusCode, "Cannot find /humans")

		res, _ = test.RunRequest("GET", fmt.Sprintf("%s/humans.txt", ts.URL))

		assert.Equal(t, http.StatusOK, res.StatusCode, "Cannot find /humans.txt")
	})
}

func Test_Prism_Forbidden(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		home := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.root")
		home.Name = "Homepage"
		home.Access = []string{"node:prism:render", "IS_AUTHENTICATED_ANONYMOUSLY"}

		manager.Save(home, false)

		raw := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.raw")
		raw.Name = "Humans.txt"
		raw.Slug = "humans.txt"
		raw.Access = []string{}

		manager.Save(raw, false)
		manager.Move(raw.Uuid, home.Uuid)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/humans.txt", ts.URL))

		assert.Equal(t, http.StatusForbidden, res.StatusCode)
	})
}
