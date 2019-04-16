#!/usr/bin/python
# -*- coding:utf-8 -*-

#go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

from gencomm import *

def GenDefine():
	with open("%s/frame/auto_define.go" %srcroot,"r") as f:
		old=f.readlines()
		write = False
		cmds = []
		cmds_flg = False
		with open("%s/define.lua" %luaroot,"w") as fw:
			for i, v in enumerate(old):
				if v.find("const (") != -1:
					cmds_flg = True
				if cmds_flg:
					if v.find("=") != -1:
						cmds.append(v[:v.rindex("=")].strip())
					elif v.find(")") != -1:
						break
			for i, v in enumerate(old):
				if v.find("func ") != -1:
					break
				if v.find("const (") != -1:
					c = old[i+1]
					if write:
						fw.write("}")
					if write == False:
						fw.write("\n\n%s = \n{\n" %"GAME_CMD")
					else:
						for cmd in cmds:
							if c.find(cmd) != -1:
								fw.write("\n\n%s = \n{\n" %cmd)
								break
					write = True
				elif v.find("_") != -1:
					if v.find("//") != -1:
						fw.write("\t%s\n" %v.replace("//", ",  --").strip())
					else:
						fw.write("\t%s\n" %v.replace("\n", ",\n").strip())
					
			fw.write("}\n\n")

def GenError():
	with open("%s/frame/auto_error.go" %srcroot,"r") as f:
		old=f.readlines()
		with open("%s/error.lua" %luaroot,"w") as fw:
			fw.write("ERROR = \n{\n")
			for i, v in enumerate(old):
				if v.find("//自动生成 复制于antnet错误码") != -1:
					break
				if v.find("=") != -1 and v.find("(") != -1:
					infos = v.split(",")
					eid = infos[1].replace(")", "")
					eid = eid.replace("\n", "")
					eid = eid.replace(" ","")
					fw.write("\t[" + eid + "] = ")
					infos = v.split("\"")
					fw.write("\"" + infos[1] + "\",\n" )
			fw.write("}")
			
def GenC2S():
	gcommitStr = "---@class C2S\n"		
	with open("%s/client.proto" %pbroot,"r") as f:
		old=f.readlines()
		with open("%s/c2s.lua" %luaroot,"w") as fw:
			fw.write("local C2S = {\n")
			funFind = False
			args = 0
			funStr = ""
			commitStr = ""
			funBody = ""
			fname = ""
			for i, v in enumerate(old):
				v = v.strip()
				if not funFind:
					if v.find("message") != -1 and v.find("C2S") != -1:
						infos = v.split(" ")
						fname = infos[1].strip().strip("{")
						gcommitStr += "---@field %s fun(" %(fname)
						funStr = "\t" + fname + " = function("
						commitStr = ""
						funBody = "\tlocal pb = GL.PB.Client.%s()" %(fname)
						funFind = True
						args = 0
						continue
				else:
					sub = 0
					if pbVersion == "3":
						sub = 1
					xinfos = v.split(" ")
					infos = []
					for s in xinfos:
						if s != "":
							infos.append(s)
					if len(infos) > 2 - sub:
						arg = infos[2 - sub].strip()
						if arg == "":
							print(v, infos)
						if args == 0:
							funStr += "%s" %(arg)
							gcommitStr += "%s:%s" %(arg, infos[1 - sub].strip())
						else:
							funStr += ", %s" %(arg)
							gcommitStr += ",%s:%s" %(arg, infos[1 - sub].strip())
						funBody += "\n\t\tif %s then pb.%s = %s end" %(arg, arg, arg)
						args = args + 1
					if v.find("}") != -1:
						gcommitStr += ")\n"
						funFind = False
						funStr += """)
	%s
		local data = pb:SerializeToString()
		GL.Net.Tcp:SendMsg(AntNet.SendData.New(GAME_CMD_C2S, GAME_CMD_C2S_AUTO, data))
	end,

""" %(funBody)
						fw.write(commitStr)
						fw.write(funStr)	
			fw.write("}\n")
			fw.write("return C2S")

	with open("%s/c2s.lua" %luaroot,"r") as fw:
		gcommitStr += fw.read()
	with open("%s/c2s.lua" %luaroot,"w") as fw:
		fw.write(gcommitStr)
		
def GenS2C():
	gcommitStr = "\n---@class S2C\n"		
	with open("%s/client.proto" %pbroot,"r") as f:
		old=f.readlines()
		
		with open("%s/s2c.lua" %luaroot,"w") as fw:
			fw.write("local S2C = {}\n")
			funFind = False
			args = 0
			funStr = ""
			commitStr = ""
			for i, v in enumerate(old):
				v = v.strip()
				if not funFind:
					if v.find("message") != -1 and (v.find("S2C") != -1 or v.find("Notify") != -1):				
						infos = v.split(" ")
						fname = infos[1].strip().strip("{")
						gcommitStr += "---@field %s %s\n" %(fname, fname)
						funStr = "S2C." + fname + " = {}"
						tName =  "S2C." + fname
						funStr += """
---@param f fun(t:%s)
function %s.Add (f)
	%s[f] = f
end
---@param f fun(t:%s)
function %s.Remove (f)
	if %s[k] then %s[k] = nil end
end				
""" %(fname, tName, tName, fname, tName, tName, tName)
							
						funStr += "\n\n"
						commitStr = "---@class %s\n" %(fname)
						funFind = True
						args = 0
						continue			
				else:
					sub = 0
					if pbVersion == "3":
						sub = 1
					xinfos = v.split(" ")
					infos = []
					for s in xinfos:
						if s != "":
							infos.append(s)
					if len(infos) > 2 - sub:					 
						arg = infos[2 - sub].strip()
						if arg == "":
							print(v, infos)
						if infos[0].strip() == "repeated":
							commitStr += "---@field %s %s[]\n" %(arg, infos[1 - sub].strip())
						else:
							commitStr += "---@field %s %s\n" %(arg, infos[1 - sub].strip())
					if v.find("}") != -1:
						funFind = False
						fw.write(commitStr)
						fw.write(funStr)	 
			fw.write("\n_s2c = S2C\n")
			fw.write("S2C.newS2C = {\n")
			fw.write("}\n")	
			fw.write("""
function S2C.DispatchError(err)
	if err == 0 then return end
	for k,v in pairs(_s2c.errors[err]) do
		if k == v and type(v) == "function" then v() end
	end
end
function S2C.RecvS2C(err, recv)
	_s2c.DispatchError(err)
	
	local pbName = GL.Net.Def[recv.head.cmd * 1024 + recv.head.act]
	if not pbName then 
		print("unknown cmd:" .. recv.head.cmd .. " act:" .. recv.head.act)
		return
	end
	
	local pb = {error = err}
	if not _s2c.newS2C[pbName] and err == 0 then
		pb = GL.PB.Client[pbName]()		
		pb:ParseFromString(recv.data)	 
	end
	local call = false
	for k,v in pairs(_s2c[pbName]) do
		if k == v and type(v) == "function" then
			v(pb)
			call = true
		end
	end
	if not call then 
		print("pb not has listener name:" .. pbName) 
	end		
end
""")
			
			fw.write("return S2C")
			
	with open("%s/s2c.lua" %luaroot,"r") as fw:
		gcommitStr += fw.read()
	with open("%s/s2c.lua" %luaroot,"w") as fw:
		fw.write(gcommitStr)
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
def GenLua():
	GenDefine()
	GenError()
	GenC2S()
	GenS2C()

