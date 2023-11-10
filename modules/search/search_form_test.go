// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SearchForm_UrlValues(t *testing.T) {
	s := &SearchForm{
		PerPage: 32,
		Page:    1,
		Type:    []*Param{NewParam("blog.type")},
	}

	assert.Equal(t, s.UrlValues().Encode(), "page=1&per_page=32&type=blog.type")
}
