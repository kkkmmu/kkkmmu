package cproxy

import (
	"fmt"
	"net"
	"regexp"
)

type Logd struct {
	address string
	filter  *regexp.Regexp
}

func NewLogd(address, filter string) (*Logd, error) {
	var fr *regexp.Regexp
	if filter != "" {
		fr = regexp.MustCompile(filter)
	}
	return &Logd{
		address: address,
		filter:  fr,
	}, nil
}

func (s *Logd) Start() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(s.address), Port: 514})
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}

		if s.filter != nil {
			match := s.filter.FindStringSubmatch(string(data[:n]))
			if len(match) > 0 {
				fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
			}
		} else {
			fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		}
	}
}
