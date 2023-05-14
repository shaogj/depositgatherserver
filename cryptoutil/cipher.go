package cryptoutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	ErrBlockSize = errors.New("error block size")
)

var iv16 = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var iv32 = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func fixIV(iv []byte, block cipher.Block) []byte {
	if len(iv) != block.BlockSize() {
		switch block.BlockSize() {
		case 16:
			return iv16
		case 32:
			return iv32
		default:
			return bytes.Repeat([]byte{0}, block.BlockSize())
		}
	}
	return iv
}

func AESCBCEncrypt(key, iv, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv = fixIV(iv, block)
	cbc := cipher.NewCBCEncrypter(block, iv)
	src = PKCS7Padding(src, block.BlockSize())
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	return dst, nil
}

func AESCBCDecrypt(key, iv, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv = fixIV(iv, block)
	cbc := cipher.NewCBCDecrypter(block, iv)
	remainder := len(src) % block.BlockSize()
	if remainder != 0 {
		return nil, ErrBlockSize
	}
	dst := make([]byte, len(src))
	cbc.CryptBlocks(dst, src)
	dst, err = PKCS7UnPadding(dst, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return dst, nil
}
