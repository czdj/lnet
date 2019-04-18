package client

import (
	"lnet"
	"lnet/iface"
	"lnet/transport"
)

type TcpClient struct {
	BaseClient
}

func NewTcpClient(netAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *TcpClient{
	tcpClient := &TcpClient{
		BaseClient:BaseClient{
			NetType:lnet.TCP,
			NetAddr:netAddr,
			transport: transport.NewTcpTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,nil),
		},
	}

	return tcpClient
}

