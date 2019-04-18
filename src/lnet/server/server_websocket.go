package server

import (
	"lnet"
	"lnet/iface"
	"lnet/transport"
)

type WebsocketServer struct {
	BaseServer
}

func NewWebsocketServer(netAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *WebsocketServer{
	WebsocketServer := &WebsocketServer{
		BaseServer:BaseServer{
			NetType:lnet.WebSocket,
			NetAddr:netAddr,
			transport: transport.NewWebsocketTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,nil),
		},
	}

	return WebsocketServer
}


