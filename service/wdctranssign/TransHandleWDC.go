package wdctranssign

import (
	"2019NNZXProj10/depositgatherserver/models"
	"time"

	//"2019NNZXProj10/depositgatherserver/wdctranssign"
	"fmt"
	"github.com/mkideal/log"
	//"strconv"

	"2019NNZXProj10/depositgatherserver/config" // "strings"
	"2019NNZXProj10/depositgatherserver/proto"
	"2019NNZXProj10/depositgatherserver/KeyStore"
	"github.com/go-xorm/xorm"
	"encoding/json"
	. "shaogj/utils"
	"errors"
)

/*
* 构造签名的交易事务
* 构造签名的孵化申请事务
* 构造签名的收益事务
* 构造签名的分享收益事务
* 构造签名的收取本金事务
* 构造签名的投票事务
* 构造签名的投票撤回事务
* 构造签名的抵押事务
* 构造签名的抵押撤回事务
* 构造签名的存证事务

*/
//SETTLESTATUS 定义交易订单状态码
type SETTLESTATUS int

const (
	SETTLE_STATUS_UN       SETTLESTATUS = iota
	SETTLE_STATUS_CREATED               // 1 申请成功(用户提交申请)
	SETTLE_STATUS_PASSED                // 2 审核通过(运营审核通过)
	SETTLE_STATUS_REJECTED              // 3 审核拒绝(运营审核拒绝)
	SETTLE_STATUS_SIGNED                // 4 开始转账，需要冻结用户资产
	SETTLE_STATUS_PENDING               // 5 转账结束，把txhash通知交易系统
	SETTLE_STATUS_SUCCESS               // 6 成功(转账成功)
	SETTLE_STATUS_FAILED                // 7 失败(转账失败)

)

//var WDCNodeUrl string = "http://192.168.1.211:19585/version"
var (
	MOrmEngine *xorm.Engine = &xorm.Engine{}
)

//curaddrrec,err := service.GWdcDataStore.GetWDCAddressRec
type WdcDataStore struct {
	OrmEngine *xorm.Engine
}

var GWdcDataStore =WdcDataStore{}

func (self *WdcDataStore) GetWDCAddressRec(curaddress string) (curaddrrec *models.WdcAccountKey, err error) {

	address_info := new(models.WdcAccountKey)
	bret,err := self.OrmEngine.Table("wdc_account_key").Where("address=?", curaddress).Get(address_info)
	if err != nil {
		log.Error("GetWDCAddressRec(), curaddress is :%s,error=%v",curaddress,err)
		return nil,err
	}
	//sgj 1016 add:
	if bret == false{
		log.Error("GetWDCAddressRec(),get rows failed,curaddress =%s,exist no row",curaddress)
		return address_info,errors.New("cur address exist no row in db")
	}else{
		log.Info("GetWDCAddressRec(),get rows succ,curaddress =%s,get row=%v",curaddress,address_info)
	}

	return address_info,nil

}

//sgj 1217adding:
func (self *WdcDataStore) GetKTCAddressRec(curKtcOrm *xorm.Engine,curaddress string) (curaddrrec *models.WdcAccountKey, err error) {

	address_info := new(models.WdcAccountKey)
	//bret,err := self.OrmEngine.Table("ktc_account_key").Where("address=?", curaddress).Get(address_info)
	bret,err := curKtcOrm.Table("ktc_account_key").Where("address=?", curaddress).Get(address_info)
	if err != nil {
		log.Error("GetKTCAddressRec(), curaddress is :%s,error=%v",curaddress,err)
		return nil,err
	}
	//sgj 1016 add:
	if bret == false{
		log.Error("GetKTCAddressRec(),get rows failed,curaddress =%s,exist no row",curaddress)
		return address_info,errors.New("cur address exist no row in db")
	}else{
		log.Info("GetKTCAddressRec(),get rows succ,curaddress =%s,get row=%v",curaddress,address_info)
	}

	return address_info,nil

}

