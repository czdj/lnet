package main;

import "lnet"

func main() {
	processor := &lnet.DefProcessor{}
	protocol := &lnet.MyProtocol{}
	tcpServer := lnet.NewTcpServer(protocol,processor)
	tcpServer.Start()
}