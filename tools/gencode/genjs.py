#!/usr/bin/python
# -*- coding:utf-8 -*-
#npm install -g protobufjs
#pbjs -t static-module -w commonjs -o proto.js *.proto
#pbts -o proto.d.ts proto.js
from gencomm import *
#cocosNetRoot,cocosHandlesRoot
def GenJSNet():
	type_map = {"32":"number","64":"number","int32":"number","uint32":"number","int64":"number","uint64":"number","string":"string","bool":"boolean"}
	f2 = open("%s/netproto.js" % cocosNetRoot, "w+")
	f2.write(
"""
(function() {
    
var C2S = net.C2S;
var _EventDispatch = net._EventDispatch;

var netlogic = 
{
"""
		)
	f3 =  open("%s/AntNet.ts" % cocosNetRoot, "w+")
	f3.write(
"""
declare var net;
export default class AntNet
{
	public static get onError(){
		return net.logic.onError;
	}
	public static get onConnect(){
		return net.logic.onConnect;
	}
	public static get onClose(){
		return net.logic.onClose;
	}
	public static get onReconnect(){
		return net.logic.onReconnect;
	}
	public static get logicPing() {
		return net.logic.ping;
	}
"""
			)

	with open("%s/client.proto" %pbroot,"r") as f:
		old=f.readlines()

	s2cm =[]
	c2ss = []
	c2sm = {}
	c2sm_type = {}
	notifys = []
	found = False
	for line in old:
		if found == False:
			sss = []
			for s in line.strip().split(" "):
				if s.strip() != "":
					sss.append(s.strip())
			if len(sss) == 2 and sss[0] == "message" and sss[1].find("C2S")!= -1:
				c2ss.append(sss[1])
				c2sm[sss[1]] = []
				c2sm_type[sss[1]] = []
				found = True
			elif len(sss) == 2 and sss[0] == "message" and sss[1].find("S2C")!= -1 and sss[1].find("GamerNotify")!= -1:
				notifys.append(sss[1])
			if len(sss) == 2 and sss[0] == "message" and sss[1].find("S2C")!= -1:
				s2cm.append(sss[1])
		else:
			sss = []
			a = line.strip()
			b = line.strip().split(" ")
			for s in line.strip().strip("repeated").strip("optional").split(" "):
				if s.strip() != "":
					sss.append(s.strip())
			if sss[0] == "}":
				found = False
			elif sss[0] == "{":
				continue
			else:
				if sss[0] == "optional":
					sss = sss[1:]
				c2sm[c2ss[len(c2ss)-1]].append(sss[1])
				c2sm_type[c2ss[len(c2ss)-1]].append([sss[0],b[0]])
	#print c2sm
	for k in c2ss:
		v = c2sm[k]
		oldk = k
		k = k[0].lower() + k[1:]
		f2.write("\n\n%s: function(" %(k))
		for vv in v:
			if vv == "id":
				continue
			f2.write("%s, " %(vv))
		if len(v) > 1:
			f2.seek(-2, 1)
		if len(v) > 1:
			f2.write('){\n\tC2S("%s", {' %(oldk))
		else:
			f2.write('){\n\tC2S("%s", {' %(oldk))
		for vv in v:
			if vv == "id":
				continue
			f2.write("%s:%s, " %(vv, vv))
		if len(v) > 1:
			f2.seek(-2, 1)
		f2.write("});\n},")
		f2.write("\n%s: new _EventDispatch()," %(k.replace("C2S", "S2C")))
	for k in notifys:
		f2.write("\n%s: new _EventDispatch()," %(k[0].lower() + k[1:]))
	f2.write(
"""
};

for(var key in netlogic)
{
    net.logic[key] = netlogic[key];
}

})()

"""
		)
	f2.close()

	#print c2sm ts
	importType = []
	for k in c2ss:
		v = c2sm[k]
		v_type = c2sm_type[k]
		oldk = k
		newk = k[0].lower() + k[1:]
		#k = k[0].lower() + k[1:]
		f3.write("\n\n\tpublic static %s(" %(k))
		i = 0
		for vv in v:
			if vv == "id":
				i = i + 1
				continue
			type = ""
			if type_map.has_key(v_type[i][0]):
				type = type_map[v_type[i][0]]
			else:
				type = v_type[i][0]
				if importType.count(type) == 0:
					importType.append(type)
			if v_type[i][1] == "repeated":
				type = type+"[]"
			f3.write("%s: %s, " %(vv,type))
			i = i + 1
		if len(v) > 1:
			f3.seek(-2, 1)
		if len(v) > 1:
			f3.write('){\n\t\tnet.logic.%s(' %(newk))
		else:
			f3.write('){\n\t\tnet.logic.%s(' % (newk))
		for vv in v:
			if vv == "id":
				continue
			f3.write("%s, " %(vv))
		if len(v) > 1:
			f3.seek(-2, 1)
		f3.write(");\n\t}")

		f3.write("\n\n\tpublic static get %s(){\n\t\treturn net.logic.%s;\n\t}" %(newk.replace("C2S", "S2C"),newk.replace("C2S", "S2C")))
		####################################################################################
		f3.write("\n\n\tpublic static async Async%s(" % (k))
		i = 0
		for vv in v:
			if vv == "id":
				i = i + 1
				continue
			type = ""
			if type_map.has_key(v_type[i][0]):
				type = type_map[v_type[i][0]]
			else:
				type = v_type[i][0]
			if v_type[i][1] == "repeated":
				type = type+"[]"
			f3.write("%s: %s, " % (vv, type))
			i = i + 1
		if len(v) > 1:
			f3.seek(-2, 1)
		f3.write('):Promise<%s>{' % (k.replace("C2S", "S2C")))
		if importType.count(k.replace("C2S", "S2C")) == 0:
			importType.append(k.replace("C2S", "S2C"))
		f3.write(
"""
		return new Promise<%s>((resolve)=>{
			let s2c = {error:200} as %s;
			let timeObj = setTimeout(()=>{
				resolve(s2c);
			}, 5000);
			let fun = function(e:%s){
				clearTimeout(timeObj);
				AntNet.%s.off(fun);
				resolve(e);
			}
			AntNet.%s.on(fun);
"""%(k.replace("C2S", "S2C"),k.replace("C2S", "S2C"),k.replace("C2S", "S2C"),newk.replace("C2S", "S2C"),newk.replace("C2S", "S2C"))
		)
		f3.write('\t\t\tAntNet.%s('%(k))
		for vv in v:
			if vv == "id":
				continue
			f3.write("%s, " % (vv))
		if len(v) > 1:
			f3.seek(-2, 1)
		f3.write(");\n\t\t});\n\t}")
		###########################################################################################################
	for k in notifys:
		f3.write("\n\n\tpublic static get %s(){\n\t\treturn net.logic.%s;\n\t}" % (k[0].lower() + k[1:], k[0].lower() + k[1:]))
	f3.write("\n}")

	f3.flush()
	f3.seek(0, 0)
	lines = f3.readlines()
	#f3.truncate()
	f3.seek(0, 0)
	str = "import {"
	for v in importType:
		str += v + ", "
	str = str[:-2]
	str += '} from "./proto";\n'
	f3.write(str)
	for line in lines:
		f3.write(line)
	f3.close()


