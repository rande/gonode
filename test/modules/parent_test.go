// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"github.com/rande/goapp"
	"github.com/rande/gonode/modules/api"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Create_Parents_With_Manager(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		handlers := app.Get("gonode.handler_collection").(base.HandlerCollection)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		node1 := handlers.NewNode("default")
		manager.Save(node1, false)

		node2 := handlers.NewNode("default")
		manager.Save(node2, false)

		node3 := handlers.NewNode("default")
		manager.Save(node3, false)

		node4 := handlers.NewNode("default")
		manager.Save(node4, false)

		// WHEN
		affectedRows, err := manager.Move(node2.Uuid, node1.Uuid)
		assert.Nil(t, err)
		assert.Equal(t, affectedRows, int64(1))

		affectedRows, err = manager.Move(node3.Uuid, node2.Uuid)
		assert.Nil(t, err)
		assert.Equal(t, affectedRows, int64(1))

		affectedRows, err = manager.Move(node4.Uuid, node3.Uuid)
		assert.Nil(t, err)
		assert.Equal(t, affectedRows, int64(1))

		// cannot move a parent node into its child
		affectedRows, err = manager.Move(node1.Uuid, node4.Uuid)
		assert.Nil(t, err)
		assert.Equal(t, affectedRows, int64(0))

		// retrieve a saved node

		node := manager.Find(node4.Uuid)

		assert.Equal(t, 3, len(node.Parents))
		assert.Contains(t, node.Parents, node1.Uuid)
		assert.Contains(t, node.Parents, node2.Uuid)
		assert.Contains(t, node.Parents, node3.Uuid)
		assert.NotContains(t, node.Parents, node4.Uuid)
	})
}

func Test_Create_Parents_With_Api(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		auth := test.GetAuthHeader(t, ts)

		handlers := app.Get("gonode.handler_collection").(base.HandlerCollection)
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		node1 := handlers.NewNode("default")
		manager.Save(node1, false)

		node2 := handlers.NewNode("default")
		manager.Save(node2, false)

		node3 := handlers.NewNode("default")
		manager.Save(node3, false)

		node4 := handlers.NewNode("default")
		manager.Save(node4, false)

		res, _ := test.RunRequest("PUT", fmt.Sprintf("%s/nodes/move/%s/%s", ts.URL, node2.Uuid, node1.Uuid), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("PUT", fmt.Sprintf("%s/nodes/move/%s/%s", ts.URL, node3.Uuid, node2.Uuid), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("PUT", fmt.Sprintf("%s/nodes/move/%s/%s", ts.URL, node4.Uuid, node3.Uuid), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("PUT", fmt.Sprintf("%s/nodes/move/%s/%s", ts.URL, node1.Uuid, node4.Uuid), nil, auth)
		assert.Equal(t, 200, res.StatusCode)

		serializer := app.Get("gonode.node.serializer").(*base.Serializer)
		op := &api.ApiOperation{}
		serializer.Deserialize(res.Body, op)

		assert.Equal(t, "OK", op.Status)
		assert.Equal(t, "Node altered: 0", op.Message)

	})
}
