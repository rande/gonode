package vault

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getVaultFs(e Encrypter, d Decrypter, key []byte) Vault {
	v := &VaultFs{
		Root:      "/tmp/goapp/test/vault",
		Encrypter: e,
		Decrypter: d,
		BaseKey:   key,
	}

	os.RemoveAll(v.Root)

	return v
}

func Test_VaultFS_Test_FileExists(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, []byte(""))

	assert.False(t, v.Has("salut"))
}

func Test_VaultFS_Unsecure_Noop(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, []byte(""))

	RunTestVault(t, v)
}

func Test_VaultFS_Secure_Noop(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141"))

	RunTestVault(t, v)
}

func Test_VaultFS_Unsecure_Aes_OFB(t *testing.T) {
	v := getVaultFs(AesOFBEncrypter, AesOFBDecrypter, []byte(""))

	RunTestVault(t, v)
}

func Test_VaultFS_Secure_Aes_OFB(t *testing.T) {
	v := getVaultFs(AesOFBEncrypter, AesOFBDecrypter, []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141"))

	RunTestVault(t, v)
}

func Test_VaultFS_Secure_Aes_OFB_NoRegression(t *testing.T) {
	v := &VaultFs{
		Root:      "../test/vault/aes/ofb",
		Encrypter: AesOFBEncrypter,
		Decrypter: AesOFBDecrypter,
		BaseKey:   []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141"),
	}

	file := "The secret file"

	assert.True(t, v.Has(file))
	
	meta, err := v.Get(file)
	assert.NoError(t, err)

	assert.Equal(t, meta["foo"].(string), "bar")

	reader, err := v.GetReader(file)
	assert.NoError(t, err)

	data := bytes.NewBufferString("")
	data.ReadFrom(reader)

	assert.Equal(t, data.String(), "The secret content")
}
