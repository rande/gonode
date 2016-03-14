// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"net/http/httptest"
	"testing"

	. "github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Find_Non_Existant(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		auth := test.GetAuthHeader(t, ts)

		res, _ := test.RunRequest("GET", ts.URL+"/api/v1.0/nodes/d703a3ab-8374-4c30-a8a4-2c22aa67763b", nil, auth)

		assert.Equal(t, 404, res.StatusCode, "Delete non existant node")
	})
}
