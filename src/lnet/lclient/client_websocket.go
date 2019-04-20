package lclient

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type WebsocketClient struct {
	BaseClient
}

func NewWebsocketClient(localAddr string, msgHandle iface.IMsgHandle) *WebsocketClient {
	WebsocketClient := &WebsocketClient{
		BaseClient: BaseClient{
			NetType:   lnet.WebSocket,
			LocalAddr: localAddr,
			transport: ltransport.NewWebsocketTransport(localAddr, lnet.DefMsgTimeout, msgHandle, nil, nil),
			msgHandle: msgHandle,
		},
	}

	return WebsocketClient
}
