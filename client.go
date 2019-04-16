package main

import (
	"lnet"
	"time"
)

func main1() {
	lnet.MsgTypeInfo.Register(11,lnet.MessageTest{})

	processor := &lnet.DefProcessor{}
	protocol := &lnet.GobProtocol{}
	//client :=  lnet.NewTcpClient("127.0.0.1:9000", protocol,processor)
	client :=  lnet.NewWebsocketClient("127.0.0.1:9000", protocol,processor)
	client.Connect()

	msg := &lnet.MessageTest{Data:"zzzzz"}
	//msg1 := "bbbbbbbb"
	for {
		client.Send(11,msg)
		//transport.Send(12,msg1)
		time.Sleep(1 * time.Second)
	}

}