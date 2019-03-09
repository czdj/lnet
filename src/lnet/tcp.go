package lnet

import (
	"fmt"
	"io"
	"net"
	"time"
	"unsafe"
)

type TcpTransport struct {
	DefTransport
	Conn net.Conn
}

func NewTcpTransport(netAddr string, timeout int, protocol  IProtocol, processor IProcessor,conn net.Conn) *TcpTransport{
	return  &TcpTransport{DefTransport:DefTransport{NetAddr:netAddr,StopFlag:0,cwrite:make(chan *[]byte,64),timeout:timeout,lastTick:time.Now().Unix(),protocol:protocol,processor:processor},Conn:conn}
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
			this.StopFlag = 1
			return err
		}
		tcpTransport := NewTcpTransport(this.NetAddr,DefMsgTimeout, this.protocol,this.processor,conn)
		this.OnNewConnect(tcpTransport)
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

	this.OnNewConnect(this)

	return nil
}

func (this *TcpTransport) read(){
	defer func() {
		this.OnClosed()
	}()

	for !this.IsStop(){
		head := PakgeHead{}
		headData := make([]byte, unsafe.Sizeof(PakgeHead{}))
		_, err := io.ReadFull(this.Conn, headData)
		if err != nil {
			fmt.Println("IO Read Err:%v",err)
			break
		}
		pakgeHead := (*PakgeHead)(unsafe.Pointer(&headData[0]))
		head.Len = pakgeHead.Len
		head.Tag = pakgeHead.Tag

		data := make([]byte,head.Len)
		_, err = io.ReadFull(this.Conn, data)
		if err != nil {
			fmt.Println("IO Read Err:%v",err)
			break
		}

		this.lastTick = time.Now().Unix()

		msg := this.protocol.Decode(head.Tag, data)
		this.processor.Process(this,msg)
	}
}

func (this *TcpTransport) write(){
	defer func() {
		this.Conn.Close()
		this.OnClosed()

		if err := recover(); err != nil {
			fmt.Println("Write panic:%v",err)
			return
		}
	}()

	var data *[]byte
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.IsStop(){
		select {
		case data = <-this.cwrite:
		case <-tick.C:
			if this.isTimeout(tick){
				this.OnClosed()
			}
		}

		if data == nil{
			continue
		}

		_,err := this.Conn.Write(*data)
		if err != nil{
			fmt.Println("Write Err:%v",err)
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *TcpTransport)Send(tag uint16, msg interface{})error{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Send panic:%v",err)
			return
		}
	}()

	if this.IsStop(){
		fmt.Println("Transport has been closed!!!")
		return nil
	}

	data := this.protocol.Encode(tag, msg)
	select {
	case this.cwrite <- &data:
	default:
		fmt.Println("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *TcpTransport) Close(){
	this.Conn.Close()
	this.StopFlag = 1
}


type TcpServer struct {
	DefServer
}

func NewTcpServer(netAddr string, protocol  IProtocol, processor IProcessor) *TcpServer{
	tcpServer := &TcpServer{DefServer:DefServer{NetType:TCP,NetAddr:netAddr,transport: NewTcpTransport(netAddr,DefMsgTimeout,protocol,processor,nil)}}

	return tcpServer
}