//20190109 ,add for BTC address:
func (self *WdcDataStore) GetBTCAddressRec(curaddress string) (curaddrrec *models.WdcAccountKey, err error) {

	address_info := new(models.WdcAccountKey)
	bret,err := self.OrmEngine.Table("btc_account_key").Where("address=?", curaddress).Get(address_info)
	if err != nil {
		log.Error("GetBTCAddressRec(), curaddress is :%s,error=%v",curaddress,err)
		return nil,err
	}
	//sgj 1016 add:
	if bret == false{
		log.Error("GetBTCAddressRec(),get rows failed,curaddress =%s,exist no row",curaddress)
	}else{
		log.Info("GetBTCAddressRec(),get rows succ,curaddress =%s,get row=%v",curaddress,address_info)
	}

	return address_info,nil

}



//11.26 adding
func (self *WdcDataStore) UtxoWdcAccount(seriid int,addressid string,curNewPrivKey string) (bool,error) {

	//保存到数据库
	addrec := models.WdcAccountKey{
		PrivKey :curNewPrivKey,
		//GenerateTime: tm.Format("2006-01-02 03:04:05 PM"),
	}
	log.Debug("exec UtxoWdcAccount()==44444 ,seriid is:%d,addressid is:%s,curNewPrivKey is %v, rec info is :%v",seriid,addressid,curNewPrivKey,addrec)
	//rows, err := self.OrmEngine.Table("address_utxo").Where("txcurid=?", curNewPrivKey).Cols("vintxid_status").Update(rec)
	rows, err := self.OrmEngine.Table("wdc_account_key").Where("address=?", addressid).Cols("priv_key").Update(addrec)
	if err != nil {
		log.Error("exec err1!,UtxoWdcAccount()==5555{seriid=%s,curNewPrivKey=%s,err=%s}",seriid ,curNewPrivKey,err.Error())
		return false,err
	}
	if rows == 0 {
		log.Debug("exec UtxoWdcAccount()==5555 finished! exist no same address. address=%s",addressid)
		return true,nil
	}else{
		log.Debug("exec UtxoWdcAccount()==5555,cur update addressid' rec succ! seriid is:%d,curNewPrivKey is %v, rec info is :%v",seriid,curNewPrivKey,addrec)
	}
	return true,nil

}

//1103add
func (self* WdcDataStore)SaveTranRecord(coinType string,fromAddr,toAddress string,settleid int64,txhash string ,curamount float64,state string,errcode int,desc string,strRaw string) error  {

	log.Info("cur to exec SaveTranRecord(),Insert transrecord!,,curamount is:%v",curamount)
	tm := time.Unix(time.Now().Unix(), 0)

	acc := models.WdcTranRecord{
		Coincode:coinType,
		From:fromAddr,
		To:toAddress,
		Amount:curamount,
		Settleid:settleid,
		Txhash:txhash,
		Status:	state,
		Errcode:int64(errcode),
		TimeCreate:	tm.Format("2006-01-02 03:04:05 PM"),
		Desc:	desc,
		Raw:strRaw,
	}
	rows, err := self.OrmEngine.Table(models.TableGGEXTranRecord).Insert(acc)
	if err != nil {
		log.Error("SaveTranRecord(),Insert err!row = :%v,rows=%d,err is-:%v \n", acc,rows,err)
		return err
	}else{
		log.Info("SaveTranRecord(),Insert succ!,row = :%v,rows=%d,err is-:%v \n", acc,rows,err)

	}
	return nil

}


var GWDCTransHandle =  WDCTransHandle{}

type WDCTransHandle struct {
	CoinType string  //币种类型
	WDCTransUrl string
	HtClient CHttpClientEx
	WdcRpcClient	*WdcRpcClient

	//sgj 1113add from DepositGather
	GatherLimit		float64
	//sgj 1114adding,总归集的地址数量
	GatherAddrCount	int

}
func (self *WDCTransHandle) Init(wdcTransUrl string,nodeconf *config.WDCNodeConf){
	self.HtClient.Init()
	self.HtClient.HeaderSet("Content-Type", "application/json;charset=utf-8")
	if wdcTransUrl != ""{
		self.WDCTransUrl = wdcTransUrl
	}else{
		self.WDCTransUrl = TestWDCTransUrl
	}
	self.WdcRpcClient =NewWdcRpcClient(nodeconf)
	//self.WdcRpcClient.HtClient.Init()
	self.GatherAddrCount = 0
}

//1029
//存证交易事务ClientToTransferProve
var TestWDCTransUrl string = "http://192.168.1.190:8089/wallet/TxUtility"

