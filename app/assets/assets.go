// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package assets

var rootDir = ""

func UpdateRootDir(path string) {

	if len(path) == 0 {
		return
	}

	rootDir = path
}
