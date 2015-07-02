package api

import (
	"encoding/json"
	"github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/extra"
	"github.com/rande/gonode/handlers"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func GetPager(app *goapp.App, res *extra.Response) *nc.ApiPager {
	p := &nc.ApiPager{}

	serializer := app.Get("gonode.node.serializer").(*nc.Serializer)
	serializer.Deserialize(res.Body, p)

	// the Element is a [string]interface so we need to convert it back to []byte
	// and then unmarshal again with the correct structure
	for k, v := range p.Elements {
		raw, _ := json.Marshal(v)

		n := nc.NewNode()
		json.Unmarshal(raw, n)

		p.Elements[k] = n
	}

	return p
}

func CheckNoResults(t *testing.T, p *nc.ApiPager) {
	assert.Equal(t, uint64(32), p.PerPage)
	assert.Equal(t, uint64(1), p.Page)
	assert.Equal(t, 0, len(p.Elements))
	assert.Equal(t, uint64(0), p.Next)
	assert.Equal(t, uint64(0), p.Previous)
}

func Test_Search_Basic(t *testing.T) {
	urls := []string{
		"/nodes",
		"/nodes?type=core.user",
		"/nodes?type=core.user&data.login=user12",
		"/nodes?type=core.user&data.login=user12&data.login=user13",
		"/nodes?&page=-1&page=1", // the last occurrence erase first values
	}

	for _, url := range urls {
		extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
			// WITH
			file, _ := os.Open("../fixtures/new_user.json")
			extra.RunRequest("POST", ts.URL+"/nodes", file)

			// WHEN
			res, _ := extra.RunRequest("GET", ts.URL+url, nil)

			p := GetPager(app, res)

			// THEN
			assert.Equal(t, uint64(32), p.PerPage)
			assert.Equal(t, uint64(1), p.Page)
			assert.Equal(t, 1, len(p.Elements))
			assert.Equal(t, uint64(0), p.Next)
			assert.Equal(t, uint64(0), p.Previous)

			n := p.Elements[0].(*nc.Node)

			assert.Equal(t, "core.user", n.Type)
			assert.False(t, n.Deleted)
		})
	}
}

func Test_Search_NoResult(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// WITH
		file, _ := os.Open("../fixtures/new_user.json")
		extra.RunRequest("POST", ts.URL+"/nodes", file)

		// WHEN
		res, _ := extra.RunRequest("GET", ts.URL+"/nodes?type=other", nil)

		p := GetPager(app, res)

		// THEN
		CheckNoResults(t, p)
	})
}

func Test_Search_Invalid_Pagination(t *testing.T) {
	urls := []string{
		"/nodes?per_page=-1",
		"/nodes?per_page=-1&page=-1",
		"/nodes?per_page=256",
		"/nodes?per_page=256&page=1",
		"/nodes?per_page=127&page=1&page=-1",
		// "/nodes?per_page=127&page=-1&page=1", // the last occurrence erase first values
	}

	for _, url := range urls {
		extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
			file, _ := os.Open("../fixtures/new_user.json")
			extra.RunRequest("POST", ts.URL+"/nodes", file)

			// WHEN
			res, _ := extra.RunRequest("GET", ts.URL+url, nil)

			assert.Equal(t, 412, res.StatusCode, "url: "+url)
		})
	}
}

func Test_Search_Invalid_OrderBy(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		// seems goji or golang block this request
		res, _ := extra.RunRequest("GET", ts.URL+"/nodes?order_by=\"1 = 1\"; DELETE * FROM node,ASC", nil)
		assert.Equal(t, 400, res.StatusCode, "url: /nodes?order_by=\"1 = 1\"; DELETE * FROM node,ASC")

		// seems goji or golang block this request
		res, _ = extra.RunRequest("GET", ts.URL+"/nodes?order_by=DELETE%20*%20FROM%20node,ASC", nil)
		assert.Equal(t, 412, res.StatusCode, "url: /nodes?order_by=DELETE%20*%20FROM%20node,ASC")
	})
}

func Test_Search_OrderBy(t *testing.T) {

	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		manager := app.Get("gonode.manager").(*nc.PgNodeManager)
		collection := app.Get("gonode.handler_collection").(nc.Handlers)

		// WITH 3 nodes
		node := collection.NewNode("core.user")
		node.Name = "User A"
		node.Weight = 1
		node.Data.(*handlers.User).FirstName = "User"
		node.Data.(*handlers.User).LastName = "A"
		node.Data.(*handlers.User).Login = "user-a"
		manager.Save(node)

		node = collection.NewNode("core.user")
		node.Name = "User AA"
		node.Weight = 2
		node.Data.(*handlers.User).FirstName = "User"
		node.Data.(*handlers.User).LastName = "AA"
		node.Data.(*handlers.User).Login = "user-aa"
		manager.Save(node)

		node = collection.NewNode("core.user")
		node.Name = "User B"
		node.Weight = 1
		node.Data.(*handlers.User).FirstName = "User"
		node.Data.(*handlers.User).LastName = "B"
		node.Data.(*handlers.User).Login = "user-b"
		manager.Save(node)

		// TESTING ASC ORDERING
		res, _ := extra.RunRequest("GET", ts.URL+"/nodes?order_by=name,ASC", nil)

		assert.Equal(t, 200, res.StatusCode, "url: /nodes?order_by=name,ASC")

		p := GetPager(app, res)

		assert.Equal(t, 3, len(p.Elements))
		assert.Equal(t, "User A", p.Elements[0].(*nc.Node).Name)
		assert.Equal(t, "User AA", p.Elements[1].(*nc.Node).Name)
		assert.Equal(t, "User B", p.Elements[2].(*nc.Node).Name)

		// TESTING DESC ORDERING
		res, _ = extra.RunRequest("GET", ts.URL+"/nodes?order_by=name,DESC", nil)

		assert.Equal(t, 200, res.StatusCode, "url: /nodes?order_by=name,DESC")

		p = GetPager(app, res)

		assert.Equal(t, 3, len(p.Elements))
		assert.Equal(t, "User B", p.Elements[0].(*nc.Node).Name)
		assert.Equal(t, "User AA", p.Elements[1].(*nc.Node).Name)
		assert.Equal(t, "User A", p.Elements[2].(*nc.Node).Name)

		// TESTING DESC ORDERING
		res, _ = extra.RunRequest("GET", ts.URL+"/nodes?order_by=name,DESC", nil)

		assert.Equal(t, 200, res.StatusCode, "url: /nodes?order_by=name,DESC")

		p = GetPager(app, res)

		assert.Equal(t, 3, len(p.Elements))
		assert.Equal(t, "User B", p.Elements[0].(*nc.Node).Name)
		assert.Equal(t, "User AA", p.Elements[1].(*nc.Node).Name)
		assert.Equal(t, "User A", p.Elements[2].(*nc.Node).Name)

		// TESTING WITH 2 ORDERING OPTION
		res, _ = extra.RunRequest("GET", ts.URL+"/nodes?order_by=weight,DESC&order_by=name,ASC", nil)

		assert.Equal(t, 200, res.StatusCode, "url: /nodes?order_by=weight,DESC&order_by=name,ASC")

		p = GetPager(app, res)

		assert.Equal(t, 3, len(p.Elements))
		assert.Equal(t, "User AA", p.Elements[0].(*nc.Node).Name)
		assert.Equal(t, "User A", p.Elements[1].(*nc.Node).Name)
		assert.Equal(t, "User B", p.Elements[2].(*nc.Node).Name)
	})
}
