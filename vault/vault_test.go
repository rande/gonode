package vault

import (
	"bytes"
	"github.com/stretchr/testify/assert"
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

	message := []byte("Another invalid message with the same key")
	// test overwrite
	written, err = v.Put(file, meta, bytes.NewBuffer(message))
	assert.Error(t, err)
	assert.Equal(t, written, 0)

	// get metadata
	meta, err = v.GetMeta(file)
	assert.NoError(t, err)
	assert.Equal(t, meta["foo"].(string), "bar")

	// get file
	writer := bytes.NewBuffer([]byte(""))
	_, err = v.Get(file, writer)
	assert.NoError(t, err)
	assert.Equal(t, []byte("Comment ca va ??"), writer.Bytes())

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
