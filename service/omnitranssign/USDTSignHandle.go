//package usdtsign
package omnitranssign

import (
	"bytes"
	"encoding/hex"
	"fmt"
	//"time"
	"strconv"
	// "log"
	// "net/rpc"
	//"os"
	// "github.com/ybbus/jsonrpc"
	// "encoding/base64"
	"github.com/icloudland/btcdx/omnijson"
	//omnirpcclient "github.com/icloudland/btcdx/rpcclient"
	"2019NNZXProj10/depositgatherserver/proto"
	//sgj 0612 add:
	"2019NNZXProj10/depositgatherserver/service/omnitranssign/omnirpc"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	//"log"
	"2019NNZXProj10/depositgatherserver/netparams"

	//"github.com/mkideal/log"
	//"errors"
	"github.com/go-spew/spew"
	//"encoding/json"
	"github.com/mkideal/log"

	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcwallet/waddrmgr"
	// "github.com/btcsuite/btcwallet/wallet/txrules"
	"strings"
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	//1017 add:
	"2019NNZXProj10/depositgatherserver/service/wdctranssign"


)

var activeNet = &netparams.TestNet3Params

//sgj1120adding:
//var UtxoRPCClient = new(wdctranssign.WdcRpcClient)

var (
	Omni_MOrmEngine *xorm.Engine = &xorm.Engine{}
)
type ErrorCode int

const (
	maxProtocolVersion = 70002
)

var m_curpkScript []byte
var m_txHex string

//sgj 0608 add:
//after: omni_createrawtx_reference,
//omni_createrawtx_change
//get m_txUSDTHex
var m_txUSDTHex string

var m_address btcutil.Address

const saltSize = 32

type DeserializationError struct {
	error
}
type InvalidParameterError struct {
	error
}
type response struct {
	result []byte
	err    error
}

func defaultInt(ptr *int, dft int) {
	if *ptr == 0 {
		*ptr = dft
	}
}

func defaultString(ptr *string, dft string) {
	if len(*ptr) == 0 {
		*ptr = dft
	}
}

var G_UsdtSignHandle = UsdtSignHandle{}

type UsdtSignHandle struct {
	omniserver string
}

//sgj 0524 add:
func InitUSDTNet(usdtTestNet3 int){
	if usdtTestNet3 == 1{
		activeNet = &netparams.TestNet3Params
	}else{
		activeNet = &netparams.MainNetParams
	}
	//fmt.Printf("in InitLTCNet(),cur ltc parames activeNet is :%v \n", activeNet)
	
}

// ChangeSource provides P2PKH change output scripts for transaction creation.
type ChangeSource func() ([]byte, error)

// rpcDecodeHexError is a convenience function for returning a nicely formatted
// RPC error which indicates the provided hex string failed to decode.
func (ser *UsdtSignHandle) rpcDecodeHexError(gotHex string) *btcjson.RPCError {
	return btcjson.NewRPCError(btcjson.ErrRPCDecodeHexString,
		fmt.Sprintf("Argument must be hexadecimal string (not %q)",
			gotHex))
}

type SignatureError struct {
	InputIndex uint32
	Error      error
}

// messageToHex serializes a message to the wire protocol encoding using the
// latest protocol version and returns a hex-encoded string of the result.
func (ser *UsdtSignHandle) messageToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", ser.internalRPCError(err.Error(), context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

//获取USDT地址的私钥
func GetAddrPrivkeyUSDT(curaddress string) (addrPrikey string,err error){
	
	engineread:= Omni_MOrmEngine
	//get address 's [privkey]
	selectsql := "select * from  btc_account_key where address = '"  + curaddress + "'"
	addr_accountinfo, err := engineread.Query(selectsql)
	if err != nil || len(addr_accountinfo) <= 0{
		log.Info("when GetAddrPrivkey(),curaddress' is:%s ,get privkey error: %v", curaddress,err)
		return "",err
	}

	curaddrprivkey := string(addr_accountinfo[0]["priv_key"])
	log.Info("when GetAddrPrivkey(),curaddress is :%s,get addr's privkey succ ,info is: %v", curaddress,curaddrprivkey)
	return curaddrprivkey,nil
}

//请求数字币,USDT of btc layer，未花费的排序过的balance,按ut.value 倒排
func GetAddrBalanceUnspentUSDT(curaddress string) (getaddrutxos []proto.AddrBalanceUnspent,status int,balance float64,err error){
	
	engineread:= Omni_MOrmEngine
	var totalbalance float64
	//同address包含的txid所属各amount数进行匹配及合并
	//comutxosql := "select ut.id,ut.addrdisp,ut.value,ut.tx_id,ut.vout,ut.txcurid,floor(ut.block_id/10000) as blockid ,generate_time from address_utxo ut,outputs o where ut.tx_id=floor(o.id/4096) and ut.value = o.value and o.tx_id is null and addrdisp= '"  + curaddress + "'" + " order by ut.value desc"
	//sgj 0830 updating:
	comutxosql := "select DISTINCT ut.id,ut.addrdisp,ut.value,ut.tx_id,ut.vout,ut.txcurid,ut.block_id from outputs_multy o,address_utxo ut where ut.id =floor(o.out_id/4096) and ut.value =o.value and addrdisp = '"  + curaddress + "'" + " and  o.txcurid < ' '" + " order by ut.value desc";

	addr_outputtxlist, err := engineread.Query(comutxosql)
	if err != nil {
		fmt.Println("when GetAddrUTXO(),get curaddress is:%v,error: %v", curaddress,err)
		return nil, proto.StatusDataSelectErr,0,err
	}
	if len(addr_outputtxlist) <= 0 {
		fmt.Println("when GetAddrUTXO(),get curaddress is:%v,no find relevent record!", curaddress)
		return nil, proto.StatusLackUTXO,0,nil
	}
	curAddrPayInfo := new(proto.AddrBalanceUnspent)

	//返回数组添加统计数目
	for i := 0; i < len(addr_outputtxlist); i++ {
		curAddrPayInfo.AddrDisp = string(addr_outputtxlist[i]["addrdisp"])
		curAddrPayInfo.TxidHex = string(addr_outputtxlist[i]["txcurid"])
		curAddrPayInfo.Vout, _ = strconv.Atoi(string(addr_outputtxlist[i]["vout"]))
		//sgj 0330 update amount:
		curAddrPayInfo.Amount, _ = strconv.ParseFloat(string(addr_outputtxlist[i]["value"]),64)
		//curAddrPayInfo.Amount = curAddrPayInfo.Amount / 100000000
		getaddrutxos = append(getaddrutxos, *curAddrPayInfo)
		totalbalance += curAddrPayInfo.Amount
	}
	return getaddrutxos,0,totalbalance,nil

}

