package main

import (
	"flag"
	"lnet"
	"lnet/lclient"
	"lnet/lmsghandle"
	"lnet/lprotocol"
	"lnet/lserver"
	"proto/pb"
	"time"
)

var (
	runSelect = flag.String("s", "s", "run select")
)

func serverStart() {
	protocol := &lprotocol.PbProtocol{}
	msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	//protocol := &lprotocol.GobProtocol{}
	msgHandle.RegisterMsg(12, pb.GameItem{})

	//server := lserver.NewTcpServer("127.0.0.1:9000", msgHandle)
	server := lserver.NewWebsocketServer("127.0.0.1:9000", msgHandle)

	server.Start()
}

func clientStart() {
	protocol := &lprotocol.PbProtocol{}
	//protocol := &lprotocol.GobProtocol{}

	msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	msgHandle.RegisterMsg(12, pb.GameItem{})

	//client := lclient.NewTcpClient("127.0.0.1:9000", msgHandle)
	client := lclient.NewWebsocketClient("ws://127.0.0.1:9000/ws", msgHandle)
	client.Connect()

	msg := &pb.GameItem{Id: 1, Type: 2, Count: 3}
	//msg1 := "bbbbbbbb"
	for {
		client.Send(msg)
		//transport.Send(12,msg1)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	flag.Parse()
	lnet.Logger = lnet.InitLogger("./logs/log.log", "")

	if *runSelect == "s" {
		serverStart()
	} else {
		clientStart()
	}

	ch := make(chan int32)
	<-ch
}
