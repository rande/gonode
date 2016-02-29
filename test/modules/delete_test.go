// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_Delete_Non_Existant_Node(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("DELETE", ts.URL+"/api/v1.0/nodes/d703a3ab-8374-4c30-a8a4-2c22aa67763b", nil, auth)

		assert.Equal(t, 404, res.StatusCode, "Delete non existant node")
	})
}

func Test_Delete_Existant_Node(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		file, _ := os.Open("../fixtures/new_user.json")
		res, _ := test.RunRequest("POST", ts.URL+"/api/v1.0/nodes", file, auth)

		assert.Equal(t, 201, res.StatusCode, "Node created")

		node := base.NewNode()

		serializer := app.Get("gonode.node.serializer").(*base.Serializer)
		serializer.Deserialize(res.Body, node)

		assert.Equal(t, "core.user", node.Type)

		res, _ = test.RunRequest("DELETE", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString(), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		serializer.Deserialize(res.Body, node)

		assert.Equal(t, node.Deleted, true)

		// test if we can delete a deleted node ...
		res, _ = test.RunRequest("DELETE", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString(), nil, auth)
		assert.Equal(t, 410, res.StatusCode)
	})
}

func Test_Delete_Find_Filter(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)
		nodes := InitSearchFixture(app)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/nodes", nil, auth)
		p := test.GetPager(app, res)

		assert.Equal(t, 4, len(p.Elements))

		res, _ = test.RunRequest("DELETE", ts.URL+"/api/v1.0/nodes/"+nodes[0].Uuid.CleanString(), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("GET", ts.URL+"/api/v1.0/nodes", nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		p = test.GetPager(app, res)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, 3, len(p.Elements))
	})
}
