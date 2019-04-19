package lserver

import (
	"lnet"
	"lnet/iface"
)

type BaseServer struct {
	Name string
	NetType lnet.NetType
	LocalAddr string //IP:Port

	transport iface.ITransport
	transportManager iface.ITransportManager
}

func (this *BaseServer) SetTransport(transport iface.ITransport){
	this.transport = transport
}

func (this *BaseServer) GetTransportMgr() iface.ITransportManager{
	return this.transportManager
}

//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *BaseServer) Start(){
	go this.transport.Listen()
}
