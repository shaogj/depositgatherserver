//调用节点的请求api：
// 1.java的rpc接口，用来取得信息，发送交易
// 2.java封装好的http 接口，用来构建交易，签名交易的各类事务函数调用

package wdctranssign

import (
	"2019NNZXProj10/depositgatherserver/config"
	"2019NNZXProj10/depositgatherserver/proto"
	"backend/support/libraries/loggers"
	"encoding/json"
	//"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	. "shaogj/utils"
	"strconv"
	"strings"

	"github.com/mkideal/log"
)

//var WdcRPCClient WdcRpcClient

//1030doing
type WdcRpcClient struct {
	HtClient CHttpClientEx
	config   *config.WDCNodeConf
}

var WDCNodeUrl string = "http://192.168.1.138:19585"

var WDCJavaSDKUrl = "http://192.168.1.190:8088/wallet/WalletUtility"

//20200608add
var GWGCTransUrl = "http://18.176.110.109:8090/wallet/TxUtility"

func NewWdcRpcClient(nodeconf *config.WDCNodeConf) *WdcRpcClient {
	curWdcRpcClient := &WdcRpcClient{
		config:   nodeconf,
		HtClient: CHttpClientEx{},
	}
	curWdcRpcClient.HtClient.Init()
	curWdcRpcClient.HtClient.HeaderSet("Content-Type", "application/json;charset=utf-8")
	return curWdcRpcClient
}

func (cur *WdcRpcClient) SendBalancePostFormNode(curpubkeyhash string) (curbalance float64, err error, errmsg string) {
	data := url.Values{}
	data.Set("pubkeyhash", curpubkeyhash)
	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "sendBalance")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err, ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return 0, err, ""
	}
	log.Info("post send success-----55,get res is :%v", string(content))
	getResp := &proto.NodeResponse{}
	err = json.Unmarshal(content, getResp)
	if nil != err {
		log.Error("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		return 0, err, ""
	}
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code {
		log.Error("SendBalancePostFormNode(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode, getResp.Message)
		return 0, errors.New(getResp.Message), "sendBalance，get结果异常"
	}
	log.Info("SendBalancePostFormNode()-----66,success,get Resp jsondata is:%v,balanceval is :%f,getResp.StatusCode is:%d", getResp, getResp.Data.(float64), getResp.StatusCode)
	return getResp.Data.(float64), nil, ""
}

//sendNonce

//to add,getNonce := 344,sendtrans := 344

func (cur *WdcRpcClient) SendNonce(curpubkeyhash string) (curbalance float64, err error, errmsg string) {
	data := url.Values{}
	data.Set("pubkeyhash", curpubkeyhash)
	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "sendNonce")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err, ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return 0, err, ""
	}
	log.Info("post send success-----66,curpubkeyhash is:%s,get res is :%v", curpubkeyhash, string(content))
	getResp := &proto.NodeResponse{}
	err = json.Unmarshal(content, getResp)
	if nil != err {
		log.Error("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		return 0, err, ""
	}
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code {
		log.Error("SendNonce(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode, getResp.Message)
		return 0, errors.New(getResp.Message), "SendNonce，get结果异常"
	}
	log.Info("SendNonce()-----66,success,get Resp jsondata is:%v,nonceval is :%d,getResp.StatusCode", getResp, getResp.Data.(float64), getResp.StatusCode)
	return getResp.Data.(float64), nil, ""
}

//1104 add,,WdcTxBlock
func (cur *WdcRpcClient) GetTransactionHeightOld(curheight int) (getBlockTrans interface{}, err error, errmsg string) {
	data := url.Values{}

	//整型转换成字符串
	curheightstr := strconv.Itoa(curheight)

	data.Set("height", curheightstr)
	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "getTransactionHeight")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err, ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return nil, err, ""
	}
	log.Info("post send success-----66,curheight is:%d,get res is :%v", curheight, string(content))
	getResp := &proto.NodeResponse{}

	curWdcTxBlock := make([]proto.WdcTxBlock, 3)
	getResp.Data = &curWdcTxBlock
	/*2 method
	txsecHexStrBef :=resSDKAccount.Data.(string)
	getResp :=&proto.WdcTxBlock{}
	err=json.Unmarshal([]byte(txsecHexStrBef),getResp)
	txHexStr = getResp.Message

	*/
	err = json.Unmarshal(content, getResp)
	//sgj 1105 watching:
	getWdcTxBlock1 := getResp.Data.(*[]proto.WdcTxBlock)
	log.Info("post send success-----77,curheight is:%d,err is:%v,get getWdcTxBlock1 is :%v", curheight, err, getWdcTxBlock1)

	if nil != err {
		log.Error("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		return nil, err, ""
	}
	log.Info("post send success---------8888,get blocktranslen is:%d", len(*getWdcTxBlock1))
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code {
		log.Error("SendNonce(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode, getResp.Message)
		return nil, errors.New(getResp.Message), "SendNonce，get结果异常"
	}
	log.Info("SendNonce()-----66,success,get Resp jsondata is:%v,nonceval is :%d,getResp.StatusCode", getResp, getResp.Data.(*[]proto.WdcTxBlock), getResp.StatusCode)
	return getResp.Data.(*[]proto.WdcTxBlock), nil, ""
}

//20200120new Node v9.0,,WdcTxBlock
func (cur *WdcRpcClient) GetTransactionHeight(curheight int) (getBlockTrans interface{}, err error, errmsg string) {
	data := url.Values{}

	//整型转换成字符串
	curheightstr := strconv.Itoa(curheight)

	data.Set("height", curheightstr)
	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "getTransactionHeight")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err, ""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.Error.Printf("Fatal error ", err.Error())
		return nil, err, ""
	}
	loggers.Info.Printf("post send getTransactionHeight success-----66,curheight is:%d,get res is :%v", curheight, string(content))
	//getResp :=&proto.NodeResponse{}

	//sgj 0220updating
	curWdcTxBlock := make([]proto.WdcTxBlock, 3)
	//getResp.Data = &curWdcTxBlock

	//err=json.Unmarshal(content,getResp)
	//sgj 20200220doing: watching:,WDC Node 9.0ver
	if string(content) == "[ ]" {
		loggers.Info.Printf("post send success,,cur blockheight is :%d,get Node New WdcTxBlock info is:%s", curheight, string(content))
		return nil, nil, "emptyBlock succ"

	}
	err = json.Unmarshal(content, &curWdcTxBlock)
	//sgj 1105 watching:
	//getWdcTxBlock1 :=getResp.Data.(*[]proto.WdcTxBlock)
	if nil != err {
		loggers.Error.Printf("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		return nil, err, ""
	}

	//sgjj 0220 adding
	getWdcTxBlock1 := &curWdcTxBlock

	loggers.Info.Printf("post send success-----77,curheight is:%d,err is:%v,get getWdcTxBlock1 is :%v", curheight, err, getWdcTxBlock1)
	loggers.Info.Printf("post send success---------8888,get blocktranslen is:%d", len(*getWdcTxBlock1))

	//return getResp.Data.(*[]proto.WdcTxBlock),nil,""
	return &curWdcTxBlock, nil, ""
}

