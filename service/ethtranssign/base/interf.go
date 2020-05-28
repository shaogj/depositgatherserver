package base

import (
	models "2019NNZXProj10/depositgatherserver/proto"
	"math/big"
)

type StructTransactionReceipt struct {
	Status            uint64
	CumulativeGasUsed uint64
	TxHash            string
	TransactionIndex  uint64

	ContractAddress string
	GasUsed         uint64
}
type EtherTranInfo struct {
	SubType  string //子类型
	From     string
	To       string
	Amount   float64
	GasPrice int64  //根据FeeAmount寄宿计算所得
	GasLimit uint64 //取自设置中
	UNonce   uint64
	Private  string //私钥
}
type EtherTranBinary struct {
	Hash string
	Raw  string
}

//sgjadd

type EtherIntf struct {
	EtherTranInfo
	IType      int     //交易类型
	GasFee     float64 //请求参数中获取的
	Balance    big.Int //钱包金额
	Nonce      uint64  //返回值
	Private    string  /*私钥*/
	PrivatePwd string  /*私钥 密码*/
	TxHash     string  /*请求 返回值 交易hash*/
	State      string
	ActBlock   int
	Err        models.ErrorInfo /*返回值 */
	Ret        chan int
}

func (self *EtherIntf) GetEtherTranInfo() EtherTranInfo {
	rec := EtherTranInfo{
		SubType:  self.SubType,
		From:     self.From,
		To:       self.To,
		Amount:   self.Amount,
		GasPrice: self.GasPrice,
		GasLimit: self.GasLimit,
		UNonce:   self.UNonce,
	}
	return rec

}
func (self *EtherIntf) GetEthTranState() models.EthTranState {

	rec := models.EthTranState{
		Txhash:   self.TxHash,
		From:     self.From,
		To:       self.To,
		Amount:   self.Amount,
		Gasfee:   self.GasFee,
		Gaslimit: int64(self.GasLimit),
		Gasprice: self.GasPrice,
		SubType:  self.SubType,
	}
	return rec
}

func (self *EtherIntf) Error(code int, desc string) {
	self.Err.Code = code
	self.Err.Desc = desc
}

type EtherClientHandle interface {
	//根据交易值 获取区块号
	TransactionBlockNumber(strHash string, retval *big.Int) error
	//获了
	TransactionReceipt(strHash string, res *StructTransactionReceipt) error

	NonceAt(strAddress string) (uint64, error)
	PendingNonceAt(strAddress string) (uint64, error)
	BalanceAt(string, string, *big.Int) models.ErrorInfo
}
