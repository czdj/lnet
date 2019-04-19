package lclient

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type TcpClient struct {
	BaseClient
}

func NewTcpClient(netAddr string, protocol iface.IProtocol, processor iface.IProcessor) *TcpClient {
	tcpClient := &TcpClient{
		BaseClient: BaseClient{
			NetType:   lnet.TCP,
			NetAddr:   netAddr,
			transport: ltransport.NewTcpTransport(netAddr, lnet.DefMsgTimeout, protocol, processor, nil, nil),
		},
	}

	return tcpClient
}
