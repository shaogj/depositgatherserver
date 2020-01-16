package ethtranssign

import (
	//"2017opensource/beepay/config"

	"fmt"

	"strings"

	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/rlp"

	"2019NNZXProj10/depositgatherserver/config"
	"2019NNZXProj10/depositgatherserver/proto"

	//"2017opensource/beepay/models"
	//"2017opensource/beepay/service/base"
	"2019NNZXProj10/depositgatherserver/service/ethtranssign/base"
	Log "github.com/mkideal/log"

	//"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	MaxChanTrans =200
	SleepMilliSecond  =2000
	LEN_TXHASH_STR_LEN =  66
	CoinType  string ="ETH"
)

var   HServer EtherSever=EtherSever{EtherBase:base.EtherBase{Coin:config.CoinEthereum}}

type EtherSever struct {
	base.EtherBase
}


func (self * EtherSever)Host(itype  int,strHost string, iPort int)  bool {
	if base.TypeLocal == itype {
		self.StrRawUrl = "\\\\.\\pipe\\geth.ipc"
	} else if base.TypeWebSocket == itype {
		self.StrRawUrl = fmt.Sprintf("ws://%s:%d", strHost, iPort)
	} else if base.TypeHttp == itype {
		self.StrRawUrl = fmt.Sprintf("http://%s:%d", strHost, iPort)
	} else {
		Log.Error("Parameter  error  ")
		return  false
	}
	return  true
}
//周期性和节点上同步费用信息
func (self* EtherSever)GoEstimateGas() {

	var conn * ethclient.Client
	var err error
	for {
		conn, err = ethclient.Dial(self.StrRawUrl)
		if err == nil {
			Log.Info("Success Connect  %s",self.StrRawUrl)
			break
		}
		Log.Error("Failed Connect %s ::: %v",self.StrRawUrl, err)
		time.Sleep(time.Duration(config.SleepMilliSecond)*time.Millisecond)
	}

	sttPreriod:=config.GetGas(self.Coin,config.CoinEstimatePeriod,1.0,1.0)
	var iSleep int64=10
	iPreriod:=sttPreriod.Limit/iSleep
	if iPreriod<1 {
		iPreriod=1
	}
	Log.Info("%s.SuggestGasPrice  Period{counter=%d, per=%d second}",self.Coin,iPreriod,iSleep)
	var iter int64=0
	for  {
		if iter>=iPreriod {
			var vprice  * big.Int=nil
			ctx, _ := context.WithTimeout(context.Background(), 120 * time.Second)
			vprice,err=conn.SuggestGasPrice(ctx)
			if err!=nil  {
				Log.Error("-----   %v",err)
			}
			if nil!=vprice {
				sttFee:=config.GetGas(self.Coin,config.CoinEtherTrans,1.0,1.0)
				sttFee.Price=(vprice.Int64()*15/10)
				Log.Info("%s.SuggestGasPrice=%d[%d]",self.Coin,sttFee.Price,vprice.Int64())
				config.SetGas(self.Coin,config.CoinEtherTrans,sttFee)
			}
			iter=0
		}
		iter+=1
		time.Sleep(time.Duration(iPreriod)*time.Second)
	}
}
//获取最新的块号
//#-----


//生成二进制交易数据
//subType string,uNonce uint64,from ,to string, famount float64,gasPrice int64,gasLimit uint64) (string,string,ErrorInfo)
func (self* EtherSever)DoTransactionPre(tran base.EtherTranInfo) (base.EtherTranBinary, proto.ErrorInfo)  {

	var signedTx *types.Transaction
	var rawTx *types.Transaction
	var err error
	raw:=base.EtherTranBinary{}
	var  gasPricebig  big.Int

	gasPricebig.SetInt64(tran.GasPrice)
	Log.Error("------subType=%s --------",tran.SubType)
	for {
		if ""==tran.SubType{
			bigAmount:=base.EtherToWei(tran.Amount)
			rawTx = types.NewTransaction(tran.UNonce, common.HexToAddress(tran.To), &bigAmount, tran.GasLimit,&gasPricebig, []byte(""))
			break
		}
		//ti:=TokenInfo{}

	}
	//签名交易
	if true {
		var ecdsaPriv *ecdsa.PrivateKey
		ecdsaPriv,err=crypto.HexToECDSA(tran.Private)
		if nil!=err {
			return  raw,proto.ErrorInfo{5084,"error private string"}
		}
		signer:=types.HomesteadSigner{}
		signature, err := crypto.Sign(signer.Hash(rawTx).Bytes(), ecdsaPriv)
		if err != nil {
			return  raw,proto.ErrorInfo{5085,"error crypto.Sign"}
		}
		signedTx, err =rawTx.WithSignature(signer, signature)
		if err != nil {
			return  raw,proto.ErrorInfo{5086,"error rawTx.WithSignature"}
		}
		if nil!=err{
			Log.Error("--SignTxWithPassphrase-- {tx=%v,err=%v}",signedTx,err)
			return  raw,proto.ErrorInfo{5085,"SignTxWithPassphrase failure"}
		}
		if nil==signedTx || nil!=err{
			Log.Error("--SignTxWithPassphrase-- {tx=%v,err=%v}",signedTx,err)
			return  raw,proto.ErrorInfo{5085,"SignTxWithPassphrase failure"}
		}
	}
	Log.Info("---SignTxWithPassphrase--- tx=%+v",signedTx)
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return  raw,proto.ErrorInfo{5086,string(err.Error())}
	}
	raw.Raw=common.ToHex(data)
	raw.Hash=signedTx.Hash().String()
	return  raw,proto.ErrorSuccess
}

