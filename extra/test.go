package extra

import (
	. "github.com/rande/goapp"
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

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	return &Response{Response: resp}, err
}

func RunHttpTest(t *testing.T, f func(t *testing.T, ts *httptest.Server, app *App)) {
	app := GetApp("../config_test.toml")
	mux := app.Get("goji.mux").(*web.Mux)

	ts := httptest.NewServer(mux)

	defer func() {
		RunRequest("PUT", ts.URL+"/uninstall", nil)
		ts.Close()
	}()

	RunRequest("PUT", ts.URL+"/uninstall", nil)
	RunRequest("PUT", ts.URL+"/install", nil)

	f(t, ts, app)
}
