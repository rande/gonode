// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/rande/gonode/core/embed/fixtures"
	"github.com/stretchr/testify/assert"
)

func Test_Template_With_Custom_Path(t *testing.T) {

	embeds := NewEmbeds()
	embeds.Add("testmodule", fixtures.GetTestEmbedFS())

	data, err := embeds.ReadFile("testmodule", "templates/layout.html")

	assert.Nil(t, err)
	assert.NotNil(t, data)
	tpl, err := template.New("default").Parse(string(data))
	assert.NoError(t, err)

	var buf bytes.Buffer

	err = tpl.Execute(&buf, nil)
	assert.Nil(t, err)

	tpl = template.New("default")
	err = ConfigureTemplates(tpl, embeds)
	assert.Nil(t, err)

	err = tpl.ExecuteTemplate(&buf, "testmodule:content.html", nil)
	assert.Nil(t, err)

	snaps.MatchSnapshot(t, buf.String())
}
