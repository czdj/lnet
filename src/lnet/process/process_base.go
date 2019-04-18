package process

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"lnet/transport"
)

type BaseProcessor struct {

}

func (this *BaseProcessor)Process(itransport iface.ITransport, msg interface{}){
	t := itransport.(*transport.TcpTransport)
	lnet.Logger.Info("process msg",zap.Any("RemoteAddr",t.Conn.RemoteAddr()),zap.Any("msg",msg))
	t.Send(msg)
}
