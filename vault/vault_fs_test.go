package vault

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func getVaultFs(e Encrypter, d Decrypter, key []byte) Vault {
	v := &VaultFs{
		Root:      "/tmp/goapp/test/vault",
		Encrypter: e,
		Decrypter: d,
		Key:       key,
	}

	os.RemoveAll(v.Root)

	return v
}

func Test_VaultFS_Test_FileExists(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, []byte(""))

	assert.False(t, v.Has("salut"))
}

func Test_VaultFS_Unsecure_Aes(t *testing.T) {
	v := getVaultFs(AesOFBEncrypter, AesOFBDecrypter, []byte(""))

	RunTestVault(t, v)
}

func Test_VaultFS_Unsecure_Noop(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, []byte(""))

	RunTestVault(t, v)
}

func Test_VaultFS_Secure_Noop(t *testing.T) {
	v := getVaultFs(NoopEncrypter, NoopDecrypter, generateKey())

	RunTestVault(t, v)
}

func Test_VaultFS_Secure_Aes(t *testing.T) {
	v := getVaultFs(AesOFBEncrypter, AesOFBDecrypter, generateKey())

	RunTestVault(t, v)
}