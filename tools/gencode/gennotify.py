#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

def GenNotify():
	with open("%s/redis_notify.proto" %pbroot,"r") as f:
		old=f.readlines()
	
	pbs = []
	find = False
	for line in old:
		if not find:
			sss = []
			for s in line.strip().split(" "):
				if s.strip() != "":
					sss.append(s)
			if len(sss) == 2 and sss[0] == "message" and sss[1] == "RedisNotify":
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
			
	
	with open("%s/frame/auto_redis_notify.go" %srcroot,"w") as f2:
		f2.write(
'''package frame
import (
	"antnet"
	"proto/pb"
)
''')
		
		for pb in pbs:
			if pbVersion == "2":
				xpb = pb[2][0].upper() + pb[2][1:]
				xpb2 = pb[1]
			else:
				xpb = pb[1][0].upper() + pb[1][1:]
				xpb2 = pb[0]
			f2.write(
'''
type %sEventFunc func(ppb *pb.%s)
type %sEvent struct {
 	l []*%sEventFunc
}
func (r *%sEvent) AddListener(fun %sEventFunc)  {
	r.l = append(r.l, &fun)
}
func (r *%sEvent) RemoveListener(fun %sEventFunc){
	
}
func (r *%sEvent) invoke(channel string, ppb *pb.%s){
	if len(r.l) == 0{
		antnet.LogInfo("redis notify not have handler channel:%%v msg:%s", channel)
	} else {
		for _, fun := range r.l {
			(*fun)(ppb)
		}
	}
}
''' %(xpb, xpb2, xpb, xpb, xpb, xpb, xpb, xpb, xpb, xpb2, xpb))
		f2.write(
'''
var RedisNotify = &struct {''')

		for pb in pbs:
			if pbVersion == "2":
				xpb = pb[2][0].upper() + pb[2][1:]
				f2.write('\n\t%s %sEvent' %(xpb, xpb))
			else:
				xpb = pb[1][0].upper() + pb[1][1:]
				f2.write('\n\t%s %sEvent' %(xpb, xpb))
		f2.write(
'''
}{}''')

		f2.write(
'''
func RedisNotifyHandlerFunc(channel string, ppb *pb.RedisNotify) {
''')
		for pb in pbs:
			if pbVersion == "2":
				xpb = pb[2][0].upper() + pb[2][1:]
			else:
				xpb = pb[1][0].upper() + pb[1][1:]
			f2.write(
'''
	if ppb.%s != nil {
		RedisNotify.%s.invoke(channel, ppb.%s)
		return
	}
''' %(xpb, xpb, xpb))
		f2.write(
'''
}
''')

		f2.write(
'''
func GetRedisNotifyString(ppb *pb.RedisNotify) string {
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