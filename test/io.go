package test

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
