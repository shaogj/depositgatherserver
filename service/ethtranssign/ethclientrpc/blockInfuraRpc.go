//package ethtranssign
package ethclientrpc

import (
	//"backend/support/config"
	//"backend/support/domain"
	"2019NNZXProj10/depositgatherserver/proto"

	mylog "github.com/mkideal/log"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var (
	infuraPostMethos map[string]bool
	//serverConf       *domain.ServerConf
)
type EtherServerConf struct {
	EtherScanApiKey    string
	EtherumNetwork    string
}

var serverConf EtherServerConf

//sgj 1127update
func InitNetwork(etherumNetwork string,etherScanApiKey string) {
	//conf := config.GetServerConf("SettleCenter")
	if etherumNetwork == "" || etherScanApiKey == "" {
		panic(errors.New("No etherumNetwork confing"))
	}
	//serverConf = conf
	serverConf.EtherScanApiKey = etherScanApiKey
	serverConf.EtherumNetwork = etherumNetwork

	infuraPostMethos = map[string]bool{
		"eth_sendRawTransaction": true,
		"eth_estimateGas":        true,
		"eth_submitWork":         true,
		"eth_submitHashrate":     true,
		//sgj 1127 add:
		"eth_getBalance":         true,
		"eth_getTransactionCount":     true,
		//1202 add:
		"eth_getTransactionReceipt" :true,

	}
}

func InfuraJsonRPC(network, method string, params interface{}, response interface{}) error {
	if _, ok := infuraPostMethos[method]; ok {
		return infuraPost(network, method, params, response)
	}
	return infuraGet(network, method, params, response)
}

func infuraGet(network, method string, params interface{}, response interface{}) error {
	client := http.Client{}
	paramsJosn, err := json.Marshal(params)
	if err != nil {
		return err
	}
	URL := fmt.Sprintf("https://api.infura.io/v1/jsonrpc/%s/%s?params=%s", network, method, url.QueryEscape(string(paramsJosn)))
	req, err := http.NewRequest("GET", URL, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		mylog.Debug("infuraGet parse error:%s", err.Error())
		return err
	}

	return nil
}

func infuraPost(network, method string, params interface{}, response interface{}) error {
	client := http.Client{}
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      57386342,
		"method":  method,
		"params":  params,
	}
	//URL := fmt.Sprintf("https://api.infura.io/v1/jsonrpc/%s", network)
	//sgj 1127,升级调用URL
	//sgj 0107 mark for ETH access api token:3aee9923486b4a49a8ec44db929e823c

	URL := "https://mainnet.infura.io/v3/3aee9923486b4a49a8ec44db929e823c"
	reqBody, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		mylog.Debug("infuraPost parse error:%s", err.Error())
		return err
	}
	return nil
}

//sgj 1130 fock from settlecenter module:
//method 2 form EtherscanJsonRPC:

func EtherscanJsonRPC(network string, params *map[string]string, response interface{}) error {
	var host string
	if network != "mainnet" {
		host = fmt.Sprintf("%s.etherscan.io", network)
	} else {
		host = "api.etherscan.io"
	}
	client := http.Client{}
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "api",
	}
	query := u.Query()
	for k, v := range *params {
		query.Set(k, v)
	}
	query.Set("ApiKeyToken", serverConf.EtherScanApiKey)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		mylog.Debug("EtherscanJsonRPC parse error:%s", err.Error())
		return err
	}
	return nil
}

func BlockNumer() (int64, error) {
	method := "eth_blockNumber"
	network := serverConf.EtherumNetwork
	params := []interface{}{}
	var response proto.BlockNumerResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}
	return strconv.ParseInt(response.Result[2:], 16, 64)
}

func Block(num int64) (*proto.Block, error) {
	method := "eth_getBlockByNumber"
	network := serverConf.EtherumNetwork
	params := []interface{}{
		fmt.Sprintf("0x%x", num),
		true,
	}
	var response proto.BlockResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}
	if response.Result == nil {
		return nil, errors.New("NOT FOUND")
	}
	return response.Result, nil
}

func Transaction(txHash string) (*proto.Transaction, error) {
	method := "eth_getTransactionByHash"
	network := serverConf.EtherumNetwork
	params := []interface{}{
		txHash,
	}
	var response proto.TransactionResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}
	return response.Result, nil
}
//sgj 11.28adding,for nonce;

func GetTransactionAccount(address string) (int64, error) {
	method := "eth_getTransactionCount"
	network := serverConf.EtherumNetwork
	//sgj add
	statusLabel := "latest"
	params := []interface{}{
		address,
		statusLabel,
	}
	var response proto.TransCountResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}
	//return response.Result, nil
	return strconv.ParseInt(response.Result[2:], 16, 64)

}

func GetTransactionPendingNonce(address string) (int64, error) {
	method := "eth_getTransactionCount"
	network := serverConf.EtherumNetwork
	//sgj add
	statusLabel := "pending"
	params := []interface{}{
		address,
		statusLabel,
	}
	var response proto.TransCountResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}
	//return response.Result, nil
	return strconv.ParseInt(response.Result[2:], 16, 64)

}

//获取余额
//add 2)POST https://<network>.infura.io/v3/YOUR-PROJECT-ID
func GetBalance(address string,statusLabel string) (int64, error) {
	method := "eth_getBalance"
	network := serverConf.EtherumNetwork
	params := []interface{}{
		address,
		statusLabel,
	}
	var response proto.TransCountResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}
	return strconv.ParseInt(response.Result[2:], 16, 64)
}

//add 3)eth_sendRawTransaction
//向节点提交一个已签名的交易以便广播到以太坊网络中
//add 2)POST https://<network>.infura.io/v3/YOUR-PROJECT-ID
func SendRawTransaction(txHexstr string) (string, error) {
	method := "eth_sendRawTransaction"
	network := serverConf.EtherumNetwork
	params := []interface{}{
		txHexstr,
	}
	var response proto.TransSendTransResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", errors.New(response.Error.Message)
	}
	return response.Result, nil

}


//end 11.28

func EtherscanTranscation(txHash string) (*proto.Transaction, error) {
	//module=proxy&action=eth_getTransactionByHash&txhash=0x1e2910a262b1008d0616a0beb24c1a491d78771baa54a33e66065e03b1f46bc1&apikey=YourApiKeyToken
	params := map[string]string{
		"module": "proxy",
		"action": "eth_getTransactionByHash",
		"txhash": txHash,
	}
	var response proto.TransactionResponse
	if err := EtherscanJsonRPC(serverConf.EtherumNetwork, &params, &response); err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}
	return response.Result, nil
}

func TransactionReceiptEthersan(txHash string) (*proto.TransactionReceipt, error) {
	params := map[string]string{
		"module": "proxy",
		"action": "eth_getTransactionReceipt",
		"txhash": txHash,
	}
	var response proto.TransactionReceiptResponse
	if err := EtherscanJsonRPC(serverConf.EtherumNetwork, &params, &response); err != nil {
		return nil, err
	}
	if response.Error != nil {
		mylog.Warn("TransactionReceiptEthersan error:%s", response.Error.Message)
		return nil, errors.New(response.Error.Message)
	}
	if response.Result == nil {
		mylog.Warn("TransactionReceiptEthersan %s TransactionReceipt is nil", txHash)
		return nil, nil
	}
	return response.Result, nil
}

func TransactionReceiptInfura(txHash string) (*proto.TransactionReceipt, error) {
	network := "mainnet"
	method := "eth_getTransactionReceipt"
	params := []interface{}{
		txHash,
	}

	var response proto.TransactionReceiptResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		mylog.Warn("TransactionReceiptInfura error:%s", err.Error())
		return nil, err
	}
	if response.Result == nil {
		mylog.Warn("TransactionReceiptInfura %s TransactionReceipt is nil", txHash)
		return nil, nil
	}

	return response.Result, nil
}

func TransactionReceipt(txhash string) *proto.TransactionReceipt {
	//if trans, _ := TransactionReceiptEthersan(txhash); trans != nil {
	//sgj1202,1224 add:
	if trans, err  := TransactionReceiptEthersan(txhash); trans != nil {
		mylog.Info("TransactionReceiptEthersan hash is: %s ,TransactionReceipt no nil,err is :%v", txhash,err)
		trans.ConfirmPlatform = "Etherscan"
		return trans
	}else{
		mylog.Error("TransactionReceiptEthersan hash is: %s ,TransactionReceipt is nil,err is :%v", txhash,err)
	}
	//if trans, _ := TransactionReceiptInfura(txhash); trans != nil {
	if trans, err2 := TransactionReceiptInfura(txhash); trans != nil {
		trans.ConfirmPlatform = "Infura"
		mylog.Info("TransactionReceiptInfura hash is: %s ,TransactionReceiptInfura is %v,err is :%v", txhash,trans,err2)
		return trans
	}else{
		mylog.Error("TransactionReceiptInfura hash is: %s ,TransactionReceiptInfura is nil,err is :%v", txhash,err2)

	}
	return nil
}

func SendTransactionRaw(signStr string) (string, error) {
	method := "eth_sendRawTransaction"
	network := serverConf.EtherumNetwork
	params := []interface{}{
		signStr,
	}
	var response proto.SendTransactionResponse
	if err := InfuraJsonRPC(network, method, params, &response); err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", errors.New(response.Error.Message)
	}
	return response.Result, nil
}
