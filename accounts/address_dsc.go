package accounts

//var chainParamsBCH = &chaincfg.MainNetParams

func AddressGenerateDSC() (getprikey string, getaddrpubkey string, getaddress string, err error) {

	//sgj 1104 tmpskip
	return
	/*
	privKeyDSC1, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println("for BCH: private key generation error: %s", err)
		return
	}
	serializedKey := privKeyDSC1.PubKey().SerializeCompressed()
	//addr3, err := cashutil.NewAddressPubKeyHash(serializedKey, chainParams)
	//sgj 1017doing for dsc
	//pubKeyAddrStruct, err := cashutil.NewAddressPubKey(serializedKey, chainParamsBCH)
	//err != nil {
		//log.Error("for BCH: pubKeyAddr key generation error: %v\n", err)

	//do: 公钥转地址
	curAddressDecoder:= dashaddress.NewDSCAddressDecoder()
	pubKeyAddr, err := curAddressDecoder.PublicKeyToAddress(serializedKey,false)
	if err != nil {
		log.Error("for DSC: pubKeyAddr key generation error: %v\n", err)
		return
	}

	//tonewaddressbch :=pubKeyAddr.AddressPubKeyHash()
	//sgj 1017BB
	//tonewaddressbch :="AddressPubKeyHash()---07"
	wif1, err := cashutil.NewWIF(privKeyDSC1, chainParamsBCH, false)
	if err != nil {
		fmt.Printf("NewWIF(), get wif1 is err,wifi :%v \n", wif1)
	}
	//privKeyDSCstr := (*cashutil.WIF)(wif1).String()

	//var strserializedPriKey string = string(privKeyDSC1.Serialize()[:])
	privKeyDSCstr, err := curAddressDecoder.PrivateKeyToWIF(privKeyDSC1.Serialize(),false)

	if err != nil {
		fmt.Println("for DSC: private key generation error: %s", err)
		return
	}

	//log.Info("get cur BCH' privKey is :%v\n",privKeyDSCstr)
	log.Info("get cur DSC' Prikey is :%v,pubKeyAddr is :%v\n",privKeyDSC1.Serialize(),pubKeyAddr)
	log.Info("get cur DSC' PubKey is:%v\n",privKeyDSC1.PubKey().SerializeUncompressed())

	//sgj 0817 add:

	var strserializedKey string = string(serializedKey[:])

	//return privKeyDSCstr,pubKeyAddr.String(),tonewaddressbch.String(),nil
	//strserializedPriKey,strserializedKey,pubKeyAddrStruct
	return privKeyDSCstr,strserializedKey,pubKeyAddr,nil

	*/
}
