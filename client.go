package main

import (
	"lnet"
	protocol2 "lnet/protocol"
	"proto/pb"
	"time"
)

func main1() {
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,pb.GameItem{})

	processor := &lnet.DefProcessor{}
	protocol := &protocol2.PbProtocol{}
	//client :=  lnet.NewTcpClient("127.0.0.1:9000", protocol,processor)
	client :=  lnet.NewWebsocketClient("ws://127.0.0.1:9000/ws", protocol,processor)
	client.Connect()

	msg := &pb.GameItem{Id:1,Type:2,Count:3}
	//msg1 := "bbbbbbbb"
	for {
		client.Send(msg)
		//transport.Send(12,msg1)
		time.Sleep(1 * time.Second)
	}

}