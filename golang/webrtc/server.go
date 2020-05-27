package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	//"time"

	"github.com/kr/pty"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/examples/util"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
	"strconv"
)

type Ctrl struct {
	Type  string
	Value string
}

func ParseTerminalCtrl(ctrl string) []int {
	fields := strings.Split(ctrl, ",")

	result := make([]int, 0, 4)

	for _, f := range fields {
		v, err := strconv.ParseInt(f, 10, 32)
		if err != nil {
			panic(err)
		}

		result = append(result, int(v))
	}

	return result
}

var Connections = make([]*webrtc.RTCPeerConnection, 0, 1)
var DCs = make([]*webrtc.RTCDataChannel, 0, 1)
var ptmx *os.File

func main() {
	addr := flag.String("address", ":50000", "Address to host the HTTP server on.")
	Cmd := flag.String("cmd", "ssh 10.71.1.3 -l tsl", "command to execute.")
	flag.Parse()

	// Everything below is the pion-WebRTC API! Thanks for using it ❤️.

	// Exchange the offer/answer via HTTP
	offerChan, answerChan := mustSignalViaHTTP(*addr)
	// Wait for the remote SessionDescription
	go func(offerOut chan webrtc.RTCSessionDescription, answerIn chan webrtc.RTCSessionDescription) {
		for offer := range offerChan {

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

			// Set the handler for ICE connection state
			// This will notify you when the peer has connected/disconnected
			peerConnection.OnICEConnectionStateChange(func(connectionState ice.ConnectionState) {
				fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
			})

			// Register data channel creation handling
			peerConnection.OnDataChannel(func(d *webrtc.RTCDataChannel) {
				fmt.Printf("New DataChannel %s %d\n", d.Label, d.ID)

				DCs = append(DCs, d)
				// Register channel opening handling
				d.OnOpen(func() {
					fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label, d.ID)
				})

				// Register message handling
				d.OnMessage(func(payload datachannel.Payload) {
					switch p := payload.(type) {
					case *datachannel.PayloadString:
						var ctrl Ctrl
						err := json.Unmarshal(p.Data, &ctrl)
						if err != nil {
							panic(err)
						}

						size := ParseTerminalCtrl(ctrl.Value)
						ws, err := pty.GetsizeFull(ptmx)
						if err != nil {
							panic(err)
						}
						ws.Rows = uint16(size[1])
						ws.Cols = uint16(size[2])

						if len(size) >= 5 {
							ws.X = uint16(size[3])
							ws.Y = uint16(size[4])
						}

						if err := pty.Setsize(ptmx, ws); err != nil {
							panic(err)
						}
						return

					case *datachannel.PayloadBinary:
						_, err := ptmx.Write(p.Data)
						if err != nil {
							log.Println(err)
						}
					default:
						fmt.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), d.Label)
					}
				})
			})

			// Sets the LocalDescription, and starts our UDP listeners
			err = peerConnection.SetRemoteDescription(offer)
			util.Check(err)

			answer, err := peerConnection.CreateAnswer(nil)
			util.Check(err)

			// Send the answer
			answerChan <- answer

			Connections = append(Connections, peerConnection)
		}
	}(offerChan, answerChan)

	go func() {
		cmds := strings.Split(*Cmd, " ")
		cmd := exec.Command(cmds[0], cmds[1:]...)
		var err error
		ptmx, err = pty.Start(cmd)
		if err != nil {
			log.Println(err)
			return
		}

		go func() {
			if _, err = io.Copy(ptmx, os.Stdin); err != nil {
				log.Println(err)
			}
		}()
		buf := make([]byte, 8192)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					log.Println(err)
				}
				panic(err)
			}

			/* Send to local console. */
			if _, err = os.Stdout.Write(buf[0:n]); err != nil {
				log.Println(err)
				panic(err)
			}

			for _, dc := range DCs {
				if err = dc.Send(datachannel.PayloadBinary{Data: buf[0:n]}); err != nil {
					log.Println(err)
					panic(err)
				}
			}
		}
	}()

	// Block forever
	select {}
}

// mustSignalViaHTTP exchange the SDP offer and answer using an HTTP server.
func mustSignalViaHTTP(address string) (offerOut chan webrtc.RTCSessionDescription, answerIn chan webrtc.RTCSessionDescription) {
	offerOut = make(chan webrtc.RTCSessionDescription)
	answerIn = make(chan webrtc.RTCSessionDescription)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var offer webrtc.RTCSessionDescription
		err := json.NewDecoder(r.Body).Decode(&offer)
		util.Check(err)

		offerOut <- offer
		answer := <-answerIn

		err = json.NewEncoder(w).Encode(answer)
		util.Check(err)

	})

	go http.ListenAndServe(address, nil)
	fmt.Println("Listening on", address)

	return
}
