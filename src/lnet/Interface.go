package lnet

import "fmt"

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
	data []byte
}

type Message struct {
	data string
}

//处理每一个连接的数据接受
type Connect interface {
	Read()
	Write()
	Close()
	IsStop() bool
}

type DefConnect struct {
	id int
	netType NetType
	protocol Protocol
	processor Processor
}

func (this *DefConnect) Read(){

}

func (this *DefConnect) Write(){

}

func (this *DefConnect) Close(){

}

func (this *DefConnect) IsStop() bool{
	return false
}

//监听类，负责接收连接
type Transport interface {
	Listen() error
	Read()
	Write(msg *Message)
	Close()
	IsStop() bool
}

type DefTransport struct{
	id int
	netType NetType
	Net  string
	Addr string
	PeerAddr string
	StopFlag bool

	protocol Protocol
	processor Processor
}

func (this *DefTransport) Listen() error{
	return nil
}

func (this *DefTransport) Read(){

}

func (this *DefTransport) Write(msg *Message){

}

func (this *DefTransport) Close(){

}

func (this *DefTransport)IsStop() bool{
	return false
}


//负责解析协议
type Protocol interface {
	Encode(msg *Message) *Pakge
	Decode(pakge *Pakge) *Message
}

//负责业务处理
type Processor  interface {
	Process(msg *Message)
}

type DefProcessor struct {

}

func (this *DefProcessor)Process(msg *Message){
	fmt.Printf("process:%v",msg)
}

type Server interface {
	Start()
}

type DefServer struct {
	netType NetType
	netAddr string

	transport Transport
	protocol  Protocol
	processor Processor
}
//接受连接，每个连接对应一个结构，每个连接开一个goroution，每一个连接里处理读写消息
func (this *DefServer) Start(){
	this.transport.Listen()
}