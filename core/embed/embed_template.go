// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

var (
	ErrRootTemplateNotFound = errors.New("root template not found")
)

type TemplateLoader struct {
	Embeds    *Embeds
	BasePath  string
	Templates map[string]*template.Template
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

func (l *TemplateLoader) Execute(path string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	fmt.Printf("Execute: %s\n", path)
	if tpl, ok := l.Templates[path]; !ok {
		return nil, ErrRootTemplateNotFound
	} else if err := tpl.ExecuteTemplate(&buf, path, data); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}
