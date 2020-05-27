package main

import (
	"cproxy"
	"flag"
	"fmt"
)

var IP = flag.String("ip", "localhost", "ip address of proxy server")
var Port = flag.String("port", "8000", "port of proxy server")
var CMD = flag.String("cmd", "ssh 10.71.1.3 -l tsl", "The remote command you need proxy")
var Mode = flag.String("mode", "client", "the proxy mode (server/client/logd)")
var Filter = flag.String("filter", "", "syslog filter regular expression")

func main() {
	flag.Parse()
	if *Mode == "server" {
		s, err := cproxy.NewServer(*IP, *Port, *CMD)
		if err != nil {
			panic(err)
		}
		s.Start()
	} else if *Mode == "client" {
		c, err := cproxy.NewClient(*IP, *Port)
		if err != nil {
			panic(err)
		}
		c.Start()
	} else if *Mode == "logd" {
		l, err := cproxy.NewLogd(*IP, *Filter)
		if err != nil {
			panic(err)
		}
		l.Start()
	} else {
		fmt.Println("Invalid mode: ", *Mode)
	}
}
