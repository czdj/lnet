package lnet

import (
	"fmt"
	"reflect"
)

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
	head PakgeHead
	data interface{}
}

type Message struct {
	Data string
}

//监听类，负责接收连接
type Transport interface {
	Listen() error
	Connect() error
	Read()
	Write(tag uint16, msg interface{})
	Close()
	IsStop() bool
}

type DefTransport struct{
	Id int
	NetType NetType
	NetAddr string
	PeerAddr string
	StopFlag bool

	protocol Protocol
	processor Processor
}

func (this *DefTransport) Listen() error{
	return nil
}

func (this *DefTransport) Connect() error{
	return nil
}

func (this *DefTransport) Read(){

}

func (this *DefTransport) Write(tag uint16, msg interface{}){

}

func (this *DefTransport) Close(){

}

func (this *DefTransport)IsStop() bool{
	return false
}


//负责解析协议
type Protocol interface {
	Encode(tag uint16, msg interface{}) []byte
	Decode(tag uint16, data []byte) interface{}
}

//负责业务处理
type Processor  interface {
	Process(msg interface{})
}

type DefProcessor struct {

}

func (this *DefProcessor)Process(msg interface{}){
	fmt.Println("process:%v",msg)
}

type Server interface {
	Start()
}

type DefServer struct {
	NetType NetType
	NetAddr string

	transport Transport
}
//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *DefServer) Start(){
	this.transport.Listen()
}