package main

import (
	"cproxy"
	"flag"
)

var Cmd = flag.String("cmd", "ssh 10.71.1.3 -l tsl", "The console that you want to proxy")
var Mode = flag.String("mode", "client", "The mode (server/client) that you want to run.")
var Address = flag.String("address", ":50000", "The server address that you want the client to connect.")

func main() {
	flag.Parse()
	if *Mode == "server" {
		cs, _ := cproxy.NewCServer(*Address, *Cmd)
		cs.Start()
	} else {
		cc, _ := cproxy.NewCClient(*Address)
		cc.Start()
	}
}
