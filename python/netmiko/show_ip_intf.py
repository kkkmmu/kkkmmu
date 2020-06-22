#!/usr/bin/env python
from netmiko import Netmiko
# from getpass import getpass

cisco1 = {
    "host": "10.55.193.124",
    "username": "admin",
    # "password": getpass(),
    "device_type": "cisco_ios",
}

net_connect = Netmiko(**cisco1)
command = "show ip int brief"

print()
print(net_connect.find_prompt())
output = net_connect.send_command(command)
net_connect.disconnect()
print(output)
print()
