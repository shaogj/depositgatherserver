package service

import (
	"2019NNZXProj10/depositgatherserver/config"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	//"strings"
	//"math/big"
	//"strconv"
	"bytes"
	"encoding/json"

	"2019NNZXProj10/depositgatherserver/accounts"
	transproto "2019NNZXProj10/depositgatherserver/proto"
	//"2019NNZXProj10/depositgatherserver/service/wdctranssign"

	"2019NNZXProj10/depositgatherserver/cryptoutil"
	"github.com/mkideal/log"
	//. "shaogj/utils"
)

type ReturnInfo struct {
	//Cmd         string      `json:"cmd"`      // 命令名,具有协议类型的作用
	InvokeResultCode    int         `json:"invokeResultCode"`    // 返回码(参见枚举 ReturnStatus)
	InvokeResultMessage string      `json:"invokeResultMessage"` // 返回码描述
	Data                interface{} `json:"data"`                // 协议数据
}

// protocol: 返回: 生成数字支付地址
type GenerateAddressRes struct {
	Count         int64    `json:"count"`
	GeneratedAddr []string `json:"getNewAddr"` // 生成地址
	CoinType      string   `json:"coinType"`
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func JSONResponseWithStatus(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	JSONResponse(w, data)
}

func GeneJsonResultFin(w http.ResponseWriter, r *http.Request, protostruct interface{}, status int, description string) {

	res := ReturnInfo{}
	res.InvokeResultMessage = description
	res.InvokeResultCode = status
	//res.Cmd = cmdname
	res.Data = protostruct
	buf := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(buf)
	err := jsonEncoder.Encode(res)
	if err != nil {
		fmt.Fprintln(w, "command %s:  result to json error: %v", res, err)
		w.Write([]byte(`{"invokeResultCode":999999,"invokeResultMessage":""}`))
	} else {
		//w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		JSONResponseWithStatus(w, res, http.StatusOK)
	}
}

const StatusNewAddressErr = 201 //  生成账号地址错误

var GSettleAccessKey string

// 跨域访问
func HttpExCrossDomainAccess(w *http.ResponseWriter) {
	// 允许跨域访问
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
}

// 请求数据json
// 请求数据转为json数据
func HttpExRequestJson(w http.ResponseWriter, r *http.Request, v interface{}) (string, transproto.ErrorInfo) {
	HttpExCrossDomainAccess(&w)
	if r.Method != "POST" {
		return "", transproto.ErrorHttpost
	}
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return "", sttErr
	}
	err = json.Unmarshal(result, v)
	if nil != err {
		sttErr := transproto.ErrorRequest
		sttErr.Desc = fmt.Sprintf("%s %s", transproto.ErrorRequest.Desc, err)
		return string(result), sttErr
	}
	return string(result), transproto.ErrorSuccess
}

//sgj 1017add for RemoteSignCreateAddress
/*
desc 服务器: 成生新的地址(批量)
请求报文
curl -d '{"cointype": "CoinDSC","accountNumber":4,"IsReturnList": 1}'  http://192.168.1.166:3377/remote/getnewaddress
*/

// 一个默认的合作商Key
var GCurGetKeyStr = []byte("1234567812345678")