//sgj 1127adding:
func (self* EtherSever)DoTransactionPreNew(tran base.EtherTranInfo) (*types.Transaction, proto.ErrorInfo)  {

	var signedTx *types.Transaction
	var rawTx *types.Transaction
	var err error
	var  gasPricebig  big.Int

	gasPricebig.SetInt64(tran.GasPrice)
	Log.Error("------subType=%s --------",tran.SubType)
	if ""==tran.SubType{
		bigAmount:=base.EtherToWei(tran.Amount)
		rawTx = types.NewTransaction(tran.UNonce, common.HexToAddress(tran.To), &bigAmount, tran.GasLimit,&gasPricebig, []byte(""))
	}
	//签名交易
	if true {
		var ecdsaPriv *ecdsa.PrivateKey
		ecdsaPriv,err=crypto.HexToECDSA(tran.Private)
		if nil!=err {
			return  nil,proto.ErrorInfo{5084,"error private string"}
		}
		signer:=types.HomesteadSigner{}
		signature, err := crypto.Sign(signer.Hash(rawTx).Bytes(), ecdsaPriv)
		if err != nil {
			return  nil,proto.ErrorInfo{5085,"error crypto.Sign"}
		}
		signedTx, err =rawTx.WithSignature(signer, signature)
		if err != nil {
			return  nil,proto.ErrorInfo{5086,"error rawTx.WithSignature"}
		}
		if nil!=err{
			Log.Error("--SignTxWithPassphrase-- {tx=%v,err=%v}",signedTx,err)
			return  nil,proto.ErrorInfo{5085,"SignTxWithPassphrase failure"}
		}
		if nil==signedTx || nil!=err{
			Log.Error("--SignTxWithPassphrase-- {tx=%v,err=%v}",signedTx,err)
			return  nil,proto.ErrorInfo{5085,"SignTxWithPassphrase failure"}
		}
	}
	Log.Info("---SignTxWithPassphrase--- tx=%+v",signedTx)
	//data, err := rlp.EncodeToBytes(signedTx)

	return  signedTx,proto.ErrorSuccess
}

//返回错误码
func IsReturnError(err error)  int {
	if nil==err {
		return 0
	}
	strErr:=err.Error()
	if true==strings.HasPrefix(strErr,"nsufficient funds for gas * price + value") {   //insufficient funds for gas * price + value
		return 3202
	}else if true==strings.HasPrefix(strErr,"ransaction nonce is too low") {
		//Transaction nonce is too low. Try incrementing the nonce.
		return 3403
	}else if true==strings.Contains(strErr,"nonce too low") {
		return 3503
	}
	//else if true==strings.HasPrefix(strErr,"known transaction") {
	//	//return 3203
	//	return 0
	//}
	return  0
}
func (self* EtherSever)DoTransactionSuf(client base.EtherClientHandle,tran base.EtherTranBinary,bWait bool) ( int, error)  {
	//ethclient,bOk:=client.(* EtherumWrapper)
	//if bOk {
	if true {
		Log.Error("%v",client.(* EtherumWrapper))
	}
	//conn:=ethclient.client
	const  imax int =100
	iter:=0
	for ;iter<100;iter+=1 {
		//ctx, _ := context.WithTimeout(context.Background(), time.Duration(base.ETH_RPC_TIMEOUT*60)* time.Second)
		//err := conn.SendTransactionRaw(ctx, tran.Raw)
		//sgj 1127duong
		//err := conn.SendTransaction(ctx, tran.Raw)
		//做监测所用，目前不需要
		var err error = nil
		if nil ==err {
			return  proto.ErrorSuccess.Code,nil
		}
		Log.Error("%s.Tx Error: %v, txhash=%s, raw=%s",self.Coin, err,tran.Hash,tran.Raw)
		strErr:=err.Error()
		//geth error:*known transaction*
		if true==strings.Contains(strErr,"known transaction"){
			return  proto.ErrorSuccess.Code,err
		}
		//2018-09-26
		//parity error:Transaction with the same hash was already imported.
		if true==strings.Contains(strErr,"same hash"){
			return  proto.ErrorSuccess.Code,err
		}
		/*sgj 1127doing
		if false==bWait {
			return  proto.ErrorOnlyDtalk_START+57,errors.New("ErrorServerStatus::"+err.Error())
		}
		*/
		iErr:=IsReturnError(err)
		if 0!=iErr {
			return  iErr,err
		}
		time.Sleep(time.Duration(8)*time.Second)
	}
	/*sgj 1124
	if iter>=imax {
		return  ErrorOnlyDtalk_START+63,errors.New(fmt.Sprintf("%s.Too  Times  to SendTransactionRaw",self.Coin))
	}
	*/
	return proto.ErrorSuccess.Code,nil
}

func (self* EtherSever)NewClient(bWait bool) (base.EtherClientHandle,error) {
	xx:=&EtherumWrapper{}
	var err error
	for {
		xx.client, err = ethclient.Dial(self.StrRawUrl)
		if err == nil {
			Log.Info("Success Connect  %s",self.StrRawUrl)
			break
		}
		if true!=bWait {
			return nil,err
		}
		Log.Error("Failed Connect %s ::: %v",self.StrRawUrl, err)
		time.Sleep(time.Duration(config.SleepMilliSecond)*time.Millisecond)
	}
	return     xx,nil
}