//广播事务
func (cur *WdcRpcClient) SendTransactionPostForm(curtraninfo string) (resdata interface{}, err error, errcode int, errmsg string) {
	data := url.Values{}
	data.Set("traninfo", curtraninfo)
	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "sendTransaction")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err, proto.ErrorRequestWDCNode.Code, proto.ErrorRequestWDCNode.Desc
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return 0, err, 3333, "sendTransactionPostForm，get结果异常"
	}
	log.Info("post sendTransactionPostForm success-----55,get res is :%v", string(content))
	getResp := &proto.NodeResponse{}
	err = json.Unmarshal(content, getResp)
	if nil != err {
		log.Error("resp=%s,url=%s,err=%v", string(content), UrlVerify, err.Error())
		return 0, err, proto.ErrorRequestWDCNodeJust.Code, proto.ErrorRequestWDCNodeJust.Desc
	}
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code {
		//res is :{"message":"Nonce is too small","data":null,"code":5000},,etc..
		//写入提现错误记录到db,message
	} else {
		log.Info("sendTransactionPostForm()-----66,success,get Resp jsondata is:%v,,getResp.StatusCode", getResp, getResp.StatusCode)
	}
	return getResp.Data, nil, getResp.StatusCode, getResp.Message

}

//3）通过地址获得公钥哈希
func (cur *WdcRpcClient) GetAddressPubHash(address string) (pubHashstr string, err error) {
	accountAddress := proto.AddressToPubkeyHash{}
	accountAddress.Address = address
	resSDKAccount := proto.JavaSDKResponse{}

	UrlVerify := fmt.Sprintf("%s/%s", WDCJavaSDKUrl, "addressToPubkeyHash")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return "", err
	}
	log.Info("WdcRpcClient.addressToPubkeyHash,get statusCode is :%s,res=%s", statusCode, strRes)
	curPubkeyHashStr := ""
	if statusCode == 200 {
		curPubkeyHashStr = resSDKAccount.Data.(string)
		log.Info("WdcRpcClient. get addressToPubkeyHash succ,value is:%v", curPubkeyHashStr)
	} else {
		log.Error("WdcRpcClient. get addressToPubkeyHash error!,value is:%v,statusCode is:%s", curPubkeyHashStr, statusCode)
		return "", err
	}
	return curPubkeyHashStr, nil

}

