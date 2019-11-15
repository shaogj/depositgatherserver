
//调用节点的请求api：
// 1.java的rpc接口，用来取得信息，发送交易
// 2.java封装好的http 接口，用来构建交易，签名交易的各类事务函数调用

package wdctranssign

import (
	"2019NNZXProj10/abitserverDepositeGather/config"
	"2019NNZXProj10/abitserverDepositeGather/proto"
	"encoding/json"
	//"bytes"
	"fmt"
	"github.com/mkideal/log"
	"io/ioutil"
	"net/http"
	"net/url"
	. "shaogj/utils"
	"strconv"
	"strings"
	"errors"
)

//var WdcRPCClient WdcRpcClient

//1030doing
type WdcRpcClient struct {
	HtClient CHttpClientEx
	config 	*config.WDCNodeConf
}

var WDCNodeUrl string = "http://192.168.1.138:19585"

var WDCJavaSDKUrl = "http://192.168.1.190:8088/wallet/WalletUtility"

func NewWdcRpcClient(nodeconf *config.WDCNodeConf) *WdcRpcClient{
	curWdcRpcClient := &WdcRpcClient{
		config:nodeconf,
		HtClient: CHttpClientEx{},
	}
	curWdcRpcClient.HtClient.Init()
	curWdcRpcClient.HtClient.HeaderSet("Content-Type", "application/json;charset=utf-8")
	return curWdcRpcClient
}


func (cur *WdcRpcClient) SendBalancePostFormNode(curpubkeyhash string) (curbalance float64,err error,errmsg string){
	data := url.Values{}
	data.Set("pubkeyhash",curpubkeyhash)
	UrlVerify:= fmt.Sprintf("%s/%s", WDCNodeUrl, "sendBalance")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0,err,""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return 0,err,""
	}
	log.Info("post send success-----55,get res is :%v",string(content))
	getResp :=&proto.NodeResponse{}
	err=json.Unmarshal(content,getResp)
	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(content),UrlVerify,err.Error())
		return 0,err,""
	}
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code{
		log.Error("SendBalancePostFormNode(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode,getResp.Message)
		return 0,errors.New(getResp.Message),"sendBalance，get结果异常"
	}
	log.Info("SendBalancePostFormNode()-----66,success,get Resp jsondata is:%v,balanceval is :%f,getResp.StatusCode is:%d",getResp,getResp.Data.(float64),getResp.StatusCode)
	return getResp.Data.(float64),nil,""
}
//sendNonce

//to add,getNonce := 344,sendtrans := 344

func (cur *WdcRpcClient) SendNonce(curpubkeyhash string) (curbalance float64,err error,errmsg string){
	data := url.Values{}
	data.Set("pubkeyhash",curpubkeyhash)
	UrlVerify:= fmt.Sprintf("%s/%s", WDCNodeUrl, "sendNonce")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return 0,err,""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return 0,err,""
	}
	log.Info("post send success-----66,curpubkeyhash is:%s,get res is :%v",curpubkeyhash,string(content))
	getResp :=&proto.NodeResponse{}
	err=json.Unmarshal(content,getResp)
	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(content),UrlVerify,err.Error())
		return 0,err,""
	}
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code{
		log.Error("SendNonce(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode,getResp.Message)
		return 0,errors.New(getResp.Message),"SendNonce，get结果异常"
	}
	log.Info("SendNonce()-----66,success,get Resp jsondata is:%v,nonceval is :%d,getResp.StatusCode",getResp,getResp.Data.(float64),getResp.StatusCode)
	return getResp.Data.(float64),nil,""
}

//1104 add,,WdcTxBlock
func (cur *WdcRpcClient) GetTransactionHeight(curheight int) (getBlockTrans interface{},err error,errmsg string){
	data := url.Values{}

	//整型转换成字符串
	curheightstr:=strconv.Itoa(curheight)

	data.Set("height",curheightstr)
	UrlVerify:= fmt.Sprintf("%s/%s", WDCNodeUrl, "getTransactionHeight")
	client := &http.Client{}
	r, _ := http.NewRequest("POST", UrlVerify, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;param=value")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err,""
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Fatal error ", err.Error())
		return nil,err,""
	}
	log.Info("post send success-----66,curheight is:%d,get res is :%v",curheight,string(content))
	getResp :=&proto.NodeResponse{}

	curWdcTxBlock :=make([]proto.WdcTxBlock,3)
	getResp.Data = &curWdcTxBlock
	/*2 method
		txsecHexStrBef :=resSDKAccount.Data.(string)
		getResp :=&proto.WdcTxBlock{}
		err=json.Unmarshal([]byte(txsecHexStrBef),getResp)
		txHexStr = getResp.Message

	*/
	err=json.Unmarshal(content,getResp)
	//sgj 1105 watching:
	getWdcTxBlock1 :=getResp.Data.(*[]proto.WdcTxBlock)
	log.Info("post send success-----77,curheight is:%d,err is:%v,get getWdcTxBlock1 is :%v",curheight,err,getWdcTxBlock1)

	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(content),UrlVerify,err.Error())
		return nil,err,""
	}
	log.Info("post send success---------8888,get blocktranslen is:%d",len(*getWdcTxBlock1))
	if getResp.StatusCode != proto.ErrorNodeRPCSuccess.Code{
		log.Error("SendNonce(),get NodeResponse error!,get StatusCode is:%d,msg is:%s ", getResp.StatusCode,getResp.Message)
		return nil,errors.New(getResp.Message),"SendNonce，get结果异常"
	}
	log.Info("SendNonce()-----66,success,get Resp jsondata is:%v,nonceval is :%d,getResp.StatusCode",getResp,getResp.Data.(*[]proto.WdcTxBlock),getResp.StatusCode)
	return getResp.Data.(*[]proto.WdcTxBlock),nil,""
}

