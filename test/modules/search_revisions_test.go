// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rande/goapp"
	"github.com/rande/gonode/core/helper"
	"github.com/rande/gonode/modules/base"
	"github.com/rande/gonode/modules/user"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Search_Revision_Basic(t *testing.T) {

	var err error

	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {

		u := app.Get("gonode.handler_collection").(base.HandlerCollection).NewNode("core.user")
		manager := app.Get("gonode.manager").(*base.PgNodeManager)

		data := u.Data.(*user.User)
		data.Email = "test@example.org"
		data.Enabled = true
		data.FirstName = "Thomas"
		data.LastName = "Rxxxx"
		data.NewPassword = "ZePassword"
		data.Username = "rande"
		u.Name = "Title 1"
		u.Access = []string{"node:api:master"}

		meta := u.Meta.(*user.UserMeta)
		meta.PasswordCost = 1 // save test time

		assert.Equal(t, 1, u.Revision, "revision should be 1")

		// create node
		u, err = manager.Save(u, true)
		u.Name = "Title 2"
		assert.Equal(t, 1, u.Revision, "revision should be 1")
		helper.PanicOnError(err)

		u, err = manager.Save(u, true)
		assert.Equal(t, 2, u.Revision, "revision should be 2")
		helper.PanicOnError(err)

		u, err = manager.Save(u, true)
		assert.Equal(t, 3, u.Revision, "revision should be 3")
		helper.PanicOnError(err)

		u, err = manager.Save(u, true)
		helper.PanicOnError(err)
		assert.Equal(t, 4, u.Revision, "revision should be 4")

		u, err = manager.Save(u, true)
		helper.PanicOnError(err)
		assert.Equal(t, 5, u.Revision, "revision should be 5")

		baseUrl := fmt.Sprintf("%s/api/v1.0/nodes/%s/revisions", ts.URL, u.Uuid.CleanString())

		values := []struct {
			Url string
			Len int
		}{
			{fmt.Sprintf("%s", baseUrl), 5},
			{fmt.Sprintf("%s?per_page=2", baseUrl), 2},
		}

		auth := test.GetDefaultAuthHeader(ts)

		for _, v := range values {
			// WHEN
			res, _ := test.RunRequest("GET", v.Url, nil, auth)

			p := test.GetPager(app, res)

			// THEN
			assert.Equal(t, uint64(1), p.Page, "Page: "+v.Url)
			assert.Equal(t, int(v.Len), len(p.Elements), "Len: "+v.Url)
		}

		// existing node
		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/2", baseUrl), nil, auth)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		node := test.GetNode(app, res)
		assert.Equal(t, 2, node.Revision, "Invalid revision number")
		assert.Equal(t, "Title 2", node.Name, "Invalid name")

		// existing node
		res, _ = test.RunRequest("GET", fmt.Sprintf("%s/100", baseUrl), nil, auth)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
