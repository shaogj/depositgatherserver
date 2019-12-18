package ktctranssign

import (
	//"2019NNZXProj10/TiggerServerPlatform/service/ktctranssign"
	//"2019NNZXProj10/TiggerServerPlatform/service/ktctranssign"
	"2019NNZXProj10/depositgatherserver/service/wdctranssign"
	"encoding/base64"
	"runtime"
	"time"

	"github.com/mkideal/log"
	//"strconv"

	"2019NNZXProj10/depositgatherserver/KeyStore"
	"2019NNZXProj10/depositgatherserver/config" // "strings"
	"2019NNZXProj10/depositgatherserver/proto"
	"encoding/json"
	"errors"
	. "shaogj/utils"
	"2019NNZXProj10/depositgatherserver/cryptoutil"
	//sgj1217 add
	"2019NNZXProj10/depositgatherserver/service/ktctranssign/ktcrpc"
	"github.com/btcsuite/btcd/btcjson"



)

//一个默认的合作商Key
var GCurGetKeyStr =[]byte("1234567812345678")

func WithdrawsDepositGatherKTC(offset, limit uint, cointype string)(addressCount int,bsucc bool){
	var reqDepositInfo proto.DepositeAddresssReq

	ht := CHttpClientEx{}
	//sgj add
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//reqInfo.MaxVol = 0
	reqDepositInfo.Limit = int(limit)
	//reqDepositInfo.Status = transproto.SETTLE_STATUS_PASSED
	reqDepositInfo.CoinCode = cointype
	reqDepositInfo.Nonce = time.Now().Unix()
	reqDepositInfo.Offset = int(offset)
	opercount,bRet := GKTCDepositHandle.DepositesAddrGatter(&reqDepositInfo)

	log.Info("DepositesAddrGatter,handle succ!,reqDepositInfo is :%v,return is :%v", reqDepositInfo,bRet)
	time.Sleep(time.Second * 2)
	return opercount,bRet

}

//sgj 1217adding
//var GWDCTransHandle =  WDCTransHandle{}
var GKTCDepositHandle =  KTCDepositHandle{}
//KTCDepositHandle

var GKtcDataStore =wdctranssign.WdcDataStore{}

type KTCDepositHandle struct {
	CoinType string  //币种类型
	WDCTransUrl string
	HtClient CHttpClientEx
	//WdcRpcClient	*WdcRpcClient
	KTCSignHandle

	//sgj 1113add from DepositGather
	GatherLimit		float64
	//sgj 1114adding,总归集的地址数量
	GatherAddrCount	int

}