func (ser *UsdtSignHandle) internalRPCError(errStr, context string) *btcjson.RPCError {
	logStr := errStr
	if context != "" {
		logStr = context + ": " + errStr
	}
	//sgj update 11.30--rpcsLog.Error(logStr)
	fmt.Printf(logStr)
	return btcjson.NewRPCError(btcjson.ErrRPCInternal.Code, errStr)
}

//sgj 0611 add:
//getopreturnval := getsplit3str(linehex," ")

func getsplit3str(originstr string, splitchr string) (descstr string) {
	cursplit := " "
	if len(splitchr) > 0 {
		cursplit = splitchr
	}
	colsvalue := strings.Split(originstr, cursplit)
	//if len(colsvalue) != 3 {
	//fmt.Printf("now!,0607---linehex index num is :%d;colsvalue value is:%s\n",len(colsvalue),colsvalue)
	fmt.Printf("now!,0607---linehex index num is :%d\n", len(colsvalue))
	//
	var getopreturnval string
	if len(colsvalue) == 3 {
		getopreturnval = colsvalue[0]
	} else if len(colsvalue) == 4 {
		getopreturnval = colsvalue[1]
	}
	// tmpnval2 := colsvalue[2]
	//fmt.Printf("succ!,0607---linehex ' get reference value->colsvalue[2] is :%s:\n",tmpnval3)
	return getopreturnval
}

//sgj 0104 add
//从地址生成脚本公钥的工具函数
func (ser *UsdtSignHandle) GenscriptPubKeyFormAddr(encodedAddr string) (string, error) {
	// Decode the provided address.
	testparams := activeNet.Params
	addr, err := btcutil.DecodeAddress(encodedAddr, testparams)
	if err != nil {
		return "", &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Invalid address or key: " + err.Error(),
		}
	}
	// Ensure the address is one of the supported types and that
	// the network encoded with the address matches the network the
	// server is currently on.
	switch addr.(type) {
	case *btcutil.AddressPubKeyHash:
	case *btcutil.AddressScriptHash:
	default:
		return "", &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Invalid address or key",
		}
	}
	if !addr.IsForNet(testparams) {
		return "", &btcjson.RPCError{
			Code: btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Invalid address: " + encodedAddr +
				" is for the wrong network",
		}
	}

	// Create a new script which pays to the provided address.
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		context := "Failed to generate pay-to-address script"
		return "", ser.internalRPCError(err.Error(), context)
	}
	//输出二进制：
	//spew.Dump(pkScript)
	cur_getpkScript := []byte(hex.EncodeToString(pkScript))
	fmt.Printf("NewHashFromStr(),make encodedAddr is :%s,get generated encode pkScript is-----A4 :%s \n", encodedAddr, cur_getpkScript)
	return string(cur_getpkScript), nil

}

// decodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when
// number of bytes when calling hex.Decode.
func (ser *UsdtSignHandle) decodeHexStr(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDecodeHexString,
			Message: "Hex string decode failed: " + err.Error(),
		}
	}
	return decoded, nil
}

//解析交易结构的工具函数2
func (ser *UsdtSignHandle) PrepareSignRawTransactionTx(tx *wire.MsgTx) string {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			//return newFutureError(err)
			fmt.Printf("PrepareSignRawTransactionTx(),err info is--%v \n", err)
		}
		//sgj add:
		//m_txHexEncode = string(buf.Bytes())

		txHex = hex.EncodeToString(buf.Bytes())
		//fmt.Printf("PrepareSignRawTransactionTx(),after Encode,txHex info is--%s \n", txHex)
		return txHex
	}
	return ""
}

