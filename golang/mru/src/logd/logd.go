package logd

import (
	"fmt"
	"net"
)

type Logd struct {
	address string
}

func New(address string) *Logd {
	return &Logd{
		address: address,
	}
}

func (s *Logd) Start() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(s.address), Port: 514})
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
	}
}
