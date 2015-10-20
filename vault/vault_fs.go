package vault

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RemoveIfExists(path string) {
	if _, error := os.Stat(path); error != nil {
		return
	}

	os.Remove(path)
}

type VaultFs struct {
	Root    string
	Algo    string
	BaseKey []byte
}

func (v *VaultFs) Has(name string) bool {
	path := filepath.Join(v.Root, GetVaultPath(GetVaultKey(name)))

	if _, error := os.Stat(path); error != nil {
		return false
	}

	return true
}

func (v *VaultFs) GetMeta(name string) (vm VaultMetadata, err error) {
	var ve *VaultElement
	var data []byte

	vm = NewVaultMetadata()

	// need to get the vaultelement
	vaultname := GetVaultKey(name)
	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))
	metafile := binfile + ".meta"

	// load vault element
	if ve, err = v.getVaultElement(vaultname); err != nil {
		return
	}

	// load metadata
	if data, err = ioutil.ReadFile(metafile); err != nil {
		return
	}

	if err = Unmarshal(v.Algo, ve.Key, data, &vm); err != nil {
		return
	}

	return
}

func (v *VaultFs) Get(name string, w io.Writer) (written int64, err error) {
	var ve *VaultElement
	var fb *os.File

	vaultname := GetVaultKey(name)

	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))

	if ve, err = v.getVaultElement(vaultname); err != nil {
		return
	}

	// load binary stream
	if fb, err = os.Open(binfile); err != nil {
		return
	}

	return Decrypt(v.Algo, ve.Key, fb, w)
}

func (v *VaultFs) Put(name string, meta VaultMetadata, r io.Reader) (written int64, err error) {
	var ve *VaultElement
	var fb *os.File
	var data []byte

	if v.Has(name) {
		return written, VaultFileExistsError
	}

	vaultname := GetVaultKey(name)

	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))
	vaultfile := binfile + ".vault"
	metafile := binfile + ".meta"

	// create vault element
	if ve, err = v.createVaultElement(vaultname); err != nil {
		return
	}

	if data, err = Marshal(v.Algo, ve.Key, meta); err != nil {
		return
	}

	if err = ioutil.WriteFile(metafile, data, 0600); err != nil {
		RemoveIfExists(vaultfile)
		RemoveIfExists(metafile)

		return
	}

	if fb, err = os.Create(binfile); err != nil {
		RemoveIfExists(vaultfile)
		RemoveIfExists(metafile)

		return
	}

	defer fb.Close()

	// Copy the input stream to the encryted stream.
	if written, err = Encrypt(v.Algo, ve.Key, r, fb); err != nil {
		RemoveIfExists(vaultfile)
		RemoveIfExists(metafile)
		RemoveIfExists(binfile)

		return
	}

	return
}

func (v *VaultFs) Remove(name string) error {
	binfile := filepath.Join(v.Root, GetVaultPath(GetVaultKey(name)))

	RemoveIfExists(binfile)
	RemoveIfExists(binfile + ".vault")
	RemoveIfExists(binfile + ".meta")

	return nil
}

func (v *VaultFs) createVaultElement(namekey []byte) (ve *VaultElement, err error) {
	var data []byte

	ve = NewVaultElement()
	ve.Algo = v.Algo

	vaultfile := filepath.Join(v.Root, GetVaultPath(namekey)) + ".vault"

	// create base folder
	path := filepath.Dir(vaultfile)
	if err = os.MkdirAll(path, 0700); err != nil {
		return
	}

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))
		data, err = Marshal(v.Algo, key, ve)
	} else {
		data, err = json.Marshal(ve)
	}

	if err != nil {
		return
	}

	if err = ioutil.WriteFile(vaultfile, data, 0600); err != nil {
		RemoveIfExists(vaultfile)

		return
	}

	return
}

func (v *VaultFs) getVaultElement(namekey []byte) (ve *VaultElement, err error) {
	// load vault element
	var data []byte

	ve = NewVaultElement()

	// generate the key from the Vault.Base + namekey
	vaultfile := filepath.Join(v.Root, GetVaultPath(namekey)) + ".vault"

	if data, err = ioutil.ReadFile(vaultfile); err != nil {
		return
	}

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))
		err = Unmarshal(v.Algo, key, data, ve)
	} else {
		err = json.Unmarshal(data, ve)
	}

	return ve, err
}
