package main;

import (
	"lnet"
	"lnet/process"
	"lnet/protocol"
	"lnet/server"
	"proto/pb"
)

func main() {
	lnet.Logger = lnet.InitLogger("./logs/log.log","")
	processor := &process.BaseProcessor{}
	protocol := &protocol.PbProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,pb.GameItem{})

	server := server.NewTcpServer("127.0.0.1:9000",protocol,processor)
	server.Start()

	ch := make(chan int32)
	<- ch
}