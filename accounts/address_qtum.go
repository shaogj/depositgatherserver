package accounts

import (
	"fmt"
	//"strings"
	"2019NNZXProj10/abitserverDepositeGather/qtumrpc"
	//"log"
   "github.com/mkideal/log"
)

func AddressGenerateDSCaa() (getprikey string, getaddrpubkey string, getaddress string, err error) {
	
	tonewaddressqtum := qtumrpc.Getnewaddress()
	fmt.Printf("qtumrpc Getnewaddress() exec succ!,get tonewaddress info is :%v\n",tonewaddressqtum)

	//sgj 0817 add:
	curaddrPrikey := qtumrpc.GetAddprivkey(tonewaddressqtum)
	log.Info("===GGEX--GetAddprivkey() exec succ!,get curaddrPrikey info is :%v\n",curaddrPrikey)
	curPubkey := curaddrPrikey + "tmppubkey"
	return curaddrPrikey,curPubkey,tonewaddressqtum,nil
}
//getprivkey, getpubKey, pubKeyAddr, err := accounts.AddressGenerateLTC()
