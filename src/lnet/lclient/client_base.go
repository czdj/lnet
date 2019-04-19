package lclient

import (
	"errors"
	"go.uber.org/zap"
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
	encodeData, err := this.msgHandle.GetProtocol().Marshal(msg)
	if err != nil {
		lnet.Logger.Error("msg marshal err", zap.Any("err", err))
		return errors.New("msg marshal err")
	}

	dp := lnet.NewDataPack()
	tag := this.msgHandle.GetMsgTag(msg)
	data, err := dp.Pack(lnet.NewMsgPackage(tag, encodeData))
	if err != nil {
		lnet.Logger.Error("数据打包错误", zap.Uint32("tag", tag), zap.Any("err", err))
		return errors.New("数据打包错误")
	}

	return this.transport.Send(data)
}
