//package btcsignlocal
package ktctranssign

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	//"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcwallet/netparams"
	"github.com/go-spew/spew"
	//add
	"github.com/go-xorm/xorm"
	//"strconv"

	//1209-remove--"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/mkideal/log"

	"encoding/json"

	"2019NNZXProj10/depositgatherserver/proto"
	"2019NNZXProj10/depositgatherserver/service/ktctranssign/ktcrpc"
)

//0507testing:
//var activeNet = &netparams.MainNetParams

//sgj 0330 for testing net
var activeNet = &netparams.TestNet3Params

var (
	KTC_MOrmEngine *xorm.Engine = &xorm.Engine{}
)
var GKTCSignHandle = KTCSignHandle{}

//sgj 0427 add;签名处理的实体类
type KTCSignHandle struct {

	//sgj 1217add from DepositGather
	GatherLimit		float64
	//sgj 1217adding,总归集的地址数量
	GatherAddrCount	int
}

//0911 add:
type KTCPreTx struct {
	Txid         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	ScriptPubKey string `json:"scriptPubKey"`
	//sgj 1210 add:
	RedeemScript string `json:"redeemScript"`
	//Value        float64 `json:"value"`
	Amount float64 `json:"amount"`
}

type ErrorCode int

const (
	maxProtocolVersion = 70002
)

var m_curpkScript []byte
var m_txHex string

var m_address btcutil.Address

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

//sgj add at 12 05, form :/go/src/github.com/btcsuite/btcwallet/internal/legacy/keystore/
// newScriptAddress initializes and returns a new P2SH address.
// iv must be 16 bytes, or nil (in which case it is randomly generated).

// ChangeSource provides P2PKH change output scripts for transaction creation.
type ChangeSource func() ([]byte, error)

// rpcDecodeHexError is a convenience function for returning a nicely formatted
// RPC error which indicates the provided hex string failed to decode.
func (ser *KTCSignHandle) rpcDecodeHexError(gotHex string) *btcjson.RPCError {
	return btcjson.NewRPCError(btcjson.ErrRPCDecodeHexString,
		fmt.Sprintf("Argument must be hexadecimal string (not %q)",
			gotHex))
}

