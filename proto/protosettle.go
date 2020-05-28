package proto

import (
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type SETTLE_TYPE int

type SETTLE_STATUS int

const (
	_                      SETTLE_STATUS = iota
	SETTLE_STATUS_CREATED                // 1 申请成功(用户提交申请)
	SETTLE_STATUS_PASSED                 // 2 审核通过(运营审核通过)
	SETTLE_STATUS_REJECTED               // 3 审核拒绝(运营审核拒绝)
	SETTLE_STATUS_SIGNED                 // 4 签名完成(生成转账signstr完成)
	SETTLE_STATUS_PENDING                // 5 打包中(待确认链上是否转账成功)
	SETTLE_STATUS_SUCCESS                // 6 成功(转账成功)
	SETTLE_STATUS_FAILED                 // 7 失败(转账失败)
)
const (
	StatusNewAddressErr = 201 //  生成账号地址错误
)

//1217add

//11.19add:
// enum: ReturnStatus
// 返回状态码
const (
	//StatusSuccess = 0 // 成功
	StatusEncodeJSONFail  = 102 // 编码json出错
	StatusDecodeJSONFail  = 103 // 解码json出错
	StatusParseError      = 104 // 解析出错
	StatusCommandNotFound = 105 // 命令未找到

	StatusSuccess = 200 // 调用成功

	StatusDataSelectErr   = 205 //  查询数据库错误
	SStatusAddrNotFound   = 206 // 用户地址不存在
	StatusNewAddrNotExist = 207 // 此类型新地址不存在

	StatusLackBalance             = 301 // 地址余额不足
	StatusLackUTXO                = 303 // 地址充值余额信息未发现
	StatusAccountNotExisted       = 304 // 地址不存在
	StatusAccountPrikeyNotExisted = 305 //地址私钥不存在

	StatusSignError       = 400 // 交易签名失败
	StatusInvalidArgument = 401 // 无效输入参数
	//StatusSignError       		= 400 // 交易签名错误
	//交易发送
	StatusShareOrderThemeErr = 500 //交易广播错误
	//sgj 1018 add:
	StatusUtxoTxWaiting = 707 // 存在当前地址的交易；等待新的Txid推送返回

)

type UnfreezeFundsStatus int

const HActionSign = "EX-ActionSign"

//sgj 1112 fix add:
const HActionAbitSign = "EX-ActionSign"

//1217
const HActionTMexSign = "EX-ActionSign"

//sgj 0106add
const HActionEETHSign = "EX-actionSign"

//1.函数名:	GenerateMultiAddress
// protocol: 请求: 生成数字支付地址
type GenerateAddressReq struct {
	//CoinType    int   `json:"coinType"`     // 数字币类型
	//Number 		int   `json:"number"`  // 生成数量
	CoinType     string `json:"coinType"`
	Count        int64  `json:"accountNumber"`
	IsReturnList int    `json:"IsReturnList"`
}

//11.14,开放充值归集指令的接口，面向opsTools：
// Encrypt
type DepositAddressGatherReq struct {
	EncryptPemTxt string `json:"encryptPemTxt"` // 归集指令私钥密文
	KeyText       string `json:"keyText"`       // 归集指令私钥秘钥
	CoinType      string `json:"coinType"`
}

// protocol: 返回: 执行充值归集指令的结果
type DepositAddressGatherRes struct {
	Count    int64  `json:"count"`
	CoinType string `json:"coinType"`
}

//end add 11.14

// protocol: 返回: 生成数字支付地址
type GenerateAddressRes struct {
	Count         int64    `json:"count"`
	GeneratedAddr []string `json:"getNewAddr"` // 生成地址
	CoinType      string   `json:"coinType"`
}

type Settle struct {
	ID            int64               `json:"-" xlsx:"-"`
	SettleId      int64               `json:"settle_id,omitempty" gorm:"default:NULL" xlsx:"#"`
	AccountId     int64               `json:"account_id,omitempty" gorm:"default:NULL" xlsx:"-"`
	FromAddress   string              `json:"from_address,omitempty" gorm:"default:NULL" xlsx:"From"`
	ToAddress     string              `json:"to_address,omitempty" gorm:"default:NULL" xlsx:"To"`
	CoinCode      string              `json:"coin_code,omitempty" gorm:"default:NULL" xlsx:"货币"`
	BlockId       string              `json:"block_id,omitempty" gorm:"default:NULL" xlsx:"-"`
	BlockHash     string              `json:"block_hash,omitempty" gorm:"default:NULL" xlsx:"-"`
	TxHash        string              `json:"tx_hash,omitempty" gorm:"default:NULL"`
	SignStr       string              `json:"sign_str,omitempty" gorm:"default:NULL" xlsx:"-"`
	QueryData     string              `json:"query_data,omitempty" gorm:"default:NULL" xlsx:"-"`
	Type          SETTLE_TYPE         `json:"type,omitempty" gorm:"default:NULL" xlsx:"类型;enum:-,充值,提现,赠送,空投,转入,转出,合约化入,合约化出,OTC转入,OTC转出,合约云转入,合约云转出,冻结,销毁,盈利"`
	Status        SETTLE_STATUS       `json:"status,omitempty" gorm:"default:NULL" xlsx:"状态;enum:-,待审核,审核通过,审核拒绝,签名完成,打包中,成功,失败"`
	Vol           decimal.Decimal     `json:"vol,omitempty" gorm:"type:decimal(36,18)" xlsx:"数额"`
	Fee           decimal.Decimal     `json:"fee,omitempty" gorm:"type:decimal(36,18)" xlsx:"手续费"`
	FeeCoinCode   string              `json:"fee_coin_code,omitempty" gorm:"default:NULL" xlsx:"手续费货币"`
	UnfreezeFunds UnfreezeFundsStatus `json:"unfreeze_funds" xlsx:"冻结状态;enum:-,已解冻,冻结中"`
	Auditor       string              `json:"auditor,omitempty" xlsx:"审核"`
	RejectReason  string              `json:"reject_reason,omitempty" xlsx:"拒绝原因"`
	Error         string              `json:"error,omitempty" gorm:"default:NULL" xlsx:"错误"`
	CreatedAt     time.Time           `json:"created_at,omitempty" gorm:"default:NULL" xlsx:"提交时间"`
	UpdatedAt     time.Time           `json:"updated_at,omitempty" gorm:"default:NULL" xlsx:"-"`
	Memo          string              `json:"memo,omitempty" gorm:"default:NULL" xlsx:"memo"`
}

type WithdrawsQuery struct {
	Status        SETTLE_STATUS       `json:"status,omitempty"`
	SettleId      int64               `json:"settle_id,omitempty"`
	Limit         int                 `json:"limit,omitempty"`
	Offset        int                 `json:"offset,omitempty"`
	CoinCode      string              `json:"coin_code"`
	UnfreezeFunds UnfreezeFundsStatus `json:"unfreeze_funds"`
	MaxVol        decimal.Decimal     `json:"max_vol"`
	MinVol        decimal.Decimal     `json:"min_vol"`
}

type WithdrawsQueryReq struct {
	WithdrawsQuery
	Nonce int64 `json:"nonce,omitempty"`
}

type WithdrawsQueryResp struct {
	Total     int      `json:"total,int"`
	Withdraws []Settle `json:"withdraws,omitempty"`
	Nonce     int64    `json:"nonce,omitempty"`
}

type WithdrawsUpdateReq struct {
	Withdraws []Settle `json:"withdraws,omitempty"`
	Nonce     int64    `json:"nonce,omitempty"`
}

type WithdrawsUpdateResp struct {
	Nonce int64 `json:"nonce,omitempty"`
}

///sgj 1112 add for Deposit of server:
type DepositeAddresssReq struct {
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
	CoinCode string   `json:"coin_code,omitempty"`
	Start    int64    `json:"start,omitempty"`
	End      int64    `json:"end,omitempty"`
	Nonce    int64    `json:"nonce,omitempty"`
	Ignore   []string `json:"-"`
}
type DepositeAddresssResp struct {
	Total    int      `json:"total,int"`
	Addresss []string `json:"addresss,omitempty"`
	Nonce    int64    `json:"nonce,omitempty"`
}

type WITHDRAWAL_CONFIG_STATUS int32

///sgj 1113 add for DepositConfig of server:

type WithDrawConfigReq struct {
	Nonce int64 `json:"nonce,omitempty"`
}
type WithDrawalConfig struct {
	ID                 int64                    `json:"-" gorm:"primary_key"`
	CoinName           string                   `json:"coin_name"`
	Fee                decimal.Decimal          `json:"fee" gorm:"type:decimal(36,18)"`
	FeeCoinName        string                   `json:"fee_coin_name"`
	MinVol             decimal.Decimal          `json:"min_vol" gorm:"type:decimal(36,18)"`
	DepositMinVol      decimal.Decimal          `json:"deposit_min_vol" gorm:"type:decimal(36,18)"`
	KYCMaxVol          decimal.Decimal          `json:"kyc_max_vol" gorm:"type:decimal(36,18)"` // 没有通过KYC认证的最大提现限制
	Status             WITHDRAWAL_CONFIG_STATUS `json:"-"`
	CoinPrecision      int64                    `json:"coin_precision"`
	WalletMaxVol       decimal.Decimal          `json:"wallet_max_vol" gorm:"type:decimal(36,18);DEFAULT:0"`        //向大钱包充值的最大值
	MiddleWalletMaxVol decimal.Decimal          `json:"middle_wallet_max_vol" gorm:"type:decimal(36,18);DEFAULT:0"` //中间归集钱包(OKEX钱包)的最大值
	CreatedAt          *time.Time               `json:"-"`
	UpdatedAt          *time.Time               `json:"-"`
}

type CoinDetailConfig struct {
	WithDrawalConfig
	CoinGroup     string `json:"coin_group"`
	TokenAddress  string `json:"token_address"`
	TokenDecimals int32  `json:"token_decimals"`
}
type WithdrawConfigResp struct {
	Total   int32               `json:"total"`
	Configs []*CoinDetailConfig `json:"configs"`
}

type Response struct {
	HTTPCode     int           `json:"-"`
	Code         string        `json:"errno"`
	Msg          string        `json:"message"`
	Header       http.Header   `json:"-"`
	Data         interface{}   `json:"data,omitempty"`
	IsGZip       bool          `json:"-"`
	IsResetToken bool          `json:"-"`
	MsgData      []interface{} `json:"-"`
}

type ErrorInfo struct {
	Code int
	Desc string
}

var (
	//请求参数错误
	ErrorRequest = ErrorInfo{Code: 1001, Desc: " 请求参数无效"}

	ErrorRequestWDCNodeRPC = ErrorInfo{Code: 7000, Desc: "请求节点RPC参数无效(WDC)"}
	//1030add

	ErrorCoinType = ErrorInfo{Code: 1004, Desc: "数字货币类型错误"}

	ErrorNodeRPCSuccess = ErrorInfo{Code: 2000, Desc: "调用节点RPC成功"}
	ErrorRequestWDCNode = ErrorInfo{Code: 5000, Desc: "请求Node错误(WDC)"}

	ErrorRequestWDCNodeJust = ErrorInfo{Code: 7000, Desc: "请求Node校验错误(WDC)"}
	/*
			2000 正确
		    2100 已确认
		    2200 未确认
		    5000 错误
		    6000 格式错误
		    7000 校验错误
		    8000 异常
	*/

	ErrorAddress = ErrorInfo{Code: 701, Desc: "无效的地址码"}
	ErrorHttpost = ErrorInfo{Code: 801, Desc: "http请求必须为POST方式"}

	ErrorSuccess       = ErrorInfo{Code: 200, Desc: "调用成功"}
	ErrorRequestWDCSDK = ErrorInfo{Code: 500, Desc: "调用SDK参数无效(WDC)"}

	ErrorGetPrivateKey = ErrorInfo{Code: 604, Desc: "无法获取用户私钥"}

	ErrDecryptFail    = ErrorInfo{Code: 3333, Desc: "解密json串失败"}
	ErrVerifyListFail = ErrorInfo{Code: 9189, Desc: "人工审核获取订单列表失败"}

	//1129 add
	ErrorRequestInfuraETHNode = ErrorInfo{Code: 7004, Desc: "请求Node错误(ETH)"}

	ErrorRequestInfuraETHSend = ErrorInfo{Code: 7005, Desc: "请求Node错误(ETH)"}
)

func Success(status ErrorInfo) bool {
	if ErrorSuccess.Code == status.Code {
		return true
	}
	return false

}

//sgj 20200109--add :
//https://blockchain.info/unspent?active=
type BtcUtxoInfo struct {
	TxHash          string `json:"tx_hash"`
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxOutputN       int64  `json:"tx_output_n"`
	Script          string `json:"script"`

	Value         int64  `json:"value"`
	ValueHex      string `json:"value_hex"`
	Confirmations int64  `json:"confirmations"`
	TxIndex       int64  `json:"tx_index"`
}

//
type BTCUnspentOutputs struct {
	CurBtcUtxoInfo []BtcUtxoInfo `json:"unspent_outputs"`
}
