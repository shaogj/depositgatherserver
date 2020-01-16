// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//package KTCtranssign
package ktcrpc

import (
	"bytes"
	"fmt"
	//"io/ioutil"
	//"path/filepath"
	"time"

	//sgj 0514 set:
	//"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcutil"
	//"2019NNZXProj10/depositgatherserver/config"
	"2019NNZXProj10/depositgatherserver/config"

	"encoding/hex"

	"github.com/ltcsuite/ltcd/wire"
	//gsgj 0612 add:
	"github.com/ybbus/jsonrpc"
    "encoding/base64"
	//"github.com/btcsuite/btcutil"
	//"encoding/json"
 	"github.com/mkideal/log"
	"2019NNZXProj10/depositgatherserver/proto"
	"encoding/json"


)

var KTCRPCClient KTCRpcClient

//连接bitcoind的RPC客户端
type KTCRpcClient struct {
	RpcClient jsonrpc.RPCClient
	config 	*config.USDTConf
}

//func NewKTCRpcClient(conf *LTCConf) *KTCRpcClient {
func NewKTCRpcClient(conf *config.USDTConf) *KTCRpcClient {
	c := &KTCRpcClient{
		config: conf,
	}
	KTCRPCClient = *c
	//go c.Run()
	return c
}

func (cur *KTCRpcClient) RpcConnect() (*KTCRpcClient, error) {

	//sgj 0611 add
	//KTCRpcClient = jsonrpc.NewClientWithOpts("http://192.168.10.232:8335", &jsonrpc.RPCClientOpts{
	//curRpcClient 
	log.Info("===0614==before KTC RpcConnect() ,RPC config info is %v", cur.config)
	cur.RpcClient = jsonrpc.NewClientWithOpts(cur.config.RPCHostPort, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			//"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("KTCcorerpc"+":"+"123456")),
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(cur.config.RPCUser+":"+cur.config.RPCPassWord)),
			},
		})
		
	log.Info("===0614==after KTC RpcConnect() ,RpcClient info is %v", cur.RpcClient)
	//cur.RpcClient = curRpcClient
	//response, err := cur.RpcClient.Call("getinfo")
    //fmt.Printf("--001---getinfo: ------001--get response info is--001: %v;;err is :%v\n", response,err)

	curblockinfo,err := cur.RpcClient.Call("getblockchaininfo")
	if err != nil {
		log.Error("KTC GetBlockChainInfo () failrue! errinfo: %v,curblockinfo is :%v", err,curblockinfo)
		//退出进程
		//log.Fatal(err)
		return nil, err
	}
	log.Info("KTC GetBlockChainInfo (),curblockinfo is: %v", curblockinfo)
	return cur,nil
}

func (cur *KTCRpcClient) GetHexMsg(txHex string) (*ltcutil.Tx, error) {
	//sgj 0409 add
	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, err
	}
	return ltcutil.NewTx(&msgTx), nil

}
//func Getinfo(){---getinfo
func GetBlockChainInfo(){
    response, err := KTCRPCClient.RpcClient.Call("getblockchaininfo")
    fmt.Printf("--001---getblockchaininfo:  get response info is--001: %v;;err is :%v\n", response,err)

} 
func (cur *KTCRpcClient) SendTransaction(txHexStr string) (txid string, err error) {

	cursendstr := txHexStr
	var rawTxHex string

	//response, err := cur.RpcClient.Call("sendrawtransaction",curTx.MsgTx(),true)
    log.Info("before KTCRpc: sendrawtransaction():params txHexStr info is: %v\n", txHexStr)
	response, err := cur.RpcClient.Call("sendrawtransaction",cursendstr,true)
    log.Info("--001---sendrawtransaction: ------001--get response info is--001: %v;;err is :%v\n", response,err)
	if err != nil {
		log.Error("SendRawTransaction() err happend!,err is :%v", err)
		return "", err
	}
	//{ <nil> -26:dust (code 64) 0};;err is :<nil>
	//sgj 1217 enhance: //&&  len(response) > 40
	if response != nil  && response.Result != nil{
		rawTxHex =response.Result.(string)
	}
	log.Info("SendRawTransaction() success!,created Txid is----->:%s", rawTxHex)
	return rawTxHex, nil
}
//sgj 0911 add sign:signmessagewithprivkey
func (cur *KTCRpcClient) SignMessageWithPrivkey(privkey string,txHexStr string) (signedrawHex string, err error) {

	var rawTxHex string
	response, err := cur.RpcClient.Call("signmessagewithprivkey",privkey,txHexStr)
    log.Info("invoke RPC signmessagewithprivkey(),is--001: get response is %v;err is :%v\n", response,err)
	if err != nil {
		log.Info("signmessagewithprivkey() err happend!,err is :%v", err)
		return "", err
	}
	if response != nil {
		rawTxHex =response.Result.(string)
	}
	log.Info("signmessagewithprivkey() success!,get rawTxHex is----->:%s", rawTxHex)
	//rawTxHex := fmt.Sprintf("%v",response)
	return rawTxHex, nil
}

