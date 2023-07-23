package dashboard

import (
	"io"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/rande/goapp"
	"github.com/rande/gonode/core/helper"
	"github.com/zenazn/goji/web"
)

// prototype of a view response, this is used to return a response from a view
// function and this will probably replace the current response system from the prism
// module.
type ViewResponse struct {
	StatusCode  int
	ContentType string
	Template    string
	Context     pongo2.Context
	Location    string
	Body        io.Reader
	Headers     map[string]string
}

func (r *ViewResponse) Add(name string, v interface{}) *ViewResponse {
	r.Context[name] = v

	return r
}

func HtmlResponse(code int, template string) *ViewResponse {
	return &ViewResponse{
		StatusCode:  code,
		Template:    template,
		Context:     pongo2.Context{},
		Headers:     map[string]string{},
		ContentType: "text/html; charset=UTF-8",
	}
}

func JsonResponse(code int, body io.Reader) *ViewResponse {
	return &ViewResponse{
		StatusCode:  code,
		ContentType: "application/json",
		Body:        body,
		Headers:     map[string]string{},
	}
}

func RedirectResponse(code int, location string) *ViewResponse {
	return &ViewResponse{
		StatusCode: code,
		Location:   location,
		Headers:    map[string]string{},
	}
}

func StreamedResponse(code int, contentType string, body io.Reader) *ViewResponse {
	return &ViewResponse{
		StatusCode:  code,
		ContentType: contentType,
		Body:        body,
		Headers:     map[string]string{},
	}
}

type ViewHandlerInterface func(c web.C, res http.ResponseWriter, req *http.Request) *ViewResponse

func InitView(app *goapp.App, creator func(app *goapp.App) ViewHandlerInterface) web.HandlerFunc {
	pongo := app.Get("gonode.pongo").(*pongo2.TemplateSet)

	handler := creator(app)

	return func(c web.C, res http.ResponseWriter, req *http.Request) {

		view := handler(c, res, req)

		// response manually handled
		if view == nil {
			return
		}

		if view.Location != "" {
			http.Redirect(res, req, view.Location, view.StatusCode)
			return
		}

		for k, v := range view.Headers {
			// we explicitly set the content type later
			// in the code
			if k == "Content-Type" {
				continue
			}

			res.Header().Set(k, v)
		}

		res.Header().Set("Content-Type", view.ContentType)

		res.WriteHeader(view.StatusCode)

		if view.Template != "" {
			view.Context.Update(pongo2.Context{
				"request": req,
			})

			tpl, err := pongo.FromFile(view.Template)

			helper.PanicOnError(err)

			data, err := tpl.ExecuteBytes(view.Context)

			helper.PanicOnError(err)

			res.Write(data)
			return
		}

		if view.Body != nil {
			_, err := io.Copy(res, view.Body)
			helper.PanicOnError(err)
		}
	}
}
