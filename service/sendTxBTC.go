// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package service

import (
	//"backend/services/base/settlecenter/worker/wdcproto"
	//"backend/services/base/settlecenter/worker/wdcproto"
	"bytes"
	"fmt"
	"time"

	//"io/ioutil"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"

	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"2019NNZXProj10/depositgatherserver/config"
	"github.com/btcsuite/btcd/wire"
	mylog "github.com/mkideal/log"
	"encoding/json"
	//1211a dd
	"github.com/btcsuite/btcd/btcjson"


)

//连接bitcoind的RPC客户端
var BtcClient BtcRpcClient


//var m_RpcClient *rpcclient.Client
type BTCConf struct{
	RPCPort		int
	RPCHostPort		string
	RPCUser		string
	RPCPassWord		string
}

type BtcRpcClient struct {
	rpcClient *rpcclient.Client
	//config 	*BTCConf
	config 	*config.USDTConf

}
func NewBtcRPCClient(conf *config.USDTConf) *BtcRpcClient {
	c := &BtcRpcClient{
		config: conf,
	}
	BtcClient = *c
	//go c.Run()
	return c
}
func (cur *BtcRpcClient) GetHexMsg(txHex string) (*btcutil.Tx, error) {
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
	return btcutil.NewTx(&msgTx), nil

}
func (cur *BtcRpcClient) RpcConnect() (*BtcRpcClient, error) {
	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}
	//节点上钱包进程的公钥：
	/*
	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := btcutil.AppDataDir("btcwallet", false)
	//sgj 0403 add:
	fmt.Printf("certHomeDir info is----------- :%s\n", certHomeDir)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	*/
	//0423 add NEWIP of cfg:
	connCfg := &rpcclient.ConnConfig{
		//Host: "127.0.0.1:8332",
		Host: cur.config.RPCHostPort,
		//Endpoint:     "ws",
		//sgj 0426 do:
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		//User:         "shaogj2017","123456"
		User:         cur.config.RPCUser,
		Pass:         cur.config.RPCPassWord,
		//sgj add 0423
		DisableTLS: true,
		//Certificates: certs,
	}

	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("connCfg info is------i999999999----- :%v;err info is :%v\n", connCfg,err)
		//sgj 0423 add
		return nil, err
	}
	cur.rpcClient = client
	curblockinfo,err := client.GetBlockChainInfo()
	if err != nil {
		mylog.Info("BTC,GetBlockChainInfo ()errinfo: %v", curblockinfo)
		//退出进程
		//log.Fatal(err)
	}
	mylog.Info("BTC,GetBlockChainInfo (),getinfo: %v", curblockinfo)
	// Get the list of unspent transaction outputs (utxos) that the
	// connected wallet has at least one private key for.
	blockCount, err := client.GetBlockCount()
	//unspent, err := client.ListUnspent()
	if err != nil {
		mylog.Info("blockCount (number): %d", blockCount)
		log.Fatal(err)
	}
	mylog.Info("BTC,blockCount (number): %d", blockCount)

	return cur,nil
}

//1211add
func (cur *BtcRpcClient) GetBlockRange() (int64, error) {
	end, err := cur.rpcClient.GetBlockCount()
	if err != nil {
		mylog.Error("BTC get total block num error:%v", err)
		return 0, err
	}
	mylog.Info("BTC,blockCount (number): %d", end)
	return end,nil

}

func (cur *BtcRpcClient) GetBlock(hashChainStr string) (*wire.MsgBlock, error) {
	//blockHash *chainhash.Hash
	KTChash,err := chainhash.NewHash([]byte(hashChainStr))
	//(*btcutil.WIF)(wif1).String()
	//end, err := cur.rpcClient.GetBlock(&(*chainhash.Hash)KTChash)

	blockTrans, err := cur.rpcClient.GetBlock(KTChash)
	if err != nil {
		mylog.Error("BTC get total block error:%v", err)
		return nil, err
	}
	//mylog.Info("BTC,get block (): %d", end)
	return blockTrans,nil

}

