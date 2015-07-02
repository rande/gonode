package api

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/extra"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Find_Non_Existant(t *testing.T) {
	extra.RunHttpTest(t, func(t *testing.T, ts *httptest.Server, app *App) {
		res, _ := extra.RunRequest("GET", ts.URL+"/nodes/d703a3ab-8374-4c30-a8a4-2c22aa67763b", nil)

		assert.Equal(t, 404, res.StatusCode, "Delete non existant node")
	})
}
