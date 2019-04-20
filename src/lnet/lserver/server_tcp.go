package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type TcpServer struct {
	BaseServer
}

func NewTcpServer(localAddr string, msgHandle iface.IMsgHandle) *TcpServer {
	tcpServer := &TcpServer{
		BaseServer: BaseServer{
			NetType:          lnet.TCP,
			LocalAddr:        localAddr,
			transport:        nil,
			transportManager: lnet.NewTransportManager(),
			msgHandle:        msgHandle,
		},
	}
	t := ltransport.NewTcpTransport(localAddr, lnet.DefMsgTimeout, msgHandle, tcpServer, nil)
	tcpServer.SetTransport(t)

	return tcpServer
}
