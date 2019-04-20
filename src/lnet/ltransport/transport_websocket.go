package ltransport

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"lnet"
	"lnet/iface"
	"net/http"
	"time"
)

type WebsocketTransport struct {
	BaseTransport
	Conn     *websocket.Conn
	Upgrader *websocket.Upgrader
}

func NewWebsocketTransport(localAddr string, timeout int, msgHandle iface.IMsgHandle, server iface.IServer, conn *websocket.Conn) *WebsocketTransport {
	re := &WebsocketTransport{
		BaseTransport: *NewBaseTransport(localAddr, timeout, msgHandle, server),
		Conn:          conn,
		Upgrader:      &websocket.Upgrader{},
	}
	if conn != nil {
		re.RemoteAddr = conn.RemoteAddr().String()
	}
	return re
}

func (this *WebsocketTransport) websocketConnHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := this.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		lnet.Logger.Error("websocket upgrade err", zap.Any("err", err))
		return
	}

	if this.server.GetTransportMgr().Len() >= lnet.MAX_CONN {
		conn.Close()
		return
	}

	websocketTransport := NewWebsocketTransport(this.LocalAddr, lnet.DefMsgTimeout, this.msgHandle, this.server, conn)
	this.server.GetTransportMgr().Add(websocketTransport)
	this.OnNewConnect(websocketTransport)
}

func (this *WebsocketTransport) Listen() error {
	http.HandleFunc("/ws", this.websocketConnHandler)
	lnet.Logger.Info("WebsocketServer is listening", zap.Any("addr", this.LocalAddr))
	err := http.ListenAndServe(this.LocalAddr, nil)
	if err != nil {
		lnet.Logger.Error("Websocket Listen err", zap.Any("err", err))
		return err
	}

	return nil
}

func (this *WebsocketTransport) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(this.LocalAddr, nil)
	if err != nil {
		this.stopFlag = 1
		lnet.Logger.Error("Connect err", zap.Any("err", err))
		return err
	}
	lnet.Logger.Info("Connect Server", zap.Any("addr", this.LocalAddr))
	this.Conn = conn
	this.RemoteAddr = conn.RemoteAddr().String()

	this.OnNewConnect(this)

	return nil
}

func (this *WebsocketTransport) Read() {
	defer func() {
		this.OnClosed()
	}()

	for !this.IsStop() {
		dp := lnet.NewDataPack()

		_, data, err := this.Conn.ReadMessage()
		if err != nil {
			lnet.Logger.Error("IO Read Err", zap.Any("err", err))
			break
		}
		this.lastTick = time.Now().Unix()

		msgPackage, err := dp.Unpack(data)
		if err != nil {
			lnet.Logger.Error("data Unpack", zap.Any("err", err))
			break
		}

		this.msgHandle.Process(this, msgPackage)
	}
}

func (this *WebsocketTransport) Write() {
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
	dp := lnet.NewDataPack()
	tick := time.NewTimer(time.Duration(this.timeout) * time.Second)
	for !this.IsStop() {
		select {
		case msgPkg = <-this.cwrite:
		case <-tick.C:
			if this.IsTimeout(tick) {
				this.OnClosed()
			}
		}

		if msgPkg == nil {
			continue
		}

		data, err := dp.Pack(msgPkg)
		if err != nil {
			lnet.Logger.Error("Pack Err", zap.Any("err", err))
			break
		}
		err = this.Conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			lnet.Logger.Error("Write Err", zap.Any("err", err))
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
	tick.Stop()
}

func (this *WebsocketTransport) Send(msgPkg iface.IMessagePackage) error {
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
	case this.cwrite <- msgPkg:
	default:
		lnet.Logger.Info("write buf full!!!")
		this.cwrite <- msgPkg
	}
	return nil
}

func (this *WebsocketTransport) Close() {
	if this.Conn != nil {
		this.Conn.Close()
	}
	this.stopFlag = 1
}
