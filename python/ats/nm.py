from netmiko import ConnectHandler

# Telnet Connection
telnetC = {
    "host": "10.55.192.186",
    "username": "admin",
    "password": "Dasan123456",
    "device_type": "ipinfusion_ocnos_telnet",
    "global_delay_factor": 0.5,
    "fast_cli": False,
}

connect = ConnectHandler(**telnetC)
connect.enable()
print(connect.send_command("terminal length 0", delay_factor=0.5))
print(connect.send_command("show arp", delay_factor=0.5))
print(connect.send_command("show mac address-table dynamic", delay_factor=0.5))
connect.disconnect()

# SSH Connection
sshC = {
    "host": "10.55.192.186",
    "username": "admin",
    "password": "Dasan123456",
    "device_type": "ipinfusion_ocnos_telnet",
}

connect = ConnectHandler(**sshC)
connect.enable()
print(connect.send_command("terminal length 0"))
print(connect.send_command("show arp", use_textfsm=True))
print(connect.send_command("show vlan brief", use_textfsm=True))
print(connect.send_command("show mac address-table dynamic", use_textfsm=True))
print(connect.send_command("show interface status", use_textfsm=True))
print(connect.send_command("show ip route", use_textfsm=True))
print(connect.send_command("show process", use_textfsm=True))
print(connect.send_command("show cpuload", use_textfsm=True))
print(connect.send_command("show memory system", use_textfsm=True))
print(connect.send_command("show ip ospf database", use_textfsm=True))
print(connect.send_command("show ip ospf interface brief", use_textfsm=True))
print(connect.send_command("show ip ospf neighbor", use_textfsm=True))
print(connect.send_command("show ipv6 neighbor", use_textfsm=True))
print(connect.send_command("show ipv6 interface brief", use_textfsm=True))
# print(connect.send_command("show interface", use_textfsm=True))
# print(connect.send_command("show ip interface brief", use_textfsm=True))
connect.disconnect()
