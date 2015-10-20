package vault

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	//	"crypto/rand"
	//	"fmt"
)

//	message := make([]byte, 2048)
//	io.ReadFull(rand.Reader, message)

// write/encrypted file
func RunTestVault(t *testing.T, v Vault, plaintext []byte) {
	var read int64

	file := "this-is-a-test"

	meta := NewVaultMetadata()
	meta["foo"] = "bar"

	reader := bytes.NewBuffer(plaintext)

	written, err := v.Put(file, meta, reader)

	assert.NoError(t, err)
	assert.True(t, written >= int64(len(plaintext))) // some cipher might add extra data
	assert.True(t, written > 0)                      // some cipher might add extra data
	assert.True(t, v.Has(file))

	invalid := []byte("Another invalid message with the same key")

	// test overwrite
	written, err = v.Put(file, meta, bytes.NewBuffer(invalid))
	assert.Error(t, err)
	assert.Equal(t, written, 0)

	// get metadata
	meta, err = v.GetMeta(file)
	assert.NoError(t, err)
	assert.Equal(t, meta["foo"].(string), "bar")

	// get file
	writer := bytes.NewBuffer([]byte(""))
	read, err = v.Get(file, writer)
	assert.Equal(t, read, len(plaintext))
	assert.True(t, len(plaintext) > 0, "plaintext length should not be empty")
	assert.NoError(t, err)
	assert.Equal(t, plaintext, writer.Bytes())

	// remove file
	v.Remove(file)
	assert.NoError(t, err)
}

// read stored encrypted files
func RunRegressionTest(t *testing.T, v Vault) {
	file := "The secret file"

	assert.True(t, v.Has(file))

	meta, err := v.GetMeta(file)
	assert.NoError(t, err)

	assert.Equal(t, meta["foo"].(string), "bar")

	writer := bytes.NewBufferString("")
	_, err = v.Get(file, writer)
	assert.NoError(t, err)
	assert.Equal(t, writer.String(), "The secret message")
}

func Test_VaultElement(t *testing.T) {
	ve := NewVaultElement()

	assert.Equal(t, ve.Algo, "aes_ctr") // default value
	assert.NotEmpty(t, ve.Key)
}

func Test_AesEncrypt(t *testing.T) {
	ve := NewVaultElement()
	ve.Key = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	assert.Equal(t, 32, len(ve.Key))

	src := bytes.NewBuffer([]byte("Hello World!!"))
	dst := bytes.NewBuffer([]byte(""))

	AesOFBEncrypter(ve.Key, src, dst)

	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))

	AesOFBDecrypter(ve.Key, dst, decrypted)

	assert.Equal(t, []byte("Hello World!!"), decrypted.Bytes())
}
