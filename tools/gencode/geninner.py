#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

def GenInner():
	with open("%s/inner.proto" %pbroot,"r") as f:
		old=f.readlines()
	
	pbs = []
	find = False
	for line in old:
		if not find:
			sss = []
			for s in line.strip().split(" "):
				if s.strip() != "":
					sss.append(s)
			if len(sss) == 2 and sss[0] == "message" and sss[1] == "Inner":
				find = True
		else:
			if line.find("=") != -1:
				ss = line.strip().split(" ")
			elif line.find("}") != -1:
				break
			else:
				continue
			sss = []
			for s in ss:
				if s.strip() != "":
					sss.append(s)
			pbs.append(ss)
			
	
	with open("%s/frame/auto_inner.go" %srcroot,"w") as f2:
		f2.write(
'''package frame
import (
	"antnet"
	"proto/pb"
)

var Inner = &struct {''')
		for pb in pbs:
			if pbVersion == "2":
				xpb = pb[2][0].upper() + pb[2][1:]
				f2.write('\n\t%s func(msgque antnet.IMsgQue, ppb *pb.%s) interface{}' %(xpb, pb[1]))
			else:
				xpb = pb[1][0].upper() + pb[1][1:]
				f2.write('\n\t%s func(msgque antnet.IMsgQue, ppb *pb.%s) interface{}' %(xpb, pb[0]))
		f2.write(
'''
}{}
''')
		f2.write(
'''
func InnerHandlerFunc(msgque antnet.IMsgQue, ppb *pb.Inner) interface{} {
	if ppb == nil {
		return nil
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
		if Inner.%s != nil {
			return Inner.%s(msgque, ppb.%s)
		}
	}
''' %(xpb, xpb, xpb, xpb))
	
		f2.write(
'''
	antnet.LogInfo("inner msg not have handler msg:%v", GetInnerPBString(ppb))
	return true
}
''')

		f2.write(
'''
func GetInnerPBString(ppb *pb.Inner) string {
	if ppb == nil {
		return ""
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
		return "%s"
	}
''' %(xpb, xpb))
		f2.write(
'''
	return ""
}
''')
