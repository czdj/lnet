package iface

import "proto/pb"

type IRpc interface {
	Send(msg interface{}) error                             //不等待结果，不阻塞
	SendWaitResult(msg interface{}) (*pb.RpcRspInfo, error) //等待结果，阻塞
}
