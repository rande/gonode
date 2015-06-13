package extra

import (
	"fmt"
	"github.com/rande/goapp"
	nc "github.com/rande/gonode/core"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func GetLifecycle(file string) *goapp.Lifecycle {

	l := goapp.NewLifecycle()

	l.Config(func(app *goapp.App) error {
		app.Set("gonode.configuration", func(app *goapp.App) interface{} {
				return GetConfiguration(file)
			})

		return nil
	})

	l.Register(func(app *goapp.App) error {
		// configure main services
		app.Set("logger", func(app *goapp.App) interface{} {
				return log.New(os.Stdout, "", log.Lshortfile)
			})

		app.Set("goji.mux", func(app *goapp.App) interface{} {
				mux := web.New()

				//		mux.Use(middleware.RequestID)
				mux.Use(middleware.Logger)
				mux.Use(middleware.Recoverer)
				//		mux.Use(middleware.AutomaticOptions)

				return mux
			})

		return nil
	})

	ConfigureApp(l)
	ConfigureGoji(l)

	return l
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

func RunHttpTest(t *testing.T, f func(t *testing.T, ts *httptest.Server, app *goapp.App)) {

	l := GetLifecycle("../config_test.toml")

	l.Run(func(app *goapp.App) error {
		var err error
		var res *Response

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

		return nil
	})

	l.Go(goapp.NewApp())
}
