package lmsghandle

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"lnet/lrpc"
	"proto/pb"
	"time"
)

type RpcMsgHandle struct {
	BaseMsgHandle
}

func NewRpcMsgHandle(protocol iface.IProtocol) *RpcMsgHandle {
	msgHandle := &RpcMsgHandle{BaseMsgHandle: *NewBaseMsgHandle(protocol)}
	return msgHandle
}

func (this *RpcMsgHandle) Process(itransport iface.ITransport, msgPackage iface.IMessagePackage) {
	msg := this.CreateMessage(msgPackage)
	if msgPackage.GetTag() == 13 {
		rpcMsg := msg.(*pb.RpcReqData)
		lnet.Logger.Info("process msg", zap.Any("RemoteAddr", itransport.GetRemoteAddr()), zap.Any("msg", rpcMsg))
		rsp := &pb.RpcRspData{Info: &pb.RpcRspInfo{Uid: rpcMsg.Info.Uid}, Rsp: rpcMsg.Req + 1}
		msgPkg := this.CreateMessagePackage(rsp)
		time.Sleep(1 * time.Second)
		itransport.Send(msgPkg)

	} else if msgPackage.GetTag() == 14 {
		rpcMsg := msg.(*pb.RpcRspData)
		if rpcMsg.Info.Uid != "" {
			lrpc.G_asyncResult.FillAsyncResult(rpcMsg.Info.Uid, rpcMsg)
			return
		} else {
			lnet.Logger.Info("process msg", zap.Any("RemoteAddr", itransport.GetRemoteAddr()), zap.Any("msg", rpcMsg))
		}
	}
}
