#!/usr/bin/python
# -*- coding:utf-8 -*-
import os,platform,sys,shutil
import gencode

root = os.path.split(os.path.realpath(__file__))[0]
root = os.path.dirname(root)

if len(sys.argv) > 1:
	if sys.argv[1] == "2" or sys.argv[1] == "3":
		pass
	else:
		print "args failed"
		sys.exit()

if len(sys.argv) > 1:
	gencode.Init()
	gencode.pb322.PB322()
	gencode.genproto.GenProto()
	# gencode.geninner.GenInner()
	# gencode.generror.GenError()
	# gencode.gencmdact.GenCmdAct()
	# gencode.gennotify.GenNotify()
	# gencode.genlua.GenLua()
	# #gencode.genconfig.GenConfig()
	# gencode.genjs.GenJS()
	# gencode.gencs.GenCS()


if 'Windows' in platform.system():
	os.system('cmd /C "set GOPATH=%s;%s&& cd %s && go build"' %(root, root + "/deps", root))
else:
	os.system('export GOPATH=%s:%s && cd %s && go build' %(root, root + "/deps", root))
