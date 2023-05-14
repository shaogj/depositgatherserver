package models

var (
	//账户信息表 wdc表存放keystore对应的字段值表
	TableWDCAccount     = "wdc_account_key"
	TableBTCAccount     = "gjc_account_key_tb"
	TableCoinPrivateKey = "coin_private_key"
	TableGGEXTranRecord = "ggex_tran_state"
	//0116add
	TableBTCTranRecord = "btc_tran_state"

	//2023
	TableBSCAccount = "bsc_account_key_mults"
)

type GjcAccountKeyTb struct {
	Id int `json:"id" xorm:"not null pk autoincr INT(11)"`
	//Uid           int    `json:"uid" xorm:"not null index INT(11)"`
	AccountName string `json:"accountname" xorm:"default ''"`
	CoinType    string `json:"cointype"` //交易币种类
	Walletid    int64  `json:"walletId" xorm:"BIGINT(20)"`
	//0926 add
	//Wallettype string    `xorm:"CHAR(32)"`
	PrivKey   string `json:"privkey" xorm:"not null TEXT"`
	PubKey    string `json:"pubkey" xorm:"not null TEXT"`
	AddressId string `json:"addressid" xorm:"not null TEXT"`
	//Txid,is made by last txout, to pay to for next time
	Utxoid      string `json:"utxoid" xorm:"default '' TEXT"`
	CreatedTime int64  `json:"created_time" xorm:"BIGINT(20)"`
	Status      int    `json:"status" xorm:"default 0 index INT(11)"`
	UpdatedTime int64  `json:"updated_time" xorm:"BIGINT(20)"`
}

// 10128add
type WdcAccountKey struct {
	Id         int    `xorm:"not null pk autoincr comment('ID') INT(11)"`
	Wallettype string `xorm:"comment('account密码') CHAR(32)"`
	Walletid   string `xorm:"unique(walletid) CHAR(120)"`
	CoinType   string `xorm:"comment('币种类型') CHAR(40)"`
	Address    string `xorm:"comment('账户公钥hash地址') unique(address) VARCHAR(200)"`
	//Status     int       `xorm:"INT(11)"`
	PrivKey string `xorm:"comment('账户私钥') VARCHAR(1024)"`
	PubKey  string `xorm:"comment('账户公钥') VARCHAR(200)"`
	//PubKeyHash string `xorm:"comment('账户公钥hash') VARCHAR(200)"`

	TimeCreate int64 `json:"created_time" xorm:"BIGINT(20)"`
	TimeUpdate int64 `json:"created_time" xorm:"BIGINT(20)"`
}

type WdcTranRecord struct {
	Settleid  int64   `json:"orderid" `
	Txhash    string  `json:"txhash" xorm:"default ''"`
	From      string  `json:"from" xorm:"default ''"`
	To        string  `json:"to" xorm:"default ''"`
	Amount    float64 `json:"amount"`
	Amountfee float64 `json:"amountfee"`
	Coincode  string  `json:"cointype" xorm:"default ''"` //交易币种类
	Status    string  `json:"status" xorm:"default ''"`

	Verifystatus int    `json:"verifystatus" xorm:"default 0 index INT(4)"` //交易审核状态
	Errcode      int64  `json:"errorid" `
	Desc         string `json:"desc" xorm:"default ''"`
	//改为varchar类型在mysql里；
	TimeCreate string `json:"time_create" `
	TimeUpdate string `json:"time_update" `
	Raw        string `json:"raw" `
}