//1211add 2:
func (cur *BtcRpcClient) call(method string, reply interface{}, v ...interface{}) error {
	rawParams := make([]json.RawMessage, 0, len(v))
	for _, param := range v {
		marshalledParam, err := json.Marshal(param)
		if err != nil {
			return err
		}
		rawMessage := json.RawMessage(marshalledParam)
		rawParams = append(rawParams, rawMessage)
	}
	result, err := cur.rpcClient.RawRequest(method, rawParams)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, reply)
}

//1,struct:
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Vin struct {
	Coinbase  string     `json:"coinbase"`
	Txid      string     `json:"txid"`
	Vout      uint32     `json:"vout"`
	ScriptSig *ScriptSig `json:"scriptSig"`
	Sequence  uint32     `json:"sequence"`
	//sgj 1211 add
	Txinwitness      []string     `json:"txinwitness"`
}
//btcjson.Vout{}

type TxRawResult struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Hash          string `json:"hash,omitempty"`
	Size          int32  `json:"size,omitempty"`
	//sgj 1211 add
	VSize          int32         `json:"vsize"`
	Weight          int32         `json:"weight"`

	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []btcjson.Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

type GetBlockVerboseResult struct {
	Hash          string        `json:"hash"`
	Confirmations int64         `json:"confirmations"`
	Size          int32         `json:"size"`
	//Txid          string        `json:"txid"`
	//sgj 1211 add
	Strippedsize          int32         `json:"strippedsize"`
	Weight          int32         `json:"weight"`

	Height        int64         `json:"height"`
	Version       int32         `json:"version"`

	VersionHex    string        `json:"versionHex"`
	MerkleRoot    string        `json:"merkleroot"`
	Tx            []TxRawResult      `json:"tx,omitempty"`
	//RawTx         []TxRawResult `json:"rawtx,omitempty"`
	Time          int64         `json:"time"`
	Mediantime          int64         `json:"mediantime"`

	Nonce         uint32        `json:"nonce"`
	Bits          string        `json:"bits"`
	Difficulty    float64       `json:"difficulty"`
	//sgj 1211 add

	Chainwork          string        `json:"chainwork"`
	NTx          int32        `json:"nTx"`

	PreviousHash  string        `json:"previousblockhash"`
	NextHash      string        `json:"nextblockhash,omitempty"`
	CoinbaseTx          interface{}        `json:"coinbaseTx"`

}


/*btc112struct

type GetBlockVerboseResult struct {
	Hash          string        `json:"hash"`
	Confirmations int64         `json:"confirmations"`
	Size          int32         `json:"size"`
	//Txid          string        `json:"txid"`
	Height        int64         `json:"height"`
	Version       int32         `json:"version"`

	VersionHex    string        `json:"versionHex"`
	MerkleRoot    string        `json:"merkleroot"`
	Tx            []string      `json:"tx,omitempty"`
	RawTx         []TxRawResult `json:"rawtx,omitempty"`
	Time          int64         `json:"time"`
	Nonce         uint32        `json:"nonce"`
	Bits          string        `json:"bits"`
	Difficulty    float64       `json:"difficulty"`
	PreviousHash  string        `json:"previousblockhash"`
	NextHash      string        `json:"nextblockhash,omitempty"`
}

 */

