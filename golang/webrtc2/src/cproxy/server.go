package cproxy

import (
	"encoding/json"
	"errors"
	"github.com/kr/pty"
	"github.com/pions/webrtc"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

type CServer struct {
	clients map[*webrtc.RTCSessionDescription]*CSession
	tos     *terminal.State
	cmd     *exec.Cmd /* Startup command. */
	ptmx    *os.File
	sc      chan os.Signal
	offerC  chan *webrtc.RTCSessionDescription /* Channel used to receive offer. */
	answerC chan *webrtc.RTCSessionDescription /* Channel used to send anser. */
	address string
}

func NewCServer(address, cmd string) (*CServer, error) {
	cmds := strings.Split(cmd, " ")
	if len(cmds) == 0 {
		return nil, errors.New("Command is necessary for a server")
	}

	/* Address validation check is necessary .*/

	return &CServer{
		clients: make(map[*webrtc.RTCSessionDescription]*CSession, 10),
		address: address,
		cmd:     exec.Command(cmds[0], cmds[1:]...),
		offerC:  make(chan *webrtc.RTCSessionDescription, 10),
		answerC: make(chan *webrtc.RTCSessionDescription, 10),
		sc:      make(chan os.Signal),
	}, nil
}

func (cs *CServer) Clear() {
	cs.ClearClients()
	terminal.Restore(int(os.Stdin.Fd()), cs.tos)
}

func (cs *CServer) Broadcast(typ int, data []byte) {
	for _, c := range cs.clients {
		c.Send(typ, data)
	}
}

func (cs *CServer) Start() {
	go cs.Listen()
	go cs.Accept()
	go cs.SignalHandler()
	defer cs.Clear()
	cs.Serve()
}

func (cs *CServer) Serve() error {
	cs.StartCommand()
	buf := make([]byte, 10000000)
	for {
		n, err := cs.ptmx.Read(buf)
		if err != nil {
			cs.sc <- os.Interrupt
			break
		}

		/* Send to local console. */
		if _, err = os.Stdout.Write(buf[0:n]); err != nil {
			cs.sc <- os.Interrupt
		}

		cs.Broadcast(DataMessage, buf[0:n])
	}

	return errors.New("INPUT/OUPUT error")
}

func (cs *CServer) StartCommand() error {
	var err error

	cs.ptmx, err = pty.Start(cs.cmd)
	if err != nil {
		panic(err)
	}

	cs.tos, err = terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	go func() {
		if _, err := io.Copy(cs.ptmx, os.Stdin); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (cs *CServer) Accept() {
	for offer := range cs.offerC {
		news, err := NewCSession("server", offer)
		if err != nil {
			panic(err)
		}
		cs.answerC <- news.answer

		go cs.HandleClientMessage(news)
		cs.clients[offer] = news
	}
}

func (cs *CServer) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var offer webrtc.RTCSessionDescription
		err := json.NewDecoder(r.Body).Decode(&offer)
		if err != nil {
			log.Fatal(err)
			return
		}

		/* Received the offer, and wait for the answer .*/
		cs.offerC <- &offer
		answer := <-cs.answerC

		err = json.NewEncoder(w).Encode(answer)
		if err != nil {
			log.Fatal(err)
			return
		}
	})

	http.ListenAndServe(cs.address, nil)
}

func (cs *CServer) HandleClientMessage(client *CSession) {
	for {
		typ, msg := client.Recv()
		if typ == CtrlMessage {
			if string(msg) == "QUIT" {
				cs.RemoveClient(client)
			}
			continue
		}
		_, err := cs.ptmx.Write(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (cs *CServer) RemoveClient(client *CSession) {
	offer := client.GetOffer()
	delete(cs.clients, &offer)
	client.Destroy()
}

func (cs *CServer) SignalHandler() {
	cs.sc = make(chan os.Signal, 1)
	signal.Notify(cs.sc, os.Interrupt, os.Kill)
	<-cs.sc
	cs.Clear()
	os.Exit(0)
}

func (cs *CServer) ClearClients() {
	cs.Broadcast(CtrlMessage, []byte("QUIT"))
	cs.clients = nil
}