//createrawtransaction,step1,创建交易信息串函数,
//A transaction without outputs is created as basis to attach Omni related outputs and change later. Note that you may create a transaction base by one or more calls of
//0607 for USDT,same to create_input func
//step 1,函数功能:创建交易结构(包括出和找零地址的公钥脚本)
func (ser *UsdtSignHandle) handleCreateRawTransaction(cmd interface{}) (interface{}, error) {
	c := cmd.(*btcjson.CreateRawTransactionCmd)

	// Validate the locktime, if given.
	if c.LockTime != nil &&
		(*c.LockTime < 0 || *c.LockTime > int64(wire.MaxTxInSequenceNum)) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Locktime out of range",
		}
	}

	// Add all transaction inputs to a new transaction after performing
	// some validity checks.
	mtx := wire.NewMsgTx(wire.TxVersion)
	for _, input := range c.Inputs {
		txHash, err := chainhash.NewHashFromStr(input.Txid)
		if err != nil {
			return nil, ser.rpcDecodeHexError(input.Txid)
		}
		//sgj add : 12. 01

		prevOut := wire.NewOutPoint(txHash, input.Vout)
		fmt.Printf("NewHashFromStr(),get txHash info is:%v; prevOut is :%v -----A1:  \n", txHash, prevOut)
		txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_DUP}, nil)
		//txIn := wire.NewTxIn(prevOut, []byte{}, nil)
		if c.LockTime != nil && *c.LockTime != 0 {
			txIn.Sequence = wire.MaxTxInSequenceNum - 1

			//script := []byte{txscript.OP_TRUE, txscript.OP_DUP,
			//txscript.OP_DROP,txscript.OP_EQUALVERIFY}
			script3 := []byte{0x04, 0x31, 0xdc, 0x00, 0x1b, 0x01, 0x62}
			txIn.SignatureScript = script3
			//SignatureScript
		}
		mtx.AddTxIn(txIn)

		//sgj add : 12. 01
		fmt.Printf("NewHashFromStr(),get txIn info is-----A2 :%v \n", *txIn)

	}

	// Add all transaction outputs to the transaction after performing
	// some validity checks.

	//sgj update 11.30
	//params := s.cfg.ChainParams
	//sgj 0109  交易包括找零至少有两个输出！
	params := activeNet.Params
	for encodedAddr, amount := range c.Amounts {
		// Ensure amount is in the valid range for monetary amounts.
		if amount <= 0 || amount > btcutil.MaxSatoshi {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCType,
				Message: "Invalid amount",
			}
		}

		// Decode the provided address.
		addr, err := btcutil.DecodeAddress(encodedAddr, params)
		if err != nil {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address or key: " + err.Error(),
			}
		}

		// Ensure the address is one of the supported types and that
		// the network encoded with the address matches the network the
		// server is currently on.
		switch addr.(type) {
		case *btcutil.AddressPubKeyHash:
		case *btcutil.AddressScriptHash:
		default:
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address or key",
			}
		}
		if !addr.IsForNet(params) {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address: " + encodedAddr +
					" is for the wrong network",
			}
		}

		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			context := "Failed to generate pay-to-address script"
			return nil, ser.internalRPCError(err.Error(), context)
		}

		// Convert the amount to satoshi.
		satoshi, err := btcutil.NewAmount(amount)
		if err != nil {
			context := "Failed to convert amount"
			return nil, ser.internalRPCError(err.Error(), context)
		}
		txOut := wire.NewTxOut(int64(satoshi), pkScript)
		mtx.AddTxOut(txOut)
		fmt.Printf("NewHashFromStr(),get pkScript is-----A3 :%v \n", pkScript)
		//输出二进制：
		spew.Dump(pkScript)
		//sgj add :0330,重要的模块变量m_txHex
		//		m_txHex = ser.PrepareSignRawTransactionTx(mtx)
		//0607 for usdt:
		log.Info("-----USDT---==Now handleCreateRawTransaction(), composed MsgTx struct,get cur m_txHex is -A4 :%s \n", m_txHex)
		//0423,最后的找零地址为自身
		m_curpkScript = []byte(hex.EncodeToString(pkScript))
	}
	//0607 for usdt:
	m_txHex = ser.PrepareSignRawTransactionTx(mtx)
	log.Info("---USDT---==Now handleCreateRawTransaction(), composed MsgTx struct,get all m_txHex is -A4 :%s \n", m_txHex)

	// Set the Locktime, if given.
	if c.LockTime != nil {
		mtx.LockTime = uint32(*c.LockTime)
	}

	//0608 usdt watching:
	fmt.Printf("Usdt,handleCreateRawTransactionm(),get tx info' mtx is :%v \n", mtx)

	mtxHex, err := ser.messageToHex(mtx)
	if err != nil {
		return nil, err
	}
	return mtxHex, nil
}

//签名Hex字符串的测试函数
func (ser *UsdtSignHandle) CheckSendHexTx(curHextx string) (curserializedTx string, err error) {

	hexStr := curHextx
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", ser.rpcDecodeHexError(hexStr)
	}
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return "", &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX decode failed: " + err.Error(),
		}
	}
	return string(serializedTx), err

}

