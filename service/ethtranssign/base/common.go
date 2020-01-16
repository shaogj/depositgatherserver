package base

import (
	"math/big"
	//"fmt"
	//"strings"
	//"strconv"
	"2019NNZXProj10/depositgatherserver/config"
	"github.com/ethereum/go-ethereum/common"
)
const (
	TypeLocal int=1
	TypeWebSocket int=2
	TypeHttp int=3

)

const (
	StatusActive int=0
)
const (

	EthTranTransaction int=1
	EthTranBalance int=2		//获取交易余额
	EthTranStatus int=3			//
	EthNonceValue int=4		//获取交易
)
const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)
const   ETH_RPC_TIMEOUT  int=16

var (
	AddressEmpty  common.Address = common.Address{}
)
//type ErrorInfo struct {
//	Code  int
//	Desc  string
//}
func  EtherToWei(fv float64) big.Int {
	var bg big.Int
	var b1 big.Int
	var b2 big.Int
	b1.SetInt64(int64(fv*config.EthEtherPrefix))
	b2.SetInt64(config.EthEtherSuffix)
	bg.Mul(&b1,&b2)
	return bg
}
func EtherFloat64( bv big.Int) float64 {
	var b2 big.Int
	var div big.Int
	div.SetInt64(config.EthEtherSuffix)
	b2.Div(&bv,&div)
	return  float64(b2.Int64())/config.EthEtherPrefix
}




