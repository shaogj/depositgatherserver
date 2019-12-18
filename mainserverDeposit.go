package main

import (
	"2019NNZXProj10/depositgatherserver/service"
	"encoding/base64"

	//"backend/support/libraries/loggers"
	"flag"
	"fmt"
	"net/http"
	"time"

	//"2019NNZXProj10/depositgatherserver"
	"os"

	"2019NNZXProj10/depositgatherserver/accounts"
	"2019NNZXProj10/depositgatherserver/service/wdctranssign"
	//sgj 1217add

	"2019NNZXProj10/depositgatherserver/service/ktctranssign"


"2019NNZXProj10/depositgatherserver/config"

	"github.com/gorilla/mux"
	mylog "github.com/mkideal/log"

	//"2019NNZXProj10/depositgatherserver/worker"
	"2019NNZXProj10/depositgatherserver/cryptoutil"
	"2019NNZXProj10/depositgatherserver/service/ktctranssign/ktcrpc"

)
var flConfig string

func init() {
	//flag.StringVar(&flConfig, "c", "./config.conf", "config filepath")
//	flag.StringVar(&flConfig, "cmy", "/Users/gejians/go/src/2019NNZXProj10/depositgatherserver/config.conf", "config filepath")
	flag.StringVar(&flConfig, "cmy", "config.conf", "config filepath")
}

//sgj 1105 add for settlecenter:

///113testing:
//Encrypt

