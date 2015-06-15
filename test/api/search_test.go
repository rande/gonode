package api

import (
	"encoding/json"
	"github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/extra"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_Search_Basic(t *testing.T) {
	urls := []string{
		"/nodes",
		"/nodes?type=core.user",
		"/nodes?type=core.user&data.login=user12",
		"/nodes?type=core.user&data.login=user12&data.login=user13",
	}

	for _, url := range urls {
		extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
			// WITH
			file, _ := os.Open("../fixtures/new_user.json")
			extra.RunRequest("POST", ts.URL+"/nodes", file)

			// WHEN
			res, _ := extra.RunRequest("GET", ts.URL+url, nil)

			p := &nc.ApiPager{}

			serializer := app.Get("gonode.node.serializer").(*nc.Serializer)
			serializer.Deserialize(res.Body, p)

			// THEN
			assert.Equal(t, uint64(32), p.PerPage)
			assert.Equal(t, uint64(1), p.Page)
			assert.Equal(t, 1, len(p.Elements))
			assert.Equal(t, uint64(0), p.Next)
			assert.Equal(t, uint64(0), p.Previous)

			// the Element is a [string]interface so we need to convert it back to []byte
			// and then unmarshal again with the correct structure
			raw, _ := json.Marshal(p.Elements[0])

			n := nc.NewNode()
			json.Unmarshal(raw, n)

			assert.Equal(t, "core.user", n.Type)
			assert.False(t, n.Deleted)
		})
	}
}
