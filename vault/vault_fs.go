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
	Root      string
	Encrypter Encrypter
	Decrypter Decrypter
	Key       []byte
}

func (v *VaultFs) Has(name string) bool {
	path := filepath.Join(v.Root, GetVaultPath(name))

	if _, error := os.Stat(path); error != nil {
		return false
	}

	return true
}

func (v *VaultFs) Get(name string) (VaultMetadata, error) {
	// need to get the vaultelement
	binfile := filepath.Join(v.Root, GetVaultPath(name))
	vaultfile := binfile + ".vault"
	metafile := binfile + ".meta"

	// load vault element
	ve, err := v.getVaultElement(vaultfile)

	if err != nil {
		return nil, err
	}

	// load metadata
	meta := NewVaultMetada()

	// load binary stream
	fm, _ := os.Open(metafile)
	err = json.NewDecoder(v.Decrypter(ve.Key, fm)).Decode(&meta)

	return meta, err
}

func (v *VaultFs) GetReader(name string) (io.Reader, error) {
	binfile := filepath.Join(v.Root, GetVaultPath(name))
	vaultfile := binfile + ".vault"

	ve, err := v.getVaultElement(vaultfile)

	if err != nil {
		return nil, err
	}

	// load binary stream
	fb, err := os.Open(binfile)

	if err != nil {
		return nil, err
	}

	return v.Decrypter(ve.Key, fb), nil
}

func (v *VaultFs) Put(name string, meta VaultMetadata, r io.Reader) (int64, error) {
	if v.Has(name) {
		return 0, VaultFileExistsError
	}

	binfile := filepath.Join(v.Root, GetVaultPath(name))

	vaultfile := binfile + ".vault"
	metafile := binfile + ".meta"

	// create base folder
	path := filepath.Dir(binfile)
	if err := os.MkdirAll(path, 0700); err != nil {
		return 0, err
	}

	// store vault element
	ve := NewVaultElement()

	fv, err := os.Create(vaultfile)
	if err != nil {
		return 0, err
	}

	defer fv.Close()

	if len(v.Key) > 0 {
		err = json.NewEncoder(v.Encrypter(v.Key, fv)).Encode(ve)
	} else {
		err = json.NewEncoder(fv).Encode(ve)
	}

	if err != nil {
		defer RemoveIfExists(vaultfile)

		return 0, err
	}

	// store metafile
	fm, err := os.Create(metafile)
	if err != nil {
		defer RemoveIfExists(vaultfile)

		return 0, err
	}
	defer fm.Close()

	err = json.NewEncoder(v.Encrypter(ve.Key, fm)).Encode(meta)
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
	if written, err := io.Copy(v.Encrypter(ve.Key, fb), r); err != nil {
		defer RemoveIfExists(vaultfile)
		defer RemoveIfExists(metafile)
		defer RemoveIfExists(binfile)

		return 0, err
	} else {
		return written, nil
	}
}

func (v *VaultFs) Remove(name string) error {
	binfile := filepath.Join(v.Root, GetVaultPath(name))

	RemoveIfExists(binfile)
	RemoveIfExists(binfile + ".vault")
	RemoveIfExists(binfile + ".meta")

	return nil
}

func (v *VaultFs) getVaultElement(vaultfile string) (*VaultElement, error) {
	// load vault element
	var err error
	ve := NewVaultElement()
	fv, _ := os.Open(vaultfile)
	if len(v.Key) > 0 {
		err = json.NewDecoder(v.Decrypter(v.Key, fv)).Decode(ve)
	} else {
		err = json.NewDecoder(fv).Decode(ve)
	}

	return ve, err
}
