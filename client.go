package main

import (
	"lnet"
	"lnet/lclient"
	"lnet/lprocess"
	"lnet/lprotocol"
	"proto/pb"
	"time"
)

func main1() {
	lnet.Logger = lnet.InitLogger("./logs/log.log","")

	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})
	lnet.MsgTypeInfo.Register(12,pb.GameItem{})

	processor := &lprocess.BaseProcessor{}
	protocol := &lprotocol.PbProtocol{}
	client :=  lclient.NewTcpClient("127.0.0.1:9000", protocol,processor)
	//client :=  lclient.NewWebsocketClient("ws://127.0.0.1:9000/ws", protocol,processor)
	client.Connect()

	msg := &pb.GameItem{Id:1,Type:2,Count:3}
	//msg1 := "bbbbbbbb"
	for {
		client.Send(msg)
		//transport.Send(12,msg1)
		time.Sleep(1 * time.Second)
	}

}