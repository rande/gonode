// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package media

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetExif_With_Valid_Exif(t *testing.T) {
	// image from github.com/rwcarlsen/goexif test suite
	f, err := os.Open("exif/2004-01-11-22-45-15-sep-2004-01-11-22-45-15a.jpg")

	assert.NoError(t, err)

	e, err := GetExif(f)

	assert.NoError(t, err)

	assert.Equal(t, e["DateTime"], "2004:01:11 22:45:19")
	assert.Equal(t, e["Software"], "M5011S-1031")
	assert.Equal(t, e["PixelXDimension"], "1600")
}

func Test_GetExif_With_Invalid_Exif(t *testing.T) {
	// image from github.com/rwcarlsen/goexif test suite
	f, err := os.Open("exif/infinite_loop_exif.jpg")

	assert.NoError(t, err)

	_, err = GetExif(f)

	assert.Error(t, err)
}