//真正的签名交易处理
func (ser *UsdtSignHandle) SignTransaction(tx *wire.MsgTx, hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*btcutil.WIF,
	p2shRedeemScriptsByAddress map[string][]byte) ([]SignatureError, error) {

	var signErrors []SignatureError
	for i, txIn := range tx.TxIn {
		prevOutScript, ok := additionalPrevScripts[txIn.PreviousOutPoint]
		spew.Dump(prevOutScript)
		log.Info("SignTransaction() exec step 01,txIn.PreviousOutPoint is :%v;prevOutScript is :%s ; m_curpkScript is :%s,\n", txIn.PreviousOutPoint, prevOutScript, m_curpkScript)

		if !ok {
			prevHash := &txIn.PreviousOutPoint.Hash
			prevIndex := txIn.PreviousOutPoint.Index
			log.Info("address's ECPrivKey(),--watching--to get from Txstore-:get privstr' prevHash is :%v;prevIndex is :%d \n", prevHash, prevIndex)
			/*
				txDetails, err := w.TxStore.TxDetails(prevHash)
				if err != nil {
					return nil, fmt.Errorf("Cannot query previous transaction "+
						"details for %v: %v", txIn.PreviousOutPoint, err)
				}
				if txDetails == nil {
					return nil, fmt.Errorf("%v not found",
						txIn.PreviousOutPoint)
				}
				prevOutScript = txDetails.MsgTx.TxOut[prevIndex].PkScript
			*/
			//sgj 1229, 不该走，否则从库里取对应的前一个未花费的输出txid：
			prevOutScript = m_curpkScript
		}

		// Set up our callbacks that we pass to txscript so it can
		// look up the appropriate keys and scripts by address.
		getKey := txscript.KeyClosure(func(addr btcutil.Address) (
			*btcec.PrivateKey, bool, error) {

			//sgj shoule be::keys's map info is :map[12r9VvL8Fjy89iSfpdog9pZPCWrT7i4eNk:5JoBo2tpFc8Yg1wvzmYvMRQx3qpF8APwNhzaQLsherKyoSrxpaM],
			//sgj test good! 解锁使用权时的私钥
			log.Info("exec KeyClosure()--001A0, len(additionalKeysByAddress) is :%d,info is :%v\n", len(additionalKeysByAddress), additionalKeysByAddress)
			//log.Info("exec KeyClosure()--001, map additionalKeysByAddress is :%v",additionalKeysByAddress)
			if len(additionalKeysByAddress) != 0 {
				//sgj add 1206:
				for key, _ := range additionalKeysByAddress {
					btwif := additionalKeysByAddress[key]
					//sgj test good!
					addrStr := addr.EncodeAddress()
					log.Info("exec KeyClosure()--001A1, addr is :%v;addrStr is :%v,key is:%v\n", addr, addrStr, key)
					//log.Info("exec KeyClosure()--001B, key is:%v,btwif.PrivKey is :%v\n",key,btwif.PrivKey)

					// key is address's string mode; to get address
					receiver1, _ := btcutil.DecodeAddress(key, activeNet.Params)
					m_address = receiver1
					return btwif.PrivKey, btwif.CompressPubKey, nil
				}
				//sgj 0109,传人私钥，则后面不继续走；
				//闭包的addr值从哪传入：
				//addr is :14G7oyy61JGwtjNKnkr6xyJ8aZWkppzaE2
				log.Info("exec KeyClosure(),sec section to get PrivKey, addr is :%v", addr)
				addrStr := addr.EncodeAddress()
				wif, ok := additionalKeysByAddress[addrStr]
				log.Info("exec KeyClosure()--001B=sec section, wif.PrivKey is :%v,addrStr is :%v", wif.PrivKey, addrStr)
				if !ok {
					return nil, false, nil
					//sgj,,fmt.Printf("no key for address")
				}
				return wif.PrivKey, wif.CompressPubKey, nil
			}
			//update to package 1206
			//privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), k.key)
			privKey := &btcec.PrivateKey{}
			fmt.Printf("exec KeyClosure(), PrivateKey is :%v", privKey)
			return privKey, true, nil
		})
		//to check！！
		getScript := txscript.ScriptClosure(func(
			addr btcutil.Address) ([]byte, error) {
			// If keys were provided then we can only use the
			// redeem scripts provided with our inputs, too.
			var script []byte
			var ok bool
			if len(p2shRedeemScriptsByAddress) != 0 {
				//sgj add 1206:
				for curaddrstr, _ := range p2shRedeemScriptsByAddress {

					//addrStr := curaddr.EncodeAddress()
					script, ok = p2shRedeemScriptsByAddress[curaddrstr]
					//script, ok := p2shRedeemScriptsByAddress[addrStr]
					fmt.Printf("exec ScriptClosure()--0020, addrStr is :%s,p2shRedeemScriptsByAddress's script is :%v\n", curaddrstr, script)
					if !ok {
						return nil, errors.New("no script for " + "address")
					}
					break
				}
				fmt.Printf("exec ScriptClosure()--002b, script is :%s", script)
				return script, nil
			}

			fmt.Printf("exec ScriptClosure()--002c, cur m_address is :%v", m_address)

			sa, ok := m_address.(waddrmgr.ManagedScriptAddress)
			if !ok {
				//return nil, errors.New("no script for " +"address")
				log.Error("get ScriptClosure() error! no script for :%s", m_address)
				return nil, nil
			}

			//return []byte{}
			byteaa, err := sa.Script()
			fmt.Printf("exec ScriptClosure() --003, sa.Script() is :%s,err is :%v\n", byteaa, err)
			return sa.Script()
		})

		// SigHashSingle inputs can only be signed if there's a
		// corresponding output. However this could be already signed,
		// so we always verify the output.

		if (hashType&txscript.SigHashSingle) !=
			txscript.SigHashSingle || i < len(tx.TxOut) {
			//sgj --activeNet.Params
			//sgj print func's params: 12 06
			//fmt.Printf("before txscript.SignTxOutput(), param0,,SignatureScript is :%s",txIn.SignatureScript)
			fmt.Printf("two,txscript.SignTxOutput()==77, params===tx is :%v, idx is :%d,pkScript is :%s, hashType is :%v, kdb is :%v, getScript is :%v,txIn.SignatureScript is :%s\n", tx, i, prevOutScript, hashType, getKey, getScript, txIn.SignatureScript)

			//step 三，返回签名脚本
			//sgj 0330 当前生成的真正的签名脚本：
			script, err := txscript.SignTxOutput(activeNet.Params,
				tx, 0, prevOutScript, hashType, getKey,
				getScript, nil)
			// Failure to sign isn't an error, it just means that
			// the tx isn't complete.
			if err != nil {
				//sgj add 1207,test good!
				//fmt.Printf("exec txscript.SignTxOutput() err !, err info  ==0AA . is :%s",err)
				signErrors = append(signErrors, SignatureError{
					InputIndex: uint32(i),
					Error:      err,
				})
				continue
			}
			txIn.SignatureScript = script
			fmt.Printf("exec txscript.SignTxOutput() finished !, get txIn'sSignatureScript info is :%s", txIn.SignatureScript)
			//输出二进制：
			spew.Dump(script)
		}

		// Either it was already signed or we just signed it.
		// Find out if it is completely satisfied or still needs more.

		vm, err := txscript.NewEngine(prevOutScript, tx, 0,
			txscript.StandardVerifyFlags, nil, nil, 774)

		//getinfo scriptSig := tx.TxIn[txIdx].SignatureScript
		//sgj add,,test good 1207
		//fmt.Printf("exec txscript.NewEngine() err !, prevOutScript info is :%s,tx.TxIn[0].SignatureScript is :%v,,err info  ==0BA . is :%s",prevOutScript,tx.TxIn[0].SignatureScript,err)
		if err == nil {
			fmt.Printf("in while: NewEngine()'s 66666666KKKKKKK is succ done! \n")
			err = vm.Execute()
		}
		if err != nil {
			fmt.Println(err)
			//sgj add 1207
			fmt.Printf("exec txscript.NewEngine(), to Execute()err !, err info  ==0BB . is :%s", err)
			signErrors = append(signErrors, SignatureError{
				InputIndex: uint32(i),
				Error:      err,
			})
		}
		fmt.Println("==usdt===Transaction successfully signed\n")
	}

	return signErrors, nil
}

//处理签名交易的command.

