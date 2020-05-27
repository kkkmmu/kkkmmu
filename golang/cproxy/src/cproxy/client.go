package cproxy

import (
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net"
	"os"
	"os/signal"
)

type Client struct {
	ip   string
	port string
	tos  *terminal.State
	conn net.Conn
	done chan struct{}
	sc   chan os.Signal
}

func NewClient(ip, port string) (*Client, error) {
	return &Client{
		ip:   ip,
		port: port,
		done: make(chan struct{}),
		sc:   make(chan os.Signal, 1),
	}, nil
}

func (c *Client) Start() {
	var err error
	c.conn, err = net.Dial("tcp", c.ip+":"+c.port)
	if err != nil {
		log.Println(err.Error() + "\r")
		return
	}

	c.tos, err = terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		c.done <- struct{}{}
	}

	go func() {
		f := bufio.NewWriter(os.Stdout)
		buf := make([]byte, 10000000)
		for {
			nr, err := c.conn.Read(buf)
			if err != nil {
				log.Println(err.Error() + "\r")
				c.done <- struct{}{}
				break
			}
			f.Write(buf[:nr])
			f.Flush()
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			nr, err := os.Stdin.Read(buf)
			if err != nil {
				log.Println(err.Error() + "\r")
				c.done <- struct{}{}
				break
			}
			if buf[0] == '`' {
				c.done <- struct{}{}
				break
			}

			_, err = c.conn.Write(buf[0:nr])
			if err != nil {
				log.Println(err.Error() + "\r")
				c.done <- struct{}{}
				break
			}
		}
	}()

	signal.Notify(c.sc, os.Interrupt, os.Kill)

	select {
	case <-c.done:
		c.conn.Close()
		terminal.Restore(int(os.Stdin.Fd()), c.tos)
	case <-c.sc:
		c.conn.Close()
		terminal.Restore(int(os.Stdin.Fd()), c.tos)
	}
	os.Exit(0)
}
