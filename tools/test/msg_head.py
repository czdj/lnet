#!/usr/bin/env python
# -*- coding:utf-8 -*-
'''
Len   uint32 //数据长度
Error uint16 //错误码
Cmd   uint8  //命令
Act   uint8  //动作
Index uint16 //序号
Flags uint16 //标记
'''
import struct

def Gen(l, cmd, act, err = 0, index = 0, flags = 0):
	return struct.pack("IHBBHH", l, err, cmd, act, index, flags)
	
def Unpack(s):
	return struct.unpack("IHBBHH", s)