func (ser *UsdtSignHandle) signRawTransaction(icmd interface{},execcompletedflag bool) (interface{}, error) {
	cmd := icmd.(*btcjson.SignRawTransactionCmd)

	//0727,,若Go2Sever，参数中得到m_txHex的值：
	if execcompletedflag == true {
		m_txHex = cmd.RawTx
	}	
	//serializedTx2, err := ser.decodeHexStr(m_txHex)
	//sgj 0612 update
	serializedTx2, err := ser.decodeHexStr(cmd.RawTx)

	fmt.Printf("--usdt0--watching==signRawTransaction() step 0,m_txHex info is :%s, \n", m_txHex)
	fmt.Printf("--usdt02--watching==signRawTransaction() step 0,cmd.RawTx info is :%s, \n", cmd.RawTx)

	if err != nil {
		return nil, err
	}
	fmt.Printf("--usdt--watching==signRawTransaction() step 1,serializedTx2 is :%s, \n", serializedTx2)

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(serializedTx2))
	if err != nil {
		e := errors.New("TX decode failed")
		return nil, DeserializationError{e}
	}

	var hashType txscript.SigHashType
	switch *cmd.Flags {
	case "ALL":
		hashType = txscript.SigHashAll
	case "NONE":
		hashType = txscript.SigHashNone
	case "SINGLE":
		hashType = txscript.SigHashSingle
	case "ALL|ANYONECANPAY":
		hashType = txscript.SigHashAll | txscript.SigHashAnyOneCanPay
	case "NONE|ANYONECANPAY":
		hashType = txscript.SigHashNone | txscript.SigHashAnyOneCanPay
	case "SINGLE|ANYONECANPAY":
		hashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	default:
		e := errors.New("Invalid sighash parameter")
		return nil, InvalidParameterError{e}
	}
	fmt.Printf("signRawTransaction() step 2,hashType is :%v, \n", hashType)

	// TODO: really we probably should look these up with btcd anyway to
	// make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []btcjson.RawTxInput
	if cmd.Inputs != nil {
		cmdInputs = *cmd.Inputs
	}
	log.Info("get cmdInputs len is :%d,info is :%v", len(cmdInputs), cmdInputs)

	for _, rti := range cmdInputs {
		//sgj --Txid,由上一个交易产生的交易ID：
		inputHash, err := chainhash.NewHashFromStr(rti.Txid)
		if err != nil {
			log.Error("get NewHashFromStr() error! cur rti.Txid is :%s, err info is :%v", rti.Txid, err)
			return nil, DeserializationError{err}
		}
		//从上个接收交易中对应出地址address;第二步从地址函数转换出ScriptPubKey(即为PreviousOutPoint的address)
		preScriptPubKey := rti.ScriptPubKey
		script, err := ser.decodeHexStr(preScriptPubKey)
		if err != nil {
			log.Error("get cmdInputs error!  cur preScriptPubKey is :%s,err info is :%v", preScriptPubKey, err)
			return nil, err
		}

		// redeemScript is only actually used iff the user provided
		// private keys. In which case, it is used to get the scripts
		// for signing. If the user did not provide keys then we always
		// get scripts from the wallet.
		// Empty strings are ok for this one and hex.DecodeString will
		// DTRT.
		if cmd.PrivKeys != nil && len(*cmd.PrivKeys) != 0 {
			redeemScript, err := ser.decodeHexStr(rti.RedeemScript)
			if err != nil {
				return nil, err
			}

			addr, err := btcutil.NewAddressScriptHash(redeemScript,
				activeNet.Params)
			if err != nil {
				return nil, DeserializationError{err}
			}
			scripts[addr.String()] = redeemScript
			fmt.Printf("signRawTransaction() ==0109==after step 2,get redeemScript info is :%s; AddressScriptHash is :%v \n", redeemScript, addr)

		}
		inputs[wire.OutPoint{
			Hash:  *inputHash,
			Index: rti.Vout,
		}] = script

		log.Info("signRawTransaction() step 3,script  info is :%s;scripts info is :%v \n", script, scripts)
		//sgj 12 .09:
		//输出二进制：
		//0109 good!!
		spew.Dump(script)
	}
	//sgj 0109 check should be right!!

	log.Info("signRawTransaction() step 3total,inputs struct info is :%v, \n", inputs)

	// Now we go and look for any inputs that we were not provided by
	// querying btcd with getrawtransaction. We queue up a bunch of async
	// requests and will wait for replies after we have checked the rest of
	// the arguments.
	//realy ,it's chan response type---1205,,sgj

	requested := make(map[wire.OutPoint]response)
	//sgj==0105 again to check!
	///	requested := make(map[wire.OutPoint]btcrpcclient.FutureGetTxOutResult)
	for _, txIn := range tx.TxIn {
		// Did we get this outpoint from the arguments?
		if _, ok := inputs[txIn.PreviousOutPoint]; ok {
			continue
		}

		// Asynchronously request the output script.
		//chainClient.||sgj update 1205
		requested[txIn.PreviousOutPoint] = response{}
		/*
			requested[txIn.PreviousOutPoint] = GetTxOutAsync(
				&txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index,
				true)
		*/
	}

	// Parse list of private keys, if present. If there are any keys here
	// they are the keys that we may use for signing. If empty we will
	// use any keys known to us already.

	//sgj add: get map address info to prikey of user;
	var keys map[string]*btcutil.WIF
	if cmd.PrivKeys != nil {
		keys = make(map[string]*btcutil.WIF)

		for _, key := range *cmd.PrivKeys {
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				return nil, DeserializationError{err}
			}

			if !wif.IsForNet(activeNet.Params) {
				s := "key network doesn't match wallet's"
				return nil, DeserializationError{errors.New(s)}
			}
			//sgj 0109 add good!!
			wif.CompressPubKey = true
			//good end!!

			addr, err := btcutil.NewAddressPubKey(wif.SerializePubKey(),
				activeNet.Params)
			if err != nil {
				return nil, DeserializationError{err}
			}
			keys[addr.EncodeAddress()] = wif
			//watch: right!!
			//5KD3sSntucZFRzDNUJusRjEhVDADir1xfDQqFxqoi7djDX5k81b
			log.Info("==sgjwatching==exec signRawTransaction()--000F2, generated correspingding [addr is :%v],it's cmd Params PrivKeys is:%v\n", addr.EncodeAddress(), key)

		}
	}

	// We have checked the rest of the args. now we can collect the async
	// txs. TODO: If we don't mind the possibility of wasting work we could
	// move waiting to the following loop and be slightly more asynchronous.

	//sgj 0421,此段code函数正常下不走
	for outPoint, _ := range requested {
		//调用decode输出的：scriptPubKey结构体---//testTx.TxOut[0].PkScript

		//m_ScriptPubKeyInfo.Vout[0].ScriptPubKey.Hex = string(m_curpkScript)
		script, err := hex.DecodeString(string(m_curpkScript))
		//script, err := hex.DecodeString(m_ScriptPubKeyInfo.Vout[0].ScriptPubKey.Hex)
		fmt.Printf("SignTransaction exec to map requested,,outPoint info is :%v,script is:%s \n", outPoint, script)
		if err != nil {
			return nil, err
		}
		inputs[outPoint] = script
	}

	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.

	//fmt.Printf("SignTransaction invoke,,----inputs's map info is :%v, \n", inputs)
	log.Info("SignTransaction invoke,,----keys's map info is :%v, \n", keys)

	//sgj 0109 真正的签名部分：！！
	if execcompletedflag == true {
		signErrs, err := ser.SignTransaction(&tx, hashType, inputs, keys, scripts)
		log.Info("SignTransaction finished!!==007,,result signErrs info is :%v,err is :%v \n", signErrs, err)
		if err != nil {
			return nil, err
		}
		//sgj 0612 add:
		log.Info("SignTransaction successed!!======0612gooding")
		var buf bytes.Buffer
		buf.Grow(tx.SerializeSize())

		curserializedTx, err := ser.CheckSendHexTx(string(buf.Bytes()))
		log.Info("SignTransaction CheckSendHexTx!!==008,,result curserializedTx is :%s,err is :%v \n", curserializedTx, err)
		//begin
		if err = tx.Serialize(&buf); err != nil {
			log.Error("SignTransaction Serialize() to panic! tx is :%v,err is :%v \n", tx, err)
			//panic(err)
		}
		//sgj 0113 add:
		Hexstrinfo := hex.EncodeToString(buf.Bytes())
		//fmt.Printf("SignTransaction finished!!==0072, tx struct to SendRawTransaction(),is :%v \n", tx)

		log.Info("SignTransaction tx's Hexstrinfo==sign last info is :%v,err is :%v \n", Hexstrinfo, err)
		signErrors := make([]btcjson.SignRawTransactionError, 0, len(signErrs))
		for _, e := range signErrs {
			input := tx.TxIn[e.InputIndex]
			signErrors = append(signErrors, btcjson.SignRawTransactionError{
				TxID:      input.PreviousOutPoint.Hash.String(),
				Vout:      input.PreviousOutPoint.Index,
				ScriptSig: hex.EncodeToString(input.SignatureScript),
				Sequence:  input.Sequence,
				Error:     e.Error.Error(),
			})
		}

		return btcjson.SignRawTransactionResult{
			Hex:      hex.EncodeToString(buf.Bytes()),
			Complete: len(signErrors) == 0,
			Errors:   signErrors,
		}, nil
	}else{
		//
		//cmdPartSignRes := icmd.(*btcjson.SignRawTransactionCmd)
		//sgj 0716 add:返回部分交易的结构数据
		var buf bytes.Buffer
		buf.Grow(tx.SerializeSize())
		if err = tx.Serialize(&buf); err != nil {
			log.Error("SignTransaction Serialize() to panic! tx is :%v,err is :%v \n", tx, err)
			//panic(err)
		}
		Hexstrinfo := hex.EncodeToString(buf.Bytes())
		//sgj 0920 fix usdt bug:
		//Hexstrinfo = curtxUSDTHex
		log.Info("Goserver1,SignTransaction part tx's Hexstrinfo==sign last info is :%v,err is :%v \n", Hexstrinfo, err)
		cmdPartSignRes := &btcjson.SignRawTransactionCmd{
			//in Goserver2.handle:
			//signErrs, err := ser.SignTransaction(&tx, hashType, inputs, keys, scripts)
			RawTx: Hexstrinfo,
			Inputs: cmd.Inputs,
			Flags: cmd.Flags,
		
		}
		return cmdPartSignRes,nil
	}

}

