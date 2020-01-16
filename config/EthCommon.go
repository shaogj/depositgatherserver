package config

import (
	"sync"
	"github.com/robfig/config"
	//. "dingjm/utils"
	//stringex "dingjm/utils/stringex"
	"os"
	"strings"
	"strconv"
	mylog "github.com/mkideal/log"

)

const (
	PasswordDefaultString string ="32z91xa"
	MaxChanTrans =200
	EthGwei=1000000000 				//1000,000,000
	EthEther=1000000000000000000   //1000,000,000,000,000,000
	EthEtherPrefix=1000000000   //用于保留浮点值的精度
	EthEtherSuffix=1000000000   //将其变成长整行
	SleepMilliSecond  =2000
	LEN_ADDRESS_STR_LEN =  42
	LEN_TXHASH_STR_LEN =  66

)
//sgj 1127 add
const (
	CoinEtherTrans string ="TRAN" //操作类型
	CoinEstimatePeriod string ="Period" //操作类型
)

//const bRemoteSign bool=true
var IsRemoteSign bool=false  //开启远程签名服务器 签名


type  SttFee   struct  {
	Limit int64
	Price int64
}

var (
	///GAS_LIMIT int64=21000
	//GAS_PRICE int64=18*EthGwei
	MapGas map[string] SttFee=map[string] SttFee{
		CoinEthereum+"."+CoinEtherTrans:{25000,22*EthGwei},
	}
	mtGas sync.Mutex
)
//add   by dingjianmin   18-8-28 下午9:11
//加载交易费用  limit and  price    ETH and  ETC
func LoadGas(conf * config.Config,sec string)    {

	if false==LoadGasItem(conf,sec,CoinEthereum ,CoinEstimatePeriod) {
		os.Exit(0)
	}

	if false==LoadGasItem(conf,sec,CoinEthereum ,CoinEtherTrans) {
		os.Exit(0)
	}
	//LoadGasItem(conf,"gas",CoinEthereum ,CoinEtherTrans)
	return
}
//func LoadEstimatePeriod(conf * config.Config,sec,cointype,fdes string) bool    {
//
//}
func LoadGasItem(conf * config.Config,sec,cointype,fdes string) bool    {

	key:=cointype+"."+fdes
	stret:=CfgGetKey(conf,sec,key,"",true)
	if ""==stret {
		mylog.Error("Can't Load Fee  %s.%s",cointype,fdes)
		return  false
	}
	strs:=strings.Split(stret,",")
	if 2!=len(strs)  {
		mylog.Error("Can't Load Fee  %s.%s",cointype,fdes)
		return  false
	}
	iva11,err1:=strconv.ParseInt(strs[0],10,64)
	iva12,err2:=strconv.ParseInt(strs[1],10,64)
	if nil!=err1 || iva11<1 {
		mylog.Error("Can't Load Fee  Value %s.%s",cointype,fdes)
		return  false
	}
	if nil!=err2 || iva12<1 {
		mylog.Error("Can't Load Fee Value %s.%s",cointype,fdes)
		return  false
	}
	sttFee:=SttFee{iva11,iva12}
	SetGas(cointype ,fdes,sttFee)
	mylog.Info("config gas key=%s value=%+v",key,sttFee)
	return true
}
//设置费用
func SetGas(cointype string,fdes string,vfee SttFee)   {
	mtGas.Lock()
	defer  mtGas.Unlock()
	str:=cointype+"."+fdes
	MapGas[str]=vfee
}
func GetGas(cointype string,fdes string ,rateLimit,ratePrice float64)  SttFee  {
	mtGas.Lock()
	defer  mtGas.Unlock()
	str:=cointype+"."+fdes
	sttFee,bok:=MapGas[str]
	if true!=bok  {
		return SttFee{0,0}
	}
	if 1.0!=rateLimit {
		sttFee.Limit=int64(rateLimit*float64(sttFee.Limit))
	}
	if 1.0!=ratePrice {
		sttFee.Price=int64(ratePrice*float64(sttFee.Price))
	}
	return sttFee
}

func CfgGetKeyInt(cfg * config.Config,sec,key string,def int, bPanic bool) int {
	ival:=CfgGetKeyInt64(cfg,sec,key,int64(def), bPanic)
	return  int(ival)
}
func CfgGetKeyInt64(cfg * config.Config,sec,key string,def int64, bPanic bool) int64 {
	val, err := cfg.String(sec, key)
	if err != nil {
		if true==bPanic {
			mylog.Error("Can't Get value {sec=%s,key=%s} error: %v",sec,key, err)
			panic("")
		}else {
			mylog.Info("Can't Get value {sec=%s,key=%s} error: %v",sec,key, err)
		}
		return  def
	}
	var ival int64
	ival,err=strconv.ParseInt(val,10,64)
	if err != nil {
		mylog.Error("Can't Get value {sec=%s,key=%s} error: %v",sec,key, err)
		panic("")
		return  def
	}
	return ival
}
