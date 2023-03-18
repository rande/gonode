// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fmt"
	"os"
)

// this is just a test to validata how the aws sdk behave
func Test_Vault_Driver_Fs(t *testing.T) {
	// init vault
	v := &DriverFs{
		Root: fmt.Sprintf("%s/gonode/%d", os.TempDir(), time.Now().Nanosecond()),
	}

	key := "test/assd"

	assert.False(t, v.Has(key))

	data := []byte("foobar et foo")

	w, err := v.GetWriter(key)

	assert.Nil(t, err)

	w.Write(data)
	w.Close()

	// should have the file
	assert.True(t, v.Has(key))

	r, err := v.GetReader(key)
	assert.Nil(t, err)

	data, err = ioutil.ReadAll(r)

	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar et foo"), data)

	err = v.Remove(key)
	assert.Nil(t, err)

	err = v.Remove(key)
	assert.NotNil(t, err)
}
