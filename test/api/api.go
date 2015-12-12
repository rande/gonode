// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/plugins/api"
	"github.com/rande/gonode/plugins/user"
	"github.com/rande/gonode/test"
)

func GetPager(app *goapp.App, res *test.Response) *api.ApiPager {
	p := &api.ApiPager{}

	serializer := app.Get("gonode.node.serializer").(*core.Serializer)
	serializer.Deserialize(res.Body, p)

	// the Element is a [string]interface so we need to convert it back to []byte
	// and then unmarshal again with the correct structure
	for k, v := range p.Elements {
		raw, _ := json.Marshal(v)

		n := core.NewNode()
		json.Unmarshal(raw, n)

		p.Elements[k] = n
	}

	return p
}

func InitSearchFixture(app *goapp.App) []*core.Node {
	manager := app.Get("gonode.manager").(*core.PgNodeManager)
	collection := app.Get("gonode.handler_collection").(core.Handlers)
	nodes := make([]*core.Node, 0)

	// WITH 3 nodes
	node := collection.NewNode("core.user")
	node.Name = "User A"
	node.Weight = 1
	node.Slug = "user-a"
	node.Data.(*user.User).FirstName = "User"
	node.Data.(*user.User).LastName = "A"
	node.Data.(*user.User).Username = "user-a"
	manager.Save(node, false)

	nodes = append(nodes, node)

	node = collection.NewNode("core.user")
	node.Name = "User AA"
	node.Weight = 2
	node.Slug = "user-aa"
	node.Data.(*user.User).FirstName = "User"
	node.Data.(*user.User).LastName = "AA"
	node.Data.(*user.User).Username = "user-aa"
	manager.Save(node, false)

	nodes = append(nodes, node)

	node = collection.NewNode("core.user")
	node.Name = "User B"
	node.Weight = 1
	node.Slug = "user-b"
	node.Data.(*user.User).FirstName = "User"
	node.Data.(*user.User).LastName = "B"
	node.Data.(*user.User).Username = "user-b"
	manager.Save(node, false)

	nodes = append(nodes, node)

	return nodes
}
