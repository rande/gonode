// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

func Test_Embeded_404(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetDefaultAuthHeader(ts)

		res, _ := test.RunRequest("GET", ts.URL+"/static/setup/hello.txt", nil, auth)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Equal(t, res.GetBody(), []byte("<html><head><title>Embed not found</title></head><body><h1>Embed not found</h1></body></html>"))
	})
}
