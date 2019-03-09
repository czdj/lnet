package main;

import "lnet"

func main() {
	processor := &lnet.DefProcessor{}
	protocol := &lnet.MyProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,"")

	tcpServer := lnet.NewTcpServer("127.0.0.1:9000",protocol,processor)
	tcpServer.Start()
}