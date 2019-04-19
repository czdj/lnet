package main

import (
	"lnet"
	"lnet/lclient"
	"lnet/lmsghandle"
	"lnet/lprotocol"
	"proto/pb"
	"time"
)

func main1() {
	lnet.Logger = lnet.InitLogger("./logs/log.log", "")

	protocol := &lprotocol.PbProtocol{}
	//protocol := &lprotocol.GobProtocol{}

	msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	msgHandle.RegisterMsg(12, pb.GameItem{})

	//client := lclient.NewTcpClient("127.0.0.1:9000", protocol, msgHandle)
	client :=  lclient.NewWebsocketClient("ws://127.0.0.1:9000/ws", protocol,msgHandle)
	client.Connect()

	msg := &pb.GameItem{Id: 1, Type: 2, Count: 3}
	//msg1 := "bbbbbbbb"
	for {
		client.Send(msg)
		//transport.Send(12,msg1)
		time.Sleep(1 * time.Second)
	}

}
