package base

type EtherBaseInterf interface {
	//GoServiceLastestBlock()
	GoEstimateGas() //周期性和节点上同步费用信息
	NewClient(bWait bool) (EtherClientHandle,error)
	DoTransactionSuf(client EtherClientHandle,tran EtherTranBinary,bWait bool) ( int, error)
	//DoTransactionPre(tran EtherTranInfo) (EtherTranBinary,proto.ErrorInfo)
	//查询交易结果
	//DoCallTransactionCheck(client EtherClientHandle,tran *EtherIntf)
}
type EtherBase struct {
	Coin string  //币种类型
	StrRawUrl string
	ChTrans    chan *EtherIntf

	Intf EtherBaseInterf
}