#客户端收消息处理
	importStr = ""
	varStr = ""
	bindStr=""
	protoNameStr=""

	for k in s2cm:
		f7 = open("%s/%sHandler.ts" %(cocosHandlesRoot,k), "w+")
		f7.write(
"""
import { %s } from "../Net/proto";

var %sHandler = function(msg: %s)
{

}

export {%sHandler}
"""%(k,k,k,k)
		)
		f7.close()
		
		protoNameStr+="\tstatic %s=\"%s\"\n"%(k,k)
		importStr+="import {%sHandler} from \"./%sHandler\";\n"%(k,k)
		varStr+="\t\t%s = %sHandler;\n"%(k,k)
		bindStr+="\t\tthis.%s.bind(this);\n"%(k)

	f4 = open("%s/ProtoHandlerList.ts" % cocosHandlesRoot, "w+")
	f4.write(
"""
%s
export default class ProtoHandlerList
{
%s
	init()
	{
%s
	}
}
"""%(importStr,varStr,bindStr)
		)

	f4.close()
	f5 = open("%s/ProtoName.ts" % cocosHandlesRoot, "w+")
	f5.write(
"""
export default class ProtoName
{
%s
}
"""%(protoNameStr)
			)

	f5.close()


def GenJS():
	GenJSNet()
	CopyFile(genpyroot + "/cocos/net.js", cocosNetRoot + "/net.js")
	CopyFile(genpyroot + "/cocos/protobuf.js", cocosNetRoot + "/protobuf.js")
	CopyFile(genpyroot + "/cocos/pako.js", cocosNetRoot + "/pako.js")
	CopyFile(genpyroot + "/web/test.html", webroot + "/test.html")
	CopyFile(genpyroot + "/web/net.js", webroot + "/net.js")
	with open("%s/net.js" % cocosNetRoot, "r") as f:
		old=f.readlines()
		write = False
		with open(webroot + "/net.js", "a") as f2:
			for line in old:
				if write:
					f2.write(line)
				if line.find("onClose:  new _EventDispatch()") == 0:
					write = True
			
	os.system('pbjs -t static-module -w commonjs -o %s/proto.js %s/auto_c2s.proto %s/auto_s2c.proto %s/client.proto %s/common.proto' %(cocosNetRoot, pb3root, pb3root, pb3root, pb3root))
	os.system('pbts -o %s/proto.d.ts %s/proto.js' %(cocosNetRoot, cocosNetRoot))
	updateProtoTsFile()
	if updateProtoJsFile() ==False:
		print "protobuf for cocos can't gen, please install protobufjs"
	CopyDir(webroot, wwwroot)
	

#laya使用	
def updateProtoTsFile():
	if os.path.exists("%s/proto.d.ts" %cocosNetRoot) == False:
		return
	type_map = {"Long":"any"}
	buf = []
	fp = open(cocosNetRoot+"/proto.d.ts","r+")
	flag = False
	for line in fp:
		tmp = line
		if line.find("*")> 0:
			#if line.startswith("import"):
			#	buf.append(line)
			continue

		if line.find("public static create")> 0:
			#buf.append(line)
			flag = True

		if flag == True:
			if line.startswith("}"):
				flag = False

		if flag == False:
			index = line.find("Long")
			if index != -1:
				tmp = line[:index]+"any"+line[index+4:]
				#line.replace("Long","any")
			buf.append(tmp)

	fp.seek(0,0)
	fp.truncate()
	for i in buf:
		fp.write(i)

def updateProtoJsFile():
	if os.path.exists("%s/proto.js" %cocosNetRoot) == False:
		return False
	with open("%s/proto.js" %cocosNetRoot,"r") as f:
		old=f.readlines()
		old[3] = 'var $protobuf = protobuf;\n'
		old.insert(11,"var proto = $root;\nvar module = {}\n")
		old.append("window.proto = $root;")
	with open("%s/proto.js" %cocosNetRoot,"w") as f:
		f.writelines(old)
	return True