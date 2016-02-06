// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Prism_Blog_Archive(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		// create a valid user into the database ...
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"

		manager.Save(node, false)

		archive := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.index")
		archive.Name = "Blog Archive"

		manager.Save(archive, false)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/prism/%s", ts.URL, archive.Uuid))

		assert.Equal(t, 200, res.StatusCode)
	})
}
