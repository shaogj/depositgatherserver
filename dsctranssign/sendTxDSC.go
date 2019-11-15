// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package dsctranssign

import (
	"bytes"
	"fmt"
	//"io/ioutil"
	"log"
	"time"

	//"github.com/davecgh/go-spew/spew"
	//sgj add:
	//"github.com/btcsuite/btcd/chaincfg/chainhash"
	"encoding/hex"
	"2019NNZXProj10/depositgatherserver/config"

	//"github.com/btcsuite/btcd/wire"
	mylog "github.com/mkideal/log"
	//sgj 0928
	//"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/wire"
	dscrpcclient "github.com/bcext/gcash/rpcclient"

)

//连接bitcoind的RPC客户端
var DSCClient DscRpcClient


//var m_RpcClient *rpcclient.Client
type BTCConf struct{
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
}

type DscRpcClient struct {
	rpcClient *dscrpcclient.Client
	config 	*config.DSCConf
}
func NewDSCRpcClient(conf *config.DSCConf) *DscRpcClient {
	c := &DscRpcClient{
		config: conf,
	}
	DSCClient = *c
	//go c.Run()
	return c
}
func (cur *DscRpcClient) GetHexMsg(txHex string) (*cashutil.Tx, error) {
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
	return cashutil.NewTx(&msgTx), nil

}
func (cur *DscRpcClient) Getnewaddress() string {
	//cur.rpcClient.GetNewaddress
	return "dfd-getnewaddress()"
}
func (cur *DscRpcClient) RpcConnect() (*DscRpcClient, error) {
	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := dscrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance cashutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}
	//节点上钱包进程的公钥：
	/*
	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := cashutil.AppDataDir("btcwallet", false)
	//sgj 0403 add:
	fmt.Printf("certHomeDir info is----------- :%s\n", certHomeDir)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	*/
	//0423 add NEWIP of cfg:
	connCfg := &dscrpcclient.ConnConfig{
		//Host: "127.0.0.1:8332",
		Host: cur.config.RPCHostPort,
		//Endpoint:     "ws",
		//sgj 0426 do:
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		//User:         "shaogj","123456"
		User:         cur.config.RPCUser,
		Pass:         cur.config.RPCPassWord,
		//sgj add 0423
		DisableTLS: true,
		//Certificates: certs,
	}
	fmt.Printf("exec connect()'s DscRpcClient!!!")
	client, err := dscrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
		mylog.Error("cur DSC RpcConnect() failed!, connCfg info is---- :%v;err info is :%v\n", connCfg,err)
		//sgj 0423 add
		return nil, err
	}else{
		//sgj 1017 add for dsc
		mylog.Error("cur DSC RpcConnect() succ!, get client is---- :%v;err info is :%v\n", client,err)
	}
	cur.rpcClient = client
	curblockinfo,err := client.GetBlockChainInfo()
	if err != nil {
		mylog.Info("BCH,GetBlockChainInfo ()errinfo: %v,err is :%v", curblockinfo,err)
		//退出进程
		//log.Fatal(err)
	}
	mylog.Info("BCH,GetBlockChainInfo (),info: %v", curblockinfo)
	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	blockCount, err := client.GetBlockCount()
	//unspent, err := client.ListUnspent()
	if err != nil {
		mylog.Info("blockCount (number): %d", blockCount)
		log.Fatal(err)
	}
	mylog.Info("BCH,blockCount (number): %d", blockCount)
	//sgj 1017 watching err:!
	//unknown address type
	/*
	getaddr,err := client.GetNewAddress("curaccount")
	if err != nil {
		fmt.Printf("GetNewAddress failed!,:getaddr is %s,err is :%s", getaddr,err)
		log.Fatal(err)
	}
	mylog.Info("GetNewAddress: %s", getaddr)
	*/
	return cur,nil
}

func (cur *DscRpcClient) SendTransaction(txHexStr string) (txid *chainhash.Hash, err error) {

	cursendstr := txHexStr
	if txHexStr == "" {
		return nil,err
	}
	//0409 add:
	curTx, err := cur.GetHexMsg(cursendstr)
	if err != nil {
		mylog.Error("getHexMsg msgTx info is :%v--err is :%v\n", curTx, err)
		return nil,err
	}
	log.Printf("getHexMsg () success!,curTx  is :%v", *curTx)
	chainHash, err := cur.rpcClient.SendRawTransaction(curTx.MsgTx(), true)
	if err != nil {
		//log.Printf("SendRawTransaction() err happend!,chainHash is :%v,err is :%v", *chainHash, err)
		mylog.Error("SendRawTransaction() err happend!,err is :%v", err)
		return nil, err
	}
	mylog.Info("SendRawTransaction() success!,BCH created Txid is----->:%s", chainHash.String())
	//return string([]byte(*chainHash)), nil
	//0424 ，把交易TXid写入数据库里，待推送；并发给java端
	return chainHash, nil
}

func (cur *DscRpcClient) RpcClose() {

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 10 seconds...")
	time.AfterFunc(time.Second*4, func() {
		log.Println("Client shutting down...")
		cur.rpcClient.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	cur.rpcClient.WaitForShutdown()
}
