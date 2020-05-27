package main

import (
	"telnet"
	//"util"
	"fmt"
)

func main() {
	c, err := telnet.New("admin", "Dasan123456", "10.71.20.182", "23")
	if err != nil {
		panic(err)
	}

	c.WriteLine("ls")
	data, n, err := c.ReadAll()
	fmt.Println(string(data))
	data, err = c.ReadUntil("#")
	fmt.Println(string(data))

	fmt.Println(n)
	fmt.Println(err)
}
