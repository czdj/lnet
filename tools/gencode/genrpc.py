#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

def GenRpc():
	with open("%s/rpc.proto" %pbroot,"r") as f:
		old=f.readlines()
	cpbs = []
	spbs = []
	pbs = []
	spbsm = {}
	cpbsm = {}
	find = False
	for line in old:
		if not find:
			sss = []
			for s in line.strip().split(" "):
				if s.strip() != "":
					sss.append(s)
			if len(sss) == 2 and sss[0] == "message" and sss[1] == "RPC":
				find = True
		else:
			if line.find("int32") != -1:
				continue
			if line.find("=") != -1:
				ss = line.strip().split(" ")
			elif line.find("}") != -1:
				break
			else:
				continue
			sss = []
			for s in ss:
				if s.strip() == "":
					continue
				sss.append(s)
			if sss[0] == "//":
				continue
			if line.find("C2S") != -1 or line.find("S2C") != -1:
				if line.find("C2S") != -1:
					if pbVersion == "2":
						cpbs.append([sss[1], sss[2], sss[1].replace("C2S", "S2C")])
						cpbsm[sss[1]] = [sss[1], sss[2], sss[1].replace("C2S", "S2C")]
					else:
						cpbs.append([sss[0], sss[1], sss[0].replace("C2S", "S2C")])
						cpbsm[sss[0]] = [sss[0], sss[1], sss[0].replace("C2S", "S2C")]
				else:
					if pbVersion == "2":
						spbs.append([sss[1], sss[2], sss[1].replace("S2C", "C2S")])
						spbsm[sss[1]] = [sss[1], sss[2], sss[1].replace("S2C", "C2S")]
					else:
						spbs.append([sss[0], sss[1], sss[0].replace("S2C", "C2S")])
						spbsm[sss[0]] = [sss[0], sss[1], sss[0].replace("S2C", "C2S")]
			else:
				pbs.append(sss)
			
	with open("%s/frame/auto_rpc.go" %srcroot,"w") as f2:
		imp ='''"antnet"
	"proto/pb"'''
		if len(cpbs) > 0:
			imp = '''"antnet"
	"proto/pb"
	"time"'''			
		f2.write(
'''package frame
import (
	%s
	"sync"
)

type IMicroServerSession interface {
	GetMicroServerId(server string, c2s *pb.C2S) int32
	GetMicroServerTimeout() int
	GetGamerId() int32
}

var rpcIndex uint16 = 1
var rpcIndexSync sync.Mutex

func GetRpcIndex() (r uint16){
	rpcIndexSync.Lock()
	if rpcIndex == 0 {
	    rpcIndex++
	}
	r = rpcIndex
	rpcIndex++
	rpcIndexSync.Unlock()
	return r
}

type rpc struct {
	NetError func(msgque antnet.IMsgQue, msg *antnet.Message) bool''' %imp)
		for pb in pbs:
			if pbVersion == "2":
				f2.write('\n\t%s func(msgque antnet.IMsgQue, msg *antnet.Message, ppb *pb.%s) bool' %(pb[1], pb[1]))
			else:
				f2.write('\n\t%s func(msgque antnet.IMsgQue, msg *antnet.Message, ppb *pb.%s) bool' %(pb[0], pb[0]))
		for pb in spbs:
			f2.write('\n\t%s func(msgque antnet.IMsgQue, msg *antnet.Message, c2s *pb.%s, s2c *pb.%s) *antnet.Error' %(pb[0], pb[2][0].upper() + pb[2][1:], pb[0]))
		f2.write(
'''
}
''')
		for pb in cpbs:
			spb = spbsm[pb[0].replace("C2S", "S2C")][1]
			spb = spb[0].upper() + spb[1:]
			f2.write('''
func (r *rpc) %s(session IMicroServerSession, ppb *pb.%s) *pb.%s {
	rpcNetMap.RLock()
	msgque, ok := rpcNetMap.M[session.GetMicroServerId("%s", &pb.C2S{%s:ppb})]
	rpcNetMap.RUnlock()
	rpb := &pb.%s{Error:pb.Int32(int32(ErrNetUnreachable.Id))}
	if !ok || !msgque.Available() {
		return rpb
	}
	ppb.Id = pb.Int32(session.GetGamerId())
	xpb := &pb.RPC{Id: ppb.Id, %s:ppb}
	c := make(chan *antnet.Message, 1)
	defer close(c)
	msg := antnet.NewMsg(GAME_CMD_RPC, GAME_CMD_RPC_AUTO, GetRpcIndex(), 0, antnet.PbData(xpb))
	re := msgque.SendCallback(msg, c)
	if re {
		ms := session.GetMicroServerTimeout()
		if ms > 0 {
			select {
			case msg := <-c:
				if msg != nil {
					rpb := msg.C2S().(*pb.RPC)
					if rpb != nil && rpb.%s != nil {
						return rpb.%s
					} else {
						antnet.LogError("parse msg failed %%#v", rpb)
					}
				}
			case <-time.After(time.Duration(ms) * time.Millisecond):
				rpb.Error = pb.Int32(int32(ErrNetTimeout.Id))
			}
		} else {
			msg := <-c
			if msg != nil {
				rpb := msg.C2S().(*pb.RPC)
				if rpb != nil && rpb.%s != nil {
					return rpb.%s
				} else {
					antnet.LogError("parse msg failed %%#v", rpb)
				}
			}
		}
	}
	
	return rpb
}
''' %(pb[0][0].upper() + pb[0][1:], pb[0][0].upper() + pb[0][1:], pb[2][0].upper() + pb[2][1:], pb[0][0].upper() + pb[0][1:], pb[0][0].upper() + pb[0][1:], pb[2][0].upper() + pb[2][1:], pb[1][0].upper() + pb[1][1:], spb, spb, spb, spb))

		f2.write('''
var RPC = &rpc{
	NetError: func(msgque antnet.IMsgQue, msg *antnet.Message) bool{
		antnet.LogError("rpc recv net error msgque:%v cmd:%v act:%v err:%v", msgque.Id(), msg.Head.Cmd, msg.Head.Act, msg.Head.Error)
		return true
	},
}
''')
		f2.write(
'''
func RPCHandlerFunc(msgque antnet.IMsgQue, msg *antnet.Message) bool {
	if msg.Head.Error > 0 {
		return RPC.NetError(msgque, msg)
	}
	ppb := msg.C2S().(*pb.RPC)
	if ppb == nil {
		return true
	}
	if ppb.ServerHello != nil{
		handleServerHello(msgque, msg, ppb.ServerHello)
	}
''')
		for pb in pbs:
			if pbVersion == "2":
				xpb = pb[2][0].upper() + pb[2][1:]
			else:
				xpb = pb[1][0].upper() + pb[1][1:]
			f2.write(
'''
	if ppb.%s != nil {
		if RPC.%s != nil {
			return RPC.%s(msgque, msg, ppb.%s)
		}
	}
''' %(xpb, xpb, xpb, xpb))

		for pb in cpbs:
			xpb = pb[1][0].upper() + pb[1][1:]
			spb = spbsm[pb[0].replace("C2S", "S2C")][1]
			spb = spb[0].upper() + spb[1:]
			f2.write(
'''
	if ppb.%s != nil {
		s2c := &pb.%s{Error:pb.Int32(0)}
		err := RPC.%s(msgque, msg, ppb.%s, s2c)
		if err != nil {
			s2c.Error = pb.Int32(int32(antnet.GetErrId(err)))
		}
		msgque.Send(antnet.NewDataMsg(antnet.PbData(&pb.RPC{%s:s2c})).CopyTag(msg))
		return true
	}
''' %(xpb, pb[0].replace("C2S", "S2C"), pb[2], xpb, spb))

		f2.write(
'''
	return true
}
''')
		f2.write(
'''
func GetRPCPBString(ppb *pb.RPC) string {
''')
		for pb in cpbs:
			xpb = pb[1][0].upper() + pb[1][1:]
			f2.write(
'''
	if ppb.%s != nil {
		return "%s"
	}
''' %(xpb, xpb))
		for pb in pbs:
			xpb = pb[1][0].upper() + pb[1][1:]
			f2.write(
'''
	if ppb.%s != nil {
		return "%s"
	}
''' %(xpb, xpb))
		f2.write(
'''
	return ""
}
''')
	return cpbsm
