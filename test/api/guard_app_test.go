package api

import (
	"fmt"
	"github.com/rande/goapp"
	"github.com/rande/gonode/test"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Guard_Error(t *testing.T) {
	test.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *goapp.App) {
		res, _ := test.RunRequest("GET", fmt.Sprintf("%s/nodes/protected", ts.URL), nil)

		assert.Equal(t, 403, res.StatusCode)
	})
}
