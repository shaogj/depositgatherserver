package cryptoutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

var (
	ErrInvalidPublicKeyType = errors.New("invalid public key type")
	ErrEmptyPEMBlock        = errors.New("empty pem block")
)

func LoadRSAPrivateKey(pemFile string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrEmptyPEMBlock
	}
	return UnmarshalPrivkey(block.Bytes)
}

func SaveRSAPrivateKey(pemFile string, privKey *rsa.PrivateKey) error {
	bytes := MarshalPrivkey(privKey)
	block := &pem.Block{
		Type:  "PRIVATE RSA KEY",
		Bytes: bytes,
	}
	f, err := os.Create(pemFile)
	if err != nil {
		return err
	}
	return pem.Encode(f, block)
}

func LoadRSAPublicKey(pemFile string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrEmptyPEMBlock
	}
	return UnmarshalPubkey(block.Bytes)
}

func SaveRSAPublicKey(pemFile string, pubKey *rsa.PublicKey) error {
	bytes, err := MarshalPubkey(pubKey)
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type:  "PUBLIC RSA KEY",
		Bytes: bytes,
	}
	f, err := os.Create(pemFile)
	if err != nil {
		return err
	}
	return pem.Encode(f, block)
}

func MarshalPubkey(pubkey *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pubkey)
}

func UnmarshalPubkey(data []byte) (*rsa.PublicKey, error) {
	p, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	pubkey, ok := p.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPublicKeyType
	}
	return pubkey, nil
}

func MarshalPrivkey(privkey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privkey)
}

func UnmarshalPrivkey(data []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(data)
}