// 0607 sgj ,btcutil.Amount is nil
func (ser *UsdtSignHandle) MakeMyTransaction(inputs []btcjson.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {
	convertedAmts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmts[addr.String()] = amount.ToBTC()
	}
	cmd := btcjson.NewCreateRawTransactionCmd(inputs, convertedAmts, lockTime)

	_, err := ser.handleCreateRawTransaction(cmd)
	//log.Info("-----usdt2---==Now handleCreateRawTransaction(), composed MsgTx struct,get m_txHex is:%s \n", m_txHex)

	return nil, err
}

//进行Usdt的签名交易---->同理按BTC的进行签名
//func (ser *UsdtSignHandle) ExecHandleProc(totalTxInputs []*btcjson.RawTxInput)(txdata interface{}, status int, err error) {
func (ser *UsdtSignHandle) ExecSignHandleProc(fromAddr string,curtxUSDTHex string,totalInputs []*omnijson.PreTx,execcompletedflag bool)(txdata interface{}, status int, err error) {
	curRawTxInputs := []btcjson.RawTxInput{}

	for _, curusdttxitem := range totalInputs {
		curTxIn := btcjson.RawTxInput{
			Txid: curusdttxitem.Txid,
			Vout: curusdttxitem.Vout,
			//sgj 1209,pubscript is last UTXO's Txid's pubscript
			//0104,ScriptPubKey,由函数计算出
			ScriptPubKey: string(curusdttxitem.ScriptPubKey),
			//0607 no setvalue
		}
		curRawTxInputs = append(curRawTxInputs, curTxIn)
	}
	PriKeys := make([]string, 0)
	//从mysql里获取私钥：
	//0611,need to add init for M_OrmEngine
	err = nil
	//从mysql里获取私钥：
	if execcompletedflag == true {
		curPrikey, err := GetAddrPrivkeyUSDT(fromAddr)
		//没取到对应私钥：
		if curPrikey == "" || err != nil {
			log.Info("command %s ,exex GetAddrPrivkeyUSDT() failue! err is: %v \n", fromAddr, err)
			return nil, proto.StatusAccountPrikeyNotExisted, err
		}
		PriKeys = append(PriKeys, curPrikey)
	}else{
		PriKeys = append(PriKeys, "")
	}

	signTransFlag := "ALL"
	var cmd = &btcjson.SignRawTransactionCmd{}
	/*
	if execcompletedflag == true {
		cmd = btcjson.NewSignRawTransactionCmd(m_txHex, &curRawTxInputs, &PriKeys, &signTransFlag)
	}else{
		//sgj 0920 fix bug:
		cmd = btcjson.NewSignRawTransactionCmd(curtxUSDTHex, &curRawTxInputs,nil, &signTransFlag)
	}
	*/
	//sjg 0608 update:
	//form m_txHex to m_txUSDTHex
	cmd = btcjson.NewSignRawTransactionCmd(curtxUSDTHex, &curRawTxInputs, &PriKeys, &signTransFlag)

	//第三步，进行签名交易
	transRes, err := ser.signRawTransaction(cmd,execcompletedflag)
	log.Info("step 2，signRawTransaction() exec finished! info return is :%v,err is :%v \n", transRes, err)
	if err != nil {
		log.Error("Failed to sign transaction，err is %v\n", err)
		//continue
		return transRes, proto.StatusSignError, err
	}
	//signedTransaction, complete, err := rpcClient.SignRawTransaction(tx.Tx)
	//第四步，to发送交易：ds
	return transRes, proto.StatusSuccess, nil

}


