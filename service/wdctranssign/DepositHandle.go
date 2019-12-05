package wdctranssign

import (
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

)

//一个默认的合作商Key
var GCurGetKeyStr =[]byte("1234567812345678")

func WithdrawsDepositGatherWDC(offset, limit uint, cointype string)(addressCount int,bsucc bool){
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
	opercount,bRet := GWDCTransHandle.DepositesAddrGatter(&reqDepositInfo)

	log.Info("DepositesAddrGatter,handle succ!,reqDepositInfo is :%v,return is :%v", reqDepositInfo,bRet)
	time.Sleep(time.Second * 2)
	return opercount,bRet

}

//sgj 1113,,小于此值，不进行归集处理
var minWDCLimit = 0.05

//大账户最大额度限制
var threshold = 150
func (self *WDCTransHandle) QueryWDCDepositesAddr(reqQueryInfo *proto.DepositeAddresssReq) (Address []string, succflag bool) {
	var signData string
	curQueryInfo := proto.DepositeAddresssReq{}
	curQueryInfo = *reqQueryInfo

	resDepositQuerySign := proto.Response{}
	transInfo := proto.DepositeAddresssResp{}
	resDepositQuerySign.Data = &transInfo
	UrlVerify := config.GbConf.SettleApiDepositQuery

	log.Info("QueryWDCDepositesAddr.UrlVerify is:%s,reqInfo is:%v", UrlVerify, curQueryInfo)

	getAddress := make([]string,0)

	//WDCGatterConfigUrl
	reqBody, err := json.Marshal(&curQueryInfo)
	if nil != err {
		log.Error("when QueryWDCDepositesAddr,Marshal to json error:%s", err.Error())
		return getAddress,false
	}

	if signData, err = auth.KSign(reqBody, config.GbConf.SettleAccessKey.AccessPrivKey); err != nil {
		log.Error("In QueryWDCDepositesAddr(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return getAddress,false
	}
	//step 2
	log.Info("QueryWDCDepositesAddr,auth.KSign succ!,signData is :%v", signData)

	ht := CHttpClientEx{}
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//req.Header.Set("abit-actionsign", signData)
	//ht.HeaderSet(proto.HActionSign, signData)
	//1112 update for abit
	ht.HeaderSet(proto.HActionAbitSign, signData)
	//1204 tmp doing:
	//ht.HeaderSet(proto.HActionSign, signData)


	log.Info("QueryWDCDepositesAddr.transInfo url=%s,cur reqdata = %v", UrlVerify, curQueryInfo)
	strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 9000, &curQueryInfo, &resDepositQuerySign)
	if nil != err {
		log.Error("QueryWDCDepositesAddr,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return getAddress,false
	} else {
		if resDepositQuerySign.Msg== "Success" {

			log.Info("QueryWDCDepositesAddr good!,get Msg' len(strRes) is:%d,resDepositQuerySign is:%v", len(strRes), resDepositQuerySign)
			log.Info("QueryWDCDepositesAddr info is:%v", transInfo)

			getAddress = transInfo.Addresss
			return getAddress,true
		}else{
			//errcode is:%v", resUpdateSign.Code)
			log.Error("QueryWDCDepositesAddr res err! err resUpdateSign is:%v", resDepositQuerySign)
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
func (self *WDCTransHandle) QueryDepositGroupConfig(group string) (getDepositConfig []DepositCoinConfig, succflag bool) {
	var signData string
	var getDepositCoinConfig = make([]DepositCoinConfig,0)

	curQueryInfo := proto.WithDrawConfigReq{}
	curQueryInfo.Nonce = time.Now().Unix()

	resDepositQuerySign := proto.Response{}
	transInfo := proto.WithdrawConfigResp{}
	resDepositQuerySign.Data = &transInfo

	UrlVerify := config.GbConf.WDCGatterConfigUrl

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
	ht.HeaderSet(proto.HActionAbitSign, signData)


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
func (self *WDCTransHandle) DepositesAddrGatter(reqQueryInfo *proto.DepositeAddresssReq) (opercount int,is bool) {


	var threshold float64 = 22;
	//fix 初始化count
	self.GatherAddrCount = 0
	//1205 fix add offset:
	var TotalAddressList = make([]string,0)
	reqQueryInfo.Offset = 0
	//循环取出所用充值地址：
	for {
		//end 1205
		curAddressList, bsucc := self.QueryWDCDepositesAddr(reqQueryInfo)
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
		log.Info("QueryWDCDepositesAddr good! get len is :%d,curAddressList is:%v", len(curAddressList), curAddressList)
		for _, getAddr := range curAddressList {
			TotalAddressList = append(TotalAddressList, getAddr)
		}
		reqQueryInfo.Offset += len(curAddressList)
	}
	//end 1205.1
	//var threshold;
	//从settlecenter测，获取配置的大账户归集限额
	configs,bsucc := self.QueryDepositGroupConfig("WDC")

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
	log.Info("exec QueryDepositGroupConfig(),get WDC GroupConfig for threshold succ ,threshold values is %.8f\n",threshold)

	//threshold = parseFloat(configs[0].threshold);
	//threshold = 250
	//WDCGatterToAddress
	curGatterToAddress := config.GbConf.WDCGatterToAddress
	curaddrrec,err := GWdcDataStore.GetWDCAddressRec(curGatterToAddress)
	if err != nil{
		log.Error("GetWDCAddressRec(),get rows for fromaddress record failed!,WDCTransProc() exec to return.curGatteraddress =%s",curGatterToAddress)
		return 0,false
	}

	//获取大账户余额	curGatterToAddress,
	addrtotalAmount,err,errmsg :=self.WdcRpcClient.SendBalancePostFormNode(curaddrrec.PubKeyHash)
	addrtotalAmount = addrtotalAmount /100000000
	if err !=nil{
		log.Error("DepositesAddrGatter.SendBalance() fail, get err=%v,errinfo :%s,cur fromAddress is: %v,getPubKeyHash is:%s", err,errmsg,curGatterToAddress,curaddrrec.PubKeyHash)
	}
	var limit = threshold - addrtotalAmount;

	//需要归集的最大额度数量
	self.GatherLimit = limit
	log.Info("curGatterToAddress(%s),GetBalance is %.8f\n",curGatterToAddress, addrtotalAmount)


	//wdcbalance :=244
	if (addrtotalAmount >= threshold) {
		log.Info("sufficient wdc balance cur value is %.8f, wdc threshold is :%f",addrtotalAmount,threshold);
		//prvkeyGatter.setPemMemory('');
		//reject('sufficient usdt');
		return 0,false;
	}

	log.Info("WithdrawsDeposites res succ! to gather limit is:%f,get TotalAddressList is:%v", limit,TotalAddressList)
	//sgj 1114checking
	//return

	for ino, curAddrItem := range TotalAddressList {

		_,gettxid := self.WDCGatherTransProc(int64(ino),curAddrItem,curGatterToAddress)
		log.Info("cur WDCGatherTransProc() finished, curAddrItem is %s, curGatterToAddress is:%s,gettxid is:%s,the rest wdc GatherLimit is :%f",curAddrItem,curGatterToAddress,gettxid,self.GatherLimit);
		//var hash = await _omnisend(addrList[i], balance, fee);
		if (self.GatherLimit <= 0 ){
			break;
		}

	}
	return self.GatherAddrCount,true


}

//归集转账过程
var curWDCFee = 0.002
//(errinfo transproto.ErrorInfo,ival uint, retval []interface{}){
func(self *WDCTransHandle) WDCGatherTransProc(iseno int64,fromaddress string, toGatherAddr string) (opsuccflag bool, tid string) {
	var ret bool = false
	log.Info("WDCGather transfer %s => %s ,coin_type %s\n", fromaddress,toGatherAddr, "WDC")
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

	curaddrrec,err := GWdcDataStore.GetWDCAddressRec(getfromAddress)
	if err != nil{
		log.Error("GetWDCAddressRec(),get rows for fromaddress record failed!,WDCTransProc() exec to return.curaddress =%s",getfromAddress)
		return false,""
	}
	getAddressPub := curaddrrec.PubKey
	getencrptedAddressPriv := curaddrrec.PrivKey

	//sgj 1115 add for encrypto
	// 对 params 进行 base64 解码
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
	//sgj 1115 end add
	//获取账户余额	getfromAddress,
	fromMount,err,errmsg :=self.WdcRpcClient.SendBalancePostFormNode(curaddrrec.PubKeyHash)
	if err !=nil{
		log.Error("WDCTransProc.SendBalance() fail, get err=%v,errinfo :%s,cur fromAddress is: %v,getPubKeyHash is:%s", err,errmsg,getfromAddress,curaddrrec.PubKeyHash)
	}
	fromMount = fromMount /100000000
	log.Info("fromAddress(%s),GetBalance.Aft is %.8f. to gather to bigaccount!\n",getfromAddress, fromMount)
	//curAmount,_:= cursettle.Vol.Float64()
	curGatherAmount := fromMount - curWDCFee

	//1114 add,满足归集最大上限为止
	if curGatherAmount > self.GatherLimit {
		curGatherAmount = self.GatherLimit
	}

	var totalNeeds float64 = (minWDCLimit + curWDCFee)	// * 100000000
	/*
	fromMount = fromMount * 100000000
	*/
	//余额不够最小归集限额,停止此比交易
	if  totalNeeds > fromMount {
		log.Info("balance is too low ignore. WDC Trans is insufficient!,cur balance is %.8f,cursettle need is:%.8f\n", fromMount,minWDCLimit)
		/*
		reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_FAILED
		reqUpdateInfo.Withdraws[0].Error = "当前余额不够"
		if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			log.Error("WDCGatherTransProc.WDCTransProc() fail, exec compare balance failed!,curid is:%d,curbalance is:%.8f,totalNeeds amount is: %.8f,cur trans break!", iseno,fromMount,totalNeeds)
		}
		*/
		return true, ""
	}
	log.Info("cur WDC Trans amount info: cur balance is %f,cursettle need is:%.8f, curFee is:%.8f\n", fromMount,totalNeeds,0.02)

	//获取账户Nonce,var getNonce int64
	time.Sleep(time.Second * 4)
	curNonce,err,errmsg :=self.WdcRpcClient.SendNonce(curaddrrec.PubKeyHash)
	if err !=nil{
		log.Error("WDCTransProc.SendNonce() fail, get err=%v,errinfo :%s,cur fromAddress is: %v,getPubKeyHash is:%s", err,errmsg,getfromAddress,curaddrrec.PubKeyHash)
	}
	getNonce := int64(curNonce)

	getToPubHashStr,err :=self.WdcRpcClient.GetAddressPubHash(toGatherAddr)
	if err !=nil{
		log.Error("WDCTransProc.GetAddressPubHash() fail, get err=%v,cur toGatherAddr = %v,cur trans break!", err,toGatherAddr)

		//1107add,转账参数不规范，通知交易系统失败：
		//reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_FAILED
		return false,""
	}else{
		log.Info("WDCTransProc.GetAddressPubHash() succ, get toGatherAddr is:%s,getToPubHashStr is:%s", toGatherAddr,getToPubHashStr)
	}

	txid,txHexStr, err,errmsg := self.ClientToTransferAccount(getAddressPub,getToPubHashStr,curGatherAmount,getAddressPriv,int64(getNonce))
	if err !=nil || errmsg !="" {
		log.Error("WDCTransProc.ClientToTransferAccount() fail, get err=%v,cur errmsg = %v,cur trans break!", err,errmsg)
		return false,""

	}else{
		log.Info("WDCTransProc.ClientToTransferAccount() succ, gettxid is:%s, txHexStr=%s,cur errmsg = %v", txid,txHexStr,errmsg)

	}

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
	GWdcDataStore.SaveTranRecord("WDCGather",getfromAddress,toGatherAddr,iseno,txid,curGatherAmount,"curstatusing",errcode,errmsg,"")

	self.GatherLimit -= curGatherAmount
	self.GatherAddrCount +=1
	return ret, txid
}

