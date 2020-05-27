import socket

def retBanner(ip, port):
	try:
		socket.setdefaulttimeout(1)
		s = socket.socket()
		s.connect((ip, port))
		banner = s.recv(1024)
		return banner
	except:
		return

def main():
	portlist = [x for x in range(1,100)]
	for x in range (2,255):
		ip = "10.71.1." + str(x)
		for port in portlist:
			banner = retBanner(ip, port)
			if banner:
				print("[+] " + ip + ": " + banner)
		
if __name__ == '__main__':
	main()
