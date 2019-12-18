package proto


//12.09add:
/*
参数：
https://mistydew.github.io/blog/2018/09/bitcoin-rpc-command-listunspent.html
1.minconf（数字，可选，默认为 1）要过滤的最小确认数。
2.maxconf（数字，可选，默认为 9999999）要过滤的最大确认数。
3.addresses（字符串）要过滤的比特币地址的 json 数组。
*/

type KtcUtxoReq struct {
	//Minconf    int64   `json:"minconf"`     //
	//Maxconf    int64  `json:"maxconf"`
	QueryAddrList          []string
	//`json:"addresses"`            //

}

type KtcUtxoAddrReq struct {
	//Minconf    int64   `json:"minconf"`     //
	//Maxconf    int64  `json:"maxconf"`
	//Address         string  `json:"address"`            //
	Address         string            //

}
//sgj 1209add :listunspent,更新于：
//http://cw.hubwiz.com/card/c/bitcoin-json-rpc-api/1/7/33/

//https://blockchain.info/unspent?active=
type CurKtcUtxoInfo struct {
	Txid string `json:"txid"`
	//TxOutputN int64  `json:"tx_output_n"`
	Vout int64  `json:"vout"`
	//TxHashBigEndian string `json:"tx_hash_big_endian"`
	Address string `json:"address"`
	Label string `json:"label"`
	//Script string `json:"script"`
	ScriptPubKey string `json:"scriptPubKey"`

	Amount float64  `json:"amount"`
	//ValueHex string `json:"value_hex"`

	RedeemScript string `json:"redeemScript"`
	Confirmations int64  `json:"confirmations"`
	Spendable bool  `json:"spendable"`
	Solvable bool  `json:"solvable"`
	Safe bool  `json:"safe"`


}

//
type KTCUnspentOutputs struct {
	CurKtcUtxoInfo []CurKtcUtxoInfo `json:"unspent_outputs"`

}

//1217 add
//11/19 add:
//0331--utxo ,未消费记录：
type AddrBalanceUnspent struct {
	TxidHex        string	`json:"txid" `   //交易id
	Vout         int		`json:"vout" `   //交易中序号
	AddrDisp        string `json:"address" `   //数字支付地址
	Amount         float64		`json:"amount" `   //交易金额
}

type TransactionAddressInfo struct {
	Address    string  `json:"address"`
	Amount    float64  `json:"amount"`
}

type SignTransactionReq struct {
	Froms        []TransactionAddressInfo   `json:"fromAddr"`            // 付款人地址
	Tos          []TransactionAddressInfo  `json:"toAddr"`            // 收款人地址
	RemainAddr    string  `json:"remainAddr"`            // 找零地址
	//不再用 Amount 		float64   `json:"amount"`  // 转账金额(比特专用)(以太坊不用)
	CoinType    string   `json:"coinType"`     // 数字币类型
	SubType    string   `json:"subType"`     // 数字币子类型
	OrderId    int64   `json:"orderId"`     // 订单号
	FeeAmount  float64  `json:"feeAmount"`
	EthLimit    int64 	`json:"ethLimit"`
	EthPrice 	int64  `json:"ethPrice"`
}