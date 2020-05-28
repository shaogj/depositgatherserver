//package btcsignlocal
package service

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"2019NNZXProj10/depositgatherserver/netparams"
	"errors"

	//"github.com/btcsuite/btcwallet/netparams"
	"github.com/davecgh/go-spew/spew"
	//add
	"2019NNZXProj10/depositgatherserver/proto"

	"2019NNZXProj10/depositgatherserver/service/wdctranssign"

	"github.com/btcsuite/btcwallet/waddrmgr"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/mkideal/log"
)

//0507testing:
//var activeNet = &netparams.MainNetParams

//sgj 0330 for testing net
var activeNet = &netparams.TestNet3Params

//sgj 1120add:

//sgj 11.19 skip gooding
//get request to get utxo of btc trans
//var UtxoRPCClient = new(wdctranssign.WdcRpcClient)

type ErrorCode int

const (
	maxProtocolVersion = 70002
)

var m_curpkScript []byte
var m_txHex string

var m_address btcutil.Address

//sgj 11.19 add:
var GBtcSignHandle BtcSignHandle

const saltSize = 32

type response struct {
	result []byte
	err    error
}
type SignatureError struct {
	InputIndex uint32
	Error      error
}

type InvalidParameterError struct {
	error
}

type DeserializationError struct {
	error
}

//sgj 0427 add;签名处理的实体类
type BtcSignHandle struct {
	GatherLimit float64
	//sgj 0116adding,总归集的地址数量
	GatherAddrCount int
}

//sgj 0524 add:
func InitBTCNet(lctTestNet3 int) {
	if lctTestNet3 == 1 {
		activeNet = &netparams.TestNet3Params
	} else {
		activeNet = &netparams.MainNetParams
	}
	fmt.Printf("in InitBTCNet(),cur btc activeNet params is :%v \n", activeNet)

}

//sgj add at 12 05, form :/go/src/github.com/btcsuite/btcwallet/internal/legacy/keystore/
// newScriptAddress initializes and returns a new P2SH address.
// iv must be 16 bytes, or nil (in which case it is randomly generated).

// ChangeSource provides P2PKH change output scripts for transaction creation.
type ChangeSource func() ([]byte, error)

// rpcDecodeHexError is a convenience function for returning a nicely formatted
// RPC error which indicates the provided hex string failed to decode.
func (ser *BtcSignHandle) rpcDecodeHexError(gotHex string) *btcjson.RPCError {
	return btcjson.NewRPCError(btcjson.ErrRPCDecodeHexString,
		fmt.Sprintf("Argument must be hexadecimal string (not %q)",
			gotHex))
}

func (ser *BtcSignHandle) internalRPCError(errStr, context string) *btcjson.RPCError {
	logStr := errStr
	if context != "" {
		logStr = context + ": " + errStr
	}
	//sgj update 11.30--rpcsLog.Error(logStr)
	fmt.Printf(logStr)
	return btcjson.NewRPCError(btcjson.ErrRPCInternal.Code, errStr)
}

// messageToHex serializes a message to the wire protocol encoding using the
// latest protocol version and returns a hex-encoded string of the result.
func (ser *BtcSignHandle) messageToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", ser.internalRPCError(err.Error(), context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

