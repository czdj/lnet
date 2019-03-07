package main;

import "lnet"

func main() {
	processor := &lnet.DefProcessor{}
	protocol := &lnet.MyProtocol{}
	lnet.MsgTypeInfo.Register(11,lnet.Message{})

	tcpServer := lnet.NewTcpServer("127.0.0.1:9000",protocol,processor)
	tcpServer.Start()
}