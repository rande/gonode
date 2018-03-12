// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parial_Reader_Incomplet_Data(t *testing.T) {

	f, _ := ioutil.TempFile(os.TempDir(), "gonode_test")
	f.Write([]byte("hello world"))
	f.Seek(0, 0)

	defer f.Close()

	r := PartialReader{
		Reader: f,
		Size:   20,
	}

	data := make([]byte, 50)

	n, err := r.Read(data)

	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, []byte("hello world"), data[:n])
	assert.Equal(t, []byte("hello world"), r.Data)
}

func Test_Parial_Reader_Complet_Data(t *testing.T) {

	f, _ := ioutil.TempFile(os.TempDir(), "gonode_test")
	f.Write([]byte("hello world"))
	f.Seek(0, 0)

	defer f.Close()

	r := PartialReader{
		Reader: f,
		Size:   5,
	}

	data := make([]byte, 50)

	n, err := r.Read(data)

	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, []byte("hello world"), data[:n])
	assert.Equal(t, []byte("hello"), r.Data)
}
