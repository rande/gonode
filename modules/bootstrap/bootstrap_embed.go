// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bootstrap

import (
	"embed"
)

//go:embed all:static
var content embed.FS

func GetEmbedFS() embed.FS {
	return content
}
