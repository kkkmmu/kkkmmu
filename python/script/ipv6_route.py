from telnetlib import Telnet
tn = Telnet('10.55.192.25', 23)   # connect to finger port
tn.write(b'admin\r\n')
print(tn.read_some())
tn.write(b'Dasan123456\r\n')
print(tn.read_some())
tn.write(b'enable\r\n')
print(tn.read_some())
tn.write(b'configure terminal\r\n')
print(tn.read_some())

for i in range(1, 2048):
    tn.write("ipv6 route 4000:%x::/64 2001:db8:11::11\r\n" % i)
    print(tn.read_until("#"))
