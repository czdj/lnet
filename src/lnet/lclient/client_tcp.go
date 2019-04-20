package lclient

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type TcpClient struct {
	BaseClient
}

func NewTcpClient(localAddr string, msgHandle iface.IMsgHandle) *TcpClient {
	tcpClient := &TcpClient{
		BaseClient: BaseClient{
			NetType:   lnet.TCP,
			LocalAddr: localAddr,
			transport: ltransport.NewTcpTransport(localAddr, lnet.DefMsgTimeout, msgHandle, nil, nil),
			msgHandle: msgHandle,
		},
	}

	return tcpClient
}