//需获取nonce
func (self *WDCTransHandle) ClientToTransferAccount(fromPubkeyStr,toPubkeyHashStr string, amount float64,prikeyStr string,nonce int64 ) (gettxid string, gettxmesage string,err error,errmsg string) {

	UrlVerify:= fmt.Sprintf("%s/%s", self.WDCTransUrl, "ClientToTransferAccount")
	curToTransferAccount :=proto.ClientToTransferAccountParams{}
	curToTransferAccount.FromPubkeyStr = fromPubkeyStr
	curToTransferAccount.ToPubkeyHashStr = toPubkeyHashStr
	curToTransferAccount.Amount = amount
	curToTransferAccount.PrikeyStr = prikeyStr
	curToTransferAccount.Nonce = nonce
	var txid,txHexStr string

	resSDKAccount := proto.JavaSDKResponse{    }

	//resTxData := proto.JavaSDKResponse{    }
	//resSDKAccount.Data =&resTxData
	log.Info("transserver.ClientToTransferProve,watching---- cur input params is:%v",curToTransferAccount)

	strRes, statusCode, errorCode, err := self.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &curToTransferAccount, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  statuscode111=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return "","",err,resSDKAccount.Message
	}
	if statusCode == 200{
		//二次对node rpc' 内部封装的解析
		log.Info("transserver.ClientToTransferProve,get statusCode is :%s,res=%s,get cur txHexStr is:%d",statusCode, strRes,txHexStr)

		//sgj 1031 testing
		txsecHexStrBef :=resSDKAccount.Data.(string)
		getResp :=&proto.NodeRPCResponse{}
		err=json.Unmarshal([]byte(txsecHexStrBef),getResp)
		if  nil!=err  {
			log.Error("resp=%s,url=%s,err=%v",string(txsecHexStrBef),UrlVerify,err.Error())
			//return byResp,statusCode,6016,err
		}else{
			log.Info("ht.ClientToTransferAccount , getResp=%v ",getResp)
		}
		txHexStr = getResp.Message
		txid = getResp.Data.(string)
		log.Info(",get cur txid is:%v,txAftMessage is :%s",txid, txHexStr)

	}else if statusCode == proto.ErrorRequestWDCSDK.Code {
		errmsg =proto.ErrorRequestWDCSDK.Desc

	}
	log.Info("transserver.ClientToTransferProve,get statusCode is :%s,txid=%s,get cur txHexStr is:%d",statusCode, txid,txHexStr)
	return txid,txHexStr,err,errmsg
}



