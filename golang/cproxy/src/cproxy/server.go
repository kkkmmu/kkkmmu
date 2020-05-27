package cproxy

import (
	"errors"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type client chan<- []byte // an outgoing message channel

type Server struct {
	ip        string
	port      string
	cmd       *exec.Cmd
	ptmx      *os.File
	entering  chan client
	leaving   chan client
	messages  chan []byte
	broadcast chan []byte
	exit      chan bool
	tos       *terminal.State
	clients   map[client]bool // all connected clients
}

func NewServer(ip, port, cmd string) (*Server, error) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) == 0 {
		return nil, errors.New("Command is necessary for a server")
	}

	return &Server{
		ip:        ip,
		port:      port,
		cmd:       exec.Command(cmds[0], cmds[1:]...),
		entering:  make(chan client),
		leaving:   make(chan client),
		messages:  make(chan []byte),        // all incoming client messages
		broadcast: make(chan []byte, 10000), // all incoming client messages
		exit:      make(chan bool),
		clients:   make(map[client]bool), // all connected clients
	}, nil
}

func (s *Server) Start() {
	var err error
	s.ptmx, err = pty.Start(s.cmd)
	if err != nil {
		panic(err)
	}

	s.tos, err = terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer terminal.Restore(int(os.Stdin.Fd()), s.tos)
	defer s.ptmx.Close()

	go func() {
		if _, err := io.Copy(s.ptmx, os.Stdin); err != nil {
			s.exit <- true
			panic(err)
		}
	}()

	go func() {
		buf := make([]byte, 1000000000)
		for {
			n, err := s.ptmx.Read(buf)
			if err != nil {
				s.exit <- true
				return
			}

			if buf[0] == '`' {
				s.exit <- true
				return
			}

			/* Send to local console. */
			if _, err = os.Stdout.Write(buf[0:n]); err != nil {
				//cs.sc <- os.Interrupt
				s.exit <- true
				return
			}

			s.broadcast <- buf[:n]
		}
	}()

	go s.broadcaster()
	go func() {
		listener, err := net.Listen("tcp", s.ip+":"+s.port)
		if err != nil {
			log.Println(err.Error() + "\r")
			s.exit <- true
			return
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}
			go s.handleConn(conn)
		}
	}()
	select {
	case <-s.exit:
		terminal.Restore(int(os.Stdin.Fd()), s.tos)
		os.Exit(0)
	}

}

func (s *Server) broadcaster() {
	for {
		select {
		case msg := <-s.messages:
			_, err := s.ptmx.Write(msg)
			if err != nil {
				s.exit <- true
			}
		case msg := <-s.broadcast:
			for cli := range s.clients {
				cli <- msg
			}
		case cli := <-s.entering:
			s.clients[cli] = true

		case cli := <-s.leaving:
			delete(s.clients, cli)
			close(cli)
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	ch := make(chan []byte, 10000) // outgoing client messages
	go s.clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- []byte(who + " login success\n\r")
	s.broadcast <- []byte(who + " has connected\n\r")
	s.entering <- ch

	buf := make([]byte, 1024)
	for {
		nr, err := conn.Read(buf)
		if err != nil {
			break
		}
		s.messages <- buf[:nr]
	}

	s.leaving <- ch
	s.broadcast <- []byte(who + " has disconnected\n\r")
	conn.Close()
}

func (s *Server) clientWriter(conn net.Conn, ch <-chan []byte) {
	for msg := range ch {
		conn.Write(msg) // NOTE: ignoring network errors
	}
}