//大账户最大额度限制
var threshold = 150
func (self *KTCDepositHandle) QueryKTCDepositesAddr(reqQueryInfo *proto.DepositeAddresssReq) (Address []string, succflag bool) {
	var signData string
	curQueryInfo := proto.DepositeAddresssReq{}
	curQueryInfo = *reqQueryInfo

	resDepositQuerySign := proto.Response{}
	transInfo := proto.DepositeAddresssResp{}
	resDepositQuerySign.Data = &transInfo
	UrlVerify := config.GbConf.SettleApiDepositQuery

	log.Info("QueryKTCDepositesAddr.UrlVerify is:%s,reqInfo is:%v", UrlVerify, curQueryInfo)

	getAddress := make([]string,0)

	//KTCGatterConfigUrl
	reqBody, err := json.Marshal(&curQueryInfo)
	if nil != err {
		log.Error("when QueryKTCDepositesAddr,Marshal to json error:%s", err.Error())
		return getAddress,false
	}

	if signData, err = auth.KSign(reqBody, config.GbConf.SettleAccessKey.AccessPrivKey); err != nil {
		log.Error("In QueryKTCDepositesAddr(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return getAddress,false
	}
	//step 2
	log.Info("QueryKTCDepositesAddr,auth.KSign succ!,signData is :%v", signData)

	ht := CHttpClientEx{}
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//req.Header.Set("abit-actionsign", signData)
	//ht.HeaderSet(proto.HActionSign, signData)
	//1112 update for abit
	//1217befing
	//ht.HeaderSet(proto.HActionAbitSign, signData)
	ht.HeaderSet(proto.HActionSign, signData)
	//1204 tmp doing:
	//ht.HeaderSet(proto.HActionSign, signData)


	log.Info("QueryKTCDepositesAddr.transInfo url=%s,cur reqdata = %v", UrlVerify, curQueryInfo)
	strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 9000, &curQueryInfo, &resDepositQuerySign)
	if nil != err {
		log.Error("QueryKTCDepositesAddr,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return getAddress,false
	} else {
		if resDepositQuerySign.Msg== "Success" {

			log.Info("QueryKTCDepositesAddr good!,get Msg' len(strRes) is:%d,resDepositQuerySign is:%v", len(strRes), resDepositQuerySign)
			log.Info("QueryKTCDepositesAddr info is:%v", transInfo)

			getAddress = transInfo.Addresss
			return getAddress,true
		}else{
			//errcode is:%v", resUpdateSign.Code)
			log.Error("QueryKTCDepositesAddr res err! err resUpdateSign is:%v", resDepositQuerySign)
			return getAddress,false
		}
	}
}

//sgj 1113 adding
type DepositCoinConfig struct{
	CoinName string
	TokenAddress string
	Threshold float64

}
var GDepositCoinConfig []DepositCoinConfig

//11.14 add：查询归集的充值配置接口：
func (self *KTCDepositHandle) QueryDepositGroupConfig(group string) (getDepositConfig []DepositCoinConfig, succflag bool) {
	var signData string
	var getDepositCoinConfig = make([]DepositCoinConfig,0)

	curQueryInfo := proto.WithDrawConfigReq{}
	curQueryInfo.Nonce = time.Now().Unix()

	resDepositQuerySign := proto.Response{}
	transInfo := proto.WithdrawConfigResp{}
	resDepositQuerySign.Data = &transInfo

	UrlVerify := config.GbConf.KTCGatterConfigUrl

	log.Info("QueryDepositGroupConfig.UrlVerify is:%s,reqInfo is:%v", UrlVerify, curQueryInfo)
	if UrlVerify == "" {
		UrlVerify = "https://devapi.ggex.com/v1/settle/getconfig?action=query"
	}
	reqBody, err := json.Marshal(&curQueryInfo)
	if nil != err {
		log.Error("when QueryDepositGroupConfig,Marshal to json error:%s", err.Error())
		return getDepositCoinConfig,false
}

	if signData, err = auth.KSign(reqBody, config.GbConf.SettleAccessKey.AccessPrivKey); err != nil {
		log.Error("In QueryDepositGroupConfig(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return getDepositCoinConfig,false
	}
	//step 2
	log.Info("QueryDepositGroupConfig,auth.KSign succ!,signData is :%v", signData)

	ht := CHttpClientEx{}
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//req.Header.Set("abit-actionsign", signData)
	//ht.HeaderSet(proto.HActionSign, signData)
	//1112 update for abit
	//ht.HeaderSet(proto.HActionAbitSign, signData)
	//1217 update for abit
	ht.HeaderSet(proto.HActionTMexSign, signData)


	log.Info("QueryDepositGroupConfig url=%s,cur reqdata = %v", UrlVerify, curQueryInfo)
	//strRes
	_, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 9000, &curQueryInfo, &resDepositQuerySign)
	if nil != err {
		log.Error("QueryDepositGroupConfig,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		time.Sleep(time.Second * 4)
		return getDepositCoinConfig,false
	} else {
		if resDepositQuerySign.Msg== "Success" {

			log.Info("QueryDepositGroupConfig good!,get strRes is:%v,resDepositQuerySign after is:%v", "strRes", transInfo.Configs)

			curCoinDetailConfig :=transInfo.Configs

			for i := 0; i < len(curCoinDetailConfig); i++ {
				var curDepositCoinConfigItem DepositCoinConfig
				var curCoinDetailConfig proto.CoinDetailConfig = *curCoinDetailConfig[i]
				log.Info("QueryDepositGroupConfig good noid:%d, item is:%v", i,curCoinDetailConfig)

				if group == curCoinDetailConfig.CoinGroup{
					volWalletMaxVol,_ := curCoinDetailConfig.WalletMaxVol.Float64()
					curDepositCoinConfigItem.Threshold = volWalletMaxVol
					curDepositCoinConfigItem.CoinName = curCoinDetailConfig.CoinName
					curDepositCoinConfigItem.TokenAddress = curCoinDetailConfig.TokenAddress
					log.Info("QueryDepositGroupConfig good noid===777:%d, item is:%v,,group is:%s,Threshold vol is:%f", i,curCoinDetailConfig,group,volWalletMaxVol)
					//sgj 1114 add for check:WalletMaxVol
					log.Info("cur group is:%s,WalletMaxVol is:%v,GetMaxVolAft is:%f,curDepositCoinConfigItem is:%v",group,curCoinDetailConfig.WalletMaxVol,volWalletMaxVol,curDepositCoinConfigItem)

					getDepositCoinConfig = append(getDepositCoinConfig,curDepositCoinConfigItem)
					log.Info("get total cur getDepositCoinConfig is :%v",getDepositCoinConfig)
				}
			}
			log.Info("cur group is:%s,getDepositCoinConfig[0] is:%v",group,getDepositCoinConfig[0])

			return getDepositCoinConfig,true
		}else{
			time.Sleep(time.Second * 4)
			log.Error("QueryDepositGroupConfig res err! err resUpdateSign is:%v", resDepositQuerySign)
			return getDepositCoinConfig,false
		}
	}
}

//开始资金归集的流程
func (self *KTCDepositHandle) DepositesAddrGatter(reqQueryInfo *proto.DepositeAddresssReq) (opercount int,is bool) {


	var threshold float64 = 22;
	//fix 初始化count
	self.GatherAddrCount = 0
	//1205 fix add offset:
	var TotalAddressList = make([]string,0)
	reqQueryInfo.Offset = 0
	//循环取出所用充值地址：
	for {
		//end 1205
		curAddressList, bsucc := self.QueryKTCDepositesAddr(reqQueryInfo)
		if bsucc == false {
			log.Error("WithdrawsDeposites err! cur reqQueryInfo is:%v", reqQueryInfo)
			//1205:
			break
		}
		if len(curAddressList) == 0 {
			log.Info("WithdrawsDeposites finished! cur reqQueryInfo is:%v,get addrlist is 0", reqQueryInfo)
			//1205:
			break
		}
		log.Info("QueryKTCDepositesAddr good! get len is :%d,curAddressList is:%v", len(curAddressList), curAddressList)
		for _, getAddr := range curAddressList {
			TotalAddressList = append(TotalAddressList, getAddr)
		}
		reqQueryInfo.Offset += len(curAddressList)
	}
	log.Info("WithdrawsDeposites Total finished! get TotalAddressList len is :%d", len(TotalAddressList))

	//end 1205.1
	//var threshold;
	//从settlecenter测，获取配置的大账户归集限额
	configs,bsucc := self.QueryDepositGroupConfig("KTC")

	if bsucc != true {
		log.Error("QueryDepositGroupConfig res err! get configs is:empty")
		return 0,false
	}
	log.Info("QueryDepositGroupConfig res good! get configs is:%v", configs[0])

	if len(configs) > 0{
		threshold =configs[0].Threshold
	}else{
		threshold = 444

	}
	log.Info("exec QueryDepositGroupConfig(),get KTC GroupConfig for threshold succ ,threshold values is %.8f\n",threshold)

	//threshold = parseFloat(configs[0].threshold);
	//threshold = 250
	//KTCGatterToAddress
	curGatterToAddress := config.GbConf.KTCGatterToAddress
	//curaddrrec,err := GWdcDataStore.GetWDCAddressRec(curGatterToAddress)
	//12.17 adding

	curaddrrec,err := GKtcDataStore.GetKTCAddressRec(GKtcDataStore.OrmEngine,curGatterToAddress)
	// GetKTCAddressRec
	if err != nil{
		log.Error("GetKTCAddressRec(),get rows for fromaddress record failed!,KTCTransProc() exec to return.curGatteraddress =%s",curGatterToAddress)
		return 0,false
	}
	//12.17--try 获取 privkey
	log.Info("exec GetKTCAddressRec,curGatterToAddress is :%s,get curaddrrec info is: %v \n", curGatterToAddress,curaddrrec)

	//获取大账户余额	curGatterToAddress,
	//1217,get KTC banlance:
	unspentLimitAddrTotals :=make([]string,0,3)
	//fromAddr := "39QXajNbM7aqurkav6DF6vyupY1cn48a8i"
	unspentLimitAddrTotals = append(unspentLimitAddrTotals,curGatterToAddress)

	getutxoinfo,utxonum,err := ktcrpc.KTCRPCClient.GetRPCTxUnSpentLimit(1,0,unspentLimitAddrTotals)	//"1Eq8xXAea47WPY5t8zUEYDKgcWB7cptZWB")
	if err != nil {
		log.Error("curGatterToAddress: %s ,exex GetTxUnSpentLimit() failue! err is: %v \n", curGatterToAddress, err)
		//	return nil, status, err
	}
	//getResp.Data.(*[]proto.WdcTxBlock)
	log.Info("exec GetTxUnSpentLimit(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! unxonum is :%d\n", getutxoinfo, utxonum)
	var addrtotalAmount float64
	//1217,getall balance
	for _,curitem := range getutxoinfo{

		addrtotalAmount += curitem.Amount
	}

	//addrtotalAmount = addrtotalAmount /100000000
	var limit = threshold - addrtotalAmount;

	//需要归集的最大额度数量
	self.GatherLimit = limit
	log.Info("curGatterToAddress(%s),GetBalance is %.8f,GatherLimit is :%f\n",curGatterToAddress, addrtotalAmount,self.GatherLimit)


	//KTCbalance :=244
	if (addrtotalAmount >= threshold) {
		log.Info("sufficient KTC balance cur value is %.8f, KTC threshold is :%f",addrtotalAmount,threshold);
		return 0,false;
	}

	log.Info("WithdrawsDeposites res succ! to gather limit is:%f,get len(TotalAddressList) is:%d,TotalAddressList is:%v", limit,len(TotalAddressList),TotalAddressList)
	//sgj 1114checking
	//return

	for ino, curAddrItem := range TotalAddressList {

		//12.17doing
		_,gettxid := self.KtCGatherTransProc(int64(ino),curAddrItem,curGatterToAddress)
		log.Info("cur KTCGatherTransProc() finished, curAddrItem is %s, curGatterToAddress is:%s,gettxid is:%s,the rest KTC GatherLimit is :%f",curAddrItem,curGatterToAddress,gettxid,self.GatherLimit);
		//var hash = await _omnisend(addrList[i], balance, fee);
		if (self.GatherLimit <= 0 ){
			break;
		}

	}
	return self.GatherAddrCount,true


}

//归集转账过程
var curKTCFee = 0.002
//(errinfo transproto.ErrorInfo,ival uint, retval []interface{}){
func(self *KTCDepositHandle) KtCGatherTransProc(iseno int64,fromaddress string, toGatherAddr string) (opsuccflag bool, tid string) {
	var ret bool = false
	log.Info("KTCGather transfer %s => %s ,coin_type %s\n", fromaddress,toGatherAddr, "KTC")
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 1<<16)
			buf = buf[:runtime.Stack(buf, true)]
			var err error
			switch x := e.(type) {
			case error:
				err = x
			case string:
				err = errors.New(x)
			}
			log.Error("==== STACK TRACE BEGIN ====\npanic: %v\n%s\n===== STACK TRACE END =====", err, string(buf))
		}
	}()

	var getfromAddress = fromaddress
	reqUpdateInfo := proto.WithdrawsUpdateReq{}
	reqUpdateInfo.Withdraws = make([]proto.Settle,1)

	//1118 update :abit-online big address
	if toGatherAddr == ""{
		toGatherAddr = "1KVcQTbsMuU5jpZSdBRXiKcbbawrGBo9h7"
	}
	//ggex.dev.test:1HFCUeNHcL6Drf4TPwBLG6RgYVe9o41BVj
	curaddrrec,err := GKtcDataStore.GetKTCAddressRec(GKtcDataStore.OrmEngine,getfromAddress)
	//curaddrrec,err = GWdcDataStore.GetKTCAddressRec(GKtcDataStore.OrmEngine,getfromAddress)

	//curaddrrec,err := GWdcDataStore.GetWDCAddressRec(getfromAddress)
	if err != nil{
		log.Error("GetKTCAddressRec(),get rows for fromaddress record failed!,KTCTransProc() exec to return.curaddress =%s",getfromAddress)
		return false,""
	}
	//getAddressPub := curaddrrec.PubKey
	getencrptedAddressPriv := curaddrrec.PrivKey

	//sgj 1115 add for encrypto
	// 对 params 进行 base64 解码
	log.Info("GetKTCAddressRec(),get getencrptedAddressPriv =====0033---is :%s", getencrptedAddressPriv)
	dencrptedEncodeStr, err := base64.StdEncoding.DecodeString(string(getencrptedAddressPriv))
	if err != nil {
		log.Error("DecodeString text is:%s,err is----AAA:%v",dencrptedEncodeStr,err)
		//return nil, err
	}
	//fmt.Println("Decrypt get decrpteddecodeStr len is:%d,val is====44:%s,org encrpted len is:",len(dencrptedEncodeStr),dencrptedEncodeStr,len(encrpted))
	delastcrptedaft, err := cryptoutil.AESCBCDecrypt(GCurGetKeyStr, nil, []byte(dencrptedEncodeStr))
	delastcrptedaftstr :=string(delastcrptedaft)
	if err != nil {
		log.Error("delastcrptedaft is: %s: decrypt error===888: %v", delastcrptedaftstr, err)
	}
	log.Info("command %s: decrypt succ===999: %s", "AESCBCDecrypt", delastcrptedaftstr)
	getAddressPriv := delastcrptedaftstr
	log.Info("after GetKTCAddressRec(),cur getfromAddress is:%s,get getAddressPriv is:%s,\n",getfromAddress, getAddressPriv)
	//sgj 1115 end add
	//获取账户余额	getfromAddress,

	/*12.17 move to GKTCSignHandle.PaySignTransProc
	//12.17doing
	//比特币的utxo是每个txid进行归集的；
	//curGatherAmount = utxo [0].amount
	curGatherAmount = 0.003
	//1114 add,满足归集最大上限为止
	if curGatherAmount > self.GatherLimit {
		curGatherAmount = self.GatherLimit
	}

	*/
	var totalNeeds float64 = (minKTCLimit + curKTCFee)	// * 100000000
	/*
	fromMount = fromMount * 100000000
	*/
	//余额不够最小归集限额,停止此比交易
	//log.Info("cur KTC Trans amount info: cur balance is %f,cursettle need is:%.8f, curFee is:%.8f\n", fromMount,totalNeeds,0.02)
	log.Info("cur KTC Trans amount info: cursettle need is:%.8f, curFee is:%.8f\n",totalNeeds,curKTCFee)

	//time.Sleep(time.Second * 4)


	//12.17 for KTC proc:
	signInfoRes, curGatherAmount,status, err := GKTCSignHandle.PaySignTransProc(fromaddress, getAddressPriv,0,toGatherAddr,self.GatherLimit)

	log.Info("KTCTransProc.SendTransactionPostForm() succ!,toamount1 is :%f,status is:%v,getcur resdata is:%v",curGatherAmount,status,signInfoRes)

	toamount1 := curGatherAmount
	//1217adding
	//end sgj 1121ing
	if signInfoRes == nil || err != nil{
		//0502	保存签名交易错误信息to DB：status，err
		//GeneJsonResultFin(w,r,ressigndata,status,desc)//"errinfo4"
		log.Error("request Pay_SignTransaction() failure ! CoinCode is %d,amount is:%f,get status info is %d:err is :%v", "KTC",toamount1,status, err)
	} else{
		if "KTC" == config.CoinKTC{

			txdecodeinfo := signInfoRes.(btcjson.SignRawTransactionResult)
			log.Info("request Pay_SignTransaction() succ ! fromaddress is %d,amount is:%f,get txdecodeinfo is :%s,status info is %d:err is :%v", fromaddress,toamount1,txdecodeinfo, status, err)
			//Sprintf
			//log.Info("44444444444444444444=====!!!!!!!")
			log.Info("to 开始广播tmpRes.Hex,cur fromaddress is:%s----！！！",fromaddress)
			//1121add
			//开始广播交易,及后续处理
			//sgj 1121doing PM,for WAtching
			//continue
			var curTxId string
			getTxid, err := ktcrpc.KTCRPCClient.SendTransaction(txdecodeinfo.Hex)
			if getTxid != "" {
				curTxId = getTxid
			}
			if err !=nil{
				log.Error("SendTransaction() failure ! fromaddress is:%s,curGatherAmount is :%f,get gettxid is :%s,err is :%v", fromaddress,curGatherAmount,curTxId, err)
			}else{
				log.Info("handleSendPostMessage() succ ! fromaddress is :%s,cur curGatherAmount is :%f,get gettxid is :%s", fromaddress,curGatherAmount,curTxId)

			}
			//if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			/*sgj begin1121
			//1104，可把转账结构写入数据库
			GWdcDataStore.SaveTranRecord(cursettle.CoinCode,getfromAddress,getToPubHashStr,cursettle.SettleId,txid,curAmount,"curstatusing",errcode,errmsg,"")
			end sgj 1121 add*/
			if curTxId > "" {

				GKtcDataStore.SaveTranRecord("KTC",fromaddress,toGatherAddr,0,curTxId,curGatherAmount,"KTCGatherSucc!",0,"KTCGatherSucc!","RawHexstrinfo")
			}
			//curDataStore.SaveTranRecord(curcoinType,fromAddr,toAddr,cur.SettleId,curTxId,curAmount,status,errcode,descTotal,desc)

		}

	}
	/*1217 forKTCsec
	//开始广播交易：
	resdata,err,errcode,errmsg:= self.WdcRpcClient.SendTransactionPostForm(txHexStr)
	if errcode == proto.ErrorNodeRPCSuccess.Code{
		log.Info("WDCTransProc.SendTransactionPostForm() succ!,txid is:%s,getcur resdata is:%v",txid,resdata)
		//转账后通知交易系统,状态值5
		//curStatus = proto.SETTLE_STATUS_PENDING
	}else{
		log.Error("WDCTransProc.SendTransactionPostForm() fail!,errcode is:%d,res errmsg is:%v",errcode,errmsg)
	}

	//1104，可把转账结构写入数据库
	//err :
	//1202 upgrade,when save to DB: use toGatherAddr replace getToPubHashStr
	GKTCDataStore.SaveTranRecord("KTCGather",getfromAddress,toGatherAddr,iseno,txid,curGatherAmount,"curstatusing",errcode,errmsg,"")
	1217 forKTCsecend	*/


	self.GatherLimit -= curGatherAmount
	self.GatherAddrCount +=1
	txid := "23223"
	return ret, txid
}