//(errinfo transproto.ErrorInfo,ival uint, retval []interface{}){
func(self *WDCTransHandle) WDCTransProc(cursettle proto.Settle, from string, accountname string) (opsuccflag bool, tid string) {
	var ret bool = false
	var curStatus proto.SETTLE_STATUS
	log.Info("transfer %s => %s mount %s coin_type %s\n", from, cursettle.ToAddress, cursettle.Vol, cursettle.CoinCode)

	//判断币种,如果不是WDC=>返回错误.
	if cursettle.CoinCode != "WDC" {
		ret = false
		return ret, ""
	}
	var getfromAddress string
	reqUpdateInfo := proto.WithdrawsUpdateReq{}
	reqUpdateInfo.Withdraws = make([]proto.Settle,1)
	if from == ""{
		getfromAddress =cursettle.FromAddress
	}else{
		getfromAddress =from
	}
	//1104 add :
	if getfromAddress == "" &&  from == ""{
		getfromAddress = "1HFCUeNHcL6Drf4TPwBLG6RgYVe9o41BVj"
	}
	reqUpdateInfo.Withdraws[0].FromAddress = getfromAddress
	reqUpdateInfo.Withdraws[0].SettleId = int64(cursettle.SettleId)
	reqUpdateInfo.Withdraws[0].CoinCode =cursettle.CoinCode
	reqUpdateInfo.Withdraws[0].SignStr = ""
	reqUpdateInfo.Nonce = time.Now().Unix()

	//1030doing
	//STEP_01__判断地址是否有效.
	//ValidateAddress(cursettle.ToAddress)
	//if err,return false

	//to do:从db里，取得address对应的pubkey,privkey

	//1030testing
	curaddrrec,err := GWdcDataStore.GetWDCAddressRec(getfromAddress)
	if err != nil{
		log.Error("GetWDCAddressRec(),get rows for fromaddress record failed!,WDCTransProc() exec to return.curaddress =%s",getfromAddress)
		return false,""
	}
	getAddressPub := curaddrrec.PubKey
	getAddressPriv := curaddrrec.PrivKey

	//获取账户余额	getfromAddress,
	fromMount,err,errmsg :=self.WdcRpcClient.SendBalancePostFormNode(curaddrrec.PubKeyHash)
	if err !=nil{
		log.Error("WDCTransProc.SendBalance() fail, get err=%v,errinfo :%s,cur fromAddress is: %v,getPubKeyHash is:%s", err,errmsg,getfromAddress,curaddrrec.PubKeyHash)
	}
	log.Info("fromAddress(%s),GetBalance is %.8f\n",getfromAddress, fromMount)
	curAmount,_:= cursettle.Vol.Float64()
	curFee,_:= cursettle.Fee.Float64()

	var totalNeeds float64 = curAmount * 100000000
	//totalNeeds = totalNeeds + curFee * 100000000
	/*
	fromMount = fromMount * 100000000
	*/
	//余额不够,通知交易系统失败
	if totalNeeds + curFee * 100000000 > fromMount {
		log.Error("cur WDC Trans is insufficient!,cur balance is %.8f,cursettle need is:%.8f\n", fromMount,totalNeeds)

		reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_FAILED
		reqUpdateInfo.Withdraws[0].Error = "当前余额不够"
		if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			log.Error("WDCTransProc.WDCTransProc() fail, exec compare balance failed!,cur cursettle is:%d,curbalance is:%.8f,totalNeeds amount is: %.8f,cur trans break!", cursettle.SettleId,fromMount,totalNeeds)
		}
		return false, ""
	}
	log.Info("cur WDC Trans amount info: cur balance is %f,cursettle need is:%.8f, curFee is:%.8f\n", fromMount,totalNeeds,curFee)

	//获取账户Nonce,var getNonce int64
	time.Sleep(time.Second * 4)
	curNonce,err,errmsg :=self.WdcRpcClient.SendNonce(curaddrrec.PubKeyHash)
	if err !=nil{
		log.Error("WDCTransProc.SendNonce() fail, get err=%v,errinfo :%s,cur fromAddress is: %v,getPubKeyHash is:%s", err,errmsg,getfromAddress,curaddrrec.PubKeyHash)
	}
	getNonce := int64(curNonce)
	toAddress :=cursettle.ToAddress

	getToPubHashStr,err :=self.WdcRpcClient.GetAddressPubHash(cursettle.ToAddress)
	if err !=nil{
		log.Error("WDCTransProc.GetAddressPubHash() fail, get err=%v,cur toAddress = %v,cur trans break!", err,toAddress)

		//1107add,转账参数不规范，通知交易系统失败：
		reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_FAILED

		if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			log.Error("WDCTransProc.GetAddressPubHash() fail, exec WithdrawsUpdate failed!,cur cursettle is:%d,toAddress is: %v,cur trans break!", cursettle.SettleId,toAddress)
		}
		return false,""
	}else{
		log.Info("WDCTransProc.GetAddressPubHash() succ, get toAddress is:%s,getToPubHashStr is:%s", toAddress,getToPubHashStr)
	}
	//,当前两次update，一次冻结，开始转账.
	//withdrawStatus.Status = SETTLE_STATUS_SIGNED

	//开始转账类型交易.
	//获取地址的私钥，调用rpc接口，进行转账：
	//sgj 1107add,签名完成状态不需要设置：SETTLE_STATUS_SIGNED
	//reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_SIGNED
	//if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
	txid,txHexStr, err,errmsg := self.ClientToTransferAccount(getAddressPub,getToPubHashStr,curAmount,getAddressPriv,int64(getNonce))
	if err !=nil || errmsg !="" {
	log.Error("WDCTransProc.ClientToTransferAccount() fail, get err=%v,cur errmsg = %v,cur trans break!", err,errmsg)
		return false,""

	}else{
		log.Info("WDCTransProc.ClientToTransferAccount() succ, gettxid is:%s, txHexStr=%s,cur errmsg = %v", txid,txHexStr,errmsg)

	}

	//开始广播交易：
	resdata,err,errcode,errmsg:= self.WdcRpcClient.SendTransactionPostForm(txHexStr)
	if errcode == proto.ErrorNodeRPCSuccess.Code{
		//resdata.(),txid,good!
		log.Info("WDCTransProc.SendTransactionPostForm() succ!,txid is:%s,getcur resdata is:%v",txid,resdata)
		//转账后通知交易系统,状态值5
		curStatus = proto.SETTLE_STATUS_PENDING

	}else{
		log.Error("WDCTransProc.SendTransactionPostForm() fail!,errcode is:%d,res errmsg is:%v",errcode,errmsg)
		curStatus = proto.SETTLE_STATUS_FAILED
	}
	//1104，可把转账结构写入数据库
	//err :=
	//sgj 1202 upgrade:
	//when save to DB:cursettle.ToAddress replace getToPubHashStr
	GWdcDataStore.SaveTranRecord(cursettle.CoinCode,getfromAddress,cursettle.ToAddress,cursettle.SettleId,txid,curAmount,"curstatusing",errcode,errmsg,"")

	reqUpdateInfo.Withdraws[0].Status = curStatus
	reqUpdateInfo.Nonce = time.Now().Unix()
	//1104 add: SignStr info:
	reqUpdateInfo.Withdraws[0].SignStr = txid
		for {
		if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			break
		} else {
			time.Sleep(time.Second * 10)
		}
	}
	//更新交易状态
	/*
		返回数据格式
		{"errno":"OK","message":"Success","data":{"nonce":144150261677}}
	*/
	return ret, "txid"
}

