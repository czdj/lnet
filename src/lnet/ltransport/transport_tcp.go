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
		this.stopFlag = 1
		lnet.Logger.Error("tcp addr err", zap.Any("err", err))
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		this.stopFlag = 1
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
			lnet.LogStack()
			return
		}
	}()

	var msgPkg iface.IMessagePackage
	var data []byte
	var err error
	dp := lnet.NewDataPack()
	writeNum := 0
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.IsStop() {
		if msgPkg == nil {
			select {
			case msgPkg = <-this.cwrite:
				if msgPkg != nil {
					data, err = dp.Pack(msgPkg)
					if err != nil {
						lnet.Logger.Error("Pack Err", zap.Any("err", err))
						return
					}
				}
			case <-tick.C:
				if this.IsTimeout(tick) {
					this.OnClosed()
				}
			}
		}

		if msgPkg == nil {
			continue
		}

		n, err := this.Conn.Write(data[writeNum:])
		if err != nil {
			lnet.Logger.Error("Write Err", zap.Any("err", err))
			break
		}
		writeNum += n
		if writeNum == len(data) {
			msgPkg = nil
			writeNum = 0
		}

		this.lastTick = time.Now().Unix()
	}
	tick.Stop()
}

func (this *TcpTransport) Send(msgPkg iface.IMessagePackage) error {
	defer func() {
		if err := recover(); err != nil {
			lnet.Logger.Error("Send panic", zap.Any("err", err))
		}
	}()

	if this.IsStop() {
		lnet.Logger.Info("Transport has been closed!!!")
		return nil
	}

	select {
	case this.cwrite <- msgPkg:
	default:
		lnet.Logger.Info("write buf full!!!")
		this.cwrite <- msgPkg
	}
	return nil
}

func (this *TcpTransport) Close() {
	if this.Conn != nil {
		this.Conn.Close()
	}
	this.stopFlag = 1
}
