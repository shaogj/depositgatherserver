package accounts

import (
	//."dingjm/utils"
	"encoding/hex"
	//"encoding/json"
	"crypto/ecdsa"      //ecdsa
	crand "crypto/rand" //ecdsa
	"github.com/ethereumproject/go-ethereum/crypto"
	"github.com/ethereumproject/go-ethereum/crypto/secp256k1"
	mylog "github.com/mkideal/log"
)

func EtherNewAccount() (string, string, error) {

	rand := crand.Reader
	privateKeyECDSA, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		mylog.Error("%s", err.Error())
		return "", "", err
	}
	Address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	strPrivateKey := hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA))
	return Address.Hex(), strPrivateKey, nil
}
func PublicAddress(stPriv string) (string, error) {
	vecdsa, err := crypto.HexToECDSA(stPriv)
	if nil != err {
		return "", err
	}
	Address := crypto.PubkeyToAddress(vecdsa.PublicKey)
	return Address.Hex(), nil
}
