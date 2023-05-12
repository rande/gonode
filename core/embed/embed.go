// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"embed"
	"errors"
)

var (
	UnableToFindEmbedError      = errors.New("Unable to find the embed file")
	ModuleDoesNotExistError     = errors.New("Module does not exist")
	InvalidPongoRefError        = errors.New("Invalid pongo reference name")
)

func NewEmbeds() *Embeds {
	return &Embeds{
		fs: make(map[string][]embed.FS),
	}
}

type Embeds struct {
	fs map[string][]embed.FS
}

func (a *Embeds) Add(module string, fs embed.FS) {
	if (a.fs[module] == nil) {
		a.fs[module] = make([]embed.FS, 0);
	}

	a.fs[module] = append(a.fs[module], fs)
}

func (a *Embeds) ReadFile(module string, path string) ([]byte, error) {
	if (a.fs[module] == nil) {
		return nil, ModuleDoesNotExistError
	}

	for i := range a.fs[module] {
		buff, _ := a.fs[module][i].ReadFile(path)

		if (buff != nil) {
			return buff, nil
		}
	}

	return nil, UnableToFindEmbedError
}