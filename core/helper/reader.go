// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"io"
)

type PartialReader struct {
	Reader io.ReadCloser
	Size   int
	Data   []byte
	Len    int
}

func (r *PartialReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	toRead := r.Size - len(r.Data)

	if toRead > 0 {
		if n < toRead {
			toRead = n
		}

		dst := make([]byte, toRead)
		n := copy(dst, p)

		r.Data = append(r.Data, dst[0:n]...)
	}

	return
}

func (r *PartialReader) Close() error {
	return r.Reader.Close()
}