//广播事务
func (cur *WdcRpcClient) SendTransactionPostForm(curtraninfo string) (resdata interface{},err error,errcode int,errmsg string) {
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
func (cur *WdcRpcClient) GetAddressPubHash(address string) (pubHashstr string,err error) {
	accountAddress := proto.AddressToPubkeyHash{}
	accountAddress.Address = address
	resSDKAccount := proto.JavaSDKResponse{    }

	UrlVerify := fmt.Sprintf("%s/%s", WDCJavaSDKUrl, "addressToPubkeyHash")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return "",err
	}
	log.Info("WdcRpcClient.addressToPubkeyHash,get statusCode is :%s,res=%s",statusCode, strRes)
	curPubkeyHashStr := ""
	if statusCode == 200{
		curPubkeyHashStr =resSDKAccount.Data.(string)
		log.Info("WdcRpcClient. get addressToPubkeyHash succ,value is:%v", curPubkeyHashStr)
	}else{
		log.Error("WdcRpcClient. get addressToPubkeyHash error!,value is:%v,statusCode is:%s", curPubkeyHashStr,statusCode)
		return "",err
	}
	return curPubkeyHashStr,nil

}

//1106 add,获取区块高度:
//RequestResponse
func (cur *WdcRpcClient) GetBlockHeight() (blockHeight float64,err error) {
	resNodeRet := proto.NodeResponse{    }

	UrlVerify := fmt.Sprintf("%s/%s", WDCNodeUrl, "height")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestResponseJson(UrlVerify, nil,5000, &resNodeRet)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return 0,err
	}
	log.Info("WdcRpcClient.GetBlockHeight,get statusCode is :%s,res=%s",statusCode, strRes)
	var curBlockHeight float64
	//resNodeRet
	if resNodeRet.StatusCode == proto.ErrorNodeRPCSuccess.Code{
	//if statusCode == proto.ErrorNodeRPCSuccess.Code{
		curBlockHeight =resNodeRet.Data.(float64)
		log.Info("WdcRpcClient. get GetBlockHeight succ,value is:%v", curBlockHeight)
	}else{
		log.Error("WdcRpcClient. get GetBlockHeight error!,value is:%v,statusCode is:%s", curBlockHeight,statusCode)
		return 0,err
	}
	return curBlockHeight,nil

}


//1105add,pubkeyHashToAddress,解析区块所用
//4）通过公钥哈希获得地址
func (cur *WdcRpcClient) GetPubkeyHashToAddress(pubkeyHash string) (pubHashstr string,err error) {
	accountAddress := proto.PubkeyHashToAddress{}
	accountAddress.PubkeyHashStr = pubkeyHash
	resSDKAccount := proto.JavaSDKResponse{    }

	UrlVerify := fmt.Sprintf("%s/%s", WDCJavaSDKUrl, "pubkeyHashToAddress")

	strRes, statusCode, errorCode, err := cur.HtClient.RequestJsonResponseJson(UrlVerify, 5000, &accountAddress, &resSDKAccount)
	if nil != err {
		log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
		return "",err
	}
	log.Info("WdcRpcClient.pubkeyHashToAddress,get statusCode is :%s,res=%s",statusCode, strRes)
	curPubHashAddrStr := ""
	if statusCode == 200{
		curPubHashAddrStr =resSDKAccount.Data.(string)
		log.Info("WdcRpcClient. get pubkeyHashToAddress succ,value is:%v", curPubHashAddrStr)
	}else{
		if statusCode == 500 {
			errors.New(resSDKAccount.Message)
		}
		log.Error("WdcRpcClient. get pubkeyHashToAddress error!,value is:%v,statusCode is:%s", curPubHashAddrStr,statusCode)
		return "",err
	}
	return curPubHashAddrStr,nil

}


//1030 add,form 方式获取nonce：sendNonce


//根据区块高度获取事务列表
/*方法: getTransactionHeight(POST)

*/