func main() {
	err := config.NewConfigTools(flConfig)
	if nil != err {
		mylog.Error("Can't load config error=%v", err)
		os.Exit(0)
	}
	cfgtools := &config.HConf

	err = config.InitConfigInfo()
	if nil != err {
		mylog.Error("from config.json,get json conf err!")
		os.Exit(0)
	}
	gbConf := &config.GbConf
	mylog.Info("--sgj==>get conf info is %v", cfgtools.CurDSCConf)

	err = service.InitMysqlDB(gbConf.MySqlCfg)
	if nil != err {
		mylog.Error("cur InitMysqlDB() to conn err!,err is :%v", err)
		os.Exit(0)

	}
	//1217add
	err = service.InitMysqlDBKTC(gbConf.MySqlCfgKTC)
	if nil != err {
		mylog.Error("cur InitMysqlDB() to conn err!,err is :%v", err)
		os.Exit(0)

	}

	accounts.GJavaSDKUrl = gbConf.JavaSDKUrl
	wdctranssign.GWDCTransHandle.Init(gbConf.WDCTransUrl, &gbConf.WDCConf)
	//wdctranssign.MOrmEngine=service.GXormMysql
	wdctranssign.GWdcDataStore.OrmEngine = service.GXormMysql

	//sgj 1217 add
	//GXormMysqlKTC
	ktctranssign.GKtcDataStore.OrmEngine = service.GXormMysqlKTC

	//"http://192.168.1.211:19585"
	wdctranssign.WDCNodeUrl = gbConf.WDCNodeUrl
	//"http://192.168.1.190:8088/wallet/WalletUtility"
	wdctranssign.WDCJavaSDKUrl = gbConf.JavaSDKUrl

	//1103add
	service.GSettleAccessKey = gbConf.SettleAccessKey.AccessPrivKey

	//一个默认的合作商Key
	getKeystr :=[]byte("1234567812345678")
	/**/
	//2)
	//curPrivkey := "453e53c2594c59c88c8efa629ba0af1dc1af2e7bd423eb9d4dd2fa6666661111"
	//curPrivkey := "c11104cc7bc872eba3af03293b7e1e7fcfc9aa146ff2ace523f260de35564c31"

	//12.17doing,,for encry to db:
	//curPrivkey :="L5DoNfVEtdwEuPkmTYQT11p7dLsmnsnMpKLsD4mZbn3ozEizdv37"
	curPrivkey :="KznoaLNGcSzJUWkBk7FLXFRbNBqkL21SVGn6CMZrEUE2qJRX3SFf"

	//1H9yXTAUS9ndsf9aWu18dU3Z7cafJzqM5i,,,,c11104cc7bc872eba3af03293b7e1e7fcfc9aa146ff2ace523f260de35564c31
	//11.15tesitn
	var encrptedEncodeStr string
	encrpted, err := cryptoutil.AESCBCEncrypt(getKeystr, nil, []byte(curPrivkey))
	if err != nil {
		mylog.Error("ccur Encrypt text is:%s,err is:%v", curPrivkey, err)
	}else{
		encrptedEncodeStr = base64.StdEncoding.EncodeToString(encrpted)
		mylog.Info("cur Encrypt text is:%s, get encrptedEncodeStr len is:%d,val is====44:%s",curPrivkey,len(encrptedEncodeStr),encrptedEncodeStr)
	}
	//2)
	// 对 params 进行 base64 解码
	//sgj 1118 do testing
	//encrptedEncodeStr = "pMWlJaOgTMxybuMoDeiynpUDWDnMcc68zf3HIdMxlcVuPtKCE70dbt8C32jFAKPn4K68AX/nZMBdn8iEVbhfTq19afj36QONEZw1OQCxpvQ="
	                     //MT5jdxLqt6lUhKFMSLc9/3gXTQqyywoUuHTdRTe793re0dloE8P1xHCGkaKCtRimk32Oc7Yetr55m6vIVcqBbCp+vaKk3hG6qAV7R2dFiKw=
	dencrptedEncodeStr, err := base64.StdEncoding.DecodeString(string(encrptedEncodeStr))
	if err != nil {
		mylog.Error("DecodeString text is:%s,err is----AAA:%v",encrptedEncodeStr,err)
		//return nil, err
	}
	/**/
	//fmt.Println("Decrypt get decrpteddecodeStr len is:%d,val is====44:%s,org encrpted len is:",len(dencrptedEncodeStr),dencrptedEncodeStr,len(encrpted))
	delastcrptedaft, err := cryptoutil.AESCBCDecrypt(getKeystr, nil, []byte(dencrptedEncodeStr))
	if err != nil {
		mylog.Error("delastcrptedaft is: %s: decrypt error===888: %v", delastcrptedaft, err)
	}
	delastcrptedaftstr :=string(delastcrptedaft)
	mylog.Info("command %s: decrypt succ===999: %s", "cmdName", delastcrptedaftstr)
	time.Sleep(time.Second * 3)

	//return
	//return
	//1113 测试归集的服务接口调用,gooding
	//for PART Model 2 to wangning

	/*
	gatherAddrCount, bret := wdctranssign.WithdrawsDepositGatherWDC(0, 30, "WDC")
	if bret == true {
		mylog.Info("handle WithdrawsDepositGatherWDC succeed!,total gatherAddrCount is :%d", gatherAddrCount)
	}else{
		mylog.Error("handle WithdrawsDepositGatherWDC failure!,total gatherAddrCount is :%d", gatherAddrCount)

	}
	return
	*/
	//1107doing
	/*
	fromPubkeyStr := "afa367a3b6122afa15b98236ca2bb94577587a5d912629ccbfbac8598daa111d"
	toPubkeyhashStr := "82e55a856e28a84e6d6c5ee431fe883fb237ef6f"
	//--"37ba3c11617d1c0465eb74dab6ca6a0dd81b1dee"
	amount:=0.005
	prikeyStr:= "ec6dca5699ff4d635382cb74fcc73a10e3aa682de36e3a3ff8e7108f5f085e81"
	nonce :=3447
	wdctranssign.GWDCTransHandle.ClientToTransferAccount(fromPubkeyStr,toPubkeyhashStr,float64(amount),prikeyStr,int64(nonce))
	time.Sleep(time.Second * 7)

	return
	end doing 1107*/

	testWdcRPCClient :=new(wdctranssign.WdcRpcClient)

	//1106 add:
	blockHeight,err := testWdcRPCClient.GetBlockHeight()
	if err != nil {
		mylog.Error("RPC GetBlockHeight()  failure! error is :%v,", err)
	}else{
		mylog.Info("RPC GetBlockHeight()  succ! get cur blockHeight is: %d", blockHeight)

	}


	//sgj 1217 adding
	if true!=cfgtools.Disable.DisKTCcoin{
		mylog.Info("--sgj==>KTC RpcConnect() conf info is %v", cfgtools.CurKTCConf)
		ktcrpc.NewKTCRpcClient(&cfgtools.CurKTCConf)
		curKTCClient,err2 := ktcrpc.KTCRPCClient.RpcConnect()

		if err2 != nil {
			mylog.Error("KTC RpcConnect() failure!error is %s: %v. process break!", err2)
			panic(fmt.Sprintf("RpcConnect error:%v", err2))
		}
		mylog.Info("KTC RpcConnect() success! curKTCClient info is %v", curKTCClient)
	}


	router :=mux.NewRouter().StrictSlash(true)


	//sgj 1017 adding:	签名服务器: 成生新的地址(批量)
	router.HandleFunc("/remote/getnewaddress", service.RemoteSignCreateAddress)

	//sgj 1114add for guiji:

	//sgj 1217cadd for guiji:
	router.HandleFunc("/remote/monitorwalletwdc", service.RemoteMonitorWalletAddress)

	strHost:=fmt.Sprintf(":%d",gbConf.WebPort)
	mylog.Info("strHost is :%s", strHost)

	err =http.ListenAndServe(strHost, router)
	if nil!=err {
		mylog.Error("%+v",err)
		os.Exit(0)
	}
	//select {}

}
