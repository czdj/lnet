#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

def GenCmdAct():
	with open("%s/idl/csv/cmdact.csv" %srcroot,"r") as f:
		old=f.readlines()[1:]
		
	with open("%s/frame/auto_define.go" %srcroot,"w") as f2:
		f2.write(
'''package frame
import (
	"proto/pb"
)
''')
		cmds = {}
		msgs = []
		for data in old:
			if data.strip() == "":
				continue
			es = data.strip().split(",")
			cmd = es[1]
			if int(cmd) == 0:
				msgs.append(es)
				continue
			if not cmds.has_key(cmd):
				cmds[cmd] = {'cmd':es[0]}
				cmds[cmd]['acts'] = []
			cmds[cmd]['acts'].append(es)
		
		f2.write("\n\nconst (")
		for i in range(1, 256):
			if cmds.has_key(str(i)):
				v = cmds[str(i)]
				f2.write("\n\t%s = %d" %(v['cmd'], i))
		f2.write("\n)")
		for i in range(1, 256):
			if cmds.has_key(str(i)):
				v = cmds[str(i)]
				f2.write("\n\nconst (")
				for es in v['acts']:
					f2.write("\n\t%s = %s //%s" %(es[2], es[3], es[6]))
				f2.write("\n)")
		
		funcs = []
		f2.write("\n")
		for i in range(1, 256):
			if cmds.has_key(str(i)):
				v = cmds[str(i)]
				f2.write("\n\nfunc initParser_%s (){" %v['cmd'])
				funcs.append("initParser_%s()" %v['cmd'])
				for es in v['acts']:
					if es[4].strip() != "" and es[5].strip() != "":
						f2.write("\n\tPbParser.Register(%s, %s, &pb.%s{}, &pb.%s{})" %(es[0], es[2], es[4], es[5]))
					elif es[4].strip() != "" and es[5].strip() == "":
						f2.write("\n\tPbParser.Register(%s, %s, &pb.%s{}, nil)" %(es[0], es[2], es[4]))
					elif es[4].strip() == "" and es[5].strip() != "":
						f2.write("\n\tPbParser.Register(%s, %s, nil, &pb.%s{})" %(es[0], es[2], es[5]))
				f2.write("\n}")
		f2.write("\n")
		f2.write("func initParser() {")
		if pbVersion == "3":
			f2.write('\n\tServerInfo.PbVer = pb.String("3.6")\n')
		else:
			f2.write('\n\tServerInfo.PbVer = pb.String("2.5")\n')
		for func in funcs:
			f2.write("\n\t%s" %func)
		for es in msgs:
			if es[4].strip() != "" and es[5].strip() != "":
				f2.write("\n\tPbParser.RegisterMsg(&pb.%s{}, &pb.%s{})" %(es[4], es[5]))
			elif es[4].strip() != "" and es[5].strip() == "":
				f2.write("\n\tPbParser.RegisterMsg(&pb.%s{}, nil)" %(es[4]))
			elif es[4].strip() == "" and es[5].strip() != "":
				f2.write("\n\tPbParser.RegisterMsg(nil, &pb.%s{})" %(es[5]))
		f2.write("\n}")
		
