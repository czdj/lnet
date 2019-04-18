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
			transport: transport.NewTcpTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,nil),
		},
	}

	return tcpServer
}


