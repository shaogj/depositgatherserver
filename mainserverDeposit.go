package main

import (
	"2019NNZXProj10/abitserverDepositeGather/service"
	"encoding/base64"

	//"backend/support/libraries/loggers"
	"flag"
	"fmt"
	"net/http"
	"time"

	//"2019NNZXProj10/abitserverDepositeGather"
	"os"

	"2019NNZXProj10/abitserverDepositeGather/accounts"
	"2019NNZXProj10/abitserverDepositeGather/service/wdctranssign"

	"2019NNZXProj10/abitserverDepositeGather/config"

	"github.com/gorilla/mux"
	mylog "github.com/mkideal/log"

	//"2019NNZXProj10/abitserverDepositeGather/worker"
	"2019NNZXProj10/abitserverDepositeGather/cryptoutil"
)
var flConfig string

func init() {
	//flag.StringVar(&flConfig, "c", "./config.conf", "config filepath")
//	flag.StringVar(&flConfig, "cmy", "/Users/gejians/go/src/2019NNZXProj10/abitserverDepositeGather/config.conf", "config filepath")
	flag.StringVar(&flConfig, "cmy", "config.conf", "config filepath")
}

//sgj 1105 add for settlecenter:

///113testing:
//Encrypt

func main() {
	fmt.Println("vim-go")
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
	accounts.GJavaSDKUrl = gbConf.JavaSDKUrl
	wdctranssign.GWDCTransHandle.Init(gbConf.WDCTransUrl, &gbConf.WDCConf)
	//wdctranssign.MOrmEngine=service.GXormMysql
	wdctranssign.GWdcDataStore.OrmEngine = service.GXormMysql

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
	curPrivkey := "453e53c2594c59c88c8efa629ba0af1dc1af2e7bd423eb9d4dd2fa6666661111"

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
	dencrptedEncodeStr, err := base64.StdEncoding.DecodeString(string(encrptedEncodeStr))
	if err != nil {
		mylog.Error("DecodeString text is:%s,err is----AAA:%v",encrptedEncodeStr,err)
		//return nil, err
	}
	//fmt.Println("Decrypt get decrpteddecodeStr len is:%d,val is====44:%s,org encrpted len is:",len(dencrptedEncodeStr),dencrptedEncodeStr,len(encrpted))
	delastcrptedaft, err := cryptoutil.AESCBCDecrypt(getKeystr, nil, []byte(dencrptedEncodeStr))
	if err != nil {
		mylog.Error("command %s: decrypt error===888: %v", "cmdName", delastcrptedaft)
	}
	delastcrptedaftstr :=string(delastcrptedaft)
	mylog.Info("command %s: decrypt succ===999: %s", "cmdName", delastcrptedaftstr)
	time.Sleep(time.Second * 3)
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
	//1105testing,good2!
	/*
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

	/* 2019.1114 skiping
	fromPubkeyStr := "1b3d954faa58c0cf7911596f056354136b3bbef996909167fd27386639cadbf4"
	toPubkeyhashStr := "b1348662bf564fe79e6fcaa33855feccd4adf98d"
	amount:=3.8
	prikeyStr:= "858322a0f4f4edc45c58ddd5b6420eb5f2e54273a5c81d366102a5f97fe56c14"
	nonce :=3
	wdctranssign.GWDCTransHandle.ClientToTransferAccount(fromPubkeyStr,toPubkeyhashStr,float64(amount),prikeyStr,int64(nonce))
	time.Sleep(time.Second * 4)

		*/
	//end 1107doing end
	testWdcRPCClient :=new(wdctranssign.WdcRpcClient)

	//1106 add:
	blockHeight,err := testWdcRPCClient.GetBlockHeight()
	if err != nil {
		mylog.Error("RPC GetBlockHeight()  failure! error is :%v,", err)
	}else{
		mylog.Info("RPC GetBlockHeight()  succ! get cur blockHeight is: %d", blockHeight)

	}

	/*
	//worker.NewWDCWorker( time.Minute * 5)
	curWDCWorker := worker.NewWDCWorker( time.Second * 500)
	go curWDCWorker.Start(nil)
	//1106 Mode2 ending
		return
	*/

	router :=mux.NewRouter().StrictSlash(true)

	//router.HandleFunc("/createaddress/{cointype}", service.CreateAddress)
	router.HandleFunc("/createaddress", service.CreateAddress)
	//sgj 1017 adding:	签名服务器: 成生新的地址(批量)
	router.HandleFunc("/remote/getnewaddress", service.RemoteSignCreateAddress)

	//sgj 1114add for guiji:

	router.HandleFunc("/remote/monitorwalletwdc", service.RemoteMonitorWalletAddress)

	strHost:=fmt.Sprintf(":%d",gbConf.WebPort)
	mylog.Info("strHost is :%s", strHost)

	//1025add

	//1028add:normal params to get all queue;
	//go service.WithdrawsTransTotal(0,10,"WDC")

	err =http.ListenAndServe(strHost, router)
	if nil!=err {
		mylog.Error("%+v",err)
		os.Exit(0)
	}
	//select {}

}
