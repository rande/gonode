// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"io"
	"os"
	"path/filepath"
)

type DriverFs struct {
	Root string
}

func (v *DriverFs) getFilename(name string) string {
	return filepath.Join(v.Root, name)
}

func (v *DriverFs) Has(name string) bool {
	if _, err := os.Stat(v.getFilename(name)); err != nil {
		return false
	}

	return true
}

func (v *DriverFs) GetReader(name string) (io.ReadCloser, error) {
	return os.Open(v.getFilename(name))
}

func (v *DriverFs) GetWriter(name string) (io.WriteCloser, error) {
	filename := v.getFilename(name)

	path := filepath.Dir(filename)

	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0600)
}

func (v *DriverFs) Remove(name string) error {
	return os.Remove(v.getFilename(name))
}
