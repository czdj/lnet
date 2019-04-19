package main

import (
	"lnet"
	"lnet/lmsghandle"
	"lnet/lprotocol"
	"lnet/lserver"
	"proto/pb"
)

func main() {
	lnet.Logger = lnet.InitLogger("./logs/log.log", "")
	protocol := &lprotocol.PbProtocol{}
	msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	//protocol := &lprotocol.GobProtocol{}
	msgHandle.RegisterMsg(12, pb.GameItem{})

	//server := lserver.NewTcpServer("127.0.0.1:9000", protocol, msgHandle)
	server := lserver.NewWebsocketServer("127.0.0.1:9000", protocol, msgHandle)

	server.Start()

	ch := make(chan int32)
	<-ch
}
