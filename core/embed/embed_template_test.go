// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"bytes"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/rande/gonode/core/embed/fixtures"
	"github.com/stretchr/testify/assert"
)

func Test_Template_With_Custom_Path(t *testing.T) {
	embeds := NewEmbeds()
	embeds.Add("testmodule", fixtures.GetTestEmbedFS())

	var buf bytes.Buffer

	templates := GetTemplates(embeds)
	ctx := map[string]interface{}{
		"Title": "Hello World!",
	}

	if tpl, ok := templates["testmodule:pages/index"]; !ok {
		assert.True(t, ok)
	} else {
		err := tpl.ExecuteTemplate(&buf, "testmodule:pages/index", ctx)

		assert.Nil(t, err)

		snaps.MatchSnapshot(t, buf.String())
	}

	buf = *bytes.NewBuffer([]byte{})

	if tpl, ok := templates["testmodule:pages/blog"]; !ok {
		assert.True(t, ok)
	} else {
		err := tpl.ExecuteTemplate(&buf, "testmodule:pages/blog", ctx)

		assert.Nil(t, err)

		snaps.MatchSnapshot(t, buf.String())
	}
}
