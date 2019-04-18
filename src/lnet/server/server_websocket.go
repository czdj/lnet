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
	websocketServer := &WebsocketServer{
		BaseServer:BaseServer{
			NetType:lnet.WebSocket,
			NetAddr:netAddr,
			transport: nil,
			transportManager:lnet.NewTransportManager(),
		},
	}

	t := transport.NewWebsocketTransport(netAddr,lnet.DefMsgTimeout,protocol,processor,websocketServer, nil)
	websocketServer.SetTransport(t)

	return websocketServer
}


