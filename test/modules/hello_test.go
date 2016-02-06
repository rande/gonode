// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
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

func Test_Hello(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		res, _ := test.RunRequest("GET", ts.URL+"/hello", nil)

		assert.Equal(t, res.GetBody(), []byte("Hello!"))
	})
}
