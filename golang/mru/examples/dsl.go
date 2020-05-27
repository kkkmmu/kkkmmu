package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	"rut"
)

var (
	IP       = flag.String("ip", "10.71.20.156", "IP address of the remote device")
	Port     = flag.String("port", "23", "Port to connect")
	Host     = flag.String("hostname", "SWITCH", "Host name of the remote device")
	Protocol = flag.String("prot", "telnet", "Login protocol(ssh/telnet)")
	User     = flag.String("u", "admin", "Username of the remote device")
	Password = flag.String("p", "Dasan123456", "Passwrod of the remote device")

	Process = flag.String("process", "hsl", "process name to debug")
	Bin     = flag.String("bin", "hsl", "binary file name path")
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		Protocol: *Protocol,
		IP:       *IP,
		Port:     *Port,
		Username: *User,
		Hostname: *Host,
		Password: *Password,
	})

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	ctx := context.Background()

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "q sh -l",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "ls -l",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}
