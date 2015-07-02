package api

import (
	. "github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/rande/gonode/extra"
	nh "github.com/rande/gonode/handlers"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_Create_User(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		// WITH
		file, _ := os.Open("../fixtures/new_user.json")
		res, _ := extra.RunRequest("POST", ts.URL+"/nodes", file)

		assert.Equal(t, 201, res.StatusCode)

		// WHEN
		node := nc.NewNode()
		serializer := app.Get("gonode.node.serializer").(*nc.Serializer)
		serializer.Deserialize(res.Body, node)

		// THEN
		assert.Equal(t, node.Type, "core.user")

		user := node.Data.(*nh.User)

		assert.Equal(t, user.FirstName, "User")
		assert.Equal(t, user.LastName, "12")
	})
}

func Test_Create_Media_With_Binary_Upload(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		// WITH
		file, _ := os.Open("../fixtures/new_image.json")
		res, _ := extra.RunRequest("POST", ts.URL+"/nodes", file)

		assert.Equal(t, 201, res.StatusCode)

		node := nc.NewNode()
		serializer := app.Get("gonode.node.serializer").(*nc.Serializer)
		serializer.Deserialize(res.Body, node)

		var message = "The content of the file, yep it is not an image"

		res, _ = extra.RunRequest("PUT", ts.URL+"/nodes/"+node.Uuid.CleanString()+"?raw", strings.NewReader(message))

		assert.Equal(t, 200, res.StatusCode)

		res, _ = extra.RunRequest("GET", ts.URL+"/nodes/"+node.Uuid.CleanString()+"?raw", nil)

		assert.Equal(t, message, string(res.GetBody()[:]))

		res, _ = extra.RunRequest("GET", ts.URL+"/nodes/"+node.Uuid.CleanString(), nil)
		assert.Equal(t, 200, res.StatusCode)

		node = nc.NewNode()
		serializer.Deserialize(res.Body, node)

		meta := node.Meta.(*nh.ImageMeta)
		assert.Equal(t, "media.image", node.Type)
		assert.Equal(t, "application/octet-stream", meta.ContentType)
	})
}
