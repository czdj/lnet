#!/usr/bin/python
# -*- coding:utf-8 -*-
#npm install -g protobufjs
#pbjs -t static-module -w commonjs -o proto.js *.proto
#pbts -o proto.d.ts proto.js
from gencomm import *

varMap = {
	'int32':'int',
	'string':'string',
	'bool':'bool'
}

def GenCSNet():
	return
	with open("%s/AntNetEvents.cs" %csnetroot,"w") as f2:
		f2.write('''
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;
using System.Reflection;
public partial class AntNet : MonoBehaviour {
''')
		with open("%s/client.proto" %pbroot,"r") as f:
			old=f.readlines()
			
		c2ss = []
		c2sm = {}
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
					found = True
				elif len(sss) == 2 and sss[0] == "message" and sss[1].find("S2C")!= -1 and sss[1].find("GamerNotify")!= -1:
					notifys.append(sss[1])
			else:
				sss = []
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
					c2sm[c2ss[len(c2ss)-1]].append([sss[0], sss[1]])
		#print c2sm
		for k in c2ss:
			v = c2sm[k]
			oldk = k
			k = k[0].lower() + k[1:]
			f2.write("\n\npublic %s(" %(k))
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
		f2.write('''
}
''')
		
		
def GenCS():
	GenCSNet()