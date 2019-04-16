#!/usr/bin/python
# -*- coding:utf-8 -*-

import os
import re
import shutil
import sys
import os.path
import platform
root = os.path.split(os.path.realpath(__file__))[0]
root = os.path.dirname(os.path.dirname(root))
pbVersion = "2"
if len(sys.argv) > 1:
	pbVersion = sys.argv[1]
shellroot = root + "/tools/bin/pb" + pbVersion 
csroot = root + "/tools/auto/cs"
cspbroot = root + "/tools/auto/cs/unity/Proto"
csnetroot = root + "/tools/auto/cs/unity/Net"
luaroot = root + "/tools/auto/lua"
luapbroot = root + "/tools/auto/lua/pb"
jsroot = root + "/tools/auto/js"
webroot = root + "/tools/auto/js/web"
cocosroot = root + "/tools/auto/js/cocos"
cocosNetRoot = cocosroot + "/Net"
cocosHandlesRoot = cocosroot + "/Handlers"
pythonroot = root + "/tools/auto/python"
genpyroot = root + "/tools/gencode"
deproot = root + "/deps/src"
confroot = root + "/conf"
srcroot = root + "/src"
wwwroot = root + "/www"
pb2root = root + "/src/idl/pb2"
pb3root = root + "/src/idl/pb3"
pbroot = root + "/src/idl/pb" +  pbVersion


newDirs = ['%s/proto/pb' %srcroot, pythonroot, jsroot, csroot, luaroot, luapbroot, webroot, cocosNetRoot,cocosHandlesRoot]
delDirs = ['%s/proto/pb' %srcroot, pythonroot, jsroot, csroot, luaroot, luapbroot, webroot, cocosNetRoot,cocosHandlesRoot]

def MoveFile(srcfile,dstfile):
    if not os.path.isfile(srcfile):
        print "%s not exist!"%(srcfile)
    else:
        fpath,fname=os.path.split(dstfile)    #分离文件名和路径
        if not os.path.exists(fpath):
            os.makedirs(fpath)                #创建路径
        shutil.move(srcfile,dstfile)          #移动文件

def CopyFile(srcfile,dstfile):
    if not os.path.isfile(srcfile):
        print "%s not exist!"%(srcfile)
    else:
        fpath,fname=os.path.split(dstfile)    #分离文件名和路径
        if not os.path.exists(fpath):
            os.makedirs(fpath)                #创建路径
        shutil.copyfile(srcfile,dstfile)      #复制文件

def CopyDir(sourceDir, targetDir):
	if not os.path.exists(targetDir):
		os.makedirs(targetDir)
	for f in os.listdir(sourceDir):
		sourceF = os.path.join(sourceDir, f)
		targetF = os.path.join(targetDir, f)
		if os.path.isfile(sourceF):
			open(targetF, "wb+").write(open(sourceF, "rb").read())
		if os.path.isdir(sourceF):
			CopyDir(sourceF, targetF)
		
def ListDirFile(rootdir):
	files = []
	list = os.listdir(rootdir) #列出文件夹下所有的目录与文件
	for i in range(0,len(list)):
		path = os.path.join(rootdir,list[i])
		if os.path.isfile(path):
			files.append(path)
	return files

def DelDir(dd):
	for parent, dirs, files in os.walk(dd):
		for name in files:
			os.remove(os.path.join(parent, name))
	for parent, dirs, files in os.walk(dd):
		for name in dirs:
			DelDir(os.path.join(parent, name))
			os.rmdir(os.path.join(parent, name))

def Init():
	for dd in delDirs:
		DelDir(dd)
	
	for nd in newDirs:
		if not os.path.exists(nd):
			os.makedirs(nd)