func (ser *KTCSignHandle) internalRPCError(errStr, context string) *btcjson.RPCError {
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
func (ser *KTCSignHandle) messageToHex(msg wire.Message) (string, error) {
	var buf bytes.Buffer
	if err := msg.BtcEncode(&buf, maxProtocolVersion, wire.WitnessEncoding); err != nil {
		context := fmt.Sprintf("Failed to encode msg of type %T", msg)
		return "", ser.internalRPCError(err.Error(), context)
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

//获取莱特币地址的私钥
func GetAddrPrivkeyKTC(curaddress string) (addrPrikey string, err error) {

	engineread := KTC_MOrmEngine
	//get address 's [privkey]
	selectsql := "select * from  ktc_account_key where address = '" + curaddress + "'"
	addr_accountinfo, err := engineread.Query(selectsql)
	if err != nil || len(addr_accountinfo) <= 0 {
		log.Info("when GetAddrPrivkeyKTC(),curaddress' is:%s ,get privkey error: %v", curaddress, err)
		return "", err
	}

	curaddrprivkey := string(addr_accountinfo[0]["priv_key"])
	log.Info("when GetAddrPrivkeyKTC(),curaddress is :%s,get addr's privkey succ ,info is: %v", curaddress, curaddrprivkey)
	return curaddrprivkey, nil
}


//sgj 1017 add for mult gorouting utxo;
//sgj 0104 add
//从地址生成脚本公钥的工具函数
func (ser *KTCSignHandle) GenscriptPubKeyFormAddr(encodedAddr string) (string, error) {
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
func (ser *KTCSignHandle) PrepareSignRawTransactionTx(tx *wire.MsgTx) string {
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
func (ser *KTCSignHandle) handleCreateRawTransaction(cmd interface{}) (interface{}, error) {
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

//真正的签名交易处理
//sgj 0911 add params: to handle
//addrkey 对应不上wif.PrivKey of type:additionalKeysByAddress
func (ser *KTCSignHandle) SignTransaction(tx *wire.MsgTx, hashType txscript.SigHashType,
	additionalPrevScripts map[wire.OutPoint][]byte,
	additionalKeysByAddress map[string]*btcutil.WIF,
	p2shRedeemScriptsByAddress map[string][]byte, addrkey map[string]string) ([]SignatureError, error) {

	var signErrors []SignatureError
	log.Info("wacthinging---0911====cur exec SignTransaction(): params tx is :%v;additionalPrevScripts is :%v ; additionalKeysByAddress is :%v\n", tx, additionalPrevScripts, additionalKeysByAddress)

	return signErrors, nil
}

// decodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when
// number of bytes when calling hex.Decode.
func (ser *KTCSignHandle) decodeHexStr(hexStr string) ([]byte, error) {
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

//处理签名交易的command.
//sgj 0911 add fromAddr;

//signRawTransaction() no need for KTC coin:
func (ser *KTCSignHandle) signRawTransaction(fromAddr string, icmd interface{}, execcompletedflag bool) (interface{}, error) {
	cmd := icmd.(*btcjson.SignRawTransactionCmd)
	/*

		return btcjson.SignRawTransactionResult{
			Hex:      hex.EncodeToString(buf.Bytes()),
			Complete: len(signErrors) == 0,
			Errors:   signErrors,
		}, nil
	*/
		cmdPartSignRes := &btcjson.SignRawTransactionCmd{
			//in Goserver2.handle:
			//signErrs, err := ser.SignTransaction(&tx, hashType, inputs, keys, scripts)
			RawTx:  "Hexstrinfo",
			Inputs: cmd.Inputs,
			Flags:  cmd.Flags,
		}
		return cmdPartSignRes, nil
}

/*
type TransactionInput struct {
	Txid string `json:"txid"`
	Vout uint32 `json:"vout"`
}
*/

type QtumTxIn struct {
	Txid string `json:"txid"`
	Vout uint32 `json:"vout”`
	//ScriptPubKey string `json:"scriptPubKey"`
	//Sequence string `json:”sequence"`
	//Value float64 `json:"value"`
}

//	Amounts  map[string]float64 `jsonrpcusage:"{\"address\":amount,...}"` // In BTC

type OutAmounts struct {
	Amounts map[string]float64 `jsonrpcusage:"{\"address\":amount,...}"` // In BTC
	//Amount  map[string]float64 `json:"{\"address\":amount,...}"` // In BTC
	//AmountStr  string `json:"{\"address\":amount,...}"` // In BTC

}

//1209转换手续费用d
func FeeDecimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", value), 64)
	return value
}
//Status
//签名交易的处理流程控制


//sgj 1217,,小于此值，不进行归集处理
var minKTCLimit = 0.00005
var minKTCLimitLast = 0.00001


//func (ser *KTCSignHandle) PaySignTransProc(curtransreq *proto.SignTransactionReq, cursettle proto.Settle,outAccountAddress string,execcompletedflag bool) (txdata interface{}, status int, err error) {
func (ser *KTCSignHandle) PaySignTransProc(fromAddr string,fromPrivKey string,toamount float64,gatherAddress string,curGatherLimit float64) (txdata interface{}, gatheredamount1 float64,status int, err error) {

	//sgj 1217,,toamount 为归集的所用数量，，utxo最大的去掉fee，即为要规的curtoamount的值
	ser.GatherLimit = curGatherLimit

	//good!,从地址生成公钥！
	//步骤记录： 1）获取自己地址的公钥
	//sgj 1120doing
	log.Info("exec PaySignTransProc() step 1,fromAddr is : %s, toamount is:%f,GatherLimit is:%f \n", fromAddr,toamount,ser.GatherLimit)

	toAddr := gatherAddress
	//toamount1 := cursettle * 100000000
	//toamount1,_:= cursettle.Vol.Float64()
	toamount1:= toamount
	//1210trying

	//toamount1 = 0.041
	remainAddr := fromAddr	//cursettle.FromAddress

	m_getCurPrivKeuy, err := ser.GenscriptPubKeyFormAddr(fromAddr)
	log.Info("fromAddr is : %s ,ToAddr is :%s,,it's  m_getCurPrivKeuy info is:%s \n", fromAddr, toAddr, m_getCurPrivKeuy)
	if len(fromAddr) < 32 || len(fromAddr) > 34 {
		log.Error("requset param fromAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", fromAddr, len(fromAddr))
		return nil, 0,proto.StatusInvalidArgument, nil
	}

	if len(toAddr) < 32 || len(toAddr) > 34 {
		log.Error("requset param toAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", toAddr, len(toAddr))
		return nil, 0,proto.StatusInvalidArgument, nil
	}
	//0330 to addr set is: toAddr:
	if len(remainAddr) < 32 || len(remainAddr) > 34 {
		log.Error("requset param RemainAddr is :%s,it's len is :%d, format is err.Invoke is returned \n", remainAddr, len(remainAddr))
		return nil, 0,proto.StatusInvalidArgument, nil
	}
	if toAddr == remainAddr {
		log.Error("requset param toAddr same to the RemainAddr is :%s,opearte is forbidden! \n", remainAddr)
		return nil, 0,proto.StatusInvalidArgument, nil

	}
	//1. 找出未花费的币（unspent output）：
	//最后的找零地址为自身
	//sgj 1207 add fro KTC:
	addrUtxolist := make([]proto.CurKtcUtxoInfo,0,8)
	curAddrLists := make([]string,0,3)
	curAddrLists = append(curAddrLists,fromAddr)


	getutxoinfo,utxonum,err := ktcrpc.KTCRPCClient.GetRPCTxUnSpentLimit(1,9999999,curAddrLists)	//"1Eq8xXAea47WPY5t8zUEYDKgcWB7cptZWB")
	if err != nil {
		log.Error("from address :%s ,to address :%s ,exec GetRPCTxUnSpentLimit() failure! err is: %v \n", fromAddr, toAddr,err)
		return nil, 0,status, err
	}

	log.Info("exec GetTxUnSpentLimit(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! unxonum is :%d\n", getutxoinfo, utxonum)
	//totalbalance is: %v ,getbalance
	if utxonum == 0 {
		log.Error("fromAddr %s ,exex GetTxUnSpentLimit() failue!,get addrUtxolist num is: %v。PaySignTransProc is break! \n", fromAddr, 0)
		return nil, 0,status, err

	}

	//11.20,请求比特币,未花费的排序过的balance,按ut.value倒排序：
	selAmountIndex :=0
	selectAmountVal:= getutxoinfo[0].Amount
	for iseno,curitem := range getutxoinfo{
		if curitem.Amount > selectAmountVal{
			selectAmountVal = curitem.Amount
			selAmountIndex = iseno
		}
	}
	log.Info("exec GetAddrUTXO(),form addrUtxolist,selAmountIndex is :%d，getutxoinfo is:%v\n",selAmountIndex,getutxoinfo[selAmountIndex])

	//sg 1210 watching:
	time.Sleep(time.Second * 4)
	addrUtxolist = append(addrUtxolist,getutxoinfo[selAmountIndex])

	log.Info("get selAmountIndex is:%d,real account info: addrUtxolist[0].TxidHex is:%s, vout is :%d,amount is :%f,", selAmountIndex,addrUtxolist[0].Txid, addrUtxolist[0].Vout, addrUtxolist[0].Amount)

	getbalance := addrUtxolist[0].Amount
	//sgj 1217,,to update ,为最大大一个utxo大切片,to gather for all utxoid
	//sgj 1217,,toamount 为归集的所用数量
	//toamount1 = getbalance
	log.Info("exec GetAddrUTXO(),addrUtxolist info is: %v ,exex GetAddressUtxo() finished! totalbalance is: %f \n", addrUtxolist, getbalance)

	//找零后，为旷工费：
	//AmountFee := 0.00002//0.0009
	// AmountFee := 0.00002 * 5//0.0009
	//AmountFee := 0.00012
	AmountFee := 0.00002

	curfee := AmountFee

	//curfee = relayFee.ToBTC()
	//sgj testing fee:	0.0006
	//有时高达4到8mBTC，现在的数值一般是0.2到1mBTC
	/*
	//curminamount := float64(toamount1) + 1000
	curminamount := float64(toamount1) + curfee
	if getbalance < curminamount {
		//返回余额不足；
		log.Error("real account vout getbalance is:%d, curminamount is :%v,", getbalance, curminamount)
		//0507==temp from broadcast unspent info; real proc alonw :

		return nil, 0,proto.StatusLackBalance, nil

	}
	//sgj 0824 wathing
	log.Info("step 1:=====In Watchcing-------------001")
	*/
	//ing--using---curGatherAmount := getbalance - AmountFee
	//1218checking
	//curGatherAmount = curGatherAmount - 0.00005
	//succ1--->curGatherAmount = 0.0015
	var curGatherAmount float64
	//curGatherAmount= FeeDecimal(curGatherAmount)
	//succ1--->2curGatherAmount= 0.0013
	curGatherAmount = getbalance - AmountFee - minKTCLimitLast//minKTCLimit
	curGatherAmount = FeeDecimal(curGatherAmount)

	//1217checking
	/*
	if  curGatherAmount > 0.0003{
		curGatherAmount -= 0.0003
	}
	to fix error：
	ktcrpc.CreateTransaction(),getresinfo info is :&{ <nil> -3:Invalid amount 0},err is :<nil>
	Why??
	*/
	//1114 add,满足归集最大上限为止
	if curGatherAmount > ser.GatherLimit {
		curGatherAmount = ser.GatherLimit
	}
	//toamount1 = curGatherAmount
	//1217 check no 0 value
	toamount1 = curGatherAmount// - 0.0003
	var totalNeeds float64 = (minKTCLimit + AmountFee)	// * 100000000
	/*
	fromMount = fromMount * 100000000
	*/
	//余额不够最小归集限额,停止此比交易
	if  totalNeeds > getbalance {
		log.Info("balance is too low ignore. KTC Trans is insufficient!,cur balance is %.8f,cursettle need is:%.8f,real toamount1 is:%.8f\n", getbalance,totalNeeds,toamount1)
		/*
		reqUpdateInfo.Withdraws[0].Status = proto.SETTLE_STATUS_FAILED
		reqUpdateInfo.Withdraws[0].Error = "当前余额不够"
		if isOk := self.WithdrawsUpdate(&reqUpdateInfo); isOk {
			log.Error("WDCGatherTransProc.WDCTransProc() fail, exec compare balance failed!,curid is:%d,curbalance is:%.8f,totalNeeds amount is: %.8f,cur trans break!", iseno,fromMount,totalNeeds)
		}
		*/
		return nil, 0,status, err
	}
	log.Info("cur KTC Trans amount info: cur balance is %.8f,curGatherAmount is:%.8f, curFee is:%.8f,real toamount1 is:%.8f\n", getbalance,curGatherAmount,AmountFee,toamount1)


	//2. 选择币的使用切片,不需要,不管Account与address的关系;只负责接口的支付功能)
	curinput := []btcjson.TransactionInput{}
	//3.30 add: 0421 算法选取
	//找出转出金额，所需要花费对应的Txid：map rec :Txid,Vout,amount1

	perInput := btcjson.TransactionInput{
		Txid: addrUtxolist[0].Txid,
		Vout: uint32(addrUtxolist[0].Vout),
	}

	curinput = append(curinput, perInput)
	//testing 第二个账户地址：

	var inaddrmapNew = OutAmounts{}
	inaddrmapNew.Amounts = make(map[string]float64, 0)

	//TransactionInput
	//receiver1, err := btcutil.DecodeAddress(toAddr, activeNet.Params)

	//curfee = float64(relayFee)
	toamount3 := addrUtxolist[0].Amount - toamount1 - curfee
	remainingAmount2 := toamount3
	remainingAmount2 = FeeDecimal(remainingAmount2)
	log.Info("cur watching1209------A03,getbalance is:%.8f,curGatherAmount is:%.8f,curfee is:%.8f, toAddr is :%s,toamount1 is:%.8f,toamount3 is:%.8f,remainingAmount2New is :%f,remainAddr is:%s", getbalance,curGatherAmount,curfee,toAddr,toamount1,toamount3,remainingAmount2,remainAddr)


	//sjg 1209 update:
	var locktime int64
	//toamount1 / 100000000
	inaddrmapNew.Amounts[toAddr] = toamount1
	inaddrmapNew.Amounts[remainAddr] = remainingAmount2

	cmdcreateparams := btcjson.NewCreateRawTransactionCmd(curinput, inaddrmapNew.Amounts, &locktime)

	log.Info("step 2:=====Bef Invoke RPC CreateTransaction(),,params cmd is :%v", cmdcreateparams)
	// 1209,,直接调用jsonRPC：
	response, err := ktcrpc.KTCRPCClient.RpcClient.Call("createrawtransaction", cmdcreateparams.Inputs, cmdcreateparams.Amounts)
	//getresinfo :=response.Result.(string)
	var curCreateRawTrans string
	if response.Result != nil {
		curCreateRawTrans = response.Result.(string)
	} else {
		curCreateRawTrans = ""
	}
	m_txHex := curCreateRawTrans

	log.Info("step 3:=====In ktcrpc.CreateTransaction(),getresinfo info is :%v,err is :%v", response, err)
	//m_txHex = ser.PrepareSignRawTransactionTx(mtx)

	fmt.Printf("step 3.1:=====In ktcrpc.CreateTransaction(), composed MsgTx struct,get m_txHex is:%s \n", m_txHex)

	//sgj watching1210 check tmping
	//time.Sleep(time.Second * 4)
	//return

	//第二步，签名交易,signedTransaction：
	//fromAddr ="1MrLfLBsujUBhmz5Da6Ceiq8aYPvbTPZ7i"
	//thePreAddr := fromAddr
	//preScriptPubKey, err := ser.GenscriptPubKeyFormAddr(thePreAddr)
	//log.Info("step 1 end,MakeMyTransaction() finished!!,real remainingAmount2 is :%.8f; relayFee is: %d,preScriptPubKey is:%s", toamount3, curfee,preScriptPubKey)

	curRawTxInputs := make([]btcjson.RawTxInput, 0)
	//第二次，传给结构体with：Txid
	curTxIn := btcjson.RawTxInput{
		Txid: perInput.Txid,
		Vout: perInput.Vout,
		//sgj 1209,pubscript is last UTXO's Txid's pubscript
		//0104,ScriptPubKey,由函数计算出
		ScriptPubKey :addrUtxolist[0].ScriptPubKey,
	}
	log.Info("step 1 ext info,req curTxIn info is :%v", curTxIn)
	PriKeys := make([]string, 0)

	//从mysql里获取私钥：

	//curPrikey, err := GetAddrPrivkeyKTC(fromAddr)
	curPrikey := fromPrivKey
	//没取到对应私钥：
	if curPrikey == "" {
		log.Info("exec GetAddrPrivkeyKTC() failue! fromAddr is:%s,err is: %v \n", fromAddr, err)
		return nil, 0,proto.StatusAccountPrikeyNotExisted, err
	}
	log.Info("to signRawTransaction(),cur fromAddr is:%s,it's curPrikey is: %s \n", fromAddr, curPrikey)
	PriKeys = append(PriKeys, curPrikey)

	log.Info("cur--exec GetAddrPrivkeyKTC(),get PriKeys is: %v \n", PriKeys)
	curRawTxInputs = append(curRawTxInputs, curTxIn)
	signTransFlag := "ALL"

	//cmd = btcjson.NewSignRawTransactionCmd(m_txHex, &curRawTxInputs, &PriKeys, &signTransFlag)

	//第三步，进行签名交易
	log.Info("before signRawTransaction() exec, param PriKeys is :%v,curRawTxInputs is :%v,get m_txHex is :%s \n", PriKeys, curRawTxInputs, m_txHex)
	//sgj 0109 真正的签名部分：！！
	//不用通常的结构数据方法，直接调用RPC func ：signrawtransaction

	//preScriptPubKey, err := ser.GenscriptPubKeyFormAddr(thePreAddr)
	preScriptPubKey := addrUtxolist[0].ScriptPubKey
	redeemScript := addrUtxolist[0].RedeemScript
	curKTCRawTxInputs := []*KTCPreTx{}
	curKTCTxIn := KTCPreTx{
		Txid: addrUtxolist[0].Txid,
		Vout:         uint32(addrUtxolist[0].Vout),
		ScriptPubKey: string(preScriptPubKey),
		//Amount:       addrUtxolist[0].Amount / 100000000,
		//sgj 1220 add

		RedeemScript: string(redeemScript),
		Amount:       addrUtxolist[0].Amount,
	}
	curKTCRawTxInputs = append(curKTCRawTxInputs, &curKTCTxIn)
	//sgj 1210 add checking:
	log.Info("bef PRC signrawtransaction(),cur param curKTCRawTxInputs info is:%s,preScriptPubKey is :%s", curKTCTxIn, preScriptPubKey)

	getCurSignedRawTxResp, err := ktcrpc.KTCRPCClient.RpcClient.Call("signrawtransaction", m_txHex, curKTCRawTxInputs, PriKeys, signTransFlag)
	if nil != err {
		log.Error("exec KTC's PRC signrawtransaction() err!,req getCurSignedRawTxResp info is:%s,err is :%v", getCurSignedRawTxResp, err)
		//return getCurSignedRawTxResp,proto.StatusSignError,err
	} else {
		log.Info("exec KTC's PRC signrawtransaction() succ!,req getCurSignedRawTxResp info is:%s,err is :%v", getCurSignedRawTxResp, err)
	}

	//0703 add:
	get_response, err := json.Marshal(getCurSignedRawTxResp.Result)
	if err != nil {
		log.Error("GetOmniTransaction(),response.Result err !, get_response is:%v,err is:%v", get_response, err)
		//return nil,err
	}
	getKTCRawTx := btcjson.SignRawTransactionResult{}
	err = json.Unmarshal([]byte(get_response), &getKTCRawTx)
	if err != nil {
		log.Error("KTC signrawtransaction(),Unmarshal to getKTCRawTx{} err !, get_response is:%v,err is:%v", get_response, err)
		//return nil,err
	}
	log.Info("step 2，exec KTC's PRC signrawtransaction()==second succ!,req info return getKTCRawTx's Hex is:%s,err is :%v", getKTCRawTx.Hex, err)
	return getKTCRawTx,toamount1, proto.StatusSuccess, err

}
