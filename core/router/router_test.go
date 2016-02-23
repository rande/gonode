// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func Test_Reverse_Basic_Usage(t *testing.T) {
	router := NewRouter(nil)

	router.Handle("prism", "/prism/:name/:format/end", func(http.ResponseWriter, *http.Request) {})

	cases := []struct {
		params url.Values
		url    string
	}{
		{url.Values{"name": []string{"foobar"}, "format": []string{"json"}}, "/prism/foobar/json/end"},
		{url.Values{"name": []string{"foobar"}, "format": []string{"json"}, "raw": []string{"1", "2"}}, "/prism/foobar/json/end?raw=1&raw=2"},
		{url.Values{"name": []string{"foobar"}, "format": []string{"json"}, "raw": []string{"1", "2"}}, "/prism/foobar/json/end?raw=1&raw=2"},
	}

	for _, data := range cases {
		url, err := router.GeneratePath("prism", data.params)

		assert.NoError(t, err)
		assert.Equal(t, data.url, url)
	}
}
