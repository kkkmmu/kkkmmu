package fileserver

import (
	"errors"
)

type FileServer struct {
	Protocol string
	Name     string
	IP       string
	Port     string
	Username string
	Password string
}

func IsValidProtocol(proto string) bool {
	if proto != "ftp" && proto != "tftp" && proto != "ssh" {
		return false
	}

	return true
}

func New(name, ip, port, user, pass, proto string) (*FileServer, error) {
	if !IsValidProtocol(proto) {
		return nil, errors.New("Invalid protocol: " + proto)
	}

	return &FileServer{
		Protocol: proto,
		Name:     name,
		IP:       ip,
		Port:     port,
		Username: user,
		Password: pass,
	}, nil
}
