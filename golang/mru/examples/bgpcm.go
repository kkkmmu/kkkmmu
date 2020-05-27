package main

import (
	"fmt"
	"telnet"
)

func main() {
	//s, err := telnet.New("admin", "", "10.71.20.230", "10027")
	s, err := telnet.New("admin", "", "10.71.20.212", "23")
	//s, err := telnet.New("admin", "", "10.71.20.33", "23")
	if err != nil {
		panic(err)
	}

	s.WriteLine("")
	data, err := s.ReadUntil(">")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	s.WriteLine("enable")
	data, err = s.ReadUntil("#")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	s.WriteLine("terminal length 0")
	data, err = s.ReadUntil("#")
	if err != nil {
		panic(err)
	}

	s.WriteLine("configure terminal")
	data, err = s.ReadUntil("#")
	if err != nil {
		panic(err)
	}

	for i := 1; i < 19; i++ {
		ifname := fmt.Sprintf("eth0/%d", i)
		s.WriteLine("interface " + ifname)
		data, err = s.ReadUntil("#")
		if err != nil {
			fmt.Println(string(data))
		}
		s.WriteLine("speed 10000")
		data, err = s.ReadUntil("#")
		if err != nil {
			fmt.Println(string(data))
		}
		s.WriteLine("bgp convergence monitor")
		data, err = s.ReadUntil("#")
		if err != nil {
			fmt.Println(string(data))
		}
		s.WriteLine(fmt.Sprintf("channel-group %d mode active", i))
		data, err = s.ReadUntil("#")
		if err != nil {
			fmt.Println(string(data))
		}

		s.WriteLine("exit")
		data, err = s.ReadUntil("#")
		if err != nil {
			fmt.Println(string(data))
		}
	}

	fmt.Println(string(data))
	s.WriteLine("exit")
}
