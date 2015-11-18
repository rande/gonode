// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	KeySize   = 32
	NonceSize = 24
)

var VaultFileExistsError = errors.New("Vault file already exists")

type VaultMetadata map[string]interface{}

func NewVaultMetadata() VaultMetadata {
	return make(VaultMetadata)
}

type VaultDriver interface {
	Has(key string) bool
	GetReader(key string) (io.ReadCloser, error)
	GetWriter(key string) (io.WriteCloser, error)
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
	r, _ := generateRandom(KeySize)

	return r
}

func generateNonce() []byte {
	r, _ := generateRandom(NonceSize)

	return r
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return b, err
	}

	return b, nil
}

func NewVaultElement() *VaultElement {
	return &VaultElement{
		Algo:    "aes_ctr",
		MetaKey: generateKey(),
		BinKey:  generateKey(),
	}
}

type Vault struct {
	Driver  VaultDriver
	Algo    string
	BaseKey []byte
}

func (v *Vault) Has(name string) bool {
	return v.Driver.Has(GetVaultPath(GetVaultKey(name)))
}

func (v *Vault) GetMeta(name string) (vm VaultMetadata, err error) {
	var ve *VaultElement
	var r io.ReadCloser

	vm = NewVaultMetadata()

	// need to get the vaultelement
	vaultname := GetVaultKey(name)

	binfile := GetVaultPath(vaultname)
	metafile := binfile + ".meta"

	// load vault element
	if ve, err = v.getVaultElement(vaultname); err != nil {
		return
	}

	// load metadata
	if r, err = v.Driver.GetReader(metafile); err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte(""))

	if _, err = io.Copy(buf, r); err != nil {
		return
	}

	if err = Unmarshal(v.Algo, ve.MetaKey, buf.Bytes(), &vm); err != nil {
		return
	}

	return
}

func (v *Vault) Get(name string, w io.Writer) (written int64, err error) {
	var ve *VaultElement
	var r io.ReadCloser

	vaultname := GetVaultKey(name)

	binfile := GetVaultPath(vaultname)

	if ve, err = v.getVaultElement(vaultname); err != nil {
		return
	}

	// load binary stream
	if r, err = v.Driver.GetReader(binfile); err != nil {
		return
	}

	defer r.Close()

	return Decrypt(v.Algo, ve.BinKey, r, w)
}

func (v *Vault) Put(name string, meta VaultMetadata, r io.Reader) (written int64, err error) {
	var ve *VaultElement
	var w io.WriteCloser
	var data []byte

	if v.Has(name) {
		return written, VaultFileExistsError
	}

	vaultname := GetVaultKey(name)

	binfile := GetVaultPath(vaultname)
	vaultfile := binfile + ".vault"
	metafile := binfile + ".meta"

	// create vault element
	if ve, err = v.createVaultElement(vaultname); err != nil {
		return
	}

	if data, err = Marshal(v.Algo, ve.MetaKey, meta); err != nil {
		return
	}

	buf := bytes.NewReader(data)

	if w, err = v.Driver.GetWriter(metafile); err != nil {
		v.removeIfExists(vaultfile)
		v.removeIfExists(metafile)

		return
	} else {
		defer w.Close()
	}

	if _, err = io.Copy(w, buf); err != nil {
		v.removeIfExists(vaultfile)
		v.removeIfExists(metafile)

		return
	}

	if w, err = v.Driver.GetWriter(binfile); err != nil {
		v.removeIfExists(vaultfile)
		v.removeIfExists(metafile)

		return
	} else {
		defer w.Close()
	}

	// Copy the input stream to the encryted stream.
	if written, err = Encrypt(v.Algo, ve.BinKey, r, w); err != nil {
		v.removeIfExists(vaultfile)
		v.removeIfExists(metafile)
		v.removeIfExists(binfile)

		return
	}

	return
}

func (v *Vault) Remove(name string) error {
	binfile := GetVaultPath(GetVaultKey(name))

	v.removeIfExists(binfile)
	v.removeIfExists(binfile + ".vault")
	v.removeIfExists(binfile + ".meta")

	return nil
}

func (v *Vault) createVaultElement(namekey []byte) (ve *VaultElement, err error) {
	var data []byte
	var w io.WriteCloser

	ve = NewVaultElement()
	ve.Algo = v.Algo

	vaultfile := GetVaultPath(namekey) + ".vault"

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))
		data, err = Marshal(v.Algo, key, ve)
	} else {
		data, err = json.Marshal(ve)
	}

	if err != nil {
		return
	}

	if w, err = v.Driver.GetWriter(vaultfile); err != nil {
		v.removeIfExists(vaultfile)

		return
	} else {
		defer w.Close()
	}

	if _, err = io.Copy(w, bytes.NewReader(data)); err != nil {
		v.removeIfExists(vaultfile)
	}

	return
}

func (v *Vault) getVaultElement(namekey []byte) (ve *VaultElement, err error) {
	// load vault element
	var r io.ReadCloser

	ve = NewVaultElement()

	// generate the key from the Vault.Base + namekey
	vaultfile := GetVaultPath(namekey) + ".vault"

	if r, err = v.Driver.GetReader(vaultfile); err != nil {
		return
	}

	buf := bytes.NewBuffer([]byte(""))

	if _, err = io.Copy(buf, r); err != nil {
		return
	}

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))
		err = Unmarshal(v.Algo, key, buf.Bytes(), ve)
	} else {
		err = json.Unmarshal(buf.Bytes(), ve)
	}

	return ve, err
}

func (v *Vault) removeIfExists(key string) {
	v.Driver.Remove(key)
}
