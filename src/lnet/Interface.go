package lnet

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"time"
)

var  DefMsgTimeout int  = 2

type MsgTypeMap struct {
	msgTypeMap map[uint16]reflect.Type
}
func (this *MsgTypeMap)Register(tag uint16,msg interface{}){
	msgType := reflect.TypeOf(msg)
	this.msgTypeMap[tag] = msgType
}

func (this *MsgTypeMap)NewMsg(tag uint16)interface{}{
	 msgType,err := this.msgTypeMap[tag]
	if err == false{
		fmt.Println("Msg Type Err!")
		return nil
	}
	 msg := reflect.New(msgType).Interface()
	return msg
}

var MsgTypeInfo MsgTypeMap = MsgTypeMap{msgTypeMap:make(map[uint16]reflect.Type)}


type NetType int
const (
	TCP NetType = iota
	UDP
	WebSocket
)
type PakgeHead struct {
	Tag uint16
	Len uint16
}
type Pakge struct {
	//head PakgeHead
	data []byte
}

type Message struct {
	Data string
}

//监听类，负责接收连接
type ITransport interface {
	Listen() error
	OnNewConnect(transport ITransport)
	Connect() error
	read()
	write()
	Send(tag uint16, msg interface{})error
	Close()
	OnClosed()
	IsStop() bool
	isTimeout(tick *time.Timer) bool
}

type DefTransport struct{
	Id int
	NetType NetType
	NetAddr string
	PeerAddr string
	StopFlag int32
	cwrite chan *[]byte
	timeout       int //传输超时
	lastTick      int64

	protocol IProtocol
	processor IProcessor
}

func (this *DefTransport) Listen() error{
	return nil
}

func (this *DefTransport) OnNewConnect(transport ITransport){
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

func (this *DefTransport) OnClosed(){
	if atomic.CompareAndSwapInt32(&this.StopFlag,0,1){
		close(this.cwrite)
		fmt.Println("connect closed !!")
	}
}
func (this *DefTransport)IsStop() bool{
	return this.StopFlag == 1
}

func (this *DefTransport)isTimeout(tick *time.Timer) bool{
	left := int(time.Now().Unix() - this.lastTick)
	if left < this.timeout  {
		tick.Reset(time.Second * time.Duration(this.timeout))
		return false
	}
	fmt.Println("msgque close because timeout wait:%v timeout:%v", left, this.timeout)
	return true
}

//负责解析协议
type IProtocol interface {
	Encode(tag uint16, msg interface{}) []byte
	Decode(tag uint16, data []byte) interface{}
}

//负责业务处理
type IProcessor  interface {
	Process(transport ITransport, msg interface{})
}

type DefProcessor struct {
	transport ITransport
}

func (this *DefProcessor)Process(transport ITransport, msg interface{}){
	fmt.Println("process:%v",msg)
	//transport.Send(11,msg)
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
	this.transport.Listen()
}