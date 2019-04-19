package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type TcpServer struct {
	BaseServer
}

func NewTcpServer(netAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *TcpServer{
	tcpServer := &TcpServer{
		BaseServer:BaseServer{
			NetType:lnet.TCP,
			NetAddr:netAddr,
			transport: nil,
			transportManager:lnet.NewTransportManager(),
		},
	}
	t := ltransport.NewTcpTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,tcpServer, nil)
	tcpServer.SetTransport(t)

	return tcpServer
}


