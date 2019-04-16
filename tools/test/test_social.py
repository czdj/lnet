#!/usr/bin/env python
# -*- coding:utf-8 -*-
import sys,time,os
root = os.path.split(os.path.realpath(__file__))[0]
root = os.path.dirname(root)
sys.path.insert(0, root + "/auto")
import python as pb
import socket
import urllib
import json
import urllib2
import msg_head

def SendHttp(url):
	print url
	req = urllib2.Request(url)
	res_data = urllib2.urlopen(req)
	res = res_data.read()
	print res
	return json.loads(res)

def Register(name, passwd):
	#登录验证服
	url = "http://127.0.0.1:5000/register?channel=mine&name=" + name + "&passwd=" + passwd
	re = SendHttp(url)
	if re["error"] != 0:
		return
	session = re["session"]
	if not re.has_key("roles"):
		#创建角色
		roleName = name + str(1)
		url = "http://127.0.0.1:5000/newrole?name=" + roleName + "&session=" + session + "&server=" + str(1) + "&type=" + str(1)
		re = SendHttp(url)
		if re["error"] != 0:
			return
		
	roles = re["roles"]
	role_id = roles[0]["id"]
	url = "http://127.0.0.1:5000/userole?id=" + str(role_id) + "&session=" + session
	re = SendHttp(url)
	
	#登录逻辑服
	c2s = pb.C2S()
	c2s.gamerLoginC2S.id = role_id
	data = c2s.SerializeToString()
	head = msg_head.Gen(len(data),1,1)

	ip_port = (re["server"]["ip"], re["server"]["port"])
	sk = socket.socket()
	sk.connect(ip_port)
	sk.sendall(head + data)

	server_reply = sk.recv(1024)
	if len(server_reply) >= 12:
		print msg_head.Unpack(server_reply[0:12])
	if len(server_reply) <= 12:		
		server_reply = server_reply + sk.recv(1024)

	s2c = pb.S2C()
	s2c.ParseFromString(server_reply[12:])
	print s2c

	c2s = pb.C2S()
	c2s.gamerLoginGetDataC2S.id = role_id
	data = c2s.SerializeToString()
	head = msg_head.Gen(len(data),1,1)
	sk.sendall(head + data)
	server_reply = sk.recv(1024)
	
	if len(server_reply) >= 12:
	    print msg_head.Unpack(server_reply[0:12])
	if len(server_reply) <= 12:
	    server_reply = server_reply + sk.recv(1024)
	
	s2c = pb.S2C()
	s2c.ParseFromString(server_reply[12:])
	print s2c
	
	
	
	c2s = pb.C2S()
	c2s.gamerFriendChatC2S.id = role_id
	c2s.gamerFriendChatC2S.toId = role_id
	c2s.gamerFriendChatC2S.msg = "xjsdkfjskjf"
	data = c2s.SerializeToString()
	head = msg_head.Gen(len(data),1,1)
	sk.sendall(head + data)
	server_reply = sk.recv(1024)
	
	if len(server_reply) >= 12:
	    print msg_head.Unpack(server_reply[0:12])
	if len(server_reply) <= 12:
	    server_reply = server_reply + sk.recv(1024)
	
	s2c = pb.S2C()
	s2c.ParseFromString(server_reply[12:])
	print s2c

	raw_input()


if __name__ == '__main__':
	if len(sys.argv) > 1:
		Register(sys.argv[1], "1")
	else:
		Register("tt344", "1")
	

