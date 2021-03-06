package ltransport

import (
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"sync/atomic"
	"time"
)

var transportId uint32 //传输层id

type BaseTransport struct {
	Id         uint32
	NetType    lnet.NetType
	LocalAddr  string
	RemoteAddr string
	stopFlag   int32
	cwrite     chan iface.IMessagePackage
	timeout    int //传输超时
	lastTick   int64

	msgHandle iface.IMsgHandle
	server    iface.IServer
}

func NewBaseTransport(localAddr string, timeout int, msgHandle iface.IMsgHandle, server iface.IServer) *BaseTransport {
	return &BaseTransport{
		Id:        atomic.AddUint32(&transportId, 1),
		LocalAddr: localAddr,
		stopFlag:  0,
		cwrite:    make(chan iface.IMessagePackage, 64),
		timeout:   timeout,
		lastTick:  time.Now().Unix(),
		msgHandle: msgHandle,
		server:    server}
}

func (this *BaseTransport) GetLocalAddr() string {
	return this.LocalAddr
}

func (this *BaseTransport) GetRemoteAddr() string {
	return this.RemoteAddr
}

func (this *BaseTransport) GetId() uint32 {
	return this.Id
}

func (this *BaseTransport) Listen() error {
	return nil
}

func (this *BaseTransport) Connect() error {
	return nil
}

func (this *BaseTransport) OnNewConnect(transport iface.ITransport) {
	go transport.Read()
	go transport.Write()
}

func (this *BaseTransport) Read() {

}

func (this *BaseTransport) Write() {

}

func (this *BaseTransport) Send(msgPkg iface.IMessagePackage) error {
	return nil
}

func (this *BaseTransport) Close() {

}

func (this *BaseTransport) OnClosed() {
	if atomic.CompareAndSwapInt32(&this.stopFlag, 0, 1) {
		close(this.cwrite)
		if this.server != nil {
			this.server.GetTransportMgr().Remove(this)
		}
		this.msgHandle.OnTransportClose(this)
	}
}

func (this *BaseTransport) IsStop() bool {
	return this.stopFlag == 1
}

func (this *BaseTransport) IsTimeout(tick *time.Timer) bool {
	left := int(time.Now().Unix() - this.lastTick)
	if left < this.timeout {
		tick.Reset(time.Second * time.Duration(this.timeout))
		return false
	}
	lnet.Logger.Info("msgque close because timeout", zap.Int("wait", left), zap.Int("timeout", this.timeout))
	return true
}
