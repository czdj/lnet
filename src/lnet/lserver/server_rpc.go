package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type RpcServer struct {
	BaseServer
}

func NewRpcServer(localAddr string, msgHandle iface.IMsgHandle) *RpcServer {
	rpcServer := &RpcServer{
		BaseServer: BaseServer{
			NetType:          lnet.TCP,
			LocalAddr:        localAddr,
			transport:        nil,
			transportManager: lnet.NewTransportManager(),
			msgHandle:        msgHandle,
		},
	}
	t := ltransport.NewTcpTransport(localAddr, lnet.DefMsgTimeout, msgHandle, rpcServer, nil)
	rpcServer.SetTransport(t)

	return rpcServer
}
