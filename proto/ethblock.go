package proto

import (
	"encoding/json"
	//"github.com/Appscrunch"
	"math/big"
	"strings"
	"github.com/shopspring/decimal"

	"time"
)

type ETHADDRESS string

func (addr ETHADDRESS) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(string(addr)))
}

func (addr *ETHADDRESS) UnmarshalJSON(b []byte) error {
	var address string
	if (len(b) == 42) && ((string(b[:2]) == "0x") || string(b[:2]) == "0X") {
		if _, ok := new(big.Int).SetString(string(b[2:]), 16); !ok {
			return nil
		}
		address = strings.ToLower(string(b))
	} else if len(b) == 40 {
		if _, ok := new(big.Int).SetString(string(b), 16); !ok {
			return nil
		}
		address = "0x" + strings.ToLower(string(b))
	} else if len(b) > 42 {
		index := strings.Index(string(b), "0x")
		if len(b)-index-1 != 42 {
			return nil
		}
		address = strings.ToLower(string(b)[index : index+42])
	} else {
		return nil
	}
	*addr = ETHADDRESS(address)
	return nil
}

func FormatETHAddress(addr string) (ETHADDRESS, bool) {
	var address ETHADDRESS
	address.UnmarshalJSON([]byte(addr))
	if address == "" {
		return "", false
	}
	return address, true
}

type RPCHeader struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Transaction struct {
	AccountId        int64           `json:"accountId,omitempty"`
	BlockHash        string          `json:"blockHash,omitempty"`
	BlockNumber      string          `json:"blockNumber,omitempty"`
	From             ETHADDRESS      `json:"from,omitempty"`
	Gas              string          `json:"gas,omitempty"`
	GasPrice         string          `json:"gasPrice,omitempty"`
	Hash             string          `json:"hash,omitempty"`
	Input            string          `json:"input,omitempty"`
	Nonce            string          `json:"nonce,omitempty"`
	To               ETHADDRESS      `json:"to,omitempty"`
	CoinCode         string          `json:"coin_code,omitempty"`
	RealTo           ETHADDRESS      `json:"real_to,omitempty"`
	Vol              decimal.Decimal `json:"vol,omitempty"`
	SettleType       SETTLE_TYPE     `json:"settle_type,omitempty"`
	TransactionIndex string          `json:"transactionIndex,omitempty"`
	Value            string          `json:"value,omitempty"`
	V                string          `json:"v,omitempty,omitempty"`
	R                string          `json:"r,omitempty,omitempty"`
	S                string          `json:"s,omitempty,omitempty"`
}

type TransactionReceipt struct {
	BlockHash         string     `json:"blockHash,omitempty"`
	BlockNumber       string     `json:"blockNumber,omitempty"`
	ContractAddress   ETHADDRESS `json:"contractAddress,omitempty"`
	CumulativeGasUsed string     `json:"cumulativeGasUsed,omitempty"`
	From              ETHADDRESS `json:"from,omitempty"`
	GasUsed           string     `json:"gasUsed,omitempty"`
	Logs              []struct {
		ETHADDRESS       string   `json:"address,omitempty"`
		Topics           []string `json:"topics,omitempty"`
		Data             string   `json:"data,omitempty"`
		BlockNumber      string   `json:"blockNumber,omitempty"`
		TransactionHash  string   `json:"transactionHash,omitempty"`
		TransactionIndex string   `json:"transactionIndex,omitempty"`
		BlockHash        string   `json:"blockHash,omitempty"`
		LogIndex         string   `json:"logIndex,omitempty"`
		Removed          bool     `json:"removed,omitempty"`
	} `json:"logs,omitempty"`
	LogsBloom        string                   `json:"logsBloom,omitempty"`
	Status           TransactionReceiptStatus `json:"status,omitempty"`
	To               ETHADDRESS               `json:"to,omitempty"`
	TransactionHash  string                   `json:"transactionHash,omitempty"`
	TransactionIndex string                   `json:"transactionIndex,omitempty"`
	ConfirmPlatform  string                   `json:"confirm_platform,omitempty"`
}

type TransactionReceiptStatus int

const (
	TransactionReceiptStatusFail    = 0
	TransactionReceiptStatusSuccess = 1
)

func (s *TransactionReceiptStatus) UnmarshalJSON(b []byte) error {
	if string(b) == `"0x01"` || string(b) == `"0x1"` || string(b) == `"1"` || string(b) == "1" {
		*s = TransactionReceiptStatusSuccess
	} else {
		*s = TransactionReceiptStatusFail
	}
	return nil
}

type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions,omitempty"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []string      `json:"uncles"`
}

type BlockResponse struct {
	RPCHeader
	Result *Block `json:"result,omitempty"`
}

type TransactionResponse struct {
	RPCHeader
	Result *Transaction `json:"result,omitempty"`
}

type BlockNumerResponse struct {
	RPCHeader
	Result string `json:"result,omitempty"`
}

//sgj 1128adding
type TransCountResponse struct {
	RPCHeader
	Result string `json:"result,omitempty"`
}

//add
type TransSendTransResponse struct {
	RPCHeader
	Result string `json:"result,omitempty"`
}

type TransactionReceiptResponse struct {
	RPCHeader
	Result *TransactionReceipt `json:"result,omitempty"`
}

type SendTransactionResponse struct {
	RPCHeader
	Result string `json:"result,omitempty"`
}

type ContractResponse struct {
	RPCHeader
	Result string `json:"result,omitempty"`
}


//sgj 1127add:

type SignAttachEthreum struct    {
	GasPrice int64
	GasLimit uint64
	Nonce uint64
}

type EthTranState struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Orderid    int64     `xorm:"unique BIGINT(20)"`
	CoinType   string    `xorm:"CHAR(16)"`
	SubType    string    `xorm:"CHAR(16)"`
	Txhash     string    `xorm:"VARCHAR(255)"`
	From       string    `xorm:"VARCHAR(255)"`
	To         string    `xorm:"VARCHAR(255)"`
	Amount     float64   `xorm:"DOUBLE"`
	Gasfee     float64   `xorm:"comment('通过java 请求调用传递过来的参数') DOUBLE"`
	Gasfeeused float64   `xorm:"comment('已经使用的交易费') DOUBLE"`
	Gasprice   int64     `xorm:"comment('交易价格') BIGINT(20)"`
	Gaslimit   int64     `xorm:"comment('最大交易步数') BIGINT(255)"`
	Gasused    int64     `xorm:"BIGINT(255)"`
	Nonce      uint64     `xorm:"BIGINT(20)"`
	Status     string    `xorm:"VARCHAR(255)"`
	Actblock   int       `xorm:"default 0 INT(255)"`
	TimeCreate time.Time `xorm:"comment('创建时间') TIMESTAMP"`
	TimeUpdate time.Time `xorm:"comment('更新时间') TIMESTAMP"`
	Desc       string    `xorm:"TEXT"`
	Raw        string    `xorm:"TEXT"`
}