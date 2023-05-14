package ethtranssign

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
	"2019NNZXProj10/depositgatherserver/proto"
	//"2017opensource/beepay/service/base"
	"2019NNZXProj10/depositgatherserver/service/ethtranssign/base"
)


type EtherumWrapper struct  {
	client  * ethclient.Client
}

func (self * EtherumWrapper)TransactionBlockNumber(strHash string,retval *big.Int) (error)  {

	return nil

}
func (self * EtherumWrapper)TransactionReceipt(strHash string ,res * base.StructTransactionReceipt)(error) {


	ctx2, _ := context.WithTimeout(context.Background(), time.Duration(base.ETH_RPC_TIMEOUT*10) * time.Second)
	vHash:=common.HexToHash(strHash)
	txReceipt,err:=self.client.TransactionReceipt(ctx2,vHash)
	if nil!=err {
		return   err
	}
	//sgj skip---PrintMarshal("receipt=",false,txReceipt)
	res.TxHash=txReceipt.TxHash.String()
	res.ContractAddress=txReceipt.ContractAddress.String()
	res.GasUsed=txReceipt.GasUsed
	res.Status=uint64(txReceipt.Status)
	return  nil
}
func (self * EtherumWrapper)NonceAt(strAddress string) (uint64, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(base.ETH_RPC_TIMEOUT*10) * time.Second)
	return self.client.NonceAt(ctx,common.HexToAddress(strAddress),nil)
}
func (self * EtherumWrapper)PendingNonceAt(strAddress string)(uint64, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(base.ETH_RPC_TIMEOUT*10) * time.Second)
	return  self.client.PendingNonceAt(ctx,common.HexToAddress(strAddress))
}
func (self * EtherumWrapper)BalanceAt(strAddress ,subType string,rVal * big.Int)( proto.ErrorInfo) {
	/*
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(base.ETH_RPC_TIMEOUT*10) * time.Second)
	if ""==subType {
		bigCount,err:=self.client.BalanceAt(ctx,common.HexToAddress(strAddress),nil)

	*rVal=*bigCount
	mylog.Info("balance=%s, address=%s, contraction=%s ",bigCount.String(),strAddress,ti.Contract)
	*/
	return proto.ErrorSuccess


}
