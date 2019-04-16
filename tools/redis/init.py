import socket

ip_port = ('127.0.0.1',9996)
s = socket.socket()
s.connect(ip_port)
s.sendall('*4\r\n$4\r\nhset\r\n$3\r\ndbs\r\n$1\r\n1\r\n$15\r\n@127.0.0.1:9997\r\n')
s.sendall('*4\r\n$4\r\nhset\r\n$3\r\ndbs\r\n$1\r\n2\r\n$15\r\n@127.0.0.1:9998\r\n')
s.close()