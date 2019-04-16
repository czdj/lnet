#!/usr/bin/python
# -*- coding:utf-8 -*-

#go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
from gencomm import *
import genrpc
def GenC2S():
	cpbsm = genrpc.GenRpc()
	with open("%s/client.proto" %pbroot,"r") as f:
		old=f.readlines()
		
	c2ss = []
	for line in old:
		sss = []
		for s in line.strip().split(" "):
			if s.strip() != "":
				sss.append(s.strip())
		if len(sss) == 2 and sss[0] == "message" and sss[1].find("C2S")!= -1:
			c2ss.append(sss[1])
			
	with open("%s/auto_c2s.proto" %pbroot,"w") as f2:
		f2.write('syntax = "proto%s";\n\nimport "client.proto";\n\n' %pbVersion)
		f2.write("message C2S\n{")
		if pbVersion == "2":
			f2.write("\n\toptional string key = 1;")
		else:
			f2.write("\n\tstring key = 1;")
		for i, c2s in enumerate(c2ss):
			if pbVersion == "2":
				f2.write("\n\toptional %s %s = %s;" %(c2s, c2s[0].lower() + c2s[1:], i + 2))
			else:
				f2.write("\n\t%s %s = %s;" %(c2s, c2s[0].lower() + c2s[1:], i + 2))
		f2.write("\n}")
		
	with open("%s/frame/auto_c2s.go" %srcroot,"w") as f2:
		f2.write(
'''package frame
import (
	"antnet"
	"proto/pb"
	"sync/atomic"
)

var C2S = &struct {
	NetError func(msgque antnet.IMsgQue, msg *antnet.Message) bool''')
		for i, c2s in enumerate(c2ss):
			if cpbsm.has_key(c2s):
				f2.write('\n\tBefRPC%s func(msgque antnet.IMsgQue, msg *antnet.Message, c2s *pb.%s) bool' %(c2s, c2s))
				f2.write('\n\tAftRPC%s func(msgque antnet.IMsgQue, msg *antnet.Message, c2s *pb.%s, s2c *pb.%s) bool' %(c2s, c2s, c2s.replace("C2S", "S2C")))
			else:
				f2.write('\n\t%s func(msgque antnet.IMsgQue, msg *antnet.Message, c2s *pb.%s, s2c *pb.%s) *antnet.Error' %(c2s, c2s, c2s.replace("C2S", "S2C")))

		f2.write(
'''
}{
	NetError: func(msgque antnet.IMsgQue, msg *antnet.Message) bool{
		antnet.LogError("inner recv net error msgque:%v cmd:%v act:%v err:%v", msgque.Id(), msg.Head.Cmd, msg.Head.Act, msg.Head.Error)
		return true
	},
}
var C2SProf = map[string]*Prof {''')

		for i, c2s in enumerate(c2ss):
			f2.write('\n\t"%s":&Prof{},' %(c2s))
		f2.write(
'''
}

func C2SHandlerFunc(msgque antnet.IMsgQue, msg *antnet.Message) bool {
	if msg.Head != nil {
		if msg.Head != nil && msg.Head.Error > 0 {
			return C2S.NetError(msgque, msg)
		}
	}
	ppb := msg.C2S().(*pb.C2S)
	if ppb == nil {
		return true
	}
''')
		hasms = False
		for c2s in c2ss:
			if cpbsm.has_key(c2s):
				hasms = True
				break
		
		if hasms:
			f2.write('''
	user := msgque.GetUser()
	if user == nil && ppb.GamerLoginC2S == nil {
	    return true
	}
	ms, msok := user.(IMicroServerSession)
''')
		for c2s in c2ss:
			if cpbsm.has_key(c2s):
				f2.write(
	'''
	if ppb.%s != nil {
		var bt int64 = 0
		if Config.Global.ProfSwitch {
			bt = antnet.UnixMs()
		}
		if C2S.BefRPC%s != nil {
			if !C2S.BefRPC%s(msgque, msg, ppb.%s) {
				return true
			}
		}
		s2c := &pb.%s{Error:pb.Int32(int32(ErrMincroServerNotFound.Id))}
		if msok {
			s2c = RPC.%s(ms, ppb.%s)
		}
		if C2S.AftRPC%s != nil {
			if !C2S.AftRPC%s(msgque, msg, ppb.%s, s2c) {
				return true
			}
		}
		if Config.Global.ProfSwitch {
			rt := antnet.UnixMs() - bt
			prof := C2SProf["%s"]
			atomic.AddInt64(&prof.CallCount, 1)
			atomic.AddInt64(&prof.AllTime, rt)
			if rt > prof.MaxTime {
				prof.MaxTime = rt
			}
			prof.AvgTime = int32(prof.AllTime / prof.CallCount)
		}
		msgque.Send(antnet.NewDataMsg(antnet.PbData(&pb.S2C{%s:s2c, /*Key:pb.String("%s"),*/ Error:s2c.Error})).CopyTag(msg))
		return true
	}
	''' %(c2s, c2s, c2s, c2s, c2s.replace("C2S", "S2C"), c2s, c2s, c2s, c2s, c2s, c2s, c2s.replace("C2S", "S2C"), c2s.replace("C2S", "S2C")[0].lower() + c2s.replace("C2S", "S2C")[1:]))
			else:
				f2.write(
	'''
	if ppb.%s != nil {
		if C2S.%s != nil {
			var bt int64 = 0
			if Config.Global.ProfSwitch {
				bt = antnet.UnixMs()
			}
			s2c := &pb.%s{Error:pb.Int32(0)}
			err := C2S.%s(msgque, msg, ppb.%s, s2c)
			if Config.Global.ProfSwitch {
				rt := antnet.UnixMs() - bt
				prof := C2SProf["%s"]
				atomic.AddInt64(&prof.CallCount, 1)
				atomic.AddInt64(&prof.AllTime, rt)
				if rt > prof.MaxTime {
					prof.MaxTime = rt
				}
				prof.AvgTime = int32(prof.AllTime / prof.CallCount)
			}
			if err != nil {
				s2c.Error = pb.Int32(int32(antnet.GetErrId(err)))
			}
			msgque.Send(antnet.NewDataMsg(antnet.PbData(&pb.S2C{%s:s2c, /*Key:pb.String("%s"),*/ Error:s2c.Error})).CopyTag(msg))
			return true
		}
	}
	''' %(c2s, c2s, c2s.replace("C2S", "S2C"), c2s, c2s, c2s, c2s.replace("C2S", "S2C"), c2s.replace("C2S", "S2C")[0].lower() + c2s.replace("C2S", "S2C")[1:]))
		f2.write(
'''
	antnet.LogInfo("c2s msg not have handler msg:%v", GetC2SPBString(msg, 0))
	return true
}
''')
		f2.write(
'''
func GetC2SPBString(msg *antnet.Message, gid int32) string {
	if msg.Head != nil {
		if msg.Head.Error > 0 {
			return ""
		}
	}
	ppb := msg.C2S().(*pb.C2S)
	if ppb == nil {
		return ""
	}
''')
		for c2s in c2ss:
			f2.write(
'''
	if ppb.%s != nil {
		if gid > 0 {
			ppb.%s.Id = pb.Int32(gid)
		}
		return "%s"
	}
''' %(c2s, c2s, c2s))
		f2.write(
'''
	return ""
}
''')

