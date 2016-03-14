// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"image"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/rande/goapp"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/media"
	"github.com/rande/gonode/modules/user"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Create_User(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		// WITH
		file, _ := os.Open("../fixtures/new_user.json")
		res, _ := test.RunRequest("POST", ts.URL+"/api/v1.0/nodes", file, auth)

		assert.Equal(t, 201, res.StatusCode, "unable to create an user")

		// WHEN
		node := base.NewNode()
		serializer := app.Get("gonode.node.serializer").(*base.Serializer)
		serializer.Deserialize(res.Body, node)

		// THEN
		assert.Equal(t, node.Type, "core.user")

		user := node.Data.(*user.User)

		assert.Equal(t, user.FirstName, "User")
		assert.Equal(t, user.LastName, "12")
	})
}

func Test_Create_Media_With_Binary_Upload(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		// WITH
		file, _ := os.Open("../fixtures/new_image.json")
		res, _ := test.RunRequest("POST", ts.URL+"/api/v1.0/nodes", file, auth)

		assert.Equal(t, 201, res.StatusCode)

		node := base.NewNode()
		serializer := app.Get("gonode.node.serializer").(*base.Serializer)
		serializer.Deserialize(res.Body, node)

		file, _ = os.Open("../fixtures/photo.jpg")

		res, _ = test.RunRequest("PUT", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString()+"?raw", file, auth)

		assert.Equal(t, 200, res.StatusCode)

		res, _ = test.RunRequest("GET", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString()+"?raw", nil, auth)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "image/jpeg", res.Header.Get("Content-Type"))

		res, _ = test.RunRequest("GET", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString(), nil, auth)
		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

		node = base.NewNode()
		serializer.Deserialize(res.Body, node)

		meta := node.Meta.(*media.ImageMeta)
		data := node.Data.(*media.Image)

		assert.Equal(t, "media.image", node.Type)
		assert.Equal(t, "image/jpeg", meta.ContentType)
		assert.Equal(t, 1024, meta.Width)
		assert.Equal(t, 768, meta.Height)

		assert.Equal(t, "the_first_media.jpg", data.Name)
		assert.Equal(t, "Canon", meta.Exif["Make"])
		assert.Equal(t, "Canon PowerShot G12", meta.Exif["Model"])
	})
}

func Test_Media_Resize_With_Orientation(t *testing.T) {

	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		message := fmt.Sprintf("Exif Orientation: %d", i)

		test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
			auth := test.GetAuthHeader(t, ts)

			// WITH
			file, _ := os.Open("../fixtures/new_image.json")
			res, _ := test.RunRequest("POST", ts.URL+"/api/v1.0/nodes", file, auth)

			assert.Equal(t, 201, res.StatusCode, message)

			node := base.NewNode()
			serializer := app.Get("gonode.node.serializer").(*base.Serializer)
			serializer.Deserialize(res.Body, node)

			file, _ = os.Open(fmt.Sprintf("../fixtures/exif_orientation/f%d-exif.jpg", i))

			res, _ = test.RunRequest("PUT", ts.URL+"/api/v1.0/nodes/"+node.Uuid.CleanString()+"?raw", file, auth)
			assert.Equal(t, http.StatusOK, res.StatusCode, message)

			res, _ = test.RunRequest("GET", ts.URL+"/prism/"+node.Uuid.CleanString()+".jpg?mr=20", nil, auth)
			assert.Equal(t, 200, res.StatusCode, message)
			assert.Equal(t, "image/jpeg", res.Header.Get("Content-Type"), message)

			config, format, err := image.DecodeConfig(res.Body)
			assert.Nil(t, err, message)
			assert.Equal(t, "jpeg", format, message)
			assert.Equal(t, 20, config.Width, message+" - invalid width (resized)")
			assert.Equal(t, 40, config.Height, message+" - invalid height (resized)")
		})
	}
}
