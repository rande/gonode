// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Hello(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetDefaultAuthHeader(ts)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/hello", nil, auth)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, res.GetBody(), []byte("Hello!"))
	})
}

func Test_Invalid_Request(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetDefaultAuthHeader(ts)

		// wrong version api number
		res, _ := test.RunRequest("GET", ts.URL+"/api/v01/hello", nil, auth)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