func (cur *KTCRpcClient) RpcClose() {

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Info("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*4, func() {
		log.Info("Client shutting down...")
		//cur.RpcClient.Shutdown()
		log.Info("Client shutdown complete.")
	})
	//cur.rpcClient.WaitForShutdown()
}


//sgj 1208
func (cur *KTCRpcClient) GetRPCTxUnSpentLimit(minconfnum,maxconfnum int,address []string) (getBtcUtxoInfo []proto.CurKtcUtxoInfo,count int,err error) {
	retbtcutxoinfo := make([]proto.CurKtcUtxoInfo,0,1)

	//BlockChainUrl := "https://blockchain.info/unspent?active="
	ktcUtxoReq := proto.KtcUtxoReq{}
	ktcUtxoReq.QueryAddrList = address
	//sec mode:

	ktcUtxoReqAddr:= make([]proto.KtcUtxoAddrReq,0,3)
	var curAddr proto.KtcUtxoAddrReq
	curAddr.Address = address[0]
	ktcUtxoReqAddr = append(ktcUtxoReqAddr,curAddr)

	//curAddr.Address = "3D3rnFDaNs2hC3YCkxU6XPP9gkFJuS6abA"
	//ktcUtxoReqAddr = append(ktcUtxoReqAddr,curAddr)

	log.Info("KtcRpcClient.GetRPCTxUnSpentLimit(),params:addrlen is:%d,address is:%v,cur ktcUtxoReq----00002 is :%v",len(address),address,ktcUtxoReq)

	minconfnum = 1
	if maxconfnum == 0 {
		maxconfnum = 9999999
	}
	//response, err :=cur.RpcClient.Call("listunspent",minconfnum,maxconfnum,ktcUtxoReq)
	getCurReqUtxoResp, err :=cur.RpcClient.Call("listunspent",minconfnum,maxconfnum,address)

	//why log.Info err !!!
	//log.Info("KtcRpcClient.GetRPCTxUnSpentLimit,get response is :%v,err is :%v",*getCurReqUtxoResp,err)

	if nil != err {
		log.Error("RPC listunspent(),error=%d.%v response=%v ",err, getCurReqUtxoResp)
		return retbtcutxoinfo,0,err
	}
	/**/
	//1209add

	//0703 add:
	get_response, err := json.Marshal(getCurReqUtxoResp.Result)
	if err != nil {
		log.Error("GetRPCTxUnSpentLimit(),response.Result err !, get_response is:%v,err is:%v", "get_response=infos", err)
		//return nil,err
	}else {
		log.Info("GetRPCTxUnSpentLimit(),response.Result succ !, get_response is:%v,err is:%v", "get_response=infos", err)
	}
	resNodeRetUtxo := []proto.CurKtcUtxoInfo{}
	err = json.Unmarshal([]byte(get_response), &resNodeRetUtxo)
	if err != nil {
		//log.Error("KtcRpcClient. get GetRPCTxUnSpentLimit error!,getutxo len is:%d,value is:%v", len(resNodeRetUtxo),retbtcutxoinfo)
		log.Error("KTC GetRPCTxUnSpentLimit(),Unmarshal to getKTCRawTx{} err !, get_response is:%v,err is:%v", get_response, err)
	}else{
		log.Info("KTC GetRPCTxUnSpentLimit succ,get resNodeRetUtxo len is:%d,value is:%v,", len(resNodeRetUtxo),resNodeRetUtxo)

	}

	//1209 end add
	//resNodeRetUtxo :=response.Result.([]proto.CurKtcUtxoInfo)

	return resNodeRetUtxo,len(resNodeRetUtxo),nil
	//return getResp.Data.(*[]proto.WdcTxBlock),nil,""


}