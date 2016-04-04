// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/blog"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Access_FindOne(t *testing.T) {

	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		manager := app.Get("gonode.manager").(*base.PgNodeManager)
		node := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("blog.post")
		data := node.Data.(*blog.Post)
		data.Title = "Blog Post 1"
		node.Access = []string{"no.role"}

		manager.Save(node, false)

		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/%s", ts.URL, node.Uuid.String()), nil, auth)

		assert.Equal(t, 404, res.StatusCode, "Should not find node as roles no not match")
	})
}
