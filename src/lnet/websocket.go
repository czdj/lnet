package lnet

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"time"
	"unsafe"
)


type WebsocketTransport struct {
	DefTransport
	Conn *websocket.Conn
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
	}
}

func (this *WebsocketTransport) websocketConnHandler(conn *websocket.Conn) {
	WebsocketTransport := NewWebsocketTransport(this.NetAddr,DefMsgTimeout, this.protocol,this.processor,conn)
	this.onNewConnect(WebsocketTransport)
}

func (this *WebsocketTransport) Listen() error {

	http.Handle("/ws", websocket.Handler(this.websocketConnHandler));
	err := http.ListenAndServe(this.NetAddr, nil);
	if err != nil{
		fmt.Println("Websocket Listen err:%v!",err)
		return err
	}

	fmt.Println("WebsocketServer is listening addr:%v!",this.NetAddr)

	return nil
}

func (this *WebsocketTransport) Connect() error{
	conn, err := websocket.Dial(this.NetAddr, "", "");
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

		_,err := this.Conn.Write(*data)
		if err != nil{
			fmt.Println("Write Err:%v",err)
			break
		}
		data = nil
		this.lastTick = time.Now().Unix()
	}
}

func (this *WebsocketTransport)Send(tag uint16, msg interface{})error{
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

	data := this.protocol.Encode(tag, msg)
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