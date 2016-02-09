// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package helper

import (
	"errors"
)

type PanicCallback func()

func PanicOnError(err error, pc ...PanicCallback) {
	if err != nil {

		if len(pc) > 0 {
			pc[0]()
		}

		panic(err)
	}
}

func PanicIf(b bool, message string, pc ...PanicCallback) {
	if b {
		PanicOnError(errors.New(message), pc...)
	}
}

func PanicUnless(b bool, message string, pc ...PanicCallback) {
	PanicIf(!b, message, pc...)
}
