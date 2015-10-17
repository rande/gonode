package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func GetCipher(name string) (Encrypter, Decrypter) {

	switch name {
	case "aes_ofb":
		return AesOFBEncrypter, AesOFBDecrypter
	case "aes_ctr":
		return AesCTREncrypter, AesCTRDecrypter
	case "no_op":
		return NoopEncrypter, NoopDecrypter

	default:
		panic("Unable to find the cipher")
	}
}

func GetEncryptionWriter(name string, key interface{}, w io.Writer) io.Writer {
	e, _ := GetCipher(name)

	return e(key, w)
}

func GetDecryptionReader(name string, key interface{}, r io.Reader) io.Reader {
	_, d := GetCipher(name)

	return d(key, r)
}

func NoopEncrypter(key interface{}, w io.Writer) io.Writer {
	return w
}

func NoopDecrypter(key interface{}, r io.Reader) io.Reader {
	return r
}

func GetAes(key interface{}) cipher.Block {
	block, err := aes.NewCipher(key.([]byte))
	if err != nil {
		panic(err)
	}

	return block
}

func AesOFBEncrypter(key interface{}, w io.Writer) io.Writer {
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(GetAes(key), iv[:])

	return &cipher.StreamWriter{S: stream, W: w}
}

func AesOFBDecrypter(key interface{}, r io.Reader) io.Reader {
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(GetAes(key), iv[:])

	return &cipher.StreamReader{S: stream, R: r}
}

func AesCTREncrypter(key interface{}, w io.Writer) io.Writer {
	var iv [aes.BlockSize]byte
	stream := cipher.NewCTR(GetAes(key), iv[:])

	return &cipher.StreamWriter{S: stream, W: w}
}

func AesCTRDecrypter(key interface{}, r io.Reader) io.Reader {
	var iv [aes.BlockSize]byte
	stream := cipher.NewCTR(GetAes(key), iv[:])

	return &cipher.StreamReader{S: stream, R: r}
}
