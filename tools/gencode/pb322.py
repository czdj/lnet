#!/usr/bin/python
# -*- coding:utf-8 -*-

from gencomm import *

enumFlag = False

def match(string):
	global enumFlag
	if string == "":
		return "\n" + string

	if string.startswith("syntax"):
		return 'syntax = "proto2";'
	if string.startswith("import"):
		return "\n" + string
	if string.startswith("message"):
		return "\n" + string	
	if string.startswith("enum"):
		enumFlag = True
		return "\n" + string
	if string.startswith("{"):
		return "\n" + string
	if string.startswith("}"):
		enumFlag = False
		return "\n" + string
	if string.startswith("//"):
		return "\n" + string
	if string.startswith("repeated"):
		return "\n\t" +string
		
	if enumFlag:
		return "\n\t" + string
	else:
		return "\n\toptional " + string

def PB322():
	for parent, dirs, files in os.walk(pb2root):
		for name in files:
			os.remove(os.path.join(parent, name))
	for parent, dirs, files in os.walk(pb3root):
		for name in files:
			with open(os.path.join(parent, name)) as f:
				old=f.readlines()
			with open("%s/%s" %(pb2root, name),"w+") as fw:
				for line in old:
					str = match(line.strip())
					fw.write(str)