func RemoteSignCreateAddress(w http.ResponseWriter, r *http.Request) {

	jReq := transproto.GenerateAddressReq{}
	strreq, sttErr := HttpExRequestJson(w, r, &jReq)
	log.Info("fun=RemoteSignCreateAddress() bef--,request=%v", jReq)
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		return
	}
	if jReq.Count <= 0 {
		GeneJsonResultFin(w, r, nil, 110098, "accountNumber is zero")
		return
	}
	if jReq.Count > 1000 {
		GeneJsonResultFin(w, r, nil, 110096, "accountNumber is too big")
		return
	}
	log.Info("fun=RemoteSignCreateAddress(),request=%s", strreq)
	//创建请求列表
	sttError := transproto.ErrorSuccess
	var getnewaddr []string = []string{}

	for i := int64(0); i < jReq.Count; i++ {
		strAdrress := ""
		//20200613update
		if config.CoinWDC == jReq.CoinType || "WGC" == jReq.CoinType {
			//1028add,调用json请求，创建地址
			getprivkey, getpubKey, getpubKeyHash, pubKeyAddr, err := accounts.AddressGenerateWDC()
			if err != nil {
				log.Error("cur exec AddressGenerateWDC() err! get privkey is:%v,pubKey is :%v\n,err is :%v", getprivkey, getpubKey, err)
			} else {
				//insert to db:Mysql
				log.Info(" cur exec AddressGenerateWDC() succ!,get pubKey is :%v\n", getpubKey)
				//保存记录到账户数据库
				//sgj 1115 add for encrytp:
				//11.15tesitn
				var encrptedEncodePrivkey string
				encrpted, err := cryptoutil.AESCBCEncrypt(GCurGetKeyStr, nil, []byte(getprivkey))
				if err != nil {
					log.Error("ccur Encrypt text is:%s,err is:%v", getprivkey, err)
				} else {
					encrptedEncodePrivkey = base64.StdEncoding.EncodeToString(encrpted)
					log.Info("get encrptedEncodeStr len is:%d,val is====44:%s", len(encrptedEncodePrivkey), encrptedEncodePrivkey)
				}

				//end add 1115

				err = AccountWDCSave(GXormMysql, "", accounts.GAccountPassword, jReq.CoinType, pubKeyAddr, encrptedEncodePrivkey, getpubKey, getpubKeyHash)
				if err != nil {
					log.Error("exec WDC GenerateAccount() failed! err is: %v", err)
					sttError = transproto.ErrorInfo{transproto.StatusNewAddressErr, "生成账号地址错误"}
					break
				}
			}
			strAdrress = pubKeyAddr //pubKeyAddr.String()
			//AccountWDCSave
		} else if "CoinBSC" == jReq.CoinType {
			strAdr, strPriv, err := accounts.EtherNewAccount()
			if nil != err {
				log.Error("get cur BSC' strAdr is :%v\n,err is :%v", strAdr, err)
				log.Error("get cur BSC' getprivkey is:%v,err is :%v\n", strPriv, err)
				break
			}
			strAdr = strings.ToLower(strAdr)

			err = GenerateAccountBSC(GXormMysql, "", accounts.GAccountPassword, jReq.CoinType, strAdr, strPriv, "getpubKey", "getpubKeyHash")

			//err = GenerateAccountBSC(GXormMysql, jReq.CoinType, strPriv, "getpubKey", strAdr)
			//err = GbDbMysql.RemoteSignAcccount(jReq.CoinType, strAdr, strPriv, strReqSign)
			if err != nil {
				log.Error("%s.GenerateAccount() error=%v", jReq.CoinType, err)
				sttError = transproto.ErrorInfo{transproto.StatusNewAddressErr, "生成账号地址错误"}
				break
			}
			strAdrress = strAdr

		} else if "CoinDSC" == jReq.CoinType {
			// sgj 1019add
			/*
				getprivkey, getpubKey, pubKeyAddr, err := accounts.EtherNewAccount()	//accounts.AddressGenerateDSC()
				if err != nil {
					log.Error("get cur DSC' pubKey is :%v\n,err is :%v", getpubKey, err)
					log.Error("get cur DSC' getprivkey is:%v,err is :%v\n", getprivkey, err)

				} else {
					//insert to db:GXormMysql
					log.Info(" cur exec AddressGenerateDSC() succ!,get pubKey is :%v\n", getpubKey)
					err = GenerateAccount(GXormMysql, jReq.CoinType, getprivkey, getpubKey, pubKeyAddr)
					if err != nil {
						log.Error("exec DSC GenerateAccount() failue! err is: %v", err)
						sttError = transproto.ErrorInfo{transproto.StatusNewAddressErr, "生成账号地址错误"}
						break
					}
				}
				strAdrress = pubKeyAddr //pubKeyAddr.String()
			*/
		} else {
			GeneJsonResultFin(w, r, nil, 111096, "error  is cointype")
		}
		if 1 == jReq.IsReturnList {
			getnewaddr = append(getnewaddr, strAdrress)
		}
		//sgj 0802 add
		log.Info("exec GenerateAccount() success! jReq.CoinType is :%s,generate addr is: %s", jReq.CoinType, strAdrress)

	}
	makeaddrs := transproto.GenerateAddressRes{
		Count:         int64(len(getnewaddr)),
		GeneratedAddr: getnewaddr,
		CoinType:      jReq.CoinType,
	}

	GeneJsonResultFin(w, r, makeaddrs, sttError.Code, sttError.Desc)
}

