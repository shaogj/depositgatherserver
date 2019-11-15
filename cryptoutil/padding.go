package cryptoutil

import (
	"errors"
)

var (
	ErrEmptyBlock  = errors.New("empty block")
	ErrPaddingSize = errors.New("invalid padding size")
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	remainder := len(ciphertext) % blockSize
	padByte := byte(blockSize - remainder)
	if cap(ciphertext)-len(ciphertext) < int(padByte) {
		pad := make([]byte, 0, padByte)
		for i := byte(0); i < padByte; i++ {
			pad = append(pad, padByte)
		}
		return append(ciphertext, pad...)
	} else {
		for i := byte(0); i < padByte; i++ {
			ciphertext = append(ciphertext, padByte)
		}
		return ciphertext
	}
}

func PKCS7UnPadding(plainText []byte, blockSize int) ([]byte, error) {
	length := len(plainText)
	if length == 0 {
		return nil, ErrEmptyBlock
	}
	unpadding := int(plainText[length-1])
	if length < unpadding || length-unpadding > len(plainText) {
		return nil, ErrPaddingSize
	}
	return plainText[:(length - unpadding)], nil
}
