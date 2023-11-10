// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/rande/gonode/core/embed"
)

var (
	ErrRootTemplateNotFound = errors.New("root template not found")
)

type Context map[string]interface{}
type FuncMap map[string]interface{}

type TemplateLoader struct {
	Embeds    *embed.Embeds
	BasePath  string
	Templates map[string]*template.Template
	FuncMap   map[string]interface{}
}

func (l *TemplateLoader) Execute(path string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if tpl, ok := l.Templates[path]; !ok {
		return nil, ErrRootTemplateNotFound
	} else if err := tpl.ExecuteTemplate(&buf, path, data); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}
