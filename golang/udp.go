package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("10.71.1.35"), Port: 514})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		/*
			_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
			if err != nil {
				fmt.Printf(err.Error())
			}
		*/
	}
}
