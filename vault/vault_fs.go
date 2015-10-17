package vault

import (
	"encoding/json"
	"io"
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

func (v *VaultFs) Get(name string) (VaultMetadata, error) {
	// need to get the vaultelement
	vaultname := GetVaultKey(name)
	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))
	metafile := binfile + ".meta"

	// load vault element
	ve, err := v.getVaultElement(vaultname)

	if err != nil {
		return nil, err
	}

	// load metadata
	meta := NewVaultMetadata()

	// load binary stream
	fm, _ := os.Open(metafile)
	err = json.NewDecoder(GetDecryptionReader(v.Algo, ve.Key, fm)).Decode(&meta)

	return meta, err
}

func (v *VaultFs) GetReader(name string) (io.Reader, error) {
	vaultname := GetVaultKey(name)

	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))

	ve, err := v.getVaultElement(vaultname)

	if err != nil {
		return nil, err
	}

	// load binary stream
	fb, err := os.Open(binfile)
	if err != nil {
		return nil, err
	}

	return GetDecryptionReader(v.Algo, ve.Key, fb), nil
}

func (v *VaultFs) Put(name string, meta VaultMetadata, r io.Reader) (int64, error) {
	if v.Has(name) {
		return 0, VaultFileExistsError
	}

	vaultname := GetVaultKey(name)

	binfile := filepath.Join(v.Root, GetVaultPath(vaultname))
	vaultfile := binfile + ".vault"
	metafile := binfile + ".meta"

	ve, err := v.createVaultElement(vaultname)

	if err != nil {
		return 0, nil
	}

	// store metafile
	fm, err := os.Create(metafile)
	if err != nil {
		defer RemoveIfExists(vaultfile)

		return 0, err
	}
	defer fm.Close()

	err = json.NewEncoder(GetEncryptionWriter(v.Algo, ve.Key, fm)).Encode(meta)
	if err != nil {
		defer RemoveIfExists(vaultfile)
		defer RemoveIfExists(metafile)

		return 0, err
	}

	// store binary
	fb, err := os.Create(binfile)
	defer fb.Close()

	if err != nil {
		defer RemoveIfExists(vaultfile)
		defer RemoveIfExists(metafile)

		return 0, err
	}

	// Copy the input stream to the encryted stream.
	if written, err := io.Copy(GetEncryptionWriter(v.Algo, ve.Key, fb), r); err != nil {
		defer RemoveIfExists(vaultfile)
		defer RemoveIfExists(metafile)
		defer RemoveIfExists(binfile)

		return 0, err
	} else {
		return written, nil
	}
}

func (v *VaultFs) Remove(name string) error {
	binfile := filepath.Join(v.Root, GetVaultPath(GetVaultKey(name)))

	RemoveIfExists(binfile)
	RemoveIfExists(binfile + ".vault")
	RemoveIfExists(binfile + ".meta")

	return nil
}

func (v *VaultFs) createVaultElement(namekey []byte) (*VaultElement, error) {
	vaultfile := filepath.Join(v.Root, GetVaultPath(namekey)) + ".vault"

	// create base folder
	path := filepath.Dir(vaultfile)
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}

	// store vault element
	ve := NewVaultElement()
	ve.Algo = v.Algo

	fv, err := os.Create(vaultfile)
	if err != nil {
		return nil, err
	}

	defer fv.Close()

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))

		err = json.NewEncoder(GetEncryptionWriter(v.Algo, key, fv)).Encode(ve)
	} else {
		err = json.NewEncoder(fv).Encode(ve)
	}

	if err != nil {
		defer RemoveIfExists(vaultfile)

		return nil, err
	}

	return ve, nil
}

func (v *VaultFs) getVaultElement(namekey []byte) (*VaultElement, error) {
	// load vault element
	var err error

	// generate the key from the Vault.Base + namekey
	vaultfile := filepath.Join(v.Root, GetVaultPath(namekey)) + ".vault"

	ve := NewVaultElement()
	fv, _ := os.Open(vaultfile)

	if len(v.BaseKey) > 0 {
		key := GetVaultKey(string(append(namekey, v.BaseKey...)[:]))
		err = json.NewDecoder(GetDecryptionReader(v.Algo, key, fv)).Decode(ve)
	} else {
		err = json.NewDecoder(fv).Decode(ve)
	}

	return ve, err
}