def GenS2C():
	with open("%s/client.proto" %pbroot,"r") as f:
		old=f.readlines()
		
	c2ss = []
	for line in old:
		sss = []
		for s in line.strip().split(" "):
			if s.strip() != "":
				sss.append(s.strip())
		if len(sss) == 2 and sss[0] == "message" and sss[1].find("S2C")!= -1:
			c2ss.append(sss[1])
			
	with open("%s/auto_s2c.proto" %pbroot,"w") as f2:
		f2.write('syntax = "proto%s";\n\nimport "client.proto";\n\n' %pbVersion)
		f2.write("message S2C\n{")
		if pbVersion == "2":
			f2.write("\n\toptional int32 error = 1;")
			f2.write("\n\toptional string key = 2;")
		else:
			f2.write("\n\tint32 error = 1;")
			f2.write("\n\tstring key = 2;")
		for i, c2s in enumerate(c2ss):
			if pbVersion == "2":
				f2.write("\n\toptional %s %s = %s;" %(c2s, c2s[0].lower() + c2s[1:], i + 3))
			else:
				f2.write("\n\t%s %s = %s;" %(c2s, c2s[0].lower() + c2s[1:], i + 3))
		f2.write("\n}")
	with open("%s/frame/auto_c2s.go" %srcroot,"a+") as f2:
		f2.write(
'''
func GetS2CPBString(ppb *pb.S2C) string {
''')
		for c2s in c2ss:
			f2.write(
'''
	if ppb.%s != nil {
		return "%s"
	}
''' %(c2s, c2s))
		f2.write(
'''
	return ""
}
''')


def GenProto():
	if 'Windows' in platform.system():
		for fname in ListDirFile('%s' %(pbroot)):
			os.system('set PATH=%s;%%PATH%% && protoc.exe --go_out=%s/proto/pb -I=%s %s' 
			%(shellroot, srcroot, pbroot, fname))

		for fname in ListDirFile('%s/proto/pb' %(srcroot)):
			with open(fname, "r") as f:
				old = f.readlines()
			with open(fname, "w") as f:
				for lin in old:
					if lin == "package %s\n" %os.path.basename(fname).split(".")[0]:
						old[old.index(lin)] = "package pb\n"
				f.writelines(old)
	else:
		os.system('export PATH=%s:$PATH && export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH && chmod -R a+x %s && protoc --go_out=import_path=pb:%s/proto/pb -I=%s %s/*.proto' 
		%(shellroot, shellroot, shellroot, srcroot, pbroot, pbroot))
		os.system('export PATH=%s:$PATH && export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH && chmod -R a+x %s && protoc --python_out=%s -I=%s %s/*.proto' 
		%(shellroot, shellroot, shellroot, pythonroot, pbroot, pbroot))
		os.system('export PATH=%s:$PATH && export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH && chmod -R a+x %s && protoc --js_out=import_style=commonjs,binary:%s -I=%s %s/*.proto' 
		%(shellroot, shellroot, shellroot, webroot, pbroot, pbroot))
		os.system('export PATH=%s:$PATH && export LD_LIBRARY_PATH=%s:$LD_LIBRARY_PATH && chmod -R a+x %s && protoc --csharp_out=%s -I=%s %s/*.proto'
		%(shellroot, shellroot, shellroot, cspbroot, pbroot, pbroot))

