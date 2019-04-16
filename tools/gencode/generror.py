#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

def GenError():
	errors = []
	with open("%s/idl/csv/error.csv" %srcroot,"r") as f:
		old=f.readlines()[1:]
		
	with open("%s/frame/auto_error.go" %srcroot,"w") as f2:
		f2.write(
'''package frame

import (
	"antnet"
)

var (''')
		for data in old:
			if data.strip() == "":
				continue
			es = data.strip().split(",")
			f2.write('\n\t%s =  antnet.NewError("%s", %s)' %(es[0], es[1], es[2]))
			errors.append(es[0])
		f2.write(
'''
)
''')

	with open("%s/antnet/error.go" %deproot,"r") as f:
		old=f.readlines()
	copy = False
	with open("%s/frame/auto_error.go" %srcroot,"a+") as f2:
		for data in old:		
			if copy:				
				if data.strip().strip("\n").strip("\r") == ")":
					f2.write(data)
					break
				else:
					var = data[:data.find("=")].strip()
					if len(var) > 0 and data.find("NewError") != -1:
						commit = data[data.find('"') + 1:data.rfind('"')].strip()
						f2.write("\t%s = antnet.%s  //%s\n" %(var, var, commit))
						errors.append(var)
			else:
				if data.find("var (") != -1: 
					f2.write(data.strip("\n") + "  //自动生成 复制于antnet错误码\n")
					copy = True
					
	with open("%s/frame/auto_error.go" %srcroot,"a+") as f2:
		f2.write('''
var errMap = map[uint16]string{		
''')
		for error in errors:		
			f2.write("\t%s.Id:%s.Str,\n" %(error, error))
			
		f2.write("}" )	
			
			