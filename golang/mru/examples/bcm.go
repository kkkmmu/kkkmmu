package main

import (
	"fmt"
	"os"
	"strings"
	"telnet"
	"time"
	"util"
)

func main() {
	s, err := telnet.New("admin", "", "10.71.20.230", "10027")
	//s, err := telnet.New("admin", "", "10.71.20.212", "23")
	if err != nil {
		panic(err)
	}

	os.Remove("bcm.log")

	defer func() {
		s.WriteLine("exit")
		s.WriteLine("exit")
		s.WriteLine("exit")
	}()

	s.WriteLine("exit")

	fmt.Println("=============================")
	for {
		tick60 := time.Tick(time.Second * 120)
		tick10 := time.Tick(time.Second * 5)
		s.WriteLine("\n")
		data, err := s.ReadUntil("ogin:")
		if err != nil {
			panic(err)
		}

		util.AppendToFile("bcm.log", data)

		fmt.Println(string(data))
		s.WriteLine("admin")
		data, err = s.ReadUntil("ssword:")
		if err != nil {
			panic(err)
		}

		util.AppendToFile("bcm.log", data)
		fmt.Println(string(data))
		s.WriteLine("Dasan123456")
		data, err = s.ReadUntil(">")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
		util.AppendToFile("bcm.log", data)
		s.WriteLine("enable")
		data, err = s.ReadUntil("#")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
		util.AppendToFile("bcm.log", data)
		s.WriteLine("terminal length 0")
		data, err = s.ReadUntil("#")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))
		util.AppendToFile("bcm.log", data)

		<-tick60
		s.WriteLine("show bgp neighbors 112.1.1.1")
		data, err = s.ReadUntil("M3000_3#")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		util.AppendToFile("bcm.log", data)
		data, err = s.ReadUntil("M3000_3#")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))

		util.AppendToFile("bcm.log", data)
		data, err = s.ReadUntil("M3000_3#")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		util.AppendToFile("bcm.log", data)
		<-tick10
		if !strings.Contains(string(data), "state = Established") {
			s.WriteLine("dlfldhsjfk-b")
			continue
		} else {
			fmt.Println("test", string(data))
			s.WriteLine("q sh -l")
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			fmt.Println(1, string(data))
			util.AppendToFile("bcm.log", data)
			data, err = s.ReadUntil("#")
			fmt.Println(2, string(data))
			util.AppendToFile("bcm.log", data)
			if err != nil {
				panic(err)
			}

			fmt.Println(1, string(data))
			util.AppendToFile("bcm.log", data)
			s.WriteLine("scontrol  -f /proc/swal pld read 0x35")
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			if strings.Contains(string(data), "0xff") {
				panic(data)
			}

			s.WriteLine("scontrol  -f /proc/swal pld read 0x36")
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			if strings.Contains(string(data), "0xff") {
				panic(data)
			}

			s.WriteLine("scontrol  -f /proc/swal pld read 0x37")
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			data, err = s.ReadUntil("#")
			fmt.Println(string(data))
			util.AppendToFile("bcm.log", data)
			if strings.Contains(string(data), "0xff") {
				panic(data)
			}

			<-tick10
			s.WriteLine("dlfldhsjfk-b")
		}
	}
}
