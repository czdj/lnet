package main

import (
	"flag"
	"go.uber.org/zap"
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
	//msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	msgHandle := lmsghandle.NewRpcMsgHandle(protocol)
	//protocol := &lprotocol.GobProtocol{}
	msgHandle.RegisterMsg(12, pb.GameItem{})
	msgHandle.RegisterMsg(13, pb.RpcReqInfo{})
	msgHandle.RegisterMsg(14, pb.RpcRspInfo{})

	//server := lserver.NewTcpServer("127.0.0.1:9000", msgHandle)
	//server := lserver.NewWebsocketServer("127.0.0.1:9000", msgHandle)
	server := lserver.NewRpcServer("127.0.0.1:9000", msgHandle)

	server.Start()
}

func clientStart() {
	protocol := &lprotocol.PbProtocol{}
	//protocol := &lprotocol.GobProtocol{}

	//msgHandle := lmsghandle.NewBaseMsgHandle(protocol)
	msgHandle := lmsghandle.NewRpcMsgHandle(protocol)

	msgHandle.RegisterMsg(12, pb.GameItem{})
	msgHandle.RegisterMsg(13, pb.RpcReqInfo{})
	msgHandle.RegisterMsg(14, pb.RpcRspInfo{})

	//client := lclient.NewTcpClient("127.0.0.1:9000", msgHandle)
	//client := lclient.NewWebsocketClient("ws://127.0.0.1:9000/ws", msgHandle)
	client := lclient.NewRpcClient("127.0.0.1:9000", msgHandle)
	client.Connect()

	//msg := &pb.GameItem{Id: 1, Type: 2, Count: 3}
	msg := &pb.RpcReqInfo{Test: 999}
	//msg1 := "bbbbbbbb"
	for {
		//client.Send(msg)
		r, _ := client.SendWaitResult(msg)
		lnet.Logger.Info("SendWaitResult ", zap.Any("Result", r))
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
