// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package fixtures

import (
	"embed"
)

//go:embed all:templates
var content embed.FS

func GetTestEmbedFS() embed.FS {
	return content
}
