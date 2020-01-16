//package omnitranssign
package omnirpc
 import (
     "fmt"
  "github.com/btcsuite/btcutil"
   "encoding/json"


 )

func Getnewaddress() string {
        //4API
        //0611 add sgj
       response, err := OmniRPCClient.RpcClient.Call("getnewaddress") 
       fmt.Printf("--005---getnewaddress: ------005--get response info is: %v;;err is :%v\n", response,err)
       rawGetHex := fmt.Sprintf("%v",response)
       fmt.Printf("--0k1---getnewaddress: --get rawGetHex info is: %s;\n", rawGetHex)
       return rawGetHex
}
func Getrawtransaction(txid string){
    response, err := OmniRPCClient.RpcClient.Call("getrawtransaction", txid , 1)
   fmt.Printf("--003---Getrawtransaction: ------003--get response info is--008: %v;;err is :%v\n", response,err)
 
} 
func Getblockhash(block int){
    response, err := OmniRPCClient.RpcClient.Call("getblockhash", block)
   fmt.Printf("--003---Getblockhash: ------003--get response info is--008: %v;;err is :%v\n", response,err)

} 

func Getblock(hash string){
    response, err := OmniRPCClient.RpcClient.Call("getblock", hash)
    fmt.Printf("--009D---getblock: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)
 
 } 
func Sendrawtransaction(tx string){
      response, err := OmniRPCClient.RpcClient.Call("sendrawtransaction", tx)
    fmt.Printf("--009D---sendrawtransaction: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)
 
}
func validateaddress(addr string){
    response, err := OmniRPCClient.RpcClient.Call("validateaddress", addr)
    fmt.Printf("--009D---validateaddress: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)
 
} 
func Createrawtransaction(ins string,outs string){

    response, err := OmniRPCClient.RpcClient.Call("createrawtransaction",ins,outs)
    fmt.Printf("--009D---createrawtransaction: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)

    } 
func Decoderawtransaction(rawtx string){
    response, err := OmniRPCClient.RpcClient.Call("decoderawtransaction", rawtx)
    fmt.Printf("--009D---decoderawtransaction: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)

    } 
func Omni_decodetransaction(rawtx string){
    response, err := OmniRPCClient.RpcClient.Call("omni_decodetransaction", rawtx)
    fmt.Printf("--009D---omni_decodetransaction: ------009D--get response info is--009B: %v;;err is :%v\n", response,err)

}
//blocks=4
func estimateFee(blocks int ) int64 {
    response, err := OmniRPCClient.RpcClient.Call("estimatefee", blocks)
    fmt.Printf("--0011B---estimatefee: ------00911--get response info is: %v;;err is :%v\n", response,err)
    //sgj tmp
    return 0
}
//func Gettxout(txid string,vout int ,unconfirmed bool= true) string {

func Gettxout(txid string,vout int ,unconfirmed bool) string {
    response, err := OmniRPCClient.RpcClient.Call("gettxout",txid,vout,unconfirmed)

    fmt.Printf("--0011B---unconfirmed: ------00911--get response info is: %v;;err is :%v\n", response,err)
     //sgj tmp
     return ""
   
} 
// OmniSendCmd defines the omni_send JSON-RPC command.

type OmniSendCmd struct {
	FromAddress string
	ToAddress   string
	PropertyId  int
	Amount      btcutil.Amount
}
type GetBalanceResult struct {
    Balance  string
	Reserved string
}
/*
*/
type GetAllBalanceResult struct {
	propertyid int
	Balance  string
	Reserved string
}

//omnijson.OmniGetBalanceCmd --return old val
func Getbalance_MP(addr string, propertyid int) (getbalanceinfo *GetBalanceResult,err error){
    response, err := OmniRPCClient.RpcClient.Call("omni_getbalance", addr, propertyid)
  
    //jsonReq := omnijson.OmniGetBalanceCmd{}
    jsonReq := GetBalanceResult{}
    get_response, err :=json.Marshal(response.Result)
    if err != nil{
        fmt.Printf("response.Result,,err is====> :%v",err)
    }else{
        //fmt.Printf("get_response succ,info is====> :%v\n",string(get_response))
    }
    err = json.Unmarshal(get_response, &jsonReq)
    if nil != err {
        fmt.Printf("--009B---Unmarshal:-- jsonReq is :%v;;err is :%v\n", jsonReq,err)
        return nil,err
    }
    //fmt.Printf("get last jsonReq info is====> :%v",jsonReq)
    return &jsonReq,nil
} 

//rawtx = createrawtx_reference(self.rawdata['transaction_to'], rawtx)['result']
func Getallbalancesforaddress_MP(addr string)  (getbalanceinfo *GetAllBalanceResult){
    response, err := OmniRPCClient.RpcClient.Call("getallbalancesforaddress_MP", addr)
    fmt.Printf("--009B---getallbalancesforaddress_MP: ------009B--get response info is--009B: %v;;err is :%v\n", response,err)
    return nil

} 

func Gettransaction_MP(tx string) string {
    response, err := OmniRPCClient.RpcClient.Call("gettransaction_MP", tx)
    fmt.Printf("--009B---gettransaction_MP: ------009B--get response info is--009B: %v;;err is :%v\n", response,err)
    return ""
} 
func listblocktransactions_MP(height int ){
    response, err := OmniRPCClient.RpcClient.Call("listblocktransactions_MP", height)
    fmt.Printf("--003---listblocktransactions_MP: ------008--get response info is--008: %v;;err is :%v\n", response,err)
 
} 
func Getproperty_MP(propertyid int){
    response, err := OmniRPCClient.RpcClient.Call("getproperty_MP", propertyid)
    fmt.Printf("--003---getproperty_MP: ------008--get response info is--008: %v;;err is :%v\n", response,err)
 
} 
func Listproperties_MP(){
    response, err := OmniRPCClient.RpcClient.Call("omni_listproperties")
    fmt.Printf("--003---omni_listproperties: ------008--get response info is--008: %v;;err is :%v\n", response,err)
    
}

func GetSimplesendPayload(propertyid int, amount string) (txHex string){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_simplesend", int(propertyid), amount)
    //String rawTxHex = String.format("00000000%08x%016x", currencyId.getValue(), amount.getWillets());
    //rawGetHex := fmt.Sprintf("%v",response)
    var cursimplePayloadResult string
    if response.Result != nil {
        cursimplePayloadResult =response.Result.(string)
    }else{
        cursimplePayloadResult = ""
    }
    fmt.Printf("--001API :GetSimplesendPayload: --get cursimplePayloadResult info is: %v;;err is :%v\n", cursimplePayloadResult,err)

    return cursimplePayloadResult
} 
//omni_createpayload_grant propertyid "amount" ( "memo" )
func GetgrantPayload(propertyid int,amount string){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_grant", int64(propertyid), amount)
    fmt.Printf("--003---GetgrantPayload: ------008--get response info is--008: %v;;err is :%v\n", response,err)

} 
func GetrevokePayload(propertyid int , amount, memo string){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_revoke", int(propertyid), amount, memo)
    fmt.Printf("--003---omni_createpayload_revoke: ------008--get response info is--008: %v;;err is :%v\n", response,err)
} 

func GettradePayload(propertyidforsale int, amountforsale string, propertiddesired int, amountdesired string){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_trade", int(propertyidforsale), amountforsale, int(propertiddesired), amountdesired)
    fmt.Printf("--007---omni_createpayload_trade: ------007--get response info is--006: %v;;err is :%v\n", response,err)

} 

//sgj 0801 update
func Createrawtx_opreturn(rawtx,payload string) string{
//func Createrawtx_opresponse, err :=(payload, rawtx=None){
    response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_opreturn", rawtx, payload)
    //getRawTx := fmt.Sprintf("%v",response)
    var currawtx_opreturn string
    if response.Result != nil {
        currawtx_opreturn =response.Result.(string)
    }else{
        currawtx_opreturn = ""
    }
    fmt.Printf("--0k4new---omni_createrawtx_opreturn: ---getRawTx info is: %v;;err is :%v\n", currawtx_opreturn,err)
    return currawtx_opreturn
} 
 
//step 5)
//sgj 0801 update
func Createrawtx_Reference(rawtx,destination string,amount float64) string{
    //the optional reference amount (minimal by default)
    // response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_reference", rawtx, destination,amount)
 
    //sgj 0806 update,remove amount
    response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_reference", rawtx, destination)
    //getRawTx := fmt.Sprintf("%v",response)

    var currawtx_reference string
    if response.Result != nil {
        currawtx_reference =response.Result.(string)
    }else{
        currawtx_reference = ""
    }
    fmt.Printf("--0k5B---omni_createrawtx_reference: --get currawtx_reference info is: %v;;err is :%v\n", currawtx_reference,err)
    return currawtx_reference

} 

func Createrawtx_change(rawtx, previnputs, destination string, fee float64) string{
  fmt.Printf("to omni_createrawtx_change()'s params info,rawtx is: %v; destination is :%s;fee is :%d;position is :%d\n", rawtx,destination,fee,1)
  //0609 add params: position 1:
  fmt.Printf("to watching==> 2's previnputs params info is :%v\n",previnputs)
  //previnputs = "'[{"txid":"48ae09d5441a88c7aa9635b9d82e7556cb75371dc8c964b828889918874c5d8b","vout":0,"scriptPubKey":"76a914c4d4f94d368e245d9d7369d4515c18b3563e729b88ac","value":0.07},{"txid":"48ae09d5441a88c7aa9635b9d82e7556cb75371dc8c964b828889918874c5d8b","vout":0,"scriptPubKey":"76a914c4d4f94d368e245d9d7369d4515c18b3563e729b88ac","value":0.07}]'"
    response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_change", rawtx, previnputs, destination, fee,1)
    var getRawTx string
    if response.Result != nil {
        getRawTx =response.Result.(string)
    }else{
        getRawTx = ""
    }
    fmt.Printf("--0k6B---omni_createrawtx_reference: --get getRawTx info is: %v;;err is :%v\n", getRawTx,err)
    return getRawTx
}

/*
func GetchangeissuerPayload(propertyid int){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_changeissuer", int(propertyid))
} 
//define omni_getproperty struct for response, err :=
func Getcrowdsale_MP(propertyid int){
    response, err := OmniRPCClient.RpcClient.Call("getcrowdsale_MP", propertyid)

} 
func GetissuancefixedPayload(ecosystem, divisible, previousid, category,subcategory, name, url, data, amount){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_issuancefixed", int(ecosystem), int(divisible), int(previousid), category,subcategory, name, url, data, amount)

}
func GetissuancecrowdsalePayload(ecosystem, divisible, previousid, category,subcategory, name, url, data, propertyiddesired, tokensperunit, deadline, earlybonus, issuerpercentage){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_issuancecrowdsale", int(ecosystem), int(divisible), int(previousid), category,subcategory, name, url, data, int(propertyiddesired), tokensperunit, int(deadline), int(earlybonus), int(issuerpercentage))

} 
func GetissuancemanagedPayload(ecosystem, divisible, previousid, category,subcategory, name, url, data){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_issuancemanaged", int(ecosystem), int(divisible), int(previousid), category,subcategory, name, url, data)
} 

func Getallbalancesforid_MP(propertyid int){
    response, err := OmniRPCClient.RpcClient.Call("getallbalancesforid_MP", propertyid)
} 
func Getactivecrowdsales_MP(){
    response, err := OmniRPCClient.RpcClient.Call("getactivecrowdsales_MP")
} 
func Getactivedexsells_MP(){
    response, err := OmniRPCClient.RpcClient.Call("getactivedexsells_MP")
} 
func Getdivisible_MP(propertyid){
    response, err := getproperty_MP(propertyid)['result']['divisible']

func Getgrants_MP(propertyid){
    response, err := OmniRPCClient.RpcClient.Call("getgrants_MP", propertyid)
} 
func Gettradessince_MP(){
    response, err := OmniRPCClient.RpcClient.Call("gettradessince_MP")
} 
func Gettrade(txhash){
    response, err := OmniRPCClient.RpcClient.Call("omni_gettrade", txhash)
} 
func Getsto_MP(txid){
    response, err := OmniRPCClient.RpcClient.Call("getsto_MP", txid , "*")
} 

func GetsendallPayload(ecosystem){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_sendall", int(ecosystem))
}    
func GetdexsellPayload(propertyidforsale, amountforsale, amountdesired, paymentwindow, minacceptfee, action){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_dexsell", int(propertyidforsale), amountforsale, amountdesired, int(paymentwindow), minacceptfee, int(action))
} 
func GetdexacceptPayload(propertyid, amount){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_dexaccept", int(propertyid), amount)
} 
func GetstoPayload(propertyid, amount){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_sto", int(propertyid), amount)
}

func GetclosecrowdsalePayload(propertyid){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_closecrowdsale", int(propertyid))
} 
func GetcanceltradesbypricePayload(propertyidforsale, amountforsale, propertiddesired, amountdesired){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_canceltradesbyprice", int(propertyidforsale), amountforsale, int(propertiddesired), amountdesired)
} 
func GetcanceltradesbypairPayload(propertyidforsale, propertiddesired){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_canceltradesbypair", int(propertyidforsale), int(propertiddesired))
} 
func GetcancelalltradesPayload(ecosystem){
    response, err := OmniRPCClient.RpcClient.Call("omni_createpayload_cancelalltrades", int(ecosystem))
}

func Createrawtx_multisig(payload, seed, pubkey, rawtx=None){
    response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_multisig", rawtx, payload, seed, pubkey)
}

*/
//omnicore-cli "omni_createrawtx_input" \
 //   "01000000000000000000" "b006729017df05eda586df9ad3f8ccfee5be340aadf88155b784d1fc0e8342ee" 0
 //func Createrawtx_input(txhash, index, rawtx=None) (rawtx string){
 func Createrawtx_input(txhash string, index int) (rawtx string){
    response, err := OmniRPCClient.RpcClient.Call("omni_createrawtx_input", rawtx, txhash, index)
    fmt.Printf("--003---omni_createrawtx_input: ------004--get response info is--004: %v;;err is :%v\n", response,err)
    return "omni_createrawtx_input======infoval"
}

