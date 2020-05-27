package main

import (
	"device"
	"fmt"
)

func main() {
	dev, err := device.New(&device.Device{
		Hostname: "SWITCH",
		Device:   "V8500_SFU",
		Type:     "PB",
		IP:       "10.71.20.230",
		Port:     "10024",
		Protocol: "telnet",
		Username: "admin",
		Password: "",
		SFU:      "A",
	})

	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}

	_, err = dev.WriteLine("enable")
	data, err := dev.Expect("#")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	dev.WriteLine("configure terminal")
	data, err = dev.Expect("#")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(dev.CurrentMode())

	dev.WriteLine("router ospf 1")
	data, err = dev.Expect("#")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(dev.CurrentMode())

	dev.WriteLine("exit")
	data, err = dev.Expect("#")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(dev.CurrentMode())

	dev.WriteLine("interface vlan 130")
	data, err = dev.Expect("#")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(dev.CurrentMode())
}
