package main;

import (
	"lnet"
	"lnet/lprocess"
	"lnet/lprotocol"
	"lnet/lserver"
	"proto/pb"
)

func main() {
	lnet.Logger = lnet.InitLogger("./logs/log.log","")
	processor := &lprocess.BaseProcessor{}
	protocol := &lprotocol.PbProtocol{}
	//protocol := &lprotocol.GobProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,pb.GameItem{})

	//server := lserver.NewTcpServer("127.0.0.1:9000",protocol,processor)
	server := lserver.NewWebsocketServer("127.0.0.1:9000",protocol,processor)

	server.Start()

	ch := make(chan int32)
	<- ch
}