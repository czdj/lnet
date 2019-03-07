package lnet

import (
	"fmt"
	"io"
	"net"
	"unsafe"
)

type TcpTransport struct {
	DefTransport
	Conn net.Conn
}

func NewTcpTransport(netAddr string, protocol  Protocol, processor Processor,conn net.Conn) *TcpTransport{
	return  &TcpTransport{DefTransport:DefTransport{NetAddr:netAddr,StopFlag:false,protocol:protocol,processor:processor},Conn:conn}
}

func (this *TcpTransport) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.NetAddr);
	if err != nil{
		fmt.Println("addr err:%v!",err)
		return err
	}
	listen, err := net.ListenTCP("tcp", tcpAddr);
	if err != nil{
		fmt.Println("tcp listen err:%v!",err)
		return err
	}
	fmt.Println("TcpServer is listening addr:%v!",tcpAddr)

	for !this.IsStop(){
		conn, err := listen.Accept();
		if err != nil{
			fmt.Println("tcp Accept err:%v!",err)
			this.StopFlag = true
			return err
		}
		tcpTransport := NewTcpTransport(this.NetAddr,this.protocol,this.processor,conn)
		go tcpTransport.Read()
		//go tcpConnect.Write()
	}

	return nil
}

func (this *TcpTransport) Connect() error{
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.NetAddr);
	if err != nil{
		fmt.Println("addr err:%v!",err)
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil{
		fmt.Println("Connect Server err:%v!",err)
		return err
	}
	fmt.Println("Connect Server Addr:%v!",tcpAddr)
	this.Conn = conn

	return nil
}


func (this *TcpTransport) Read(){
	for !this.IsStop(){
		head := PakgeHead{}
		headData := make([]byte, unsafe.Sizeof(PakgeHead{}))
		_, err := io.ReadFull(this.Conn, headData)
		if err != nil {
			this.StopFlag = true
			fmt.Println("IO Read Err:%v",err)
			break
		}
		pakgeHead := (*PakgeHead)(unsafe.Pointer(&headData[0]))
		head.Len = pakgeHead.Len
		head.Tag = pakgeHead.Tag

		data := make([]byte,head.Len)
		_, err = io.ReadFull(this.Conn, data)
		if err != nil {
			this.StopFlag = true
			fmt.Println("IO Read Err:%v",err)
			break
		}

		msg := this.protocol.Decode(head.Tag, data)
		this.processor.Process(msg)
	}
}

func (this *TcpTransport) Write(tag uint16, msg *Message){
	data := this.protocol.Encode(tag, msg)

	_,err := this.Conn.Write(data)
	if err != nil{
		fmt.Println("Write Err:%v",err)
	}

	//b := *((*[]byte)(unsafe.Pointer(&pakge.head)))
	//_,err := this.Conn.Write(b)
	//if err != nil{
	//	fmt.Println("Write Err:%v",err)
	//}
	//
	//_,err = this.Conn.Write(pakge.data)
	//if err != nil{
	//	fmt.Println("Write Err:%v",err)
	//}
}

func (this *TcpTransport) Close(){
	this.Conn.Close()
	this.StopFlag = true
}


type TcpServer struct {
	DefServer
}
var(

)
func NewTcpServer(netAddr string, protocol  Protocol, processor Processor) *TcpServer{
	tcpServer := &TcpServer{DefServer:DefServer{NetType:TCP,NetAddr:netAddr,transport: NewTcpTransport(netAddr,protocol,processor,nil)}}

	return tcpServer
}