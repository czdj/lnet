package lnet

import (
	"fmt"
	"io"
	"net"
	"unsafe"
)

type TcpConnect struct {
	DefConnect
	conn net.Conn
}

func (this *TcpConnect) Read(){
	for {
		pakge := Pakge{}
		head := make([]byte, unsafe.Sizeof(PakgeHead{}))
		_, err := io.ReadFull(this.conn, head)
		if err != nil {
			break
		}
		pakgeHead := (*PakgeHead)(unsafe.Pointer(&head[0]))
		pakge.head.Len = pakgeHead.Len
		pakge.head.Tag = pakgeHead.Tag

		pakge.data = make([]byte,pakgeHead.Len)
		_, err = io.ReadFull(this.conn, pakge.data)
		if err != nil {
			break
		}

		msg := this.protocol.Decode(&pakge)
		this.processor.Process(msg)
	}
}

func (this *TcpConnect) Write(msg *Message){
	pakge := this.protocol.Encode(msg)
	this.conn.Write(*(*[]byte)(unsafe.Pointer(pakge)))
}

func (this *TcpConnect) Close(){
	this.conn.Close()
}

type TcpTransport struct {
	DefTransport
	conn net.Conn
}

func (this *TcpTransport) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.Addr);
	if err != nil{
		fmt.Printf("addr err:%v!",err)
		return err
	}
	listen, err := net.ListenTCP("tcp", tcpAddr);
	if err != nil{
		fmt.Printf("tcp listen err:%v!",err)
		return err
	}

	for !this.IsStop(){
		conn, err := listen.Accept();
		if err != nil{
			fmt.Printf("tcp Accept err:%v!",err)
			this.StopFlag = true
			return err
		}

		tcpTransport := &TcpTransport{DefTransport:DefTransport{},conn:conn}
		go tcpTransport.Read()
		//go tcpConnect.Write()
	}

	return nil
}

func (this *TcpTransport) Read(){
	for !this.IsStop(){
		pakge := Pakge{}
		head := make([]byte, unsafe.Sizeof(PakgeHead{}))
		_, err := io.ReadFull(this.conn, head)
		if err != nil {
			this.StopFlag = true
			break
		}
		pakgeHead := (*PakgeHead)(unsafe.Pointer(&head[0]))
		pakge.head.Len = pakgeHead.Len
		pakge.head.Tag = pakgeHead.Tag

		pakge.data = make([]byte,pakgeHead.Len)
		_, err = io.ReadFull(this.conn, pakge.data)
		if err != nil {
			this.StopFlag = true
			break
		}

		msg := this.protocol.Decode(&pakge)
		this.processor.Process(msg)
	}
}

func (this *TcpTransport) Write(msg *Message){
	pakge := this.protocol.Encode(msg)
	this.conn.Write(*(*[]byte)(unsafe.Pointer(pakge)))
}

func (this *TcpTransport) Close(){
	this.conn.Close()
	this.StopFlag = true
}


type TcpServer struct {
	DefServer
}
var(

)
func NewTcpServer( protocol  Protocol, processor Processor) *TcpServer{
	tcpServer := &TcpServer{}
	tcpServer.transport = &TcpTransport{}
	tcpServer. protocol = protocol
	tcpServer.processor = processor

	return tcpServer
}