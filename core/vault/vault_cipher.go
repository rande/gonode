// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Encrypter func(key interface{}, r io.Reader, w io.Writer) (int64, error)

type Decrypter func(key interface{}, r io.Reader, w io.Writer) (int64, error)

func GetCipher(mode string) (Encrypter, Decrypter) {

	switch mode {
	case "aes_ofb":
		return AesOFBEncrypter, AesOFBDecrypter
	case "aes_ctr":
		return AesCTREncrypter, AesCTRDecrypter
	case "aes_cbc":
		return AesCBCEncrypter, AesCBCDecrypter
	case "aes_gcm":
		return AesGCMEncrypter, AesGCMDecrypter
	case "no_op":
		return NoopEncrypter, NoopDecrypter

	default:
		panic("Unable to find the cipher: " + mode)
	}
}

func Marshal(mode string, key interface{}, v interface{}) (data []byte, err error) {
	if data, err = json.Marshal(v); err != nil {
		return
	}

	src := bytes.NewBuffer(data)
	dst := bytes.NewBuffer([]byte(""))

	if _, err = Encrypt(mode, key, src, dst); err != nil {
		return
	}

	return dst.Bytes(), nil
}

func Unmarshal(mode string, key interface{}, data []byte, v interface{}) (err error) {
	src := bytes.NewBuffer(data)
	dst := bytes.NewBuffer([]byte(""))

	if _, err = Decrypt(mode, key, src, dst); err != nil {
		return
	}

	if err = json.Unmarshal(dst.Bytes(), v); err != nil {
		return
	}

	return
}

func Encrypt(mode string, key interface{}, r io.Reader, w io.Writer) (int64, error) {
	e, _ := GetCipher(mode)

	return e(key, r, w)
}

func Decrypt(mode string, key interface{}, r io.Reader, w io.Writer) (int64, error) {
	_, d := GetCipher(mode)

	return d(key, r, w)
}

func NoopEncrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	return io.Copy(w, r)
}

func NoopDecrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	return io.Copy(w, r)
}

func getAes(key interface{}) cipher.Block {
	block, err := aes.NewCipher(key.([]byte))
	if err != nil {
		panic(err)
	}

	return block
}

func AesOFBEncrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(getAes(key), iv[:])

	return NoopEncrypter(key, r, &cipher.StreamWriter{S: stream, W: w})
}

func AesOFBDecrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(getAes(key), iv[:])

	return NoopDecrypter(key, &cipher.StreamReader{S: stream, R: r}, w)
}

func AesCTREncrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	var iv [aes.BlockSize]byte
	stream := cipher.NewCTR(getAes(key), iv[:])

	return NoopEncrypter(key, r, &cipher.StreamWriter{S: stream, W: w})
}

func AesCTRDecrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	var iv [aes.BlockSize]byte
	stream := cipher.NewCTR(getAes(key), iv[:])

	return NoopDecrypter(key, &cipher.StreamReader{S: stream, R: r}, w)
}

// this implementation required to load all information into memory before encrypting
// data.
func AesGCMEncrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	gcm, err := cipher.NewGCM(getAes(key))
	if err != nil {
		return 0, nil
	}

	nonce, _ := generateRandom(NonceSize)

	plaintext, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, nil
	}

	out := gcm.Seal(nonce, nonce, plaintext, nil)

	written, err := w.Write(out)

	return int64(written), err
}

// this implementation required to load all information into memory before decrypting
// data.
func AesGCMDecrypter(key interface{}, r io.Reader, w io.Writer) (int64, error) {
	gcm, err := cipher.NewGCM(getAes(key))
	if err != nil {
		return 0, err
	}

	ciphertext, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	nonce := make([]byte, NonceSize)
	copy(nonce, ciphertext)

	out, err := gcm.Open(nil, nonce, ciphertext[NonceSize:], nil)
	if err != nil {
		return 0, err
	}

	written, err := w.Write(out)

	return int64(written), err
}

func pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)

	return append(data, pad...), nil
}

func unpad(data []byte, blocklen int) ([]byte, error) {

	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}

	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}

	padlen := int(data[len(data)-1])

	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("invalid padding, padlen invalid")
	}

	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padlen], nil
}

func AesCBCEncrypter(key interface{}, r io.Reader, w io.Writer) (written int64, err error) {
	var twritten, read int
	aes, err := aes.NewCipher(key.([]byte))

	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize())

	block := cipher.NewCBCEncrypter(aes, iv)

	for {
		buf := make([]byte, block.BlockSize()-1)
		if read, err = io.ReadFull(r, buf); err != nil {
			if err == io.EOF {
				return written, nil
			} else if err == io.ErrUnexpectedEOF {
				// nothing
			} else {
				return
			}
		}

		if buf, err = pad(buf[:read], block.BlockSize()); err != nil {
			return written, err
		}

		block.CryptBlocks(buf, buf)

		if twritten, err = w.Write(buf); err != nil {
			return
		} else {
			written += int64(twritten)
		}
	}
}

func AesCBCDecrypter(key interface{}, r io.Reader, w io.Writer) (written int64, err error) {
	var twritten int

	aes, err := aes.NewCipher(key.([]byte))

	if err != nil {
		return 0, err
	}

	iv := make([]byte, aes.BlockSize())

	d := cipher.NewCBCDecrypter(aes, iv)

	for {

		buf := make([]byte, d.BlockSize())

		if _, err = io.ReadFull(r, buf); err != nil {
			if err == io.EOF {
				return written, nil
			} else {
				return
			}
		}

		d.CryptBlocks(buf, buf)

		if buf, err = unpad(buf, d.BlockSize()); err != nil {
			return
		}
		if twritten, err = w.Write(buf); err != nil {
			return
		} else {
			written += int64(twritten)
		}
	}
}
