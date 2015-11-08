package vault

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

var VaultFileExistsError = errors.New("Vault file already exists")

type VaultMetadata map[string]interface{}

func NewVaultMetadata() VaultMetadata {
	return make(VaultMetadata)
}

type Vault interface {
	Has(key string) bool
	GetMeta(key string) (VaultMetadata, error)
	Get(key string, w io.Writer) (int64, error)
	Put(key string, meta VaultMetadata, r io.Reader) (int64, error)
	Remove(key string) error
}

type VaultElement struct {
	MetaKey []byte `json:"meta_key"`
	BinKey  []byte `json:"bin_key"`
	Algo    string `json:"algo"`
	Hash    string `json:"hash"`
}

func GetVaultKey(name string) []byte {
	sum := sha256.Sum256([]byte(name))

	return sum[:]
}

func GetVaultPath(sum []byte) string {
	return fmt.Sprintf("%x/%x/%x.bin", sum[0:1], sum[1:2], sum[2:])
}

func generateKey() []byte {

	b := make([]byte, 32)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return b
}

func NewVaultElement() *VaultElement {
	return &VaultElement{
		Algo:    "aes_ctr",
		MetaKey: generateKey(),
		BinKey:  generateKey(),
	}
}
