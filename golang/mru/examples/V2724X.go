package main

import (
	"command"
	"context"
	"fmt"
	"rut"
	"util"
)

func main() {
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		IP:       "10.71.20.20",
		Port:     "23",
		Username: "admin",
		Hostname: "SWITCH",
		Password: "",
	})

	if err != nil {
		panic(err)
	}

	dev.Init()

	ctx := context.Background()

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "q sh -l",
	})
	fmt.Println(data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "bcm dump soc",
	})
	fmt.Println(data)
	util.SaveToFile("soc.txt", []byte(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "bcm dump socmem diff",
	})
	fmt.Println(data)
	util.SaveToFile("socmem.txt", []byte(data))
}
