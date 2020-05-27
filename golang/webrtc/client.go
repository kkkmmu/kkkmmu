package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	//"time"

	"github.com/kr/pty"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/examples/util"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
	"golang.org/x/crypto/ssh/terminal"
)

type Ctrl struct {
	Type  string
	Value string
}

func SetUpTerminal(dc *webrtc.RTCDataChannel) {
	_, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	winSize, err := pty.GetsizeFull(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := Ctrl{
		Type:  "set_terminal",
		Value: fmt.Sprintf("%d,%d,%d,%d", winSize.Rows, winSize.Cols, winSize.X, winSize.Y),
	}

	data, err := json.Marshal(ctrl)
	if err != nil {
		panic(err)
	}
	dc.Send(datachannel.PayloadString{Data: data})
}

func main() {
	addr := flag.String("address", ":50000", "Address that the HTTP server is hosted on.")
	label := flag.String("id", "data", "Label for your connnection.")
	flag.Parse()

	// Everything below is the pion-WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.New(config)
	util.Check(err)

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel(*label, nil)
	util.Check(err)

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState ice.ConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("'%s'-'%d' open.\n", dataChannel.Label, dataChannel.ID)
		SetUpTerminal(dataChannel)
	})

	// Register the OnMessage to handle incoming messages
	dataChannel.OnMessage(func(payload datachannel.Payload) {
		switch p := payload.(type) {
		case *datachannel.PayloadString:
			fmt.Printf("Message '%s' from DataChannel '%s' payload '%s'\n", p.PayloadType().String(), dataChannel.Label, string(p.Data))
		case *datachannel.PayloadBinary:
			f := bufio.NewWriter(os.Stdout)
			f.Write(p.Data)
			f.Flush()
		default:
			fmt.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), dataChannel.Label)
		}
	})

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	util.Check(err)

	// Exchange the offer for the answer
	answer := mustSignalViaHTTP(offer, *addr)

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answer)
	util.Check(err)

	buf := make([]byte, 1024)
	for {
		nr, err := os.Stdin.Read(buf)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		err = dataChannel.Send(datachannel.PayloadBinary{Data: buf[0:nr]})
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	// Block forever
	select {}
}

// mustSignalViaHTTP exchange the SDP offer and answer using an HTTP Post request.
func mustSignalViaHTTP(offer webrtc.RTCSessionDescription, address string) webrtc.RTCSessionDescription {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(offer)
	util.Check(err)

	resp, err := http.Post("http://"+address, "application/json; charset=utf-8", b)
	util.Check(err)
	defer resp.Body.Close()

	var answer webrtc.RTCSessionDescription
	err = json.NewDecoder(resp.Body).Decode(&answer)
	util.Check(err)

	return answer
}
