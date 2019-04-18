package server

import (
	"lnet"
	"lnet/iface"
)

type BaseServer struct {
	NetType lnet.NetType
	NetAddr string

	transport iface.ITransport
}
//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *BaseServer) Start(){
	go this.transport.Listen()
}