// sgj 1114add for
// DepositAddressGatherReq

/*
func RemoteMonitorWalletAddress(w http.ResponseWriter, r *http.Request) {

	jReq := transproto.DepositAddressGatherReq{}
	strreq, sttErr := HttpExRequestJson(w, r, &jReq)
	log.Info("fun=RemoteSignCreateAddress() bef--,request=%v", jReq)
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		return
	}
	//20200614add,,for WGCFee
	if jReq.EncryptPemTxt == "EncryptPemTxt2019WDC1114val" || jReq.EncryptPemTxt == "EncryptPemTxt2020WGCFee1114val" {
		//if jReq.EncryptPemTxt != "EncryptPemTxt2019WDC1114val"
	} else {
		GeneJsonResultFin(w, r, nil, 110098, "EncryptPemTxt is nocorrect!")
		return
	}
	if jReq.KeyText != "UCt38sGmp" {
		GeneJsonResultFin(w, r, nil, 110098, "KeyText is nocorrect!")
		return
	}
	log.Info("fun=RemoteMonitorWalletAddress(),request=%s", strreq)
	//创建请求列表
	sttError := transproto.ErrorSuccess
	var gatherAddrCount int
	var gatherWGCCount, gatherWDCCount int
	var bret bool
	//20200611 add for WGC
	if "WDC" == jReq.CoinType || "WGC" == jReq.CoinType {
		//WithdrawsDepositGatherWDC
		//1113 测试归集的服务接口调用

		//1204,limit set to 50
		//0611update,,CoinType,"WDC"
		gatherAddrCount, bret = wdctranssign.WithdrawsDepositGatherWDC(0, 50, jReq.CoinType)

		if bret != true {
			log.Error("cur exec WithdrawsDepositGatherWDC() err! get gatherAddrCount is:%d\n,err is :%v", gatherAddrCount, "errinfomsgskip")
		} else {
			log.Info(" cur exec WithdrawsDepositGatherWDC() succ!,get gatherAddrCount is :%d\n", gatherAddrCount)
		}
	} else if "WGCWDCAll" == jReq.CoinType {
		//sgj 2020066 add,合并归集所有的WGC，WDC的地址；
		gatherWGCCount, gatherWDCCount, bret = wdctranssign.WithdrawsDepositGatherWGCWDCAddrAll(0, 50, jReq.CoinType)

		if bret != true {
			log.Error("cur exec WithdrawsDepositGatherWGCWDCAddrAll() err! get gatherAddrCount is:%d\n,err is :%v", gatherAddrCount, "errinfomsgskip")
		} else {
			log.Info(" cur exec WithdrawsDepositGatherWGCWDCAddrAll() succ!,get gatherAddrCount is :%d\n", gatherAddrCount)
		}

	} else if "WGCFee" == jReq.CoinType {
		//sgj 20200614 add
		gatherAddrCount, bret = wdctranssign.WithdrawsDepositGatherWDCFee(0, 50, jReq.CoinType, jReq.FeeAmount, jReq.FeeThreshold)

		if bret != true {
			log.Error("cur exec WithdrawsDepositGatherWDCFee() err! get gatherAddrCount is:%d\n,err is :%v", gatherAddrCount, "errinfomsgskip")
		} else {
			log.Info(" cur exec WithdrawsDepositGatherWDCFee() succ!,get gatherAddrCount is :%d\n", gatherAddrCount)
		}

	} else {

		GeneJsonResultFin(w, r, nil, 111096, "error  is cointype")
	}
	//sgj 0802 add
	log.Info("exec WithdrawsDepositGatherWDC() success! jReq.CoinType is :%s,gatherAddrCount is: %d", jReq.CoinType, gatherAddrCount)

	makeaddrs := transproto.DepositAddressGatherRes{
		Count:    int64(gatherAddrCount),
		CoinType: jReq.CoinType,
	}
	if "WGCWDCAll" == jReq.CoinType {
		makeaddrsAll := transproto.DepositAddressGatherResAllCount{
			WGCGatherCount: int64(gatherWGCCount),
			WDCGatherCount: int64(gatherWDCCount),
			CoinType:       jReq.CoinType,
		}
		GeneJsonResultFin(w, r, makeaddrsAll, sttError.Code, sttError.Desc)

	} else {
		GeneJsonResultFin(w, r, makeaddrs, sttError.Code, sttError.Desc)

	}
}

*/

