package client

import (
	"lnet"
	"lnet/iface"
	"lnet/transport"
)

type WebsocketClient struct {
	BaseClient
}

func NewWebsocketClient(netAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *WebsocketClient{
	WebsocketClient := &WebsocketClient{
		BaseClient:BaseClient{
			NetType:lnet.WebSocket,
			NetAddr:netAddr,
			transport: transport.NewWebsocketTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,nil),
		},
	}

	return WebsocketClient
}
