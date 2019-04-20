package lclient

import (
	"lnet"
	"lnet/iface"
)

type BaseClient struct {
	NetType   lnet.NetType
	LocalAddr string
	transport iface.ITransport
	msgHandle iface.IMsgHandle
}

//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *BaseClient) Connect() error {
	return this.transport.Connect()
}

func (this *BaseClient) Send(msg interface{}) error {
	msgPkg := this.msgHandle.CreateMessagePackage(msg)
	return this.transport.Send(msgPkg)
}
