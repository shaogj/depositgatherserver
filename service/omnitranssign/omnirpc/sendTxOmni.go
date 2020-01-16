// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//package omnitranssign
package omnirpc

import (
	"bytes"
	"fmt"
	//"io/ioutil"
	//"path/filepath"
	"time"

	//"github.com/btcsuite/btcd/rpcclient"
	//sgj 0514 set:
	//"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcutil"
	//"github.com/davecgh/go-spew/spew"
	//sgj add:
	//"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	//sgj 0524 add:
	"2019NNZXProj10/depositgatherserver/config"

	"encoding/hex"

	"github.com/ltcsuite/ltcd/wire"
	//gsgj 0612 add:
	"github.com/ybbus/jsonrpc"
    "encoding/base64"
	//"github.com/btcsuite/btcutil"
	//"encoding/json"
 	"github.com/mkideal/log"

)

var OmniRPCClient OmniUsdtRpcClient

//连接bitcoind的RPC客户端
//var LtcClient LtcRpcClient
//var m_RpcClient *rpcclient.Client
/*
type LTCConf struct{
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
}
*/
type OmniUsdtRpcClient struct {
	RpcClient jsonrpc.RPCClient
//	rpcClient *rpcclient.Client
//	config 	*LTCConf
	config 	*config.USDTConf
}

//func NewOmniUsdtRpcClient(conf *LTCConf) *OmniUsdtRpcClient {
func NewOmniUsdtRpcClient(conf *config.USDTConf) *OmniUsdtRpcClient {
	c := &OmniUsdtRpcClient{
		config: conf,
	}
	OmniRPCClient = *c
	//go c.Run()
	return c
}

func (cur *OmniUsdtRpcClient) RpcConnect() (*OmniUsdtRpcClient, error) {

	//sgj 0611 add
	//OmniRPCClient = jsonrpc.NewClientWithOpts("http://192.168.10.232:8335", &jsonrpc.RPCClientOpts{
	//curRpcClient 
	cur.RpcClient = jsonrpc.NewClientWithOpts(cur.config.RPCHostPort, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			//"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("omnicorerpc"+":"+"123456")),
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(cur.config.RPCUser+":"+cur.config.RPCPassWord)),
			},
		})
		
	log.Info("===0614==after usdt RpcConnect() ,RpcClient info is %v", cur.RpcClient)
	curblockinfo,err := cur.RpcClient.Call("getblockchaininfo")
	if err != nil {
		log.Info("Usdt GetBlockChainInfo () failrue! errinfo: %v,curblockinfo is :%v", err,curblockinfo.Result)
		//退出进程
		//log.Fatal(err)
		return nil, err
	}
	log.Info("Usdt GetBlockChainInfo (),getinfo: %v", curblockinfo)

	//sgj 1119adding
	curblockcount,err := cur.RpcClient.Call("getblockcount")
	if err != nil {
		log.Info("Usdt getblockcount() failure! errinfo: %v,curblockcount is :%v", err,curblockcount)
		//退出进程
		//log.Fatal(err)
		return nil, err
	}
	log.Info("Usdt curblockcount(),getinfo: %v", curblockcount.Result)

	return cur,nil
}

func (cur *OmniUsdtRpcClient) GetHexMsg(txHex string) (*ltcutil.Tx, error) {
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
func Getinfo(){
    response, err := OmniRPCClient.RpcClient.Call("getinfo")
    fmt.Printf("--001---getinfo: ------001--get response info is--001: %v;;err is :%v\n", response,err)

} 
func (cur *OmniUsdtRpcClient) SendTransaction(txHexStr string) (txid string, err error) {

	cursendstr := txHexStr
	var rawTxHex string
	/*if txHexStr == "" {
		cursendstr = "0100000001c5935d1b263b3de546964253a0b34728d1d7601166077e6c2c12a5ee5ccffaac020000006b483045022100d5e973686a5d03cfe378cd62f4d623e8b70f23b58737b34d784777a4de746ee0022028826a107394b2c8adb41e22f2969c06fd31a64918073b36c1ee4ffc02b28fa9012102405835f8a07f7d08259b0caa18799ed36c3b7e98e86634cc3dece14939630b75feffffff029e700000000000001976a914e4b74f540b845ae81100e03fc64a41da721fb91a88aca20c0100000000001976a914660371326d3a2e064c278b20107a65dad847e8a988ac00000000"
	}
	*/
	//response, err := cur.RpcClient.Call("sendrawtransaction",curTx.MsgTx(),true)
	response, err := cur.RpcClient.Call("sendrawtransaction",cursendstr,true)
    log.Info("Invoke OmniUsdt's RPC sendrawtransaction(): get response info is: %v;;err is :%v\n", response,err)
	if err != nil {
		//log.Printf("SendRawTransaction() err happend!,chainHash is :%v,err is :%v", *chainHash, err)
		log.Error("SendRawTransaction() err happend!,txHexStr is:%s,err is :%v", txHexStr,err)
		return "", err
	}
	//sgj 0925 update add condi:
	//response.Error != nil
	if response != nil && response.Result !=nil {
		rawTxHex =response.Result.(string)
	}
	log.Info("SendRawTransaction() success!,created Txid is----->:%s", rawTxHex)
	//0424 ，把交易TXid写入数据库里，待推送；并发给java端
	//sgj 0813 update to remove
	//rawTxHex := fmt.Sprintf("%v",response)

	return rawTxHex, nil
}

func (cur *OmniUsdtRpcClient) RpcClose() {

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
