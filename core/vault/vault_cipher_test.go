// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var encKey = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
var invalidKey = []byte("baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
var message = "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?"

func Test_GetCipher(t *testing.T) {
	encrypter, decrypter := GetCipher("aes_ofb")
	assert.NotNil(t, encrypter)
	assert.NotNil(t, decrypter)

	encrypter, decrypter = GetCipher("aes_ctr")
	assert.NotNil(t, encrypter)
	assert.NotNil(t, decrypter)

	encrypter, decrypter = GetCipher("aes_cbc")
	assert.NotNil(t, encrypter)
	assert.NotNil(t, decrypter)

	encrypter, decrypter = GetCipher("aes_gcm")
	assert.NotNil(t, encrypter)
	assert.NotNil(t, decrypter)

	encrypter, decrypter = GetCipher("no_op")
	assert.NotNil(t, encrypter)
	assert.NotNil(t, decrypter)

	assert.Panics(t, func() { GetCipher("invalid_cypher") })
}

func Test_AesOFB(t *testing.T) {
	ve := NewVaultElement()
	ve.MetaKey = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	assert.Equal(t, 32, len(ve.MetaKey))

	src := bytes.NewBuffer([]byte("Hello World!!"))
	dst := bytes.NewBuffer([]byte(""))

	_, err := AesOFBEncrypter(ve.MetaKey, src, dst)

	assert.Nil(t, err)
	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))

	_, err = AesOFBDecrypter(ve.MetaKey, dst, decrypted)

	assert.Nil(t, err)
	assert.Equal(t, []byte("Hello World!!"), decrypted.Bytes())
}

func Test_AesCBC(t *testing.T) {
	src := bytes.NewBuffer([]byte(message))
	dst := bytes.NewBuffer([]byte(""))

	_, err := AesCBCEncrypter(encKey, src, dst)
	assert.Nil(t, err)
	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))
	_, err = AesCBCDecrypter(encKey, dst, decrypted)

	assert.Nil(t, err)
	assert.Equal(t, []byte(message), decrypted.Bytes())
}

func Test_AesCTR(t *testing.T) {
	src := bytes.NewBuffer([]byte(message))
	dst := bytes.NewBuffer([]byte(""))

	_, err := AesCTREncrypter(encKey, src, dst)
	assert.Nil(t, err)
	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))
	_, err = AesCTRDecrypter(encKey, dst, decrypted)

	assert.Nil(t, err)
	assert.Equal(t, []byte(message), decrypted.Bytes())
}

func Test_AesGCM(t *testing.T) {
	src := bytes.NewBuffer([]byte(message))
	dst := bytes.NewBuffer([]byte(""))

	_, err := AesGCMEncrypter(encKey, src, dst)
	assert.Nil(t, err)
	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))
	_, err = AesGCMDecrypter(encKey, dst, decrypted)

	assert.Nil(t, err)
	assert.Equal(t, []byte(message), decrypted.Bytes())
}

func Test_AesGCM_Invalid_Key(t *testing.T) {
	src := bytes.NewBuffer([]byte(message))
	dst := bytes.NewBuffer([]byte(""))

	_, err := AesGCMEncrypter(encKey, src, dst)
	assert.Nil(t, err)
	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))
	_, err = AesGCMDecrypter(invalidKey, dst, decrypted)

	assert.NotNil(t, err)
	// assert.Equal(t, []byte(message), decrypted.Bytes())
}
