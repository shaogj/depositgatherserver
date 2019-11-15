package proto

/*
TxUtility. ClientToTransferProve()
*/
//public static JSONObject ClientToTransferProve(String fromPubkeyStr, Long nonce,byte[] payload,String prikeyStr){

type ClientToTransferProveParams struct {
	FromPubkeyStr string	 `json:"fromPubkeyStr"`
	Payload	string 			`json:"payload"`
	Nonce	int64			 `json:"nonce"`
	PrikeyStr	string 		`json:"prikeyStr"`

}

//构造签名的交易事务
//public static JSONObject ClientToTransferAccount(String fromPubkeyStr, String toPubkeyHashStr, BigDecimal amount, String prikeyStr,Long nonce){
type ClientToTransferAccountParams struct {
	FromPubkeyStr string	 `json:"fromPubkeyStr"`
	ToPubkeyHashStr	string 			`json:"toPubkeyHashStr"`
	//Amount	int64			 `json:"amount"`
	Amount	float64			 `json:"amount"`
	PrikeyStr	string 		`json:"prikeyStr"`
	Nonce	int64			 `json:"nonce"`

}
type ClientToIncubateProfitParams struct {
	FromPubkeyStr string
	ToPubkeyHashStr string
	Amount	int64
	PrikeyStr string
	Txid string
	Nonce int64
}
//访问参数，返回协议
type JavaSDKResponse struct {
	Message string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
	StatusCode   string   `json:"statusCode,omitempty"`
}

type NodeRPCResponse struct {
	Message string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
	//2000,5000,6000,值类型
	StatusCode   int   `json:"statusCode,omitempty"`
}

//
//Node的RPC的，返回协议
type NodeResponse struct {
	Message string      `json:"message"`
	Data  interface{} `json:"data,omitempty"`
	StatusCode   int   `json:"code,omitempty"`
}

//通过keystore获得地址
type KeystoreToAddress struct {
	KsJson      string     `json:"ksJson"`
	PassWord    string     `json:"password"`
}

//生成keystore文件
type AccountPassword struct {
	Password      string     `json:"password"`

}


//verifyAddress
type VerifyAddressReq struct {
	Address      string     `json:"address"`

}

//通过地址获得公钥哈希
type AddressToPubkeyHash struct {
	Address	string `json:"address"`
}


//1104add,通过地址获得公钥哈希,,RPC解析用
type PubkeyHashToAddress struct {
	PubkeyHashStr	string `json:"r1Str"`
}

type WdcTxBlock struct {
	BlockHash string `json:"block_hash"`
	Height int64  `json:"height"`

	Version int  `json:"version"`
	TxHash string `json:"tx_hash"`
	Type int64  `json:"type"`
	Nonce int64  `json:"nonce"`
	FormPubKey string `json:"form"`


	GasPrice int64  `json:"gas_price"`
	Amount int64  `json:"amount"`

	Payload string `json:"payload"`
	Signature string `json:"signature"`
	ToPubKeyHash string `json:"to"`

}

//1030

//生成keystore文件
type SendBalanceReq struct {
	PubkeyHash      string     `json:"pubkeyhash"`

}

//生成keystore文件
type SendTransactionReq struct {
	TransactionStr      string

}