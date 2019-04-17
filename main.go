package main;

import (
	"lnet"
	"proto/pb"
)

func main() {
	lnet.Logger = lnet.InitLogger("log.log","")
	processor := &lnet.DefProcessor{}
	protocol := &lnet.PbProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,pb.GameItem{})

	server := lnet.NewTcpServer("127.0.0.1:9000",protocol,processor)
	server.Start()

	ch := make(chan int32)
	<- ch
}