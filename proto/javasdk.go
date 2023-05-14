package proto

/*
TxUtility. ClientToTransferProve()
*/
//public static JSONObject ClientToTransferProve(String fromPubkeyStr, Long nonce,byte[] payload,String prikeyStr){

type ClientToTransferProveParams struct {
	FromPubkeyStr string `json:"fromPubkeyStr"`
	Payload       string `json:"payload"`
	Nonce         int64  `json:"nonce"`
	PrikeyStr     string `json:"prikeyStr"`
}

//构造签名的交易事务
//public static JSONObject ClientToTransferAccount(String fromPubkeyStr, String toPubkeyHashStr, BigDecimal amount, String prikeyStr,Long nonce){
type ClientToTransferAccountParams struct {
	FromPubkeyStr   string `json:"fromPubkeyStr"`
	ToPubkeyHashStr string `json:"toPubkeyHashStr"`
	//Amount	int64			 `json:"amount"`
	Amount    float64 `json:"amount"`
	PrikeyStr string  `json:"prikeyStr"`
	Nonce     int64   `json:"nonce"`
}

//sgj 20200604 add,,for WGC Token trans
//1.31构造签名的资产定义的转账的规则调用事务

type CreateSignToDeployforRuleTransferParams struct {
	FromPubkeyStr string  `json:"fromPubkeyStr"`
	TxHash1       string  `json:"txHash1"`
	PrikeyStr     string  `json:"prikeyStr"`
	Nonce         int64   `json:"nonce"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Value         float64 `json:"value"`
}

type ClientToIncubateProfitParams struct {
	FromPubkeyStr   string
	ToPubkeyHashStr string
	Amount          int64
	PrikeyStr       string
	Txid            string
	Nonce           int64
}

//访问参数，返回协议
type JavaSDKResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode string      `json:"statusCode,omitempty"`
}

type NodeRPCResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	//2000,5000,6000,值类型
	StatusCode int `json:"statusCode,omitempty"`
}

//sgj 20200607add
type BlockPayloadResponse struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int64  `json:"value"`
}

//
//Node的RPC的，返回协议
type NodeResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"code,omitempty"`
}

//sgj 20200604add
type BlockHeadWDCResponse struct {
	BlockSize          int         `json:"blockSize"`
	BlockHash          string      `json:"blockHash"`
	NVersion           int         `json:"nVersion"`
	HashPrevBlock      string      `json:"hashPrevBlock"`
	HashMerkleRoot     string      `json:"hashMerkleRoot"`
	HashMerkleState    string      `json:"hashMerkleState"`
	HashMerkleIncubate string      `json:"hashMerkleIncubate"`
	NHeight            int         `json:"nHeight"`
	NTime              int64       `json:"nTime"`
	NBits              string      `json:"nBits"`
	NNonce             string      `json:"nNonce"`
	BlockNotice        interface{} `json:"blockNotice,omitempty"`
	BlockBodyData      interface{} `json:"body,omitempty"`
}

//通过keystore获得地址
type KeystoreToAddress struct {
	KsJson   string `json:"ksJson"`
	PassWord string `json:"password"`
}

//生成keystore文件
type AccountPassword struct {
	Password string `json:"password"`
}

//verifyAddress
type VerifyAddressReq struct {
	Address string `json:"address"`
}

//通过地址获得公钥哈希
type AddressToPubkeyHash struct {
	Address string `json:"address"`
}

//1104add,通过地址获得公钥哈希,,RPC解析用
type PubkeyHashToAddress struct {
	PubkeyHashStr string `json:"r1Str"`
}

type WdcTxBlock struct {
	BlockHash string `json:"block_hash"`
	Height    int64  `json:"height"`

	Version    int    `json:"version"`
	TxHash     string `json:"tx_hash"`
	Type       int64  `json:"type"`
	Nonce      int64  `json:"nonce"`
	FormPubKey string `json:"form"`

	GasPrice int64 `json:"gas_price"`
	Amount   int64 `json:"amount"`

	Payload      string `json:"payload"`
	Signature    string `json:"signature"`
	ToPubKeyHash string `json:"to"`
}

//sgj 0220 add:
type WdcTxBlockNew struct {
	BlockHash   string `json:"blockHash"`
	BlockHeight int64  `json:"blockHeight"`

	//sgj update name
	TransactionHash string `json:"transactionHash"`
	Version         int    `json:"version"`
	Type            int64  `json:"type"`
	Nonce           int64  `json:"nonce"`
	From            string `json:"from"`

	GasPrice int64 `json:"gasPrice"`
	Amount   int64 `json:"amount"`

	Payload      string `json:"payload"`
	ToPubKeyHash string `json:"to"`
	Signature    string `json:"signature"`
}

//sgj 20200605 add
type PayloadStrReq struct {
	Payload string `json:"payload"`
}

//sgj 20200605add
//解析的block中的WGC的数据
type AssetPayLoadTransaction struct {
	StatusCode      string `json:"statusCode"`
	PubkeyFirstSign int64  `json:"pubkeyFirstSign"`
	PubkeyFirst     string `json:"pubkeyFirst"`
	SignFirst       int64  `json:"signFirst"`
	Data            string `json:"data"`
	Message         int64  `json:"message"`
}

//1030

//生成keystore文件
type SendBalanceReq struct {
	PubkeyHash string `json:"pubkeyhash"`
}

//生成keystore文件
type SendTransactionReq struct {
	TransactionStr string
}
