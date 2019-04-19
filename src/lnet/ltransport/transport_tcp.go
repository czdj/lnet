package ltransport

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"lnet"
	"lnet/iface"
	"net"
	"time"
)

type TcpTransport struct {
	BaseTransport
	Conn net.Conn
}

func NewTcpTransport(localAddr string, timeout int, msgHandle iface.IMsgHandle, server iface.IServer, conn net.Conn) *TcpTransport {
	re := &TcpTransport{
		BaseTransport: *NewBaseTransport(localAddr, timeout, msgHandle, server),
		Conn:          conn,
	}
	if conn != nil {
		re.RemoteAddr = conn.RemoteAddr().String()
	}
	return re
}

func (this *TcpTransport) Listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.LocalAddr)
	if err != nil {
		lnet.Logger.Error("Tcp Addr Err", zap.Any("err", err))
		return err
	}

	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		lnet.Logger.Error("tcp listen err", zap.Any("err", err))
		return err
	}
	lnet.Logger.Info("TcpServer is listening", zap.Any("addr", tcpAddr))

	for !this.IsStop() {
		conn, err := listen.Accept()
		if err != nil {
			lnet.Logger.Error("tcp Accept err", zap.Any("err", err))
			this.stopFlag = 1
			return err
		}
		///TODO:配置
		if this.server.GetTransportMgr().Len() >= lnet.MAX_CONN {
			conn.Close()
			continue
		}

		tcpTransport := NewTcpTransport(this.LocalAddr, lnet.DefMsgTimeout, this.msgHandle, this.server, conn)
		this.server.GetTransportMgr().Add(tcpTransport)

		this.OnNewConnect(tcpTransport)
	}

	return nil
}

func (this *TcpTransport) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.LocalAddr)
	if err != nil {
		lnet.Logger.Error("tcp addr err", zap.Any("err", err))
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		lnet.Logger.Error("Connect Server err", zap.Any("err", err))
		return err
	}
	lnet.Logger.Info("Connect Server ", zap.Any("Addr", tcpAddr))
	this.Conn = conn
	this.RemoteAddr = conn.RemoteAddr().String()

	this.OnNewConnect(this)

	return nil
}

func (this *TcpTransport) Read() {
	defer func() {
		this.OnClosed()
	}()

	for !this.IsStop() {
		dp := lnet.NewDataPack()

		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(this.Conn, headData)
		if err != nil {
			lnet.Logger.Error("IO Read Err", zap.Any("err", err))
			break
		}

		msgPackge, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		data := make([]byte, msgPackge.GetLen())
		_, err = io.ReadFull(this.Conn, data)
		if err != nil {
			lnet.Logger.Error("IO Read Err", zap.Any("err", err))
			break
		}
		msgPackge.SetData(data)

		this.lastTick = time.Now().Unix()

		this.msgHandle.Process(this, msgPackge)
	}
}

func (this *TcpTransport) Write() {
	defer func() {
		this.Conn.Close()
		this.OnClosed()

		if err := recover(); err != nil {
			lnet.Logger.Error("Write panic", zap.Any("err", err))
			return
		}
	}()

	var data *[]byte
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.IsStop() {
		select {
		case data = <-this.cwrite:
		case <-tick.C:
			if this.IsTimeout(tick) {
				this.OnClosed()
			}
		}

		if data == nil {
			continue
		}

		_, err := this.Conn.Write(*data)
		if err != nil {
			lnet.Logger.Error("Write Err", zap.Any("err", err))
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *TcpTransport) Send(data []byte) error {
	defer func() {
		if err := recover(); err != nil {
			lnet.Logger.Error("Send panic", zap.Any("err", err))
			return
		}
	}()

	if this.IsStop() {
		lnet.Logger.Info("Transport has been closed!!!")
		return nil
	}

	select {
	case this.cwrite <- &data:
	default:
		lnet.Logger.Info("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *TcpTransport) Close() {
	this.Conn.Close()
	this.stopFlag = 1
}
