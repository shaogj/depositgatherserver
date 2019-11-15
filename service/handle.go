package service

import (
	"2019NNZXProj10/abitserverDepositeGather/config"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strings"
	//"math/big"
	//"strconv"
	"bytes"
	"encoding/json"

	"2019NNZXProj10/abitserverDepositeGather/accounts"
	transproto "2019NNZXProj10/abitserverDepositeGather/proto"
	"2019NNZXProj10/abitserverDepositeGather/service/wdctranssign"

	"2019NNZXProj10/abitserverDepositeGather/KeyStore"

	"github.com/gorilla/mux"
	"github.com/mkideal/log"
	. "shaogj/utils"
	"time"
	"2019NNZXProj10/abitserverDepositeGather/cryptoutil"

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

//跨域访问
func HttpExCrossDomainAccess(w *http.ResponseWriter) {
	// 允许跨域访问
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
}

//请求数据json
//请求数据转为json数据
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

func CreateAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//cmdName := "CreateAddress"
	coinType := vars["cointype"]
	var status int
	var desc string
	var genPubKeyAddr string
	fmt.Println("vim-go")
	if coinType == "CoinDSC" {
		//sgj 0821 gooding:
		//getprivkey, getpubKey, pubKeyAddr, err := accounts.AddressGenerateDSC()

		//sgj 1104 temp skip:
		getprivkey := "getprivkey"
		getpubKey := "getpubKey"
		pubKeyAddr := "pubKeyAddr"
		var err error = nil
		//genPubKeyAddr = pubKeyAddr
		if err != nil {
			log.Error("exec AddressGenerateQtum() Failure!")
			status = StatusNewAddressErr
			desc = "生成账号地址错误"
		} else {
			log.Info("doing--AddressGenerateDSC() exec succ!,get pubKeyAdd is :%s,pubKey is %s,puraddrPrikey info is :%v\n", pubKeyAddr, getpubKey, getprivkey)
			//err = GenerateAccount(qtumtranssign.Qtum_MOrmEngine,"QTUM",getprivkey,getpubKey,pubKeyAddr)

		}
	}
	makeaddrs := GenerateAddressRes{
		GeneratedAddr: []string{genPubKeyAddr},
		CoinType:      coinType,
	}

	GeneJsonResultFin(w, r, makeaddrs, status, desc)
}

//sgj 1017add for RemoteSignCreateAddress
/*
desc 服务器: 成生新的地址(批量)
请求报文
curl -d '{"cointype": "CoinDSC","accountNumber":4,"IsReturnList": 1}'  http://192.168.1.166:3377/remote/getnewaddress
*/

