package client

import (
	"fmt"
	"ssh"
	"telnet"
	"time"
)

type Client interface {
	ReadUntil(delims ...string) ([]byte, error)
	ReadUntilIndex(delims ...string) ([]byte, int, error)
	ReadUntilSkip(delims, sdelims []string) ([]byte, error)
	ReadUntilContains(delims, cdelims []string) ([]byte, error)
	ReadString(delim byte) (string, error)
	Read(buf []byte) (int, error)
	ReadAll() ([]byte, int, error)
	SkipUntil(delims ...string) error
	Write(buf []byte) (int, error)
	WriteLine(data string) (int, error)
	WriteString(data string) (int, error)
	ReadTimeout(time.Duration) ([]byte, error)
}

func New(user, pass, proto, ip, port string) (Client, error) {
	if proto != "ssh" && proto != "telnet" {
		return nil, fmt.Errorf("Unknow proto: %s", proto)
	}

	if proto == "telnet" {
		return telnet.New(user, pass, ip, port)
	}

	return ssh.New(user, pass, ip, port)
}