//sgj 1218add
//DepositAddressGatherReq

// sgj 0116add for EETH gather:WithdrawsDepositGatherEETH
func RemoteMonitorWalletEETHAddress(w http.ResponseWriter, r *http.Request) {

	jReq := transproto.DepositAddressGatherReq{}
	strreq, sttErr := HttpExRequestJson(w, r, &jReq)
	log.Info("fun=RemoteSignCreateAddress() bef--,request=%v", jReq)
	if true != transproto.Success(sttErr) {
		GeneJsonResultFin(w, r, nil, sttErr.Code, sttErr.Desc)
		return
	}
	if jReq.EncryptPemTxt != "EncryptPemTxt2019WDC1114val" {
		GeneJsonResultFin(w, r, nil, 110098, "EncryptPemTxt is nocorrect!")
		return
	}
	if jReq.KeyText != "UCt38sGmp" {
		GeneJsonResultFin(w, r, nil, 110098, "KeyText is nocorrect!")
		return
	}
	log.Info("fun=RemoteMonitorWalletAddress(),request=%s", strreq)
	//创建请求列表
	sttError := transproto.ErrorSuccess
	var gatherAddrCount int
	var bret bool
	if "BTC" == jReq.CoinType || "USDT" == jReq.CoinType || "ETH" == jReq.CoinType {
		//WithdrawsDepositGatherWDC
		//1217 测试归集的服务接口调用

		//1204,limit set to 50;;;
		//1230,limit set to 50,,from 4
		//	gatherAddrCount, bret = ktctranssign.WithdrawsDepositGatherEETH(0, 50, jReq.CoinType)
		//2023.0323
		//gatherAddrCount, bret = WithdrawsDepositGatherEETH(0, 50, jReq.CoinType)

		if bret != true {
			log.Error("cur exec WithdrawsDepositGatherEETH() err! get gatherAddrCount is:%d\n,err is :%v", gatherAddrCount, "errinfomsgskip")
		} else {
			log.Info(" cur exec WithdrawsDepositGatherEETH() succ!,get gatherAddrCount is :%d\n", gatherAddrCount)
		}
	} else {
		GeneJsonResultFin(w, r, nil, 111096, "error  is cointype")
	}
	//sgj 0802 add
	log.Info("exec WithdrawsDepositGatherEETH() success! jReq.CoinType is :%s,gatherAddrCount is: %d", jReq.CoinType, gatherAddrCount)

	makeaddrs := transproto.DepositAddressGatherRes{
		Count:    int64(gatherAddrCount),
		CoinType: jReq.CoinType,
	}

	GeneJsonResultFin(w, r, makeaddrs, sttError.Code, sttError.Desc)
}

//const HActionSign     = "GGEX-ActionSign"
