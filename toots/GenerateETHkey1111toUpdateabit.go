package main

import (
    "crypto/ecdsa"
    "fmt"
    "log"
	"encoding/hex"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"
)

func main() {
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        log.Fatal(err)
    } 
    privateKeyBytes := crypto.FromECDSA(privateKey)
    //sgj 1111add:
    //var signPrivKey = "dd8bfcf42b66c994478736539ccf0d9e7bb008fa030e4e2c4e3401f938284ebb"
    fmt.Println("privateKeyBytestohex is----001:%s",privateKeyBytes)
    encodeTostr :=hex.EncodeToString(privateKeyBytes)
   
    fmt.Println("privateKeyBytestohex EncodeTostr is----001:%s",encodeTostr)
    
    //privateKeyBytes = []byte(signPrivKey)
    //sgj end add
    fmt.Println("Private: ", hexutil.Encode(privateKeyBytes)[2:])

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }   
    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
    fmt.Println("Pubclic: ", address)

}