func (self *WDCTransHandle) WithdrawsUpdate(reqUpdateInfo *proto.WithdrawsUpdateReq) (is bool) {

	var signData string
	curUpdateInfo := proto.WithdrawsUpdateReq{}
	curUpdateInfo = *reqUpdateInfo
	resUpdateSign := proto.Response{}
	transInfo := proto.WithdrawsUpdateResp{}
	resUpdateSign.Data = &transInfo

	//UrlVerify := config.GbConf.SettleApiReq.SettlApiUpdate
	UrlVerify := config.GbConf.SettleApiUpdate

	reqBody, err := json.Marshal(&curUpdateInfo)
	if nil != err {
		log.Error("when WithdrawsUpdate,Marshal to json error:%s", err.Error())
		return false
	}

	if signData, err = auth.KSign(reqBody, config.GbConf.SettleAccessKey.AccessPrivKey); err != nil {
		log.Error("In WithdrawsUpdate(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
		return false
	}
	//step 2
	log.Info("WithdrawsUpdate,auth.KSign succ!,signData is :%v", signData)

	ht := CHttpClientEx{}
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//req.Header.Set("abit-actionsign", signData)
	//ht.HeaderSet(proto.HActionSign, signData)
	//1112 update for abit
	ht.HeaderSet(proto.HActionAbitSign, signData)


	log.Info("WDCTransProc.transInfo url=%s,cur reqdata = %v", UrlVerify, curUpdateInfo)
	strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 9000, &curUpdateInfo, &resUpdateSign)
	if nil != err {
		log.Error("withdrawsUpdate,ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return false
	} else {
		if resUpdateSign.Msg== "Success" {
			log.Info("WDCTransProc.withdrawsUpdate good!,get strRes is:%v,resUpdateSign is:%v", strRes, resUpdateSign)
			return true
		}else{
			//errcode is:%v", resUpdateSign.Code)
			log.Error("WDCTransProc.withdrawsUpdate res err! err resUpdateSign is:%v", resUpdateSign)
			return false
		}
	}
}

//sgj 1112PM add for gatterDepositServer:
//finished! part1
/*

WithdrawsDeposites
reqUpdateInfo := proto.WithdrawsUpdateReq{}
	reqUpdateInfo.Withdraws = make([]proto.Settle,1)
	if from == ""{
		getfromAddress =cursettle.FromAddress
	}else{
		getfromAddress =from
	}
	//1104 add :
	if getfromAddress == "" &&  from == ""{
		getfromAddress = "1HFCUeNHcL6Drf4TPwBLG6RgYVe9o41BVj"
	}
	reqUpdateInfo.Withdraws[0].FromAddress = getfromAddress
	reqUpdateInfo.Withdraws[0].SettleId = int64(cursettle.SettleId)
	reqUpdateInfo.Withdraws[0].CoinCode =cursettle.CoinCode
	reqUpdateInfo.Withdraws[0].SignStr = ""
	reqUpdateInfo.Nonce = time.Now().Unix()

==:


*/