//1106 add,获取区块高度:
//RequestResponse
func (cur *WdcRpcClient) GetBlockHeight() (blockHeight float64, err error) {
	resNodeRet := proto.NodeResponse{}

	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "height")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil, 5000, &resNodeRet)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return 0, err
	}
	log.Info("WdcRpcClient.GetBlockHeight,get statusCode is :%s,res=%s", statusCode, strRes)
	var curBlockHeight float64
	//resNodeRet
	if resNodeRet.StatusCode == proto.ErrorNodeRPCSuccess.Code {
		//if statusCode == proto.ErrorNodeRPCSuccess.Code{
		curBlockHeight = resNodeRet.Data.(float64)
		log.Info("WdcRpcClient. get GetBlockHeight succ,value is:%v", curBlockHeight)
	} else {
		log.Error("WdcRpcClient. get GetBlockHeight error!,value is:%v,statusCode is:%s", curBlockHeight, statusCode)
		return 0, err
	}
	return curBlockHeight, nil

}

//20200109----add form utxo for:https://blockchain.info/unspent?active=1Eq8xXAea47WPY5t8zUEYDKgcWB7cptZWB
//getBlockTrans interface{},

func (cur *WdcRpcClient) GetBTCTxUnSpentLimit(address string) (getBtcUtxoInfo []proto.BtcUtxoInfo, count int, err error) {
	//resNodeRet := proto.NodeResponse{    }
	retbtcutxoinfo := make([]proto.BtcUtxoInfo, 0, 10)
	resNodeRet := proto.BTCUnspentOutputs{}
	//WDCNodeUrl
	BlockChainUrl := "https://blockchain.info/unspent?active="
	BlockChainUrlStr := fmt.Sprintf("%s%s", BlockChainUrl, address)
	UrlVerify := BlockChainUrlStr
	//fmt.Sprintf("%s/%s", BlockChainUrlStr, "height")

	//sgj 0106PMadd,,to upgrade timeout for BTC,5s,to 15s
	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil, 15000, &resNodeRet)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return retbtcutxoinfo, 0, err
	}
	log.Info("WdcRpcClient.GetTxUnSpentLimit,get statusCode is :%s,res=%s", statusCode, strRes)
	//var curBlockHeight float64
	//resNodeRet
	if len(resNodeRet.CurBtcUtxoInfo) > 0 {
		//if statusCode == proto.ErrorNodeRPCSuccess.Code{
		retbtcutxoinfo = resNodeRet.CurBtcUtxoInfo
		log.Info("WdcRpcClient. get GetTxUnSpentLimit succ,value is:%v", retbtcutxoinfo)
	} else {
		log.Error("WdcRpcClient. get GetTxUnSpentLimit error!,value is:%v,statusCode is:%s", retbtcutxoinfo, statusCode)
		return retbtcutxoinfo, 0, err
	}
	return retbtcutxoinfo, len(retbtcutxoinfo), nil
	//return getResp.Data.(*[]proto.WdcTxBlock),nil,""

}

//sgj 0116ing add

func (cur *WdcRpcClient) GetTxUnSpentLimit(address string) (getBtcUtxoInfo []proto.BtcUtxoInfo, count int, err error) {
	//resNodeRet := proto.NodeResponse{    }
	retbtcutxoinfo := make([]proto.BtcUtxoInfo, 0, 10)
	resNodeRet := proto.BTCUnspentOutputs{}
	//WDCNodeUrl
	BlockChainUrl := "https://blockchain.info/unspent?active="
	BlockChainUrlStr := fmt.Sprintf("%s%s", BlockChainUrl, address)
	UrlVerify := BlockChainUrlStr
	//fmt.Sprintf("%s/%s", BlockChainUrlStr, "height")

	//sgj 0106PMadd,,to upgrade timeout for BTC,5s,to 15s
	//sgj 0105 updating,try,,15000 to 20000
	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil, 20000, &resNodeRet)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return retbtcutxoinfo, 0, err
	}
	log.Info("WdcRpcClient.GetTxUnSpentLimit,get statusCode is :%s,res=%s", statusCode, strRes)
	//var curBlockHeight float64
	//resNodeRet
	if len(resNodeRet.CurBtcUtxoInfo) > 0 {
		//if statusCode == proto.ErrorNodeRPCSuccess.Code{
		retbtcutxoinfo = resNodeRet.CurBtcUtxoInfo
		log.Info("WdcRpcClient. get GetTxUnSpentLimit succ,value is:%v", retbtcutxoinfo)
	} else {
		log.Error("WdcRpcClient. get GetTxUnSpentLimit error!,value is:%v,statusCode is:%s", retbtcutxoinfo, statusCode)
		return retbtcutxoinfo, 0, err
	}
	return retbtcutxoinfo, len(retbtcutxoinfo), nil
	//return getResp.Data.(*[]proto.WdcTxBlock),nil,""

}

//1105add,pubkeyHashToAddress,解析区块所用
//4）通过公钥哈希获得地址
func (cur *WdcRpcClient) GetPubkeyHashToAddress(pubkeyHash string) (pubHashstr string, err error) {
	accountAddress := proto.PubkeyHashToAddress{}
	accountAddress.PubkeyHashStr = pubkeyHash
	resSDKAccount := proto.JavaSDKResponse{}

	UrlVerify := fmt.Sprintf("%s/%s", WDCJavaSDKUrl, "pubkeyHashToAddress")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return "", err
	}
	log.Info("WdcRpcClient.pubkeyHashToAddress,get statusCode is :%s,res=%s", statusCode, strRes)
	curPubHashAddrStr := ""
	if statusCode == 200 {
		curPubHashAddrStr = resSDKAccount.Data.(string)
		log.Info("WdcRpcClient. get pubkeyHashToAddress succ,value is:%v", curPubHashAddrStr)
	} else {
		if statusCode == 500 {
			errors.New(resSDKAccount.Message)
		}
		log.Error("WdcRpcClient. get pubkeyHashToAddress error!,value is:%v,statusCode is:%s", curPubHashAddrStr, statusCode)
		return "", err
	}
	return curPubHashAddrStr, nil

}


//20200109PM add for verifyAddress
func (cur *WdcRpcClient) CheckVerifyAddress(curAddr string) (VerifyRet int64,err error) {
	accountAddress := proto.VerifyAddressReq{}
	accountAddress.Address = curAddr
	resSDKAccount := proto.JavaSDKResponse{    }
	var curVerifyVal int64
	var curVerifyValStr string
	UrlVerify := fmt.Sprintf("%s/%s", WDCJavaSDKUrl, "verifyAddress")
	//add for test
	log.Info("WdcRpcClient.CheckVerifyAddress,cur UrlVerify is ---to check :%s", UrlVerify)

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return -2,err
	}
	log.Info("WdcRpcClient.CheckVerifyAddress,get statusCode is :%s,res=%s",statusCode, strRes)
	if statusCode == 200{
		curVerifyValStr =resSDKAccount.Data.(string)
		if curVerifyValStr == "-1" || curVerifyValStr == "-2"{
			curVerifyVal = -1
		}else{
			tmpVal,_ := strconv.Atoi(curVerifyValStr)
			curVerifyVal = int64(tmpVal)
		}

		log.Info("WdcRpcClient. get CheckVerifyAddress succ,curVerifyVal is:%v", curVerifyVal)
	}else{
		if statusCode == 500 {
			errors.New(resSDKAccount.Message)
		}
		curVerifyVal = -3
		log.Error("WdcRpcClient. get CheckVerifyAddress error!,value is:%v,statusCode is:%s", curVerifyVal,curVerifyVal)
		return curVerifyVal,err
	}
	return curVerifyVal,nil

}

//20200605add,Token WGC getbalance get method:
//RequestResponse, http://47.74.183.249:19585/TokenBalance/?code=WGC&address=WX1KVcQTbsMuU5jpZSdBRXiKcbbawrGBo9h7

func (cur *WdcRpcClient) GetWDCAddrTokenBalance(coinName string, address string) (curbalance float64, err error) {
	resNodeRet := proto.NodeResponse{}

	UrlVerify := fmt.Sprintf("%s/%s/?code=%s&address=%s", WDCNodeUrl, "TokenBalance", coinName, address)

	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil, 5000, &resNodeRet)

	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return 0, err
	}
	log.Info("WdcRpcClient.GetWDCAddrTokenBalance,get statusCode is :%s,res=%s", statusCode, strRes)
	var curAddrBalance float64

	if resNodeRet.StatusCode == proto.ErrorNodeRPCSuccess.Code {
		curAddrBalance = resNodeRet.Data.(float64)
		log.Info("WdcRpcClient. get GetWDCAddrTokenBalance succ,value is:%v", curAddrBalance)
	} else {
		log.Error("WdcRpcClient. get GetWDCAddrTokenBalance error!,value is:%v,statusCode is:%s", curAddrBalance, statusCode)
		return 0, err
	}
	return curAddrBalance, nil

}

//20200603 add,获取区块数据新接口:	Get方式解析区块
//RequestResponse，http://47.74.183.249:19585/block/2154049
func (cur *WdcRpcClient) GetBlockFullData(curheight int) (getBlockTrans interface{}, err error, errmsg string) {

	//sgj 0605 tmp testing
	//0607
	//return nil, nil, "test WGC trans skip succ"

	resNodeRet := proto.BlockHeadWDCResponse{}
	curWdcTxBlock := make([]proto.WdcTxBlockNew, 3)
	resNodeRet.BlockBodyData = &curWdcTxBlock
	var curBlockcontent string
	UrlVerify := fmt.Sprintf("%s/%s/%d", WDCNodeUrl, "block", curheight)
	log.Info("WdcRpcClient.cur Invoke RPC URL is:%v", UrlVerify)

	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil, 5000, &resNodeRet)
	if nil != err || statusCode != 200 {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return nil, nil, "emptyBlock succ"
	}
	log.Info("WdcRpcClient.GetBlockFullData,get statusCode is :%s,res=%s", statusCode, strRes)
	curBlockcontent = string(strRes)
	log.Info("WdcRpcClient. get GetBlockFullData succ,value is:%v", curheight)

	//sgj 20200220doing: watching:,WDC Node 9.0ver
	if string(curBlockcontent) == "[ ]" {
		log.Info("post send success,,cur blockheight is :%d,get Node New WdcTxBlock info is:%s", curheight, string(curBlockcontent))
		return nil, nil, "emptyBlock succ"

	}

	//sgjj 0220 adding
	getWdcTxBlock1 := &curWdcTxBlock

	log.Info("get blockdata finish!,curheight is:%d,err is:%v,get blocktranslen is:%d,getWdcTxBlock1 is :%v", curheight, err, len(*getWdcTxBlock1), getWdcTxBlock1)

	return &curWdcTxBlock, nil, ""

}

//20200605 add,获取区块数据Token' sub data:	Get方式解析区块
//RequestResponse，http://47.74.183.249:19585/block/2154049
func (cur *WdcRpcClient) GetPayLoadTransaction(assetPayload string) (getPayloadData interface{}, err error, statusCode int) {
	curPayload := proto.PayloadStrReq{}
	curPayload.Payload = assetPayload
	resSDKAccount := proto.JavaSDKResponse{}

	//resTransData :=proto.AssetPayLoadTransaction{}
	UrlVerify := fmt.Sprintf("%s/%s", GWGCTransUrl, "getAssetTransfer")
	getResp := proto.BlockPayloadResponse{}

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &curPayload, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return &getResp, err, statusCode
	}
	log.Info("WdcRpcClient.GetPayLoadTransaction,get statusCode is :%s,res=%s", statusCode, strRes)
	curPubStr := ""

	//if statusCode == 200{
	if resSDKAccount.StatusCode == "2000" {
		//二次对node rpc' 内部封装的解析
		txsecMessageStr := resSDKAccount.Message
		log.Info("WdcRpcClient.GetPayLoadTransaction,get statusCode is :%s,message=%s,", statusCode, txsecMessageStr)

		//sgj 0103 fix bug
		if txsecMessageStr == "" {
			log.Error("ht.GetPayLoadTransaction , getResp's Data =%v .it is emptystring!return is err!", txsecMessageStr)

		} else {
			err = json.Unmarshal([]byte(txsecMessageStr), &getResp)
			if nil != err {
				log.Error("resp=%s,url=%s,err=%v", string(txsecMessageStr), UrlVerify, err.Error())
			} else {
				log.Info("ht.GetPayLoadTransaction , getResp.From=%s,getResp.To=%s,valus is %d, ", getResp.From, getResp.To, getResp.Value)
			}

		}
	} else {

		log.Error("WdcRpcClient. get GetPayLoadTransaction error!,value is:%v,statusCode is:%s", curPubStr, statusCode)
	}
	return &getResp, nil, statusCode

}

//1030 add,form 方式获取nonce：sendNonce

//根据区块高度获取事务列表
/*方法: getTransactionHeight(POST)

 */
