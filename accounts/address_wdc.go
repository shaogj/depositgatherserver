package accounts

import (
	"2019NNZXProj10/depositgatherserver/proto"
	. "shaogj/utils"
	"fmt"

	"github.com/mkideal/log"
)

//1028,调用http接口取得json信息
//http://192.168.1.190:8088/wallet/WalletUtility/fromPassword
//1028:
/*
	//serializedKey := privKeyDSC1.PubKey().SerializeCompressed()
	1.创建地址时：
	json 的请求：
	WalletUtility.fromPassword()，返回keystore
	2）
	1.5 通过keystore获得地址
	WalletUtility.keystoreToAddress()
	参数：
	1）、keystore（String)
	2）、密码（String)

	WalletUtility.addressToPubkeyHash()
	3）
	1.3 通过地址获得公钥哈希
	WalletUtility.addressToPubkeyHash()
	参数：
	1）、地址字符串（String)
	返回类型：String（十六进制字符串）
	返回值：pubkeyHash
*/

var GAccountPassword = "1111122222"
var GJavaSDKUrl string = ""

func VerifyAddress(curAddress string) int {

	accountAddr := proto.VerifyAddressReq{}
	accountAddr.Address = curAddress
	resSDKAccount := proto.JavaSDKResponse{}
	var strUrl string
	if GJavaSDKUrl == "" {
		strUrl = "http://192.168.1.190:8088/wallet/WalletUtility"
	} else {
		strUrl = GJavaSDKUrl
	}
	UrlVerify := fmt.Sprintf("%s/%s", strUrl, "verifyAddress")
	ht := CHttpClientEx{}
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 5000, &accountAddr, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  statuscode111=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		//log.Error("transserver.transInfo res err! err is:%v", sttError)
	}
	verifyFalg := -2
	if statusCode == 200 {
		verifyFalg = resSDKAccount.Data.(int)
	}
	log.Info("transserver.verifyAddress,get statusCode is :%s,res=%s", statusCode, strRes)
	return verifyFalg

}

func AddressGenerateWDC() (getprikey string, getaddrpubkey string, getpubkeyhash string, getaddress string, err error) {

	//1028 PMadd:
	stdAccountPassword := GAccountPassword
	accountPassword := proto.AccountPassword{}
	accountPassword.Password = stdAccountPassword
	resSDKAccount := proto.JavaSDKResponse{}

	var strUrl string
	if GJavaSDKUrl == "" {
		strUrl = "http://192.168.1.190:8090/wallet/WalletUtility"
	} else {
		strUrl = GJavaSDKUrl
	}
	UrlVerify := fmt.Sprintf("%s/%s", strUrl, "fromPassword")
	ht := CHttpClientEx{}
	ht.Init()
	//ht.HeaderSet("Content-Type", "text/json")
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//1)通过pass 创建keystore,fromPassword
	strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 5000, &accountPassword, &resSDKAccount)
	if nil != err || statusCode != 200 {
		log.Error("ht.RequestResponseJsonJson  statuscode111=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		//return "","","","",err
	}
	curKeyStoreStr := ""
	if statusCode == 200 {
		curKeyStoreStr = resSDKAccount.Data.(string)

	}
	log.Info("transserver.fromPassword,get statusCode is :%s,res=%s", statusCode, strRes)
	log.Info("json.Unmarshal succ!,cur get resSDKAccount is:%v,get StatusCode is :%s", resSDKAccount, resSDKAccount.StatusCode)

	//返回：  返回数据  网络状态, 错误码  错误信息
	curKeystore := proto.KeystoreToAddress{}
	curKeystore.KsJson = curKeyStoreStr
	curKeystore.PassWord = stdAccountPassword
	log.Info("transserver. get curKeystore is----watching---001:%v", curKeystore)

	//2)通过keystore获得地址,keystoreToAddress
	UrlVerify = fmt.Sprintf("%s/%s", strUrl, "keystoreToAddress")

	strRes, statusCode, errorCode, err = ht.RequestJsonResponseJson(UrlVerify, 5000, &curKeystore, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
	}
	log.Info("transserver.keystoreToAddress,get statusCode is :%s,res=%s", statusCode, strRes)
	curAddressStr := ""
	if statusCode == 200 {
		curAddressStr = resSDKAccount.Data.(string)
	}
	log.Info("transserver. get curAddressStr is----watching---002:%v", curAddressStr)

	//3）通过地址获得公钥哈希
	accountAddress := proto.AddressToPubkeyHash{}
	accountAddress.Address = curAddressStr

	UrlVerify = fmt.Sprintf("%s/%s", strUrl, "addressToPubkeyHash")

	strRes, statusCode, errorCode, err = ht.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return
	}
	log.Info("transserver.addressToPubkeyHash,get statusCode is :%s,res=%s", statusCode, strRes)
	curPubkeyHashStr := ""
	if statusCode == 200 {
		curPubkeyHashStr = resSDKAccount.Data.(string)
		log.Info("transserver. get addressToPubkeyHash succ,value is:%v", curPubkeyHashStr)
	} else {
		log.Error("transserver. get addressToPubkeyHash error!,value is:%v,statusCode is:%s", curPubkeyHashStr, statusCode)
		//return
	}
	//4)通过keystore获得公钥
	UrlVerify = fmt.Sprintf("%s/%s", strUrl, "keystoreToPubkey")

	strRes, statusCode, errorCode, err = ht.RequestJsonResponseJson(UrlVerify, 5000, &curKeystore, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
	}
	curPubkeyStr := ""
	if statusCode == 200 {
		curPubkeyStr = resSDKAccount.Data.(string)
		log.Info("transserver. get keystoreToPubkey succ,value is:%v", curPubkeyStr)
	} else {
		log.Error("transserver. get keystoreToPubkey error!,value is:%v,statusCode is:%s", curPubkeyStr, statusCode)
	}
	//1.4 通过keystore获得私钥
	UrlVerify = fmt.Sprintf("%s/%s", strUrl, "obtainPrikey")

	strRes, statusCode, errorCode, err = ht.RequestJsonResponseJson(UrlVerify, 5000, &curKeystore, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		//log.Error("transserver.transInfo res err! err is:%v", sttError)
	}
	//log.Info("transserver.obtainPrikey,get statusCode is :%s,res=%s",statusCode, strRes)
	log.Info("json.Unmarshal succ!,cur get resSDKAccount is:%v,get StatusCode is :%s", resSDKAccount, resSDKAccount.StatusCode)
	curObtainPrikey := ""
	if statusCode == 200 {
		curObtainPrikey = resSDKAccount.Data.(string)
		log.Info("transserver. get curObtainPrikey succ,value is:%v", curObtainPrikey)
	} else {
		log.Error("transserver. get curObtainPrikey error!,value is:%v,statusCode is:%s", curObtainPrikey, statusCode)
	}

	return curObtainPrikey, curPubkeyStr, curPubkeyHashStr, curAddressStr, nil
}
