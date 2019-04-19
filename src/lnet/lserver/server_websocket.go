package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type WebsocketServer struct {
	BaseServer
}

func NewWebsocketServer(localAddr string, protocol  iface.IProtocol, processor iface.IProcessor) *WebsocketServer{
	websocketServer := &WebsocketServer{
		BaseServer:BaseServer{
			NetType:lnet.WebSocket,
			LocalAddr:localAddr,
			transport: nil,
			transportManager:lnet.NewTransportManager(),
		},
	}

	t := ltransport.NewWebsocketTransport(localAddr,lnet.DefMsgTimeout,protocol,processor,websocketServer, nil)
	websocketServer.SetTransport(t)

	return websocketServer
}


