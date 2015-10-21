package vault

import (
	//	"bytes"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	//	"fmt"
	"fmt"
)

func getVaultFs(algo string, key []byte) Vault {
	v := &VaultFs{
		Root:    "/tmp/goapp/test/vault",
		Algo:    algo,
		BaseKey: key,
	}

	os.RemoveAll(v.Root)

	return v
}

var largeMessage []byte
var smallMessage []byte
var xLargeMessage []byte

var key = []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141")

func init() {
	smallMessage = []byte("Comment ca va ??")

	largeMessage = make([]byte, 1024*1024*1+2)
	io.ReadFull(rand.Reader, largeMessage)

	fmt.Println("Start generating XLarge message")
	xLargeMessage = make([]byte, 1024*1024*10+2)
	io.ReadFull(rand.Reader, xLargeMessage)
	fmt.Println("End generating XLarge message")
}

func Test_VaultFS_Test_FileExists(t *testing.T) {
	v := getVaultFs("no_op", []byte(""))

	assert.False(t, v.Has("salut"))
}

func Test_VaultFS_Unsecure_Noop(t *testing.T) {
	v := getVaultFs("no_op", []byte(""))

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Secure_Noop(t *testing.T) {
	v := getVaultFs("no_op", key)

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Unsecure_Aes_OFB(t *testing.T) {
	v := getVaultFs("aes_ofb", []byte(""))

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Secure_Aes_OFB(t *testing.T) {
	v := getVaultFs("aes_ofb", key)

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Secure_Aes_OFB_Large(t *testing.T) {
	v := getVaultFs("aes_ofb", key)

	RunTestVault(t, v, largeMessage)
}

func Test_VaultFS_Secure_Aes_CTR(t *testing.T) {
	v := getVaultFs("aes_ctr", key)

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Secure_Aes_CTR_Large(t *testing.T) {
	v := getVaultFs("aes_ctr", key)

	RunTestVault(t, v, largeMessage)
}

func Test_VaultFS_Secure_Aes_CTR_XLarge(t *testing.T) {
	v := getVaultFs("aes_ctr", key)

	RunTestVault(t, v, xLargeMessage)
}

func Test_VaultFS_Secure_Aes_CBC(t *testing.T) {
	v := getVaultFs("aes_cbc", key)

	RunTestVault(t, v, smallMessage)
}

func Test_VaultFS_Secure_Aes_CBC_Large(t *testing.T) {
	v := getVaultFs("aes_cbc", key)

	RunTestVault(t, v, largeMessage)
}

func Test_VaultFS_Secure_Aes_CBC_XLarge(t *testing.T) {
	v := getVaultFs("aes_cbc", key)

	RunTestVault(t, v, xLargeMessage)
}

//func Test_Generate_Regression_Files(t *testing.T) {
//
//	types := []string{"aes_ofb", "aes_ctr", "aes_cbc"}
//
//	for _, v := range types {
//		v := &VaultFs{
//			Root:    "../test/vault/" + v,
//			Algo:    v,
//			BaseKey: []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141"),
//		}
//
//		file := "The secret file"
//		data := bytes.NewBufferString("The secret message")
//		meta := NewVaultMetadata()
//		meta["foo"] = "bar"
//
//		if v.Has(file) {
//			v.Remove(file)
//		}
//
//		v.Put(file, meta, data)
//	}
//}

func Test_VaultFS_Secure_Aes_OFB_NoRegression(t *testing.T) {
	v := &VaultFs{
		Root:    "../test/vault/aes_ofb",
		Algo:    "aes_ofb",
		BaseKey: key,
	}

	RunRegressionTest(t, v)
}

func Test_VaultFS_Secure_Aes_CTR_NoRegression(t *testing.T) {
	v := &VaultFs{
		Root:    "../test/vault/aes_ctr",
		Algo:    "aes_ctr",
		BaseKey: key,
	}

	RunRegressionTest(t, v)
}

func Test_VaultFS_Secure_Aes_CBC_NoRegression(t *testing.T) {
	v := &VaultFs{
		Root:    "../test/vault/aes_cbc",
		Algo:    "aes_cbc",
		BaseKey: key,
	}

	RunRegressionTest(t, v)
}
