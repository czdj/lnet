#!/usr/bin/python
# -*- coding:utf-8 -*-
import os,os.path,redis,socket
root = os.path.split(os.path.realpath(__file__))[0]
if os.path.exists('%s/redis.conf' %(root)):
	os.remove('%s/redis.conf' %(root))
if os.path.exists('%s/redis1.conf' %(root)):
	os.remove('%s/redis1.conf' %(root))
with open("%s/redis.conf" %root,"w") as f:
	f.write('''
daemonize yes
bind 127.0.0.1
port 9999

appendonly yes 
appendfsync everysec
# 只增文件的文件名称。（默认是appendonly.aof）  
appendfilename appendonly0.aof  
''')
with open("%s/redis1.conf" %root,"w") as f2:
	f2.write('''
daemonize yes
bind 127.0.0.1
port 9998


appendonly yes 
appendfsync everysec
# 只增文件的文件名称。（默认是appendonly.aof）  
appendfilename appendonly1.aof  
''')
with open("%s/redis2.conf" %root,"w") as f3:
	f3.write('''
daemonize yes
bind 127.0.0.1
port 9997


appendonly yes
appendfsync everysec
# 只增文件的文件名称。（默认是appendonly.aof）
appendfilename appendonly2.aof
''')
os.system('redis-server %s/redis.conf' %(root))
os.system('redis-server %s/redis1.conf' %(root))
os.system('redis-server %s/redis2.conf' %(root))

ip_port = ('127.0.0.1',9999)
s = socket.socket()
s.connect(ip_port)
s.sendall('*4\r\n$4\r\nhset\r\n$3\r\ndbs\r\n$1\r\n1\r\n$15\r\n@127.0.0.1:9998\r\n')
s.close()

