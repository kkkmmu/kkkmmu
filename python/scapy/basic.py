#! /usr/bin/env python
from scapy.all import *

# Generate 255 ip.
hosts=IP(dst="10.71.1.146", src="10.0.1.*")/ICMP()
for h in hosts:
    print(h.src)
sendp(hosts)

# Generate 255 ip.
hosts1=IP(dst="10.71.1.146", src="10.71.1.0/30")/ICMP()
for h in hosts1:
    print(h.src)
send(hosts1)

# Generate 255 ip.
hosts2=IP(dst="10.71.1.146", src="2.*.1.1")/ICMP()
for h in hosts2:
    print(h.src)
send(hosts2)

macs=Ether(dst=["ff:ff:ff:ff:ff:ff", "00:11:12:12:12:12", "00:11:11:11:11:11"])
for m in macs:
    print(m.dst)

# Generate 10 unicast mac address. 
macs=Ether(dst=["00:"+":".join("%02x" % b for b in [random.randrange(256) for _ in range(5)]) for _ in range(10)])
for m in macs:
    print(m.dst)


# Generate 10 multicast mac address. 
macs=Ether(dst=["01:"+":".join("%02x" % b for b in [random.randrange(256) for _ in range(5)]) for _ in range(10)])
for m in macs:
    print(m.dst)


# Sniff received ARP packet and do some action
def pkt_monitor_callback(pkt):
    for p in pkt: #who-has or is-at
        if ARP in p:
            if p[ARP].op == 1 or p[ARP].op == 2:
                print(p.sprintf("%Ether.dst% %Ether.src% %ARP.psrc% %ARP.hwsrc% %ARP.op% %ARP.pdst% %ARP.hwdst% %ARP.hwtype% %ARP.ptype% %ARP.hwlen% %ARP.plen%"))
                send(IP(dst="10.71.1.146")/ICMP()/p/"Hello from docker!")


sniff(prn=pkt_monitor_callback, store=0)
