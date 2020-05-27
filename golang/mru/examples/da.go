package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"rut"
	"strings"
)

var (
	IP             = flag.String("ip", "10.55.192.211", "IP address of the remote device")
	Port           = flag.String("port", "23", "Port to connect")
	Host           = flag.String("hostname", "SWITCH", "Host name of the remote device")
	Protocol       = flag.String("prot", "telnet", "Login protocol(ssh/telnet)")
	User           = flag.String("u", "admin", "Username of the remote device")
	Password       = flag.String("p", "", "Passwrod of the remote device")
	ServerIP       = flag.String("sip", "10.55.2.202", "SFTP server IP address of the bin/core file")
	ServerUser     = flag.String("su", "liwei", "Username of the SFTP Server")
	ServerPassword = flag.String("sp", "Lee123!@#", "Passwrod of the SFTP Server")

	Core    = flag.String("core", "nsm.core", "coredump file name(path)")
	Bin     = flag.String("bin", "nsm.unstrip", "binary file name path")
	Process = flag.String("process", "hsl", "the process name to be checked")

	CTX = context.Background()
)

func main() {
	flag.Parse()
	dev, err := rut.New(&rut.RUT{
		Name:     "M3000-210",
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
		return
	}

	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()

	_, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "q sh -l",
	})

	if err != nil {
		panic(err)
	}

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "which gdb",
	})

	if !strings.Contains(string(data), "bin") {
		fmt.Println("Cannot find gdb on this device ", *IP)
		return
	}

	if err = dev.SCPGet("/etc/"+filepath.Base(*Core), *ServerIP, *ServerUser, *ServerPassword, filepath.Dir(*Core), filepath.Base(*Core)); err != nil {
		fmt.Printf("Cannot download file: %s from %s with: %s\n", *Core, *ServerIP, err.Error())
		return
	}

	if err = dev.SCPGet("/etc/"+filepath.Base(*Bin), *ServerIP, *ServerUser, *ServerPassword, filepath.Dir(*Bin), filepath.Base(*Bin)); err != nil {
		fmt.Printf("Cannot download file: %s from %s with: %s\n", *Bin, *ServerIP, err.Error())
		return
	}

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "ls",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "gdb " + filepath.Base(*Bin) + " " + filepath.Base(*Core),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "bt",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "frame 0",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "info locals",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "frame 1",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "info locals",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "info registers",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
