package lprocess

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
)

type BaseProcessor struct {
}

func (this *BaseProcessor) Process(itransport iface.ITransport, msg interface{}) {
	lnet.Logger.Info("process msg", zap.Any("RemoteAddr", itransport.GetRemoteAddr()), zap.Any("msg", msg))
	itransport.Send(msg)
}
