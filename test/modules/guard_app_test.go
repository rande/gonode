// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
)

func Test_Guard_Error(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/api/v1.0/nodes/protected", ts.URL), nil)

		assert.Equal(t, 403, res.StatusCode)
	})
}
