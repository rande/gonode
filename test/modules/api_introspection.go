// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_API_GET_Handlers_Node(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/handlers/node", nil, auth)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

func Test_API_GET_Handlers_View(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/handlers/view", nil, auth)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}

func Test_API_GET_Services(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/services", nil, auth)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	})
}
