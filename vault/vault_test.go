package vault

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

// write/encrypted file
func RunTestVault(t *testing.T, v Vault) {
	file := "this-is-a-test"

	meta := NewVaultMetadata()
	meta["foo"] = "bar"

	reader := bytes.NewBuffer([]byte("Comment ca va ??"))

	written, err := v.Put(file, meta, reader)

	assert.NoError(t, err)
	assert.Equal(t, written, 16)
	assert.True(t, v.Has(file))

	// test overwrite
	written, err = v.Put(file, meta, bytes.NewBuffer([]byte("Another content")))
	assert.Error(t, err)
	assert.Equal(t, written, 0)

	// get metadata
	meta, err = v.Get(file)
	assert.NoError(t, err)
	assert.Equal(t, meta["foo"].(string), "bar")

	// get reader
	stream, err := v.GetReader(file)
	dst := bytes.NewBuffer([]byte(""))
	io.Copy(dst, stream)
	assert.Equal(t, []byte("Comment ca va ??"), dst.Bytes())

	// remove file
	v.Remove(file)
	assert.NoError(t, err)
}

// read stored encrypted files
func RunRegressionTest(t *testing.T, v Vault) {
	file := "The secret file"

	assert.True(t, v.Has(file))

	meta, err := v.Get(file)
	assert.NoError(t, err)

	assert.Equal(t, meta["foo"].(string), "bar")

	reader, err := v.GetReader(file)
	assert.NoError(t, err)

	data := bytes.NewBufferString("")
	data.ReadFrom(reader)

	assert.Equal(t, data.String(), "The secret message")
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

	io.Copy(AesOFBEncrypter(ve.Key, dst), src)

	assert.NotEmpty(t, dst.String())

	decrypted := bytes.NewBuffer([]byte(""))

	io.Copy(decrypted, AesOFBDecrypter(ve.Key, dst))

	assert.Equal(t, []byte("Hello World!!"), decrypted.Bytes())
}
