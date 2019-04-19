package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type TcpServer struct {
	BaseServer
}

func NewTcpServer(localAddr string, protocol iface.IProtocol, processor iface.IProcessor) *TcpServer {
	tcpServer := &TcpServer{
		BaseServer: BaseServer{
			NetType:          lnet.TCP,
			LocalAddr:        localAddr,
			transport:        nil,
			transportManager: lnet.NewTransportManager(),
		},
	}
	t := ltransport.NewTcpTransport(localAddr, lnet.DefMsgTimeout, protocol, processor, tcpServer, nil)
	tcpServer.SetTransport(t)

	return tcpServer
}
