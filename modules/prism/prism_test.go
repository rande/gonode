// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package prism

import (
	"html/template"
	"net/http"
	"net/url"
	"testing"

	"github.com/rande/gonode/core/router"
	"github.com/rande/gonode/modules/base"
	"github.com/stretchr/testify/assert"
	"github.com/zenazn/goji/web"
)

func Test_PrimPath_Node_With_Path(t *testing.T) {

	r := router.NewRouter(nil)
	r.Handle("prism_path", "/:path", func(c web.C, res http.ResponseWriter, req *http.Request) {})
	r.Handle("prism_path_format", "/:path.:format", func(c web.C, res http.ResponseWriter, req *http.Request) {})

	node := base.NewNode()
	node.Path = "/path/to/my/content"

	cases := []struct {
		node   *base.Node
		params url.Values
		url    template.HTML
	}{
		{node, url.Values{}, "/path/to/my/content"},
		{node, url.Values{"name": []string{"foobar"}}, template.HTML("/path/to/my/content?name=foobar")},
		{node, url.Values{"format": []string{"html"}}, template.HTML("/path/to/my/content.html")},
		{node, url.Values{"format": []string{"html"}, "name": []string{"foobar"}}, template.HTML("/path/to/my/content.html?name=foobar")},
	}

	f := PrismPath(r)

	for _, data := range cases {
		url := f(data.node, data.params)

		assert.Equal(t, data.url, url)
	}
}

func Test_PrimPath_Node_Without_Path(t *testing.T) {

	r := router.NewRouter(nil)
	r.Handle("prism", "/prism/:uuid", func(c web.C, res http.ResponseWriter, req *http.Request) {})
	r.Handle("prism_format", "/prism/:uuid.:format", func(c web.C, res http.ResponseWriter, req *http.Request) {})

	node := base.NewNode()

	cases := []struct {
		node   *base.Node
		params url.Values
		url    template.HTML
	}{
		{node, url.Values{}, template.HTML("/prism/11111111-1111-1111-1111-111111111111")},
		{node, url.Values{"name": []string{"foobar"}}, template.HTML("/prism/11111111-1111-1111-1111-111111111111?name=foobar")},
		{node, url.Values{"name": []string{"foobar"}, "format": []string{"html"}}, template.HTML("/prism/11111111-1111-1111-1111-111111111111.html?name=foobar")},
	}

	f := PrismPath(r)

	for _, data := range cases {
		url := f(data.node, data.params)

		assert.Equal(t, data.url, url)
	}
}
