package api

import (
	. "github.com/rande/goapp"
	"github.com/rande/gonode/extra"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_Hello(t *testing.T) {
	var res *extra.Response
	var app *App

	app = extra.GetApp("../config_test.toml")
	ts := app.Get("testserver").(*httptest.Server)

	res, _ = extra.RunRequest("GET", ts.URL+"/hello", nil)

	assert.Equal(t, res.GetBody(), []byte("Hello!"))
}
