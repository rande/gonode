// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	. "github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_Delete_Non_Existant_Node(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		res, _ := test.RunRequest("DELETE", ts.URL+"/nodes/d703a3ab-8374-4c30-a8a4-2c22aa67763b", nil)

		assert.Equal(t, 404, res.StatusCode, "Delete non existant node")
	})
}

func Test_Delete_Existant_Node(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		file, _ := os.Open("../fixtures/new_user.json")
		res, _ := test.RunRequest("POST", ts.URL+"/nodes", file)

		assert.Equal(t, 201, res.StatusCode, "Node created")

		node := nc.NewNode()

		serializer := app.Get("gonode.node.serializer").(*nc.Serializer)
		serializer.Deserialize(res.Body, node)

		assert.Equal(t, "core.user", node.Type)

		res, _ = test.RunRequest("DELETE", ts.URL+"/nodes/"+node.Uuid.CleanString(), nil)
		assert.Equal(t, 200, res.StatusCode)

		serializer.Deserialize(res.Body, node)

		assert.Equal(t, node.Deleted, true)

		// test if we can delete a deleted node ...
		res, _ = test.RunRequest("DELETE", ts.URL+"/nodes/"+node.Uuid.CleanString(), nil)
		assert.Equal(t, 410, res.StatusCode)
	})
}

func Test_Delete_Find_Filter(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		nodes := InitSearchFixture(app)

		res, _ := test.RunRequest("GET", ts.URL+"/nodes", nil)
		p := GetPager(app, res)

		assert.Equal(t, 3, len(p.Elements))

		res, _ = test.RunRequest("DELETE", ts.URL+"/nodes/"+nodes[0].Uuid.CleanString(), nil)
		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("GET", ts.URL+"/nodes", nil)
		assert.Equal(t, 200, res.StatusCode)

		p = GetPager(app, res)

		assert.Equal(t, 2, len(p.Elements))
	})
}
