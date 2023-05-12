// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"bytes"
	"io"
	"strings"
)

type PongoTemplateLoader struct {
	Embeds *Embeds
	BasePath string
}

// Abs calculates the path to a given template. Whenever a path must be resolved
// due to an import from another template, the base equals the parent template's path.
func (l *PongoTemplateLoader) Abs(base, name string) string {
	// for _, lookupPath := range l.Paths {

	// 	path := lookupPath + "/" + name

	// 	_, err := l.Asset(path)

	// 	if err == nil {
	// 		return path
	// 	}
	// }

	return name
}

// Get returns an io.Reader where the template's content can be read from.
func (l *PongoTemplateLoader) Get(path string) (io.Reader, error) {

	// blog:foo/blog.post.tpl => module=blog templates/foo/blog.post.tpl

	sections := strings.Split(path, ":")


	if (len(sections) != 2) {
		return nil, InvalidPongoRefError
	}

	data, err := l.Embeds.ReadFile(sections[0], l.BasePath + "templates/" + strings.Join(sections[1:], "/"))

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
