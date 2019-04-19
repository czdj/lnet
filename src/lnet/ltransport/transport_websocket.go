package ltransport

import (
	"errors"
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

func NewWebsocketTransport(localAddr string, timeout int, protocol iface.IProtocol, processor iface.IProcessor, server iface.IServer, conn *websocket.Conn) *WebsocketTransport {
	re := &WebsocketTransport{
		BaseTransport: *NewBaseTransport(localAddr, timeout, protocol, processor, server),
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

	WebsocketTransport := NewWebsocketTransport(this.LocalAddr, lnet.DefMsgTimeout, this.protocol, this.processor, this.server, conn)
	this.OnNewConnect(WebsocketTransport)
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

		msg := lnet.MsgTypeInfo.NewMsg(msgPackage.GetTag())
		this.protocol.Unmarshal(msgPackage.GetData(), msg)

		this.processor.Process(this, msg)
	}
}

func (this *WebsocketTransport) Write() {
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

		err := this.Conn.WriteMessage(websocket.BinaryMessage, *data)
		if err != nil {
			lnet.Logger.Error("Write Err", zap.Any("err", err))
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *WebsocketTransport) Send(msg interface{}) error {
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

	encodeData, err := this.protocol.Marshal(msg)
	if err != nil {
		lnet.Logger.Error("数据编码错误", zap.Any("err", err))
		return err
	}

	dp := lnet.NewDataPack()
	tag := lnet.MsgTypeInfo.Tag(msg)
	data, err := dp.Pack(lnet.NewMsgPackage(tag, encodeData))
	if err != nil {
		lnet.Logger.Error("数据打包错误", zap.Uint32("tag", tag), zap.Any("err", err))
		return errors.New("Pack error msg ")
	}

	select {
	case this.cwrite <- &data:
	default:
		lnet.Logger.Info("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *WebsocketTransport) Close() {
	this.Conn.Close()
	this.stopFlag = 1
}
