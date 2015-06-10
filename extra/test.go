package extra

import (
	. "github.com/rande/goapp"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	nc "github.com/rande/gonode/core"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"fmt"
)

func GetApp(file string) *App {

	app := NewApp()

	app.Set("gonode.configuration", func(app *App) interface{} {
		return GetConfiguration(file)
	})

	// configure main services
	app.Set("logger", func(app *App) interface{} {
		return log.New(os.Stdout, "", log.Lshortfile)
	})

	app.Set("goji.mux", func(app *App) interface{} {
		mux := web.New()

		//		mux.Use(middleware.RequestID)
		mux.Use(middleware.Logger)
		mux.Use(middleware.Recoverer)
		//		mux.Use(middleware.AutomaticOptions)

		return mux
	})

	ConfigureApp(app)
	ConfigureGoji(app)

	return app
}

type Response struct {
	*http.Response
	RawBody  []byte
	bodyRead bool
}

func (r Response) GetBody() []byte {
	var err error

	if !r.bodyRead {
		r.RawBody, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		r.bodyRead = true
	}

	return r.RawBody
}

func RunRequest(method string, url string, body io.Reader) (*Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)

	nc.PanicOnError(err)

	resp, err := client.Do(req)

	return &Response{Response: resp}, err
}

func RunHttpTest(t *testing.T, f func(t *testing.T, ts *httptest.Server, app *App)) {
	var err error
	var res *Response

	app := GetApp("../config_test.toml")
	mux := app.Get("goji.mux").(*web.Mux)

	ts := httptest.NewServer(mux)

	defer func() {
		ts.Close()
	}()

	res, err = RunRequest("PUT", ts.URL+"/uninstall", nil)
	nc.PanicIf(res.StatusCode != http.StatusOK, fmt.Sprintf("Expected code 200, get %d\n%s", res.StatusCode, string(res.GetBody()[:])))
	nc.PanicOnError(err)

	res, err = RunRequest("PUT", ts.URL+"/install", nil)
	nc.PanicIf(res.StatusCode != http.StatusOK, fmt.Sprintf("Expected code 200, get %d\n%s", res.StatusCode, string(res.GetBody()[:])))
	nc.PanicOnError(err)

	f(t, ts, app)
}
