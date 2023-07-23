// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dashboard

import (
	"embed"
)

//go:embed all:static all:templates
var content embed.FS

func GetEmbedFS() embed.FS {
	return content
}
