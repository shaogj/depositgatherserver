// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//package Qtumtranssign
package qtumrpc

import (
	//"bytes"
	"fmt"
	//"io/ioutil"
	"time"

	"2019NNZXProj10/depositgatherserver/config"

	"encoding/base64"
	"github.com/mkideal/log"
	"github.com/ybbus/jsonrpc"
)

var QtumRPCClient QtumRpcClient

//连接bitcoind的RPC客户端
type BTCConf struct{
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
}

type QtumRpcClient struct {
	RpcClient jsonrpc.RPCClient
	config 	*config.DSCConf
}

//func NewQtumRpcClient(conf *LTCConf) *QtumRpcClient {
func NewQtumRpcClient(conf *config.DSCConf) *QtumRpcClient {
	c := &QtumRpcClient{
		config: conf,
	}
	QtumRPCClient = *c
	//go c.Run()
	return c
}

func (cur *QtumRpcClient) RpcConnect() (*QtumRpcClient, error) {

	//sgj 0611 add
	//QtumRPCClient = jsonrpc.NewClientWithOpts("http://192.168.10.232:8335", &jsonrpc.RPCClientOpts{
	//curRpcClient 
	fmt.Print("===0614=77777777777=before Qtum RpcConnect() ,RPC config info is %v\n", cur.config)
	cur.RpcClient = jsonrpc.NewClientWithOpts(cur.config.RPCHostPort, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			//"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("Qtumcorerpc"+":"+"123456")),
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(cur.config.RPCUser+":"+cur.config.RPCPassWord)),
			},
		})
		
	fmt.Print("==88888888==after Qtum RpcConnect() ,RpcClient info is %v", cur.RpcClient)
	//cur.RpcClient = curRpcClient
	response, err := cur.RpcClient.Call("getinfo")
    fmt.Printf("--001---getinfo: ------001--get response info is--001: %v;;err is :%v\n", response,err)

	curblockinfo,err := cur.RpcClient.Call("getblockchaininfo")
	if err != nil {
		log.Info("Qtum GetBlockChainInfo () failrue! errinfo: %v,curblockinfo is :%v", err,curblockinfo)
		//退出进程
		//log.Fatal(err)
		return nil, err
	}
	log.Info("Qtum GetBlockChainInfo (),getinfo: %v", curblockinfo)
	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	//
	//blockCount, err := client.GetBlockCount()

	return cur,nil
}

func Getinfo(){
    response, err := QtumRPCClient.RpcClient.Call("getinfo")
    fmt.Printf("--001---getinfo: ------001--get response info is--001: %v;;err is :%v\n", response,err)

}
func (cur *QtumRpcClient) SendTransaction(txHexStr string) (txid string, err error) {

	cursendstr := txHexStr
	var rawTxHex string
	/*if txHexStr == "" {
		cursendstr = "0100000001c5935d1b263b3de546964253a0b34728d1d7601166077e6c2c12a5ee5ccffaac020000006b483045022100d5e973686a5d03cfe378cd62f4d623e8b70f23b58737b34d784777a4de746ee0022028826a107394b2c8adb41e22f2969c06fd31a64918073b36c1ee4ffc02b28fa9012102405835f8a07f7d08259b0caa18799ed36c3b7e98e86634cc3dece14939630b75feffffff029e700000000000001976a914e4b74f540b845ae81100e03fc64a41da721fb91a88aca20c0100000000001976a914660371326d3a2e064c278b20107a65dad847e8a988ac00000000"
	}
	*/
	//response, err := cur.RpcClient.Call("sendrawtransaction",curTx.MsgTx(),true)
    log.Info("before QtumRpc: sendrawtransaction():params txHexStr info is: %v;err is :%v\n", txHexStr,err)
	response, err := cur.RpcClient.Call("sendrawtransaction",cursendstr,true)
    log.Info("--001---sendrawtransaction: ------001--get response info is--001: %v;;err is :%v\n", response,err)
	if err != nil {
		log.Info("SendRawTransaction() err happend!,err is :%v", err)
		return "", err
	}
	if response != nil {
		rawTxHex =response.Result.(string)
	}
	log.Info("SendRawTransaction() success!,created Txid is----->:%s", rawTxHex)
	//0424 ，把交易TXid写入数据库里，待推送；并发给java端
	//rawTxHex := fmt.Sprintf("%v",response)
	return rawTxHex, nil
}
//sgj 0911 add sign:signmessagewithprivkey
func (cur *QtumRpcClient) SignMessageWithPrivkey(privkey string,txHexStr string) (signedrawHex string, err error) {

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

func (cur *QtumRpcClient) RpcClose() {

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
