package main;

import "lnet"

func main() {
	processor := &lnet.DefProcessor{}
	protocol := &lnet.GobProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,"")

	server := lnet.NewWebsocketServer("127.0.0.1:9000",protocol,processor)
	server.Start()

	ch := make(chan int32)
	<- ch
}