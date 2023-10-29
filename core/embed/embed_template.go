// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"bytes"
	"html/template"
	"io"
	"strings"
)

type TemplateLoader struct {
	Embeds   *Embeds
	BasePath string
	Template *template.Template
}

// Abs calculates the path to a given template. Whenever a path must be resolved
// due to an import from another template, the base equals the parent template's path.
func (l *TemplateLoader) Abs(base, name string) string {
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
// blog:foo/blog.post.html => module=blog templates/foo/blog.post.html
func (l *TemplateLoader) Get(path string) (io.Reader, error) {
	sections := strings.Split(path, ":")

	if len(sections) != 2 {
		return nil, ErrInvalidPongoRef
	}

	data, err := l.Embeds.ReadFile(sections[0], l.BasePath+"templates/"+strings.Join(sections[1:], "/"))

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func (l *TemplateLoader) Execute(path string, data interface{}) (string, error) {
	var buf bytes.Buffer

	err := l.Template.ExecuteTemplate(&buf, path, data)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
