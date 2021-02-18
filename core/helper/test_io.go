// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"io"
	"strings"
)

type TestCloserReader struct {
	Data io.Reader
}

func (t *TestCloserReader) Read(p []byte) (n int, err error) {
	return t.Data.Read(p)
}

func (t *TestCloserReader) Close() error {
	return nil // no nothing ...
}

func NewTestCloserReader(data string) *TestCloserReader {
	return &TestCloserReader{
		Data: strings.NewReader(data),
	}
}
