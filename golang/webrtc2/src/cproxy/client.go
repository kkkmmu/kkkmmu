package cproxy

import (
	"encoding/json"
	"fmt"
	//"github.com/kr/pty"
	"github.com/pions/webrtc"
	"golang.org/x/crypto/ssh/terminal"
	//"github.com/pions/webrtc/pkg/datachannel"
	"bufio"
	"bytes"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"strings"
)

type CClient struct {
	conn    *CSession
	tos     *terminal.State
	offerC  chan webrtc.RTCSessionDescription /* Channel used to receive offer. */
	answerC chan webrtc.RTCSessionDescription /* Channel used to send anser. */
	ready   chan bool
	address string
	sc      chan os.Signal
}

func NewCClient(server string) (*CClient, error) {
	return &CClient{
		address: server,
		offerC:  make(chan webrtc.RTCSessionDescription, 10),
		answerC: make(chan webrtc.RTCSessionDescription, 10),
		ready:   make(chan bool),
	}, nil
}

func (cc *CClient) Start() {
	go cc.Init()
	answer := <-cc.answerC
	err := cc.conn.SetAnswer(answer)
	if err != nil {
		panic(err)
	}
	cc.SetUpTerminal()
	go cc.MessageHandler()
	go cc.SignalHandler()
	go cc.Run()
	defer cc.Clear()
	select {}
}

func (cc *CClient) Clear() {
	cc.conn.Destroy()
	terminal.Restore(int(os.Stdin.Fd()), cc.tos)
}

func (cc *CClient) MessageHandler() {
	f := bufio.NewWriter(os.Stdout)
	for {
		typ, msg := cc.conn.Recv()
		if typ == CtrlMessage {
			if string(msg) == "OPEN" {
				cc.ready <- true
			} else if string(msg) == "QUIT" {
				cc.sc <- os.Interrupt
			}
			continue
		}
		f.Write(msg)
		f.Flush()
	}
}

func (cc *CClient) Run() {
	<-cc.ready

	buf := make([]byte, 1024)
	for {
		nr, err := os.Stdin.Read(buf)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		/*
			if sdata == '`' {
				cc.conn.Send(CtrlMessage, []byte("QUIT"))
				continue
			}
		*/

		cc.conn.Send(DataMessage, buf[0:nr])
	}
}

func (cc *CClient) Init() {
	news, err := NewCSession(GenerateUserID(), nil)
	if err != nil {
		panic(err)
	}
	cc.conn = news

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(news.GetOffer())
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://"+cc.address, "application/json; charset=utf-8", b)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer resp.Body.Close()

	var answer webrtc.RTCSessionDescription
	err = json.NewDecoder(resp.Body).Decode(&answer)
	if err != nil {
		panic(err)
	}

	cc.answerC <- answer
}

func (cc *CClient) SetUpTerminal() {
	var err error
	cc.tos, err = terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	/*
		winSize, err := pty.GetsizeFull(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		ctrl := CControl{
			Type:  "set_terminal",
			Value: fmt.Sprintf("%d,%d,%d,%d", winSize.Rows, winSize.Cols, winSize.X, winSize.Y),
		}

		data, err := json.Marshal(ctrl)
		if err != nil {
			panic(err)
		}
		cc.conn.Send(CtrlMessage, data)
	*/
}

func (cc *CClient) SignalHandler() {
	cc.sc = make(chan os.Signal, 1)
	signal.Notify(cc.sc, os.Interrupt, os.Kill)
	<-cc.sc
	cc.Clear()
	os.Exit(0)
}