//签名交易的处理流程控制
func (ser *UsdtSignHandle) GetVerifiedParams(fromAddrCur string,fromPrivKey string,toamountcur float64,gatherAddress string) (fromAddr,toAddr, remainAddr string,toamount float64,errcode int,err error) {

	log.Info("exec PaySignTransProc() step 1,fromAddrCur is : %v \n", fromAddrCur)

	fromAddr =fromAddrCur

	toAddr = gatherAddress
	//toamount1 := cursettle * 100000000
	////sgj 1121PMdoing:临时交易数据data：

	//toamount1,_:= cursettle.Vol.Float64()
	toamount1 :=toamountcur
	//sgj 0106doing
	//toamount1 = 0.14

	//toamount1 = toamount1 * 100000000
	remainAddr = fromAddr  //cursettle.FromAddress

	if len(fromAddr) < 32 || len(fromAddr) > 34 {
		log.Error("requset param fromAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", fromAddr, len(fromAddr))
		return "", "", "", 0, proto.StatusInvalidArgument, nil
	}

	if len(toAddr) < 32 || len(toAddr) > 34 {
		log.Error("requset param toAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", toAddr, len(toAddr))
		return "", "", "", 0, proto.StatusInvalidArgument, nil
	}
	//0330 to addr set is: toAddr:
	if len(remainAddr) < 32 || len(remainAddr) > 34 {
		log.Error("requset param RemainAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", remainAddr, len(remainAddr))
		return "", "", "", 0, proto.StatusInvalidArgument, nil
	}
	if toAddr == remainAddr {
		log.Error("requset param toAddr same to the RemainAddr is :%s,opearte is forbidden! \n", remainAddr)
		return "", "", "", 0, proto.StatusInvalidArgument, nil
		
	}
	return fromAddr,toAddr, remainAddr,toamount1,0,nil

}
//a := strconv.FormatFloat(10.010, 'f', -1, 64)
//输出：10.01
func FloatTostrwithprec(fv float64, prec int) string {
    return strconv.FormatFloat(fv, 'f', prec, 64)
}

//sgj 0116doing
var UtxoRPCClientUSDT = new(wdctranssign.WdcRpcClient)