//sgj 0104 add
//从地址生成脚本公钥的工具函数
func (ser *BtcSignHandle) GenscriptPubKeyFormAddr(encodedAddr string) (string, error) {
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

//解析交易结构的工具函数2
func (ser *BtcSignHandle) PrepareSignRawTransactionTx(tx *wire.MsgTx) string {
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

// handleCreateRawTransaction handles createrawtransaction commands.
//2017.org params 11.29:
//s *rpcServer, , closeChan <-chan struct{}

//step 1,函数功能:创建交易结构(包括出和找零地址的公钥脚本)
func (ser *BtcSignHandle) handleCreateRawTransaction(cmd interface{}) (interface{}, error) {
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
		m_txHex = ser.PrepareSignRawTransactionTx(mtx)
		fmt.Printf("Now handleCreateRawTransaction(), composed MsgTx struct,get m_txHex is-----A4 :%s \n", m_txHex)
		//0423,最后的找零地址为自身
		m_curpkScript = []byte(hex.EncodeToString(pkScript))
	}

	// Set the Locktime, if given.
	if c.LockTime != nil {
		mtx.LockTime = uint32(*c.LockTime)
	}

	// Return the serialized and hex-encoded transaction.  Note that this
	// is intentionally not directly returning because the first return
	// value is a string and it would result in returning an empty string to
	// the client instead of nothing (nil) in the case of an error.

	//from tx_test.go--A:
	//基本不用签名设置的脚本值
	fmt.Printf("handleCreateRawTransactionm(),get tx info' SignatureScript is-2 :%v \n", mtx.TxIn[0].SignatureScript)
	//fmt.Printf("handleCreateRawTransactionm(),get tx info ' PreviousOutPoint.Hash is--3 :%v \n", mtx.TxIn[1].PreviousOutPoint.Hash)

	mtxHex, err := ser.messageToHex(mtx)
	if err != nil {
		return nil, err
	}
	return mtxHex, nil
}

//step1,创建交易信息串函数
func (ser *BtcSignHandle) MakeMyTransaction(inputs []btcjson.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {
	convertedAmts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmts[addr.String()] = amount.ToBTC()
	}
	cmd := btcjson.NewCreateRawTransactionCmd(inputs, convertedAmts, lockTime)

	_, err := ser.handleCreateRawTransaction(cmd)
	//TxIn.SignatureScript
	//TxIn.PreviousOutPoint.Hash
	fmt.Printf("step first,,exec MakeMyTransaction() finished !(build pubkeyscript finish)! \n")

	return nil, err
}

//真正的签名交易处理
func (ser *BtcSignHandle) SignTransaction(tx *wire.MsgTx, hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*btcutil.WIF,
	p2shRedeemScriptsByAddress map[string][]byte) ([]SignatureError, error) {

	var signErrors []SignatureError
	//sgj add 0717 update:
	//fromAddr :=getdbfromaddr
	fromAddr := "getdbfromaddrstr"
	getCurpkScript, _ := ser.GenscriptPubKeyFormAddr(fromAddr)
	m_curpkScript := []byte(getCurpkScript)
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
			//sgj 0112 地址管理，查询DB,PubKeyAddress--to get privkey
			/*==sgh 1206 add

			 */

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
		fmt.Println("Transaction successfully signed")
	}

	return signErrors, nil
}

// decodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when
// number of bytes when calling hex.Decode.
func (ser *BtcSignHandle) decodeHexStr(hexStr string) ([]byte, error) {
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

//签名Hex字符串的测试函数
func (ser *BtcSignHandle) CheckSendHexTx(curHextx string) (curserializedTx string, err error) {

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

//处理签名交易的command.
func (ser *BtcSignHandle) signRawTransaction(icmd interface{}, execcompletedflag bool) (interface{}, error) {
	cmd := icmd.(*btcjson.SignRawTransactionCmd)

	//sgj 0720 add:
	//若Go2Sever，参数中得到m_txHex的值：
	if execcompletedflag == true {
		m_txHex = cmd.RawTx
	}
	serializedTx2, err := ser.decodeHexStr(m_txHex)
	if err != nil {
		return nil, err
	}
	log.Info("execcompletedflag is :%v,signRawTransaction() step 1,serializedTx2 is :%s, \n", execcompletedflag, serializedTx2)

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
			log.Info("signRawTransaction() ==0109==after step 2,get redeemScript info is :%s; AddressScriptHash is :%v \n", redeemScript, addr)

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
	//sgj 0719 add
	log.Info("cur params's cmd.PrivKeys info is :%v, \n", cmd.PrivKeys)

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
	log.Info("signRawTransaction() -----datawatching ------0717====1")

	if cmd.PrivKeys != nil {
		keys = make(map[string]*btcutil.WIF)

		for _, key := range *cmd.PrivKeys {
			wif, err := btcutil.DecodeWIF(key)
			if err != nil {
				return nil, DeserializationError{err}
			}
			log.Info("signRawTransaction() -----datawatching ------0717====2")

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
	log.Info("signRawTransaction() -----datawatching ------0717====3")

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
	} else {
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
		log.Info("Goserver1,SignTransaction part tx's Hexstrinfo==sign last info is :%v,err is :%v \n", Hexstrinfo, err)
		cmdPartSignRes := &btcjson.SignRawTransactionCmd{
			//in Goserver2.handle:
			//signErrs, err := ser.SignTransaction(&tx, hashType, inputs, keys, scripts)
			RawTx:  Hexstrinfo,
			Inputs: cmd.Inputs,
			Flags:  cmd.Flags,
		}
		return cmdPartSignRes, nil
	}
}

//Status
//签名交易的处理流程控制
func (ser *BtcSignHandle) PaySignTransProc(fromAddr string, fromPrivKey string, toamount float64, gatherAddress string, curGatherLimit float64) (txdata interface{}, gatheredamount1 float64, status int, err error) {
	//func (ser *BtcSignHandle) PaySignTransProc(curtransreq *proto.SignTransactionReq, cursettle proto.Settle,outAccountAddress string,execcompletedflag bool) (txdata interface{}, status int, err error) {

	//sgj 1120doing
	log.Info("exec PaySignTransProc() step 1,curtransreq of fromAddr is : %v \n", fromAddr)
	//sgj 1217,,toamount 为归集的所用数量，，utxo最大的去掉fee，即为要规的curtoamount的值
	ser.GatherLimit = curGatherLimit

	//good!,从地址生成公钥！
	//步骤记录： 1）获取自己地址的公钥
	//sgj 1120doing
	log.Info("exec PaySignTransProc() step 1,fromAddr is : %s, toamount is:%f,GatherLimit is:%f \n", fromAddr, toamount, ser.GatherLimit)

	toAddr := gatherAddress
	//toamount1 := cursettle * 100000000
	//toamount1,_:= cursettle.Vol.Float64()
	toamount1 := toamount
	//1210trying

	//toamount1 = 0.041
	remainAddr := fromAddr //cursettle.FromAddress

	//end 1120adding
	m_getCurPrivKeuy, err := ser.GenscriptPubKeyFormAddr(fromAddr)
	log.Info("fromAddr is : %s ,ToAddr is :%s,,it's  m_getCurPrivKeuy info is:%s \n", fromAddr, toAddr, m_getCurPrivKeuy)
	if len(fromAddr) < 32 || len(fromAddr) > 34 {
		log.Error("requset param fromAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", fromAddr, len(fromAddr))
		return nil, 0, proto.StatusInvalidArgument, nil

	}

	if len(toAddr) < 32 || len(toAddr) > 34 {
		log.Error("requset param toAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", toAddr, len(toAddr))
		return nil, 0, proto.StatusInvalidArgument, nil
	}
	//0330 to addr set is: toAddr:
	if len(remainAddr) < 32 || len(remainAddr) > 34 {
		log.Error("requset param RemainAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", remainAddr, len(remainAddr))
		return nil, 0, proto.StatusInvalidArgument, nil
	}
	if toAddr == remainAddr {
		log.Error("requset param toAddr same to the RemainAddr is :%s,opearte is forbidden! \n", remainAddr)
		return nil, 0, proto.StatusInvalidArgument, nil

	}
	//1. 找出未花费的币（unspent output）;可能多个输入时，循环取utxo：
	//0329 update,	//GetAddrUTXO(fromAddr)
	//testing from is: 1AJQ3jXhUF8WiisEcuVd8Xmfq4QJ7n1SdL
	/*
		addrUtxolist, status, getbalance, err := GetAddrBalanceUnspent(fromAddr)

		if err != nil {
			log.Error("command %s ,exex GetAddressUtxo() failue! err is: %v \n", fromAddr, err)
			return nil, status, err
		}
		if status != 0 {
			log.Error("from address :%s ,to address :%s ,exec GetAddressUtxo() failure! no result info!,status is:%d \n", fromAddr, toAddr,status)
			return nil, status, nil
		}
		log.Info("exec GetAddrUTXO(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! totalbalance is: %v \n", addrUtxolist, getbalance)

		//txrules.DefaultRelayFeePerKb
		log.Info("real account info: addrUtxolist[0].TxidHex is:%s, vout is :%d,amount is :%v,", addrUtxolist[0].TxidHex, addrUtxolist[0].Vout, addrUtxolist[0].Amount)
	*/

	//addrUtxolist :=[]proto.AddrBalanceUnspent{}
	addrUtxolist := make([]proto.AddrBalanceUnspent, 0, 8)
	curBtcUtxoInfo := proto.AddrBalanceUnspent{}
	//getutxoinfo,utxonum,err := UtxoRPCClient.GetTxUnSpentLimit("1Eq8xXAea47WPY5t8zUEYDKgcWB7cptZWB")
	getutxoinfo, utxonum, err := UtxoRPCClient.GetBTCTxUnSpentLimit(fromAddr)
	if err != nil {
		log.Error("fromAddr %s ,exex GetTxUnSpentLimit() failue! err is: %v \n", fromAddr, err)
		//return nil, status, err
		return nil, 0, proto.StatusLackUTXO, nil

	}

	log.Info("exec GetTxUnSpentLimit(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! unxonum is :%d\n", getutxoinfo, utxonum)
	//totalbalance is: %v ,getbalance
	if utxonum == 0 {
		log.Error("fromAddr %s ,exex GetTxUnSpentLimit() failue!,get addrUtxolist num is: %v \n", fromAddr, 0)
		return nil, 0, proto.StatusLackUTXO, nil

	}
	//11.20需要add 入排序
	selAmountIndex := 0
	selectAmountVal := getutxoinfo[0].Value
	for iseno, curitem := range getutxoinfo {
		if curitem.Value > selectAmountVal {
			selectAmountVal = curitem.Value
			selAmountIndex = iseno
		}
	}
	log.Info("exec GetAddrUTXO(),form addrUtxolist,selAmountIndex is :%d，getutxoinfo is:%v\n", selAmountIndex, getutxoinfo[selAmountIndex])

	curBtcUtxoInfo.TxidHex = getutxoinfo[selAmountIndex].TxHashBigEndian
	curBtcUtxoInfo.Vout = int(getutxoinfo[selAmountIndex].TxOutputN)
	curBtcUtxoInfo.Amount = float64(getutxoinfo[selAmountIndex].Value)
	/**/
	addrUtxolist = append(addrUtxolist, curBtcUtxoInfo)

	log.Info("get selAmountIndex is:%d,real account info: addrUtxolist[0].TxidHex is:%s, vout is :%d,amount is :%v,", selAmountIndex, addrUtxolist[0].TxidHex, addrUtxolist[0].Vout, addrUtxolist[0].Amount)

	getbalance := addrUtxolist[0].Amount
	//sgj 1120 end add
	curminamount := float64(toamount1) + 1000
	if getbalance < curminamount {
		//返回余额不足；
		log.Error("real account vout getbalance is:%d, curminamount is :%v,", getbalance, curminamount)
		//0507==temp from broadcast unspent info; real proc alonw :
		return nil, 0, proto.StatusLackBalance, nil

	}
	//2. 选择币的使用切片：credit（unspent output）
	//balance:(不需要,不管Account与address的关系;只负责接口的支付功能)
	curinput := []btcjson.TransactionInput{}
	//3.30 add: 0421 算法选取
	//找出转出金额，所需要花费对应的Txid：map rec :Txid,Vout,amount1
	perInput := btcjson.TransactionInput{
		//sgj 0108，0404 :
		Txid: addrUtxolist[0].TxidHex,
		//Txid:	"acfacf5ceea5122c6c7e07661160d7d12847b3a053429646e53d3b261b5d93c5",
		Vout: uint32(addrUtxolist[0].Vout),
	}
	curinput = append(curinput, perInput)
	//testing 第二个账户地址：

	receiver1, err := btcutil.DecodeAddress(toAddr, activeNet.Params)
	//log.Info("exec GetAddrUTXO(),addrUtxolist==watching---1,toAddr is:%s,receiver1 is :%v",toAddr,receiver1)
	if err != nil {
		log.Error("Failed to btcutil.DecodeAddress,toAddr is:%s，err is %v\n", toAddr, err)
	}
	//设置找零address:
	//receiver2,err := btcutil.DecodeAddress("1MrLfLBsujUBhmz5Da6Ceiq8aYPvbTPZ7i",activeNet.Params)
	receiver2, err := btcutil.DecodeAddress(remainAddr, activeNet.Params)
	if err != nil {
		log.Error("Failed to btcutil.DecodeAddress,remainAddr is:%s，err is %v\n", remainAddr, err)
	}
	//sgj 0402 testing for blance enought!:
	//找零后，为旷工费：
	//relayFee := txrules.DefaultRelayFeePerKb *6
	curfee := 52100.0

	//0830,add transfee:,//0920 change back
	//relayFee := txrules.DefaultRelayFeePerKb *80
	relayFee := txrules.DefaultRelayFeePerKb * 6

	//curfee = relayFee.ToBTC()
	//sgj testing fee:	0.0006
	//有时高达4到8mBTC，现在的数值一般是0.2到1mBTC

	//log.Info("exec GetAddrUTXO(),addrUtxolist==watching---i2i21,remainAddr is:%s,receiver2 is:%v",remainAddr,receiver2)
	curfee = float64(relayFee)
	toamount3 := addrUtxolist[0].Amount - toamount1 - curfee
	//toamount2 := 12000
	remainingAmount2 := toamount3
	receivers := map[btcutil.Address]btcutil.Amount{
		//sgj 0109 add
		//Value:         btcutil.Amount(value).ToBTC(),
		receiver1: btcutil.Amount(toamount1), //try amount sgj
		receiver2: btcutil.Amount(remainingAmount2),
	}
	var locktime int64 = time.Now().Unix()

	//第一步，创建交易,createTransaction：
	log.Info("step 1 bef,MakeMyTransaction() parmas curinput is:%v,receivers is :%v", curinput, receivers)
	_, err = ser.MakeMyTransaction(curinput, receivers, &locktime)
	log.Info("step 1 end,MakeMyTransaction() finished!!,real remainingAmount2 is :%d; relayFee is: %d", toamount3, relayFee)

	//第二步，签名交易,signedTransaction：
	//fromAddr ="1MrLfLBsujUBhmz5Da6Ceiq8aYPvbTPZ7i"
	thePreAddr := fromAddr
	preScriptPubKey, err := ser.GenscriptPubKeyFormAddr(thePreAddr)

	curRawTxInputs := make([]btcjson.RawTxInput, 0)
	//第二次，传给结构体with：Txid
	curTxIn := btcjson.RawTxInput{
		Txid: perInput.Txid,
		Vout: perInput.Vout,
		//sgj 1209,pubscript is last UTXO's Txid's pubscript
		//0104,ScriptPubKey,由函数计算出
		ScriptPubKey: string(preScriptPubKey),
	}
	log.Info("step 1 ext info,req curTxIn info is :%v", curTxIn)
	PriKeys := make([]string, 0)

	//从mysql里获取私钥：
	//11.20 update
	if true {
		//curPrikey, err := GetAddrPrivkey(fromAddr)
		//GWdcDataStore
		curaddrrec, err := wdctranssign.GWdcDataStore.GetBTCAddressRec(fromAddr)
		if err != nil {
			log.Error("GetBTCAddressRec(),get rows for fromaddress record failed!,GetBTCAddressRec() exec to return.curaddress =%s", fromAddr)
			return nil, 0, proto.StatusAccountPrikeyNotExisted, nil
		}
		getAddressPub := curaddrrec.PubKey
		curPrikey := curaddrrec.PrivKey
		log.Info("step 1 bef,MakeMyTransaction() parmas curinput is:%s,get curPrikey is :%s", getAddressPub, curPrikey)

		//没取到对应私钥：
		if curPrikey == "" || err != nil {
			log.Info("command %s ,exex GetAddressUtxo() failue! err is: %v \n", fromAddr, err)
			return nil, 0, proto.StatusAccountPrikeyNotExisted, err
		}
		PriKeys = append(PriKeys, curPrikey)
	} else {
		PriKeys = append(PriKeys, "")

	}
	//PriKeys = append(PriKeys,"5KD3sSntucZFRzDNUJusRjEhVDADir1xfDQqFxqoi7djDX5k81b")
	curRawTxInputs = append(curRawTxInputs, curTxIn)
	signTransFlag := "ALL"
	var cmd = &btcjson.SignRawTransactionCmd{}
	cmd = btcjson.NewSignRawTransactionCmd(m_txHex, &curRawTxInputs, &PriKeys, &signTransFlag)

	//第三步，进行签名交易
	transRes, err := ser.signRawTransaction(cmd, true)
	//if execcompletedflag == true {

	log.Info("step 2，signRawTransaction() exec finished! info return is :%v,err is :%v \n", transRes, err)
	if err != nil {
		log.Error("Failed to sign transaction，err is %v\n", err)
		//continue
		return transRes, 0, proto.StatusSignError, err
	}
	//signedTransaction, complete, err := rpcClient.SignRawTransaction(tx.Tx)
	//第四步，to发送交易：ds
	return transRes, 0, proto.StatusSuccess, nil
}
