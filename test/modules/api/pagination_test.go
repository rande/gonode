// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core"
	"github.com/rande/gonode/modules/user"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Pagination(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		values := []struct {
			Url      string
			PerPage  uint64
			Len      uint64
			Page     uint64
			Previous uint64
			Next     uint64
		}{
			{"/nodes?per_page=6", 6, 6, 1, 0, 2},
			{"/nodes?per_page=13", 13, 12, 1, 0, 0},
			{"/nodes?per_page=2&page=6", 2, 2, 6, 5, 0},
		}

		// WITH
		// create a valid user into the database ...
		manager := app.Get("gonode.manager").(*core.PgNodeManager)

		for i := 0; i < 11; i++ {
			u := app.Get("gonode.handler_collection").(core.HandlerCollection).NewNode("core.user")

			data := u.Data.(*user.User)
			data.Email = "test@example.org"
			data.Enabled = true
			data.FirstName = "Thomas"
			data.LastName = "Rxxxx"
			data.NewPassword = "ZePassword"
			data.Username = fmt.Sprintf("rande%02d", i)

			meta := u.Meta.(*user.UserMeta)
			meta.PasswordCost = 1 // save test time
			manager.Save(u, false)
		}

		auth := test.GetAuthHeader(t, ts)

		for _, v := range values {
			// paginate result ...
			res, _ := test.RunRequest("GET", ts.URL+v.Url, nil, auth)

			p := GetPager(app, res)

			fmt.Printf("%v", p.Elements)
			// THEN
			assert.Equal(t, v.PerPage, p.PerPage, "Wrong PerPage value: "+v.Url)
			assert.Equal(t, int(v.Len), len(p.Elements), "Wrong Len value: "+v.Url)
			assert.Equal(t, v.Page, p.Page, "Wrong Page value: "+v.Url)
			assert.Equal(t, v.Next, p.Next, "Wrong Next value: "+v.Url)
			assert.Equal(t, v.Previous, p.Previous, "Wrong Prvious value: "+v.Url)
		}
	})
}
