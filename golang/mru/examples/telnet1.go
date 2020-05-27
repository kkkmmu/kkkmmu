package main

import (
	"fmt"
	"telnet"
)

func main() {
	s, err := telnet.New("admin", "", "10.71.20.230", "10027")
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

	fmt.Println(string(data))
	s.WriteLine("show running-config")
	data, err = s.ReadUntil("#")
	data, err = s.ReadUntil("#")
	data, err = s.ReadUntil("#")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	s.WriteLine("exit")
}
