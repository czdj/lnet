package lclient

import (
	"lnet"
	"lnet/iface"
)

type BaseClient struct {
	NetType lnet.NetType
	NetAddr string

	transport iface.ITransport
}

//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *BaseClient) Connect() error{
	return this.transport.Connect()
}

func (this *BaseClient) Send(msg interface{})error{
	return this.transport.Send(msg)
}
