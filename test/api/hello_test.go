package api

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/extra"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Hello(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		res, _ := extra.RunRequest("GET", ts.URL+"/hello", nil)

		assert.Equal(t, res.GetBody(), []byte("Hello!"))
	})
}
