package lnet

import (
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

//负责网络相关功能的处理
type ITransport interface {
	Listen() error
	Connect() error
	onNewConnect(transport ITransport)
	read()
	write()
	Send(msg interface{})error
	Close()
	onClosed()
	isStop() bool
	isTimeout(tick *time.Timer) bool
}

type DefTransport struct{
	Id int
	NetType NetType
	NetAddr string
	PeerAddr string
	stopFlag int32
	cwrite chan *[]byte
	timeout       int //传输超时
	lastTick      int64

	protocol IProtocol
	processor IProcessor
}

func (this *DefTransport) Listen() error{
	return nil
}

func (this *DefTransport) onNewConnect(transport ITransport){
	go transport.read()
	go transport.write()
}

func (this *DefTransport) Connect() error{
	return nil
}

func (this *DefTransport) read(){

}

func (this *DefTransport) write(){

}

func (this *DefTransport)Send(tag uint16, msg interface{})error{
	return nil
}

func (this *DefTransport) Close(){

}

func (this *DefTransport) onClosed(){
	if atomic.CompareAndSwapInt32(&this.stopFlag,0,1){
		close(this.cwrite)
		Logger.Info("connect closed !")
	}
}
func (this *DefTransport)isStop() bool{
	return this.stopFlag == 1
}

func (this *DefTransport)isTimeout(tick *time.Timer) bool{
	left := int(time.Now().Unix() - this.lastTick)
	if left < this.timeout  {
		tick.Reset(time.Second * time.Duration(this.timeout))
		return false
	}
	Logger.Info("msgque close because timeout",zap.Int("wait",left),zap.Int("timeout",this.timeout))
	return true
}

//负责解析协议
type IProtocol interface {
	Marshal(msg interface{})([]byte,error)
	Unmarshal(data []byte, v interface{}) error
}

//负责业务处理
type IProcessor  interface {
	Process(transport ITransport, msg interface{})
}

type DefProcessor struct {

}

func (this *DefProcessor)Process(transport ITransport, msg interface{}){
	t := transport.(*TcpTransport)
	Logger.Info("process msg",zap.Any("RemoteAddr",t.Conn.RemoteAddr()),zap.Any("msg",msg))
	transport.Send(msg)
}

type Server interface {
	Start()
}

type DefServer struct {
	NetType NetType
	NetAddr string

	transport ITransport
}
//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *DefServer) Start(){
	go this.transport.Listen()
}

type Client interface {
	Connect() error
	Send(msg interface{})error
}

type DefClient struct {
	NetType NetType
	NetAddr string

	transport ITransport
}
//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *DefClient) Connect() error{
	return this.transport.Connect()
}

func (this *DefClient) Send(msg interface{})error{
	return this.transport.Send(msg)
}