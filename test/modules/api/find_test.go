// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Find_Non_Existant(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", ts.URL+"/nodes/d703a3ab-8374-4c30-a8a4-2c22aa67763b", nil, auth)

		assert.Equal(t, 404, res.StatusCode, "Delete non existant node")
	})
}
