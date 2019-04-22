package lrpc

import (
	"lnet/iface"
	"proto/pb"
	"time"
)

type Rpc struct {
	transport      iface.ITransport
	msgHandle      iface.IMsgHandle
	asyncResultMgr *AsyncResultMgr
}

func (this *Rpc) Send(msg interface{}) error {
	msgPkg := this.msgHandle.CreateMessagePackage(msg)
	return this.transport.Send(msgPkg)
}

func (this *Rpc) SendWaitResult(msg interface{}) (*pb.RpcRspInfo, error) {
	msgRpc := msg.(pb.RpcReqInfo)
	asyncR := this.asyncResultMgr.Add()
	msgRpc.Uid = asyncR.GetUid()
	msgPkg := this.msgHandle.CreateMessagePackage(msgRpc)
	this.transport.Send(msgPkg)

	resp, err := asyncR.GetResult(2 * time.Second)
	if err == nil {
		return resp, nil
	} else {
		//超时了 或者其他原因结果没等到
		this.asyncResultMgr.Remove(asyncR.GetUid())
		return nil, err
	}
}
