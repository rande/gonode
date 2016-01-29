// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"errors"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIf(b bool, message string) {
	if b {
		PanicOnError(errors.New(message))
	}
}

func PanicUnless(b bool, message string) {
	PanicIf(!b, message)
}
