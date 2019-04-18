package transport

import (
	"go.uber.org/zap"
	"io"
	"lnet"
	"lnet/iface"
	"net"
	"sync/atomic"
	"time"
	"unsafe"
)

type TcpTransport struct {
	BaseTransport
	Conn net.Conn
}

func NewTcpTransport(netAddr string, timeout int, protocol  iface.IProtocol, processor iface.IProcessor,server iface.IServer, conn net.Conn) *TcpTransport{
	return  &TcpTransport{
		BaseTransport:BaseTransport{
			Id:          atomic.AddUint32(&transportId, 1),
			LocalAddr:   netAddr,
			stopFlag:    0,
			cwrite:      make(chan *[]byte,64),
			timeout:     timeout,
			lastTick:    time.Now().Unix(),
			protocol:    protocol,
			processor:   processor,
			server:      server},
		Conn: conn,
	}
}

func (this *TcpTransport) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.LocalAddr);
	if err != nil{
		lnet.Logger.Error("Tcp Addr Err",zap.Any("err",err))
		return err
	}

	listen, err := net.ListenTCP("tcp", tcpAddr);
	if err != nil{
		lnet.Logger.Error("tcp listen err",zap.Any("err",err))
		return err
	}
	lnet.Logger.Info("TcpServer is listening",zap.Any("addr",tcpAddr))

	for !this.IsStop(){
		conn, err := listen.Accept();
		if err != nil{
			lnet.Logger.Error("tcp Accept err",zap.Any("err",err))
			this.stopFlag = 1
			return err
		}
		///TODO:配置
		if this.server.GetTransportMgr().Len() >= 30000{
			conn.Close()
			continue
		}

		tcpTransport := NewTcpTransport(this.LocalAddr,lnet.DefMsgTimeout, this.protocol,this.processor,this.server,conn)
		this.server.GetTransportMgr().Add(tcpTransport)

		this.OnNewConnect(tcpTransport)
	}

	return nil
}

func (this *TcpTransport) Connect() error{
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.LocalAddr);
	if err != nil{
		lnet.Logger.Error("tcp addr err",zap.Any("err",err))
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil{
		lnet.Logger.Error("Connect Server err",zap.Any("err",err))
		return err
	}
	lnet.Logger.Info("Connect Server ",zap.Any("Addr",tcpAddr))
	this.Conn = conn

	this.OnNewConnect(this)

	return nil
}

func (this *TcpTransport) Read(){
	defer func() {
		this.OnClosed()
	}()

	for !this.IsStop(){
		head := lnet.PakgeHead{}
		headData := make([]byte, unsafe.Sizeof(lnet.PakgeHead{}))
		_, err := io.ReadFull(this.Conn, headData)
		if err != nil {
			lnet.Logger.Error("IO Read Err",zap.Any("err",err))
			break
		}
		pakgeHead := (*lnet.PakgeHead)(unsafe.Pointer(&headData[0]))
		head.Len = pakgeHead.Len
		head.Tag = pakgeHead.Tag

		data := make([]byte,head.Len)
		_, err = io.ReadFull(this.Conn, data)
		if err != nil {
			lnet.Logger.Error("IO Read Err",zap.Any("err",err))
			break
		}

		this.lastTick = time.Now().Unix()

		msg := lnet.MsgTypeInfo.NewMsg(head.Tag)
		this.protocol.Unmarshal(data,msg)

		this.processor.Process(this,msg)
	}
}

func (this *TcpTransport) Write(){
	defer func() {
		this.Conn.Close()
		this.OnClosed()

		if err := recover(); err != nil {
			lnet.Logger.Error("Write panic",zap.Any("err",err))
			return
		}
	}()

	var data *[]byte
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.IsStop(){
		select {
		case data = <-this.cwrite:
		case <-tick.C:
			if this.IsTimeout(tick){
				this.OnClosed()
			}
		}

		if data == nil{
			continue
		}

		_,err := this.Conn.Write(*data)
		if err != nil{
			lnet.Logger.Error("Write Err",zap.Any("err",err))
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *TcpTransport)Send(msg interface{})error{
	defer func() {
		if err := recover(); err != nil {
			lnet.Logger.Error("Send panic",zap.Any("err",err))
			return
		}
	}()

	if this.IsStop(){
		lnet.Logger.Info("Transport has been closed!!!")
		return nil
	}

	encodeData, err := this.protocol.Marshal(msg)
	if err != nil{
		return err
	}

	head := &lnet.PakgeHead{}
	head.Tag = lnet.MsgTypeInfo.Tag(msg)
	head.Len = uint16(len(encodeData))

	data := make([]byte,unsafe.Sizeof(lnet.PakgeHead{}))
	ptr := (*lnet.PakgeHead)(unsafe.Pointer(&data[0]))
	ptr.Len = head.Len
	ptr.Tag = head.Tag
	data = append(data,encodeData...)

	select {
	case this.cwrite <- &data:
	default:
		lnet.Logger.Info("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *TcpTransport) Close(){
	this.Conn.Close()
	this.stopFlag = 1
}

