package lclient

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type WebsocketClient struct {
	BaseClient
}

func NewWebsocketClient(netAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *WebsocketClient{
	WebsocketClient := &WebsocketClient{
		BaseClient:BaseClient{
			NetType:   lnet.WebSocket,
			NetAddr:   netAddr,
			transport: ltransport.NewWebsocketTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,nil, nil),
		},
	}

	return WebsocketClient
}