//一个默认的合作商Key
var GCurGetKeyStr =[]byte("1234567812345678")

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
		if config.CoinWDC == jReq.CoinType {
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
				}else{
					encrptedEncodePrivkey = base64.StdEncoding.EncodeToString(encrpted)
					log.Info("get encrptedEncodeStr len is:%d,val is====44:%s",len(encrptedEncodePrivkey),encrptedEncodePrivkey)
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
		} else if "CoinDSC" == jReq.CoinType {
			// sgj 1019add

			getprivkey, getpubKey, pubKeyAddr, err := accounts.AddressGenerateDSC()
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

//sgj 1114add for
//DepositAddressGatherReq
func RemoteMonitorWalletAddress(w http.ResponseWriter, r *http.Request) {

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
		if "WDC" == jReq.CoinType {
			//WithdrawsDepositGatherWDC
			//1113 测试归集的服务接口调用
			gatherAddrCount, bret = wdctranssign.WithdrawsDepositGatherWDC(0, 30, "WDC")

			if bret != true {
				log.Error("cur exec WithdrawsDepositGatherWDC() err! get gatherAddrCount is:%d\n,err is :%v", gatherAddrCount,"errinfomsgskip")
			} else {
				log.Info(" cur exec WithdrawsDepositGatherWDC() succ!,get gatherAddrCount is :%d\n", gatherAddrCount)
			}
		} else {
			GeneJsonResultFin(w, r, nil, 111096, "error  is cointype")
		}
		//sgj 0802 add
		log.Info("exec WithdrawsDepositGatherWDC() success! jReq.CoinType is :%s,gatherAddrCount is: %d", jReq.CoinType, gatherAddrCount)

	makeaddrs := transproto.DepositAddressGatherRes{
		Count:         int64(gatherAddrCount),
		CoinType:      jReq.CoinType,
	}

	GeneJsonResultFin(w, r, makeaddrs, sttError.Code, sttError.Desc)
}

//const HActionSign     = "GGEX-ActionSign"

//sgj 1025doing,send request query to settlecenter
func WithdrawsTransTotal(offset, limit uint, cointype string) (uint, []interface{}) {
	var reqInfo transproto.WithdrawsQueryReq

	ht := CHttpClientEx{}
	//sgj add
	ht.Init()
	ht.HeaderSet("Content-Type", "application/json;charset=utf-8")

	//reqInfo.MaxVol = 0
	reqInfo.Limit = int(limit)
	reqInfo.Status = transproto.SETTLE_STATUS_PASSED
	reqInfo.CoinCode = cointype
	reqInfo.Nonce = time.Now().Unix()
	reqInfo.Offset = int(offset)
	//sgj 1028,,查询settle 的提现类型记录状态为：SETTLE_STATUS_PASSED
	//获取交易所需要的详细信息

	UrlVerify := config.GbConf.SettleApiQuery
	//UrlVerify := config.GbConf.SettleApiReq.SettlApiQuery
	var signData string
	for {
		//fix to this,每次变量初始化
		resQuerySign := transproto.Response{}
		transInfo := transproto.WithdrawsQueryResp{}
		resQuerySign.Data = &transInfo

		log.Info("withdrawsQuery.UrlVerify is:%s,reqInfo is:%v", UrlVerify, reqInfo)

		reqBody, err := json.Marshal(&reqInfo)
		if nil != err {
			log.Error("when withdrawsQuery,Marshal to json error:%s", err.Error())
			return 0, nil
		}
		if signData, err = auth.KSign(reqBody, GSettleAccessKey); err != nil {
			log.Error("In withdrawsQuery(),auth.KSign failed,signData is :%v,err is:%v", signData, err)
			return 0, nil
		}
		//step 2
		log.Info("withdrawsQuery,auth.KSign succ!,signData is :%v", signData)

		//req.Header.Set("abit-actionsign", signData)

		//
		//ht.HeaderSet(transproto.HActionSign, signData)
		//sgj 1112 adding,STDusing!!:
		ht.HeaderSet(transproto.HActionAbitSign, signData)

		strRes, statusCode, errorCode, err := ht.RequestJsonResponseJson(UrlVerify, 9000, &reqInfo, &resQuerySign)
		if nil != err {
			log.Error("ht.RequestResponseJsonJson  status=%d,error=%d.%v url=%s ", statusCode, errorCode, err, UrlVerify)
			time.Sleep(time.Second * 30)
			continue
		}
		log.Info("transserver.transInfo res=%s", strRes)
		log.Info("json.Unmarshal succ!,cur get resQuerySign :%v,get field code is :%s", resQuerySign, resQuerySign.Code)
		for ino, curSettItem := range transInfo.Withdraws {
			// 2 审核通过(运营审核通过)
			log.Info("get SettlApiQuery queue info, cur ino is:%d,SettItem record Status is:%v,curSettItem is :%v", ino, curSettItem.Status, curSettItem)
			curStatus := curSettItem.Status
			if curStatus != transproto.SETTLE_STATUS_PASSED {
				log.Error("get SettlApiQuery queue info, curSettItem record Status is:%s,is skiped!,curSettItem is :%v", curStatus, curSettItem)
				continue
			}
			//note,from 是大账户地址："1HFCUeNHcL6Drf4TPwBLG6RgYVe9o41BVj",
			if curSettItem.CoinCode == "WDC" {
				wdctranssign.GWDCTransHandle.WDCTransProc(curSettItem, config.GbConf.WDCTransferOutAddress, "exaaccountName")
			}

		}

		log.Info("transInfo len is:%d,finished!,to wait 30s", len(transInfo.Withdraws))

		time.Sleep(time.Second * 30)

	}
	return 1, nil

}
