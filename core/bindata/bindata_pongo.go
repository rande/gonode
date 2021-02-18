// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bindata

import (
	"bytes"
	"io"
)

type PongoTemplateLoader struct {
	Asset func(name string) ([]byte, error)
	Paths []string
}

// Abs calculates the path to a given template. Whenever a path must be resolved
// due to an import from another template, the base equals the parent template's path.
func (l *PongoTemplateLoader) Abs(base, name string) string {
	for _, lookupPath := range l.Paths {

		path := lookupPath + "/" + name

		_, err := l.Asset(path)

		if err == nil {
			return path
		}
	}

	return name
}

// Get returns an io.Reader where the template's content can be read from.
func (l *PongoTemplateLoader) Get(path string) (io.Reader, error) {
	data, err := l.Asset(path)

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
