// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package embed

import (
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
)

var (
	ErrUnableToFindEmbed  = errors.New("unable to find the embed file/directory")
	ErrModuleDoesNotExist = errors.New("module does not exist")
	ErrInvalidTemplateRef = errors.New("invalid template reference name")
)

type Embed struct {
	Module string
	Path   string
}

func NewEmbeds() *Embeds {
	return &Embeds{
		fs: make(map[string][]embed.FS),
	}
}

type Embeds struct {
	fs map[string][]embed.FS
}

func (a *Embeds) Add(module string, fs embed.FS) {
	if a.fs[module] == nil {
		a.fs[module] = make([]embed.FS, 0)
	}

	a.fs[module] = append(a.fs[module], fs)
}

func (a *Embeds) ReadFile(module string, path string) ([]byte, error) {
	if a.fs[module] == nil {
		return nil, ErrModuleDoesNotExist
	}

	for i := range a.fs[module] {
		buff, _ := a.fs[module][i].ReadFile(path)

		if buff != nil {
			return buff, nil
		}
	}

	return nil, ErrUnableToFindEmbed
}

func (a *Embeds) GetModules() []string {
	modules := make([]string, 0)

	for k := range a.fs {
		modules = append(modules, k)
	}

	return modules
}

func (a *Embeds) ReadDir(module string, path string) ([]fs.DirEntry, error) {
	if a.fs[module] == nil {
		return nil, ErrModuleDoesNotExist
	}

	for i := range a.fs[module] {
		dir, _ := a.fs[module][i].ReadDir(path)

		if dir != nil {
			return dir, nil
		}
	}

	return nil, ErrUnableToFindEmbed
}

func (a *Embeds) GetFilesByExt(ext string) []Embed {
	var embeds []Embed

	for _, module := range a.GetModules() {
		embeds = append(embeds, readDir(module, a.fs[module][0], "templates")...)
	}

	filteredEmbeds := make([]Embed, 0)
	for _, embed := range embeds {
		if filepath.Ext(embed.Path) == ext {
			filteredEmbeds = append(filteredEmbeds, embed)
		}
	}

	return filteredEmbeds
}

func readDir(module string, dir embed.FS, path string) []Embed {
	var embeds []Embed

	files, err := dir.ReadDir(path)

	if err != nil {
		return embeds
	}

	for _, file := range files {
		if file.IsDir() {
			embeds = append(embeds, readDir(module, dir, path+"/"+file.Name())...)
		}

		embeds = append(embeds, Embed{
			Module: module,
			Path:   path + "/" + file.Name(),
		})
	}

	return embeds
}
