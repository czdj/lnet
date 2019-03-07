package main

import (
	"lnet"
	"time"
)

func main() {
	lnet.MsgTypeInfo.Register(11,lnet.Message{})

	processor := &lnet.DefProcessor{}
	protocol := &lnet.MyProtocol{}
	transport := lnet.NewTcpTransport("127.0.0.1:9000",protocol,processor,nil)
	transport.Connect()
	msg := &lnet.Message{Data:"zzzzz"}
	for {
		transport.Write(11,msg)
		time.Sleep(5000)
	}

}