package lnet

import (
	"go.uber.org/zap"
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
	return  &TcpTransport{
		DefTransport:DefTransport{
			NetAddr:netAddr,
			stopFlag:0,
			cwrite:make(chan *[]byte,64),
			timeout:timeout,
			lastTick:time.Now().Unix(),
			protocol:protocol,
			processor:processor},
		Conn:conn,
	}
}

func (this *TcpTransport) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.NetAddr);
	if err != nil{
		Logger.Error("Tcp Addr Err",zap.Any("err",err))
		return err
	}

	listen, err := net.ListenTCP("tcp", tcpAddr);
	if err != nil{
		Logger.Error("tcp listen err",zap.Any("err",err))
		return err
	}
	Logger.Info("TcpServer is listening",zap.Any("addr",tcpAddr))

	for !this.isStop(){
		conn, err := listen.Accept();
		if err != nil{
			Logger.Error("tcp Accept err",zap.Any("err",err))
			this.stopFlag = 1
			return err
		}
		tcpTransport := NewTcpTransport(this.NetAddr,DefMsgTimeout, this.protocol,this.processor,conn)
		this.onNewConnect(tcpTransport)
	}

	return nil
}

func (this *TcpTransport) Connect() error{
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.NetAddr);
	if err != nil{
		Logger.Error("tcp addr err",zap.Any("err",err))
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil{
		Logger.Error("Connect Server err",zap.Any("err",err))
		return err
	}
	Logger.Info("Connect Server ",zap.Any("Addr",tcpAddr))
	this.Conn = conn

	this.onNewConnect(this)

	return nil
}

func (this *TcpTransport) read(){
	defer func() {
		this.onClosed()
	}()

	for !this.isStop(){
		head := PakgeHead{}
		headData := make([]byte, unsafe.Sizeof(PakgeHead{}))
		_, err := io.ReadFull(this.Conn, headData)
		if err != nil {
			Logger.Error("IO Read Err",zap.Any("err",err))
			break
		}
		pakgeHead := (*PakgeHead)(unsafe.Pointer(&headData[0]))
		head.Len = pakgeHead.Len
		head.Tag = pakgeHead.Tag

		data := make([]byte,head.Len)
		_, err = io.ReadFull(this.Conn, data)
		if err != nil {
			Logger.Error("IO Read Err",zap.Any("err",err))
			break
		}

		this.lastTick = time.Now().Unix()

		msg := MsgTypeInfo.NewMsg(head.Tag)
		this.protocol.Unmarshal(data,msg)

		this.processor.Process(this,msg)
	}
}

func (this *TcpTransport) write(){
	defer func() {
		this.Conn.Close()
		this.onClosed()

		if err := recover(); err != nil {
			Logger.Error("Write panic",zap.Any("err",err))
			return
		}
	}()

	var data *[]byte
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.isStop(){
		select {
		case data = <-this.cwrite:
		case <-tick.C:
			if this.isTimeout(tick){
				this.onClosed()
			}
		}

		if data == nil{
			continue
		}

		_,err := this.Conn.Write(*data)
		if err != nil{
			Logger.Error("Write Err",zap.Any("err",err))
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *TcpTransport)Send(msg interface{})error{
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("Send panic",zap.Any("err",err))
			return
		}
	}()

	if this.isStop(){
		Logger.Info("Transport has been closed!!!")
		return nil
	}

	encodeData, err := this.protocol.Marshal(msg)
	if err != nil{
		return err
	}

	head := &PakgeHead{}
	head.Tag = MsgTypeInfo.Tag(msg)
	head.Len = uint16(len(encodeData))

	data := make([]byte,unsafe.Sizeof(PakgeHead{}))
	ptr := (*PakgeHead)(unsafe.Pointer(&data[0]))
	ptr.Len = head.Len
	ptr.Tag = head.Tag
	data = append(data,encodeData...)

	select {
	case this.cwrite <- &data:
	default:
		Logger.Info("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *TcpTransport) Close(){
	this.Conn.Close()
	this.stopFlag = 1
}


type TcpServer struct {
	DefServer
}

func NewTcpServer(netAddr string, protocol  IProtocol, processor IProcessor) *TcpServer{
	tcpServer := &TcpServer{
		DefServer:DefServer{
			NetType:TCP,
			NetAddr:netAddr,
			transport: NewTcpTransport(netAddr,DefMsgTimeout,protocol,processor,nil),
		},
	}

	return tcpServer
}

type TcpClient struct {
	DefClient
}

func NewTcpClient(netAddr string, protocol  IProtocol, processor IProcessor) *TcpClient{
	tcpClient := &TcpClient{
		DefClient:DefClient{
			NetType:TCP,
			NetAddr:netAddr,
			transport: NewTcpTransport(netAddr,DefMsgTimeout,protocol,processor,nil),
		},
	}

	return tcpClient
}