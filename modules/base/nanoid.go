// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// This is a adaptation of https://github.com/matoous/go-nanoid/blob/master/gonanoid.go

package base

import (
	"crypto/rand"
)

// defaultAlphabet is the alphabet used for ID characters by default.
var defaultAlphabet = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	defaultSize = 16
)

// New generates secure URL-friendly unique ID.
// Accepts optional parameter - length of the ID to be generated (21 by default).
func NewId() string {
	bytes := make([]byte, defaultSize)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	id := make([]rune, defaultSize)
	for i := 0; i < defaultSize; i++ {
		id[i] = defaultAlphabet[bytes[i]&61]
	}
	return string(id[:defaultSize])
}
