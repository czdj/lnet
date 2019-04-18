package server

import (
	"lnet"
	"lnet/iface"
	"lnet/transport"
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
	t := transport.NewTcpTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,tcpServer, nil)
	tcpServer.SetTransport(t)

	return tcpServer
}


