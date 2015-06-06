package core

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
