package lserver

import (
	"lnet"
	"lnet/iface"
	"lnet/ltransport"
)

type WebsocketServer struct {
	BaseServer
}

func NewWebsocketServer(localAddr string, msgHandle iface.IMsgHandle) *WebsocketServer {
	websocketServer := &WebsocketServer{
		BaseServer: BaseServer{
			NetType:          lnet.WebSocket,
			LocalAddr:        localAddr,
			transport:        nil,
			transportManager: lnet.NewTransportManager(),
			msgHandle:        msgHandle,
		},
	}

	t := ltransport.NewWebsocketTransport(localAddr, lnet.DefMsgTimeout, msgHandle, websocketServer, nil)
	websocketServer.SetTransport(t)

	return websocketServer
}
