/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

//package dashcash
package dashaddress

import (
	"fmt"
	//"github.com/blocktree/bitcoin-adapter/bitcoin"
	"github.com/blocktree/dashcash-adapter/dashcash_addrdec"
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	//"github.com/blocktree/go-owcrypt"
)

func init() {

}

//var (
//	AddressDecoder = &openwallet.AddressDecoder{
//		PrivateKeyToWIF:    PrivateKeyToWIF,
//		PublicKeyToAddress: PublicKeyToAddress,
//		WIFToPrivateKey:    WIFToPrivateKey,
//	}
//)

type addressDecoder struct {
	//wm *WalletManager //钱包管理者
}

/*
//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}
*/
func NewDSCAddressDecoder() *addressDecoder {
	decoder := addressDecoder{}
	return &decoder
}
//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {

	cfg := dashcash_addrdec.DSC_mainnetPrivateWIFCompressed
	if isTestnet {
		cfg = dashcash_addrdec.DSC_testnetPrivateWIFCompressed
	}

	wif, _ := dashcash_addrdec.Default.AddressEncode(priv, cfg)

	return wif, nil

}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	dashcash_addrdec.Default.IsTestNet = isTestnet
	address, err := dashcash_addrdec.Default.AddressEncode(pub)
	if err != nil {
		return "", err
	}

	//if decoder.wm.Config.RPCServerType == bitcoin.RPCServerCore {


	return address, nil

}


//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {

	cfg := dashcash_addrdec.DSC_mainnetPrivateWIFCompressed
	if isTestnet {
		cfg = dashcash_addrdec.DSC_testnetPrivateWIFCompressed
	}

	priv, err := dashcash_addrdec.Default.AddressDecode(wif, cfg)
	if err != nil {
		return nil, err
	}

	return priv, err

}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func (decoder *addressDecoder) ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error) {
	var (
		hash []byte
	)

	cfg := addressEncoder.BTC_mainnetAddressBech32V0
	//sgj del
	if len(scriptPubKey) == 22 || len(scriptPubKey) == 34 {

		hash = scriptPubKey[2:]

		address := addressEncoder.AddressEncode(hash, cfg)

		return address, nil

	} else {
		return "", fmt.Errorf("scriptPubKey length is invalid")
	}

}