package qtumrpc

 import (
     "fmt"

    //"github.com/icloudland/btcdx/omnijson"
    //omnirpcclient "github.com/icloudland/btcdx/rpcclient"
 )

func Getnewaddress() string {
	//4API
   response, err := QtumRPCClient.RpcClient.Call("getnewaddress","") 
   //fmt.Printf("--005---getnewaddress()  QtumRPC error: get response info is: %v;;err is :%v\n", response,err)
   curnewadd :=response.Result.(string)
   //rawGetHex := fmt.Sprintf("%v",response)
   fmt.Printf("--0k1---getnewaddress: --getnewaddress get info is: %s,err is :%v\n", curnewadd,err)
   return curnewadd
}

//sgj 0807 add:
func GetAddprivkey(addr string) string {

    response, err := QtumRPCClient.RpcClient.Call("dumpprivkey",addr) 
    curaddprivkey :=response.Result.(string)
    fmt.Printf("--0k2---GetAddprivkey: curaddprivkey get info is: %v;;err is :%v\n", curaddprivkey,err)
    return curaddprivkey
}

func GetInfo() string {
	//4API
   response, err := QtumRPCClient.RpcClient.Call("getinfo","") 
   fmt.Printf("--004---getinfo()  QtumRPC error: get response info is: %v;;err is :%v\n", response,err)
   //curnewadd :=response.Result.(string)
   rawGetinfostr := fmt.Sprintf("%v",response)
   fmt.Printf("--04K---getinfo: --get curnewadd info is: %s;\n", rawGetinfostr)
   return rawGetinfostr
}

//sgj 0809 add new:
func GetAccountAddress(account string) (string, error) {
	//response, err := QtumRPCClient.RpcClient.Call("getaccountaddress", []interface{}{account})
    //sg 0814 add:
    response, err := QtumRPCClient.RpcClient.Call("getaccountaddress", []interface{}{"product.1"})
    
    //stding:   curnewadd :=response.Result.(string)
    rawGetHex := fmt.Sprintf("%v",response)
    fmt.Printf("--06API--- getaccountaddress: --get curnewadd info is: %s;\n", rawGetHex)
    if err != nil {
		return "", err
	}
	return rawGetHex, nil
	//return response.Result.(string), nil
}


func Getblockchaininfo()(interface{},error){
    response, err := QtumRPCClient.RpcClient.Call("getblockchaininfo")
    //fmt.Printf("--002API---Qtum RPC Invoke Getblockchaininfo(): get response info is--009B: %v;;err is :%v\n", response,err)
    return response,err
 } 


 //sgj 0824 add:
func CreateTransaction(createtxstr []string,addr string) (curinfo interface{},err error) {

    //addr:= "{\"MQ3Jx2krKJbDgAcxJrXvL73KXH4zVf8mWg\":11}");
    response, err := QtumRPCClient.RpcClient.Call("createrawtransaction",createtxstr,addr) 
    getresinfo :=response.Result
    //curaddprivkey :=response.Result.(string)
    fmt.Printf("--077---CreateTransaction: input params createtxstr is ;%s,CreateTransaction get response info  is: %v;;err is :%v\n", createtxstr,response,err)
    return getresinfo,err
}