//propertyid,资产ID
//fromAddr string,fromPrivKey string,toamount float64,gatherAddress string,curGatherLimit float64
func (ser *UsdtSignHandle) PaySignTransProc(propertyid int,fromAddr string,fromPrivKey string,toamount float64,gatherAddress string,curGatherLimit float64,execcompletedflag bool) (txdata interface{}, status int, err error) {

//func (ser *UsdtSignHandle) PaySignTransProc(propertyid int,curtransreq *proto.SignTransactionReq,cursettle proto.Settle,outAccountAddress string,execcompletedflag bool) (txdata interface{}, status int, err error) {

	fromAddr,toAddr, remainAddr,toamount,paramstatus,err := ser.GetVerifiedParams(fromAddr,fromPrivKey,toamount,gatherAddress)
	if paramstatus > 0 {
		return nil, proto.StatusInvalidArgument, nil
	}
	m_getCurPubKey, err := ser.GenscriptPubKeyFormAddr(fromAddr)
	log.Info("fromAddr is : %s ,ToAddr is :%s,,it's  m_getCurPrivKeuy info is:%s \n", fromAddr, toAddr, m_getCurPubKey)

	//2018---06.13 by sgj--USDT的amount只取小数点后一位
	stramount33usdt := FloatTostrwithprec(toamount,-1)
	log.Info("cur propertyid is :%d,usdt amount value, stramount33usdt val is %s\n", propertyid,toamount,stramount33usdt)
	rawGetHex := omnirpc.GetSimplesendPayload(propertyid,stramount33usdt)//"0.47"
	//String rawTxHex = String.format("00000000%08x%016x", currencyId.getValue(), amount.getWillets());
	//sgj 0727 update:
	log.Info("Invoke OmniRpc 1), GenscriptPubKeyFormAddr() info: propertyid is :%d,getrawGetHex is: %v\n", propertyid,rawGetHex)
	 /*
	 response, err := omnirpc.OmniRPCClient.RpcClient.Call("omni_createpayload_simplesend", 2, "15.7")
	 //String rawTxHex = String.format("00000000%08x%016x", currencyId.getValue(), amount.getWillets());
	 */

	addrUtxolist := make([]proto.AddrBalanceUnspent,3,8)
	getutxoinfo,utxonum,err := UtxoRPCClientUSDT.GetTxUnSpentLimit(fromAddr)

	if err != nil {
		log.Error("fromAddr %s ,exex GetTxUnSpentLimit() failue! err is: %v \n", fromAddr, err)
		return nil, status, err
	}

	log.Info("exec GetTxUnSpentLimit(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! unxonum is :%d\n", addrUtxolist, utxonum)
	//totalbalance is: %v ,getbalance
	if utxonum == 0 {
		log.Error("fromAddr %s ,exex GetTxUnSpentLimit() failue!,get addrUtxolist num is: %v \n", fromAddr, 0)
		return nil, status, err

	}
	//totalbalance is: %v ,getbalance
	if utxonum == 0 {
		log.Error("fromAddr %s ,exex GetAddressUtxo() failue!,get addrUtxolist num is: %v \n", fromAddr, 0)
		return nil, status, err

	}
	//11.20,请求比特币,未花费的排序过的balance,按ut.value 倒排
	selAmountIndex :=0
	selectAmountVal:= getutxoinfo[0].Value
	for iseno,curitem := range getutxoinfo{
		if curitem.Value > selectAmountVal{
			selectAmountVal = curitem.Value
			selAmountIndex = iseno
		}
	}
	log.Info("exec GetAddrUTXO(),form addrUtxolist,selAmountIndex is :%d，getutxoinfo is:%v\n",selAmountIndex,getutxoinfo[selAmountIndex])

	addrUtxolist[0].TxidHex = getutxoinfo[selAmountIndex].TxHashBigEndian
	addrUtxolist[0].Vout = int(getutxoinfo[selAmountIndex].TxOutputN)
	addrUtxolist[0].Amount = float64(getutxoinfo[selAmountIndex].Value)

	log.Info("get selAmountIndex is:%d,real account info: addrUtxolist[0].TxidHex is:%s, vout is :%d,amount is :%v,", selAmountIndex,addrUtxolist[0].TxidHex, addrUtxolist[0].Vout, addrUtxolist[0].Amount)

	getbalance := addrUtxolist[0].Amount
	log.Info("exec GetAddrUTXO(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! totalbalance is: %v \n", addrUtxolist, getbalance)

	//sgj 0402 if balance < amount1 ,then 返回余额不足；
	//curminamount := float64(toamount) + 1000
	//addres的BTC余额，必须大于需要的最低的BTC‘s fee：
	//0917,real need 0.00000546 BTC
	if getbalance < 0.00006{
	//if getbalance < curminamount{
		log.Error("real account vout getbalance is:%d, cur to amount is :%v,", getbalance,toamount)
		return nil, proto.StatusLackBalance, nil
	}
	//var irealamount float64
	balanceinfo, err := omnirpc.Getbalance_MP(fromAddr, propertyid)
	//balanceinfo, err := omnirpc.Getbalance_MP("myThsLLa2H9cQJ5JCD2kfxymMDoj8aFHTT", 1)
	if err != nil {
		//amount <=0 ,余额不足，返回；
		log.Error("omnirpc.Getbalance_MP(),err is====> :%v\n", err)
	} else {
		irealtokenamount,_ := strconv.ParseFloat(balanceinfo.Balance, 64)
		log.Info("omnirpc.Getbalance_MP() exec succ!,get info is :%v;irealtokenamount is :%d\n", balanceinfo,irealtokenamount)
		//sgj 0924 update,for USDT,txfee is 2,judge if lack:
		//sgj 1120 update minlimit to 0.3
		//tmp wathcing :
		//irealtokenamount = 3
		if irealtokenamount <= 0.3 {
			return nil, proto.StatusLackBalance, nil
		}
	}
	log.Info("Invoke OmniRpc 2), Getbalance_MP() info: propertyid is :%d,fromAddr is: %s\n", propertyid,fromAddr)

	curinput := []btcjson.TransactionInput{}
	//找出转出金额，所需要花费对应的Txid：map rec :Txid,Vout,amount1
	perInput := btcjson.TransactionInput{
		Txid: addrUtxolist[0].TxidHex,
		Vout: uint32(addrUtxolist[0].Vout),
	}
	curinput = append(curinput, perInput)

	//step 3),Createtransaction=no output,no amount
	msgTx, err := ser.MakeMyTransaction(curinput, nil, nil)
	log.Info("step3===>USDT, info,(msgTx) value is :%v,orign m_txHex----> info is :%v\n", msgTx, m_txHex)
	//sgj  0612,same to omni_createrawtx_input

	//m_txHex = "0100000002de95b97cf4c67ec01485fd698ec154a325ff69dd3e58435d7024bae7f69534c20000000000ffffffffb3b60aaa69b860c9bf31e742e3b37e75a2a553fd0bebf8aaf7da0e9bb07316ee0200000000ffffffff0000000000"
	payload := rawGetHex

	//sgj 0801 update
	getRawHex := omnirpc.Createrawtx_opreturn(m_txHex, payload)
	log.Info("step4==>Invoke OmniRpc 3),exec Createrawtx_opreturn(),req getRawHex info is :%v\n", getRawHex)
	//step 5)
	//申请测试币的充值地址：
	destination := toAddr

	optamount := 0.0009
	//optamount := toamount
	getNewRawHex := omnirpc.Createrawtx_Reference(getRawHex, destination, optamount)

	log.Info("step5===>Invoke OmniRpc 4),exec Createrawtx_Reference(),req getNewRawHex info is :%v\n", getNewRawHex)

	//1. 找出未花费的币（unspent output）;可能多个输入时，循环取utxo：
	curUstdRawTxInputs := []*omnijson.PreTx{}
	thePreAddr := fromAddr
	preScriptPubKey, err := ser.GenscriptPubKeyFormAddr(thePreAddr)
	//第二次，传给结构体with：Txid
	//curTxIn := btcjson.RawTxInput{
	curTxIn := omnijson.PreTx{
		Txid: perInput.Txid,
		Vout: perInput.Vout,
		ScriptPubKey: string(preScriptPubKey),
		//add for usdt:
		//Value: addrUtxolist[0].Amount,

		//sgj 0806 update Amouint Fix:
		Value: addrUtxolist[0].Amount / 100000000,
		//you need to specify txid, vout, scriptPubKey

	}
	log.Info("step 1 ext info,req curTxIn info is :%v", curTxIn)

	curUstdRawTxInputs = append(curUstdRawTxInputs, &curTxIn)

	//from step 5 and add a change output back to
	// sgj 0920 update to usdt:
	fee := 0.00006
	//sgj 20200114 update fee value
	fee = 0.000007
	fromAddr = remainAddr
	//getCurRawTxResp, err := omnirpc.OmniRPCClient.Call("omni_createrawtx_change", getNewRawHex, curUstdRawTxInputs, fromAddr, fee, 1)
	//sgj 1121 watching
	log.Info("step6 pre==>  to Invoke OmniRpc 5),exec Createrawtx_change(),fromAddr is:%s,toamount is:%f,cur curUstdRawTxInputs is:%v,getNewRawHex info is:%s,err is :%v", fromAddr,toamount,curUstdRawTxInputs,getNewRawHex,err)
	getCurRawTxResp, err := omnirpc.OmniRPCClient.RpcClient.Call("omni_createrawtx_change", getNewRawHex, curUstdRawTxInputs, fromAddr, fee, 1)
	if nil != err {
		log.Error("step6==>Invoke OmniRpc 5),exec Createrawtx_change() err!,req getCurRawTxResp info is:%s,err is :%v", getCurRawTxResp,err)
		return
		return nil, proto.StatusSignError, nil
	}
	getNewChangeRawTx :=getCurRawTxResp.Result.(string)
	//log.Info("step6======>ext Createrawtx_change() succ!,req getCurRawTx info is :%v", getCurRawTx)

	log.Info("step6==>>Invoke OmniRpc 5),exec omni_createrawtx_change(),req getNewChangeRawTx info is :%v\n", getNewChangeRawTx)
	//与BTC的签名完全一致；
	signInfoRes, status, err := ser.ExecSignHandleProc(fromAddr, getNewChangeRawTx,curUstdRawTxInputs,execcompletedflag)
	//log.Info("request Pay_SignTransaction() succ ! info is %v; status is :%d,err is :%v", signInfoRes, status, err)
	return signInfoRes, status, err
}

