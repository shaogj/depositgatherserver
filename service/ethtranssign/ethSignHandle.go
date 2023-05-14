package ethtranssign

import (
	"2019NNZXProj10/depositgatherserver/config"
	"2019NNZXProj10/depositgatherserver/service/ethtranssign/ethclientrpc"

	"fmt"
	//"2019NNZXProj10/depositgatherserver/service/ethtranssign"
	"2019NNZXProj10/depositgatherserver/service/ethtranssign/base"
	"2019NNZXProj10/depositgatherserver/proto"
	"github.com/mkideal/log"
	"github.com/go-xorm/xorm"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
	"strings"


)
func main() {
	fmt.Println("vim-go")
}

type EthSignHandle struct {
	CoinType string  //币种类型
	EthTransUrl string
	//HtClient CHttpClientEx
	//WdcRpcClient	*WdcRpcClient

}
var GEthSignHandle  EthSignHandle

var (
	ETH_MOrmEngine *xorm.Engine = &xorm.Engine{}
)
//获取ETH的地址的私钥
func GetAddrPrivkeyETH(curaddress string) (addrPrikey string,err error){

	engineread:= ETH_MOrmEngine
	//get address 's [privkey]
	selectsql := "select * from  eth_account_key where address = '"  + curaddress + "'"
	addr_accountinfo, err := engineread.Query(selectsql)
	if err != nil || len(addr_accountinfo) <= 0{
		log.Info("when GetAddrPrivkeyETH(),curaddress' is:%s ,get privkey error: %v", curaddress,err)
		return "",err
	}

	curaddrprivkey := string(addr_accountinfo[0]["priv_key"])
	log.Info("when GetAddrPrivkeyETH(),curaddress is :%s,get addr's privkey succ ,info is: %v", curaddress,curaddrprivkey)
	return curaddrprivkey,nil
}

//签名交易的处理流程控制
func (ser *EthSignHandle) PaySignTransProc(curtransreq *proto.SignTransactionReq, cursettle proto.Settle,outAccountAddress string,execcompletedflag bool) (txdata interface{}, status int, err error) {

	log.Info("exec PaySignTransProc() step 1,curtransreq of cursettle is : %v \n", cursettle)

	fromAddr :=cursettle.FromAddress
	if fromAddr == ""{
		fromAddr = outAccountAddress
	}
	toAddr := cursettle.ToAddress
	toamount1,_:= cursettle.Vol.Float64()
	toamountfee1,_:= cursettle.Fee.Float64()

	//11.29 doing---toamount1 = toamount1 * 100000000

	attachEthereum:=proto.SignAttachEthreum{}
	//目前方式，不连接geth：

	//tran.Nonce=uNonce
	enduNonce, err := ethclientrpc.GetTransactionAccount(outAccountAddress)
	if err != nil {
		log.Error("GetTransactionAccount get nonce num error:%s", err.Error())
		return nil,proto.ErrorRequestInfuraETHNode.Code,err
	}
	log.Info("GetTransactionAccount get nonce num succ!,get end is:%d", enduNonce)

	uNonce2,_:=ethclientrpc.GetTransactionPendingNonce(fromAddr)
	if uNonce2 > enduNonce {
		enduNonce = uNonce2
	}
	log.Info("uNoncePending2=%d,enduNonce=%d=====004",uNonce2,enduNonce)
	attachEthereum.Nonce = uint64(enduNonce)
	//获取余额1128:
	getNodeAmount, err := ethclientrpc.GetBalance(outAccountAddress,"latest")
	if err != nil {
		log.Error("GetBalance get nonce num error:%s", err.Error())
		return nil,proto.ErrorRequestInfuraETHNode.Code,err
	}
	getNodeAmountNew := float64(getNodeAmount)/config.EthEtherPrefix
	//sgj upgrade:
	//再次移位：
	getNodeAmountNew = getNodeAmountNew / config.EthEtherPrefix
	log.Info("GetBalance get nonce num succ!,get getNodeAmount is:%d,get getNodeAmountNew is:%f", getNodeAmount,getNodeAmountNew)
	//"当前余额不够"
	totalNeedAmount1 := toamountfee1 + toamount1
	if getNodeAmountNew < totalNeedAmount1 {
		log.Error("GetBalance get amount num is lack!")
		return nil,proto.StatusLackBalance,err

	}

	sttFee:=config.GetGas(config.CoinEthereum,config.CoinEtherTrans,1.0,1.0)
	attachEthereum.GasLimit =uint64(sttFee.Limit)
	attachEthereum.GasPrice = sttFee.Price

	//end sgj 11.28adding

	tran:=base.EtherTranInfo{
		SubType:"",
		From:fromAddr,
		To:toAddr,
		Amount:toamount1,
		GasPrice:attachEthereum.GasPrice,
		GasLimit:attachEthereum.GasLimit,
		UNonce:attachEthereum.Nonce,
	}
	//tranErrDesc :=fmt.Sprintf("nonce1=%d,nonce2=%d,nonce=%d",uNonce1,uNonce2,uNonce)
	log.Info("tran =%v",tran)

	//加载 用户私钥
	//2.生成进制的交易数据
	//vRaw:=base.EtherTranBinary{}
	if config.CoinEthereum==cursettle.CoinCode {
		curPrikey, err := GetAddrPrivkeyETH(fromAddr)
		//没取到对应私钥：
		if curPrikey == "" || err != nil {
			log.Info("command %s ,exex GetAddrPrivkeyETH() failue! err is: %v \n", fromAddr, err)
			return nil, proto.StatusAccountPrikeyNotExisted, err
		}
		tran.Private = curPrikey
		getRawTx,errErrInfo:= HServer.DoTransactionPreNew(tran)
		if proto.ErrorSuccess.Code !=errErrInfo.Code   {
			log.Error("sign Tx invalid!, Error{%d,%s} data=%+v  , url=%s ",errErrInfo.Code,errErrInfo.Desc,tran,"signInfoUrl")
		}else{
			log.Info(" sign Tx  succ!{%d,%s} data=%+v  , getRawTx=%v ",errErrInfo.Code,errErrInfo.Desc,tran,getRawTx)

		}
		//11.29 add SendTX:
		//SendTransactionRaw
		tranRaw:=base.EtherTranBinary{}
		data, err := rlp.EncodeToBytes(getRawTx)
		if err != nil {
			return   nil,5086,err
		}
		tranRaw.Raw=common.ToHex(data)
		tranRaw.Hash=getRawTx.Hash().String()

		txHash,err := ethclientrpc.SendTransactionRaw(tranRaw.Raw)
		if nil==err {
			log.Error("%s.Tx SendTransactionRaw succ!: txhash=%s, raw=%s",cursettle.CoinCode,tranRaw.Hash,tranRaw.Raw)
			return  txHash,proto.ErrorSuccess.Code,nil
		}else{
			log.Error("%s.Tx SendTransactionRaw Error!: err:= %v, txhash=%s, raw=%s",cursettle.CoinCode, err,tranRaw.Hash,tranRaw.Raw)
		}

		strErr:=err.Error()
		//geth error:*known transaction*
		if true==strings.Contains(strErr,"known transaction"){
			return  "",proto.ErrorSuccess.Code,err
		}
		//parity error:Transaction with the same hash was already imported.
		if true==strings.Contains(strErr,"same hash"){
			return  "",proto.ErrorSuccess.Code,err
		}
		return "",proto.ErrorRequestInfuraETHSend.Code,err
	}

	return nil,0 ,nil
}