//2GOOD INFO:
func (cur *BtcRpcClient) GetBlockVerbose(hashChainStr string) (*GetBlockVerboseResult,error){//(*wire.MsgBlock, error) {
	//blockHash *chainhash.Hash
	//KTChash,err := chainhash.NewHash([]byte(hashChainStr))
	//(*btcutil.WIF)(wif1).String()
	//end, err := cur.rpcClient.GetBlock(&(*chainhash.Hash)KTChash)
	//var tx domain.OmniTx
	//var tx interface{}
	var err error
	//omni_gettransaction
	//if err = cur.call("getblock", &tx, hashChainStr,2); err == nil {
	//1211 add:
	getBlockTrans := GetBlockVerboseResult{}
	if err = cur.call("getblock", &getBlockTrans, hashChainStr,2); err != nil {
		mylog.Error("BTC get GetBlockVerbose,return getBlockTrans data is,err! getBlockTrans is:%v, error:%v", getBlockTrans,err)	//tx
		return nil,err
	}
	//mylog.Info("BTC get GetBlockVerbose,return tx is:%v, error:%v", tx,err)


	//---ing-skiplog--mylog.Info("BTC get GetBlockVerbose Good!!,return getBlockTrans is:%v, error:%v", getBlockTrans,err)


	//blockTrans
	//tx, err := cur.rpcClient.GetBlock(KTChash)
	//mylog.Info("BTC,get block (): %d", end)
	//return blockTrans,nil
	return &getBlockTrans,nil

}
//3) add
func (cur *BtcRpcClient) VinNew(txVerbose *TxRawResult) []string {
	var hashChainStrTxin string
	addr := make([]string, 0)
	var ret []string
	for _, v := range txVerbose.Vin {
		hashChainStrTxin = v.Txid
		if hashChainStrTxin == "" {
			//conbase,txid = ""
			return ret
		}
		getTransInfo := TxRawResult{}
		//getblock
		var err error
		if err = cur.call("getrawtransaction", &getTransInfo, hashChainStrTxin, true); err != nil {
			mylog.Error("KTC get GetBlockVerbose,return getTransInfo is:%v, error:%v", getTransInfo, err) //tx
		}
		mylog.Info("KTC get getrawtransaction Good!!,return getTransInfo is:%v, error:%v", getTransInfo, err)

		addr = append(addr, getTransInfo.Vout[v.Vout].ScriptPubKey.Addresses...)
	}
	/*
	for _, v := range t.TxIn {
		result, err := d.KtcRpcClient.GetRawTransactionVerbose(&v.PreviousOutPoint.Hash)
		if err != nil {
			loggers.Warn.Printf("Ktc get %s prevout:%v", t.TxHash(), err)
			continue
		}
		addr = append(addr, result.Vout[v.PreviousOutPoint.Index].ScriptPubKey.Addresses...)
	}
	*/
	return addr
}

/**/
func (cur *BtcRpcClient) SendTransaction(txHexStr string) (txid *chainhash.Hash, err error) {

	cursendstr := txHexStr
	if txHexStr == "" {
		cursendstr = "0100000001c5935d1b263b3de546964253a0b34728d1d7601166077e6c2c12a5ee5ccffaac020000006b483045022100d5e973686a5d03cfe378cd62f4d623e8b70f23b58737b34d784777a4de746ee0022028826a107394b2c8adb41e22f2969c06fd31a64918073b36c1ee4ffc02b28fa9012102405835f8a07f7d08259b0caa18799ed36c3b7e98e86634cc3dece14939630b75feffffff029e700000000000001976a914e4b74f540b845ae81100e03fc64a41da721fb91a88aca20c0100000000001976a914660371326d3a2e064c278b20107a65dad847e8a988ac00000000"
	}
	//0409 add:
	curTx, err := cur.GetHexMsg(cursendstr)
	if err != nil {
		//log.Fatal(err)
		mylog.Error("getHexMsg msgTx info is :%v--err is :%v\n", curTx, err)
		return nil,err
	}
	log.Printf("getHexMsg () success!,curTx  is :%v", *curTx)
	chainHash, err := cur.rpcClient.SendRawTransaction(curTx.MsgTx(), true)
	if err != nil {
		//sgj 0502 remove:	log.Fatal(err)
		//log.Printf("SendRawTransaction() err happend!,chainHash is :%v,err is :%v", *chainHash, err)
		mylog.Error("SendRawTransaction() err happend!,err is :%v", err)
		return nil, err
	}
	mylog.Info("SendRawTransaction() success!,created Txid is----->:%s", chainHash.String())
	//return string([]byte(*chainHash)), nil
	//0424 ，把交易TXid写入数据库里，待推送；并发给java端
	return chainHash, nil
}

func (cur *BtcRpcClient) RpcClose() {

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
