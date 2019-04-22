package lclient

import (
	"lnet"
	"lnet/iface"
	"lnet/lrpc"
	"lnet/ltransport"
	"proto/pb"
	"time"
)

type RpcClient struct {
	BaseClient
	asyncResultMgr *lrpc.AsyncResultMgr
}

func NewRpcClient(localAddr string, msgHandle iface.IMsgHandle) *RpcClient {
	rpcClient := &RpcClient{
		BaseClient: BaseClient{
			NetType:   lnet.TCP,
			LocalAddr: localAddr,
			transport: ltransport.NewTcpTransport(localAddr, lnet.DefMsgTimeout, msgHandle, nil, nil),
			msgHandle: msgHandle,
		},
		asyncResultMgr: lrpc.G_asyncResult,
	}

	return rpcClient
}

func (this *RpcClient) Send(msg interface{}) error {
	msgPkg := this.msgHandle.CreateMessagePackage(msg)
	return this.transport.Send(msgPkg)
}

func (this *RpcClient) SendWaitResult(msg interface{}) (*pb.RpcRspInfo, error) {
	msgRpc := msg.(*pb.RpcReqInfo)
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
