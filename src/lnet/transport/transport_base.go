package transport

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"sync/atomic"
	"time"
)

var transportId uint32 //传输层id

type BaseTransport struct{
	Id uint32
	NetType lnet.NetType
	LocalAddr string
	PeerAddr string
	stopFlag int32
	cwrite chan *[]byte
	timeout       int //传输超时
	lastTick      int64

	protocol iface.IProtocol
	processor iface.IProcessor
	server    iface.IServer
}


func (this *BaseTransport) GetId() uint32{
	return this.Id
}

func (this *BaseTransport) Listen() error{
	return nil
}

func (this *BaseTransport) Connect() error{
	return nil
}

func (this *BaseTransport) OnNewConnect(transport iface.ITransport){
	go transport.Read()
	go transport.Write()
}

func (this *BaseTransport) Read(){

}

func (this *BaseTransport) Write(){

}

func (this *BaseTransport)Send(msg interface{})error{
	return nil
}

func (this *BaseTransport) Close(){

}

func (this *BaseTransport) OnClosed(){
	if atomic.CompareAndSwapInt32(&this.stopFlag,0,1){
		close(this.cwrite)
		lnet.Logger.Info("connect closed !")
	}
}
func (this *BaseTransport)IsStop() bool{
	return this.stopFlag == 1
}

func (this *BaseTransport)IsTimeout(tick *time.Timer) bool{
	left := int(time.Now().Unix() - this.lastTick)
	if left < this.timeout  {
		tick.Reset(time.Second * time.Duration(this.timeout))
		return false
	}
	lnet.Logger.Info("msgque close because timeout",zap.Int("wait",left),zap.Int("timeout",this.timeout))
	return true
}
