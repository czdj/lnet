package lnet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"unsafe"
)


type WebsocketTransport struct {
	DefTransport
	Conn *websocket.Conn
	Upgrader *websocket.Upgrader
}

func NewWebsocketTransport(netAddr string, timeout int, protocol  IProtocol, processor IProcessor,conn *websocket.Conn) *WebsocketTransport{
	return  &WebsocketTransport{
		DefTransport:DefTransport{
			NetAddr:netAddr,
			stopFlag:0,
			cwrite:make(chan *[]byte,64),
			timeout:timeout,
			lastTick:time.Now().Unix(),
			protocol:protocol,
			processor:processor},
		Conn:conn,
		Upgrader:&websocket.Upgrader{},
	}
}

func (this *WebsocketTransport) websocketConnHandler(w http.ResponseWriter, r *http.Request){
	conn, err := this.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	WebsocketTransport := NewWebsocketTransport(this.NetAddr,DefMsgTimeout, this.protocol,this.processor,conn)
	this.onNewConnect(WebsocketTransport)
}

func (this *WebsocketTransport) Listen() error {
	http.HandleFunc("/ws",this.websocketConnHandler);
	fmt.Println("WebsocketServer is listening addr:%v!",this.NetAddr)
	err := http.ListenAndServe(this.NetAddr, nil);
	if err != nil{
		fmt.Println("Websocket Listen err:%v!",err)
		return err
	}

	return nil
}

func (this *WebsocketTransport) Connect() error{
	conn, _,err := websocket.DefaultDialer.Dial(this.NetAddr, nil);
	if err != nil{
		fmt.Println("Connect err:%v!",err)
		return err
	}

	fmt.Println("Connect Server Addr:%v!",this.NetAddr)
	this.Conn = conn

	this.onNewConnect(this)

	return nil
}

func (this *WebsocketTransport) read(){
	defer func() {
		this.onClosed()
	}()

	for !this.isStop(){
		_, data, err := this.Conn.ReadMessage()
		if err != nil {
			fmt.Println("IO Read Err:%v",err)
			break
		}
		this.lastTick = time.Now().Unix()

		pakgeHead := (*PakgeHead)(unsafe.Pointer(&data[0]))
		tag := pakgeHead.Tag
		data = data[unsafe.Sizeof(PakgeHead{}):]

		msg := MsgTypeInfo.NewMsg(tag)
		this.protocol.Unmarshal(data,msg)

		this.processor.Process(this,msg)
	}
}

func (this *WebsocketTransport) write(){
	defer func() {
		this.Conn.Close()
		this.onClosed()

		if err := recover(); err != nil {
			fmt.Println("Write panic:%v",err)
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

		err := this.Conn.WriteMessage(websocket.BinaryMessage,*data)
		if err != nil{
			fmt.Println("Write Err:%v",err)
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *WebsocketTransport)Send(msg interface{})error{
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Send panic:%v",err)
			return
		}
	}()

	if this.isStop(){
		fmt.Println("Transport has been closed!!!")
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
		fmt.Println("write buf full!!!")
		this.cwrite <- &data
	}
	return nil
}

func (this *WebsocketTransport) Close(){
	this.Conn.Close()
	this.stopFlag = 1
}


type WebsocketServer struct {
	DefServer
}

func NewWebsocketServer(netAddr string, protocol  IProtocol, processor IProcessor) *WebsocketServer{
	WebsocketServer := &WebsocketServer{
		DefServer:DefServer{
			NetType:WebSocket,
			NetAddr:netAddr,
			transport: NewWebsocketTransport(netAddr,DefMsgTimeout,protocol,processor,nil),
		},
	}

	return WebsocketServer
}

type WebsocketClient struct {
	DefClient
}

func NewWebsocketClient(netAddr string, protocol  IProtocol, processor IProcessor) *WebsocketClient{
	WebsocketClient := &WebsocketClient{
		DefClient:DefClient{
			NetType:WebSocket,
			NetAddr:netAddr,
			transport: NewWebsocketTransport(netAddr,DefMsgTimeout,protocol,processor,nil),
		},
	}

	return WebsocketClient
}