package cproxy

import (
	"fmt"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
	"log"
	"time"
)

const (
	CtrlMessage = iota
	DataMessage
)

type CSession struct {
	Name      string
	ID        uint16
	dc        *webrtc.RTCDataChannel
	conn      *webrtc.RTCPeerConnection
	offer     *webrtc.RTCSessionDescription
	answer    *webrtc.RTCSessionDescription
	ctrlC     chan string
	mesgC     chan []byte
	ratelimit <-chan time.Time
}

func NewCSession(name string, remote *webrtc.RTCSessionDescription) (*CSession, error) {
	config := webrtc.RTCConfiguration{
		IceServers: []webrtc.RTCIceServer{
			{
				//	URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	wrc, err := webrtc.New(config)
	if err != nil {
		return nil, fmt.Errorf("Cannot create new session with: %s", err)
	}

	cs := &CSession{
		conn:      wrc,
		ctrlC:     make(chan string, 10),
		mesgC:     make(chan []byte, 1000),
		ratelimit: time.Tick(time.Millisecond * 5),
	}

	wrc.OnICEConnectionStateChange(func(connectionState ice.ConnectionState) {
		//log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	/*If off is not empty this means we are create a server, else we are create a client session.*/
	if remote != nil {
		wrc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
			cs.dc = dc
			dc.OnOpen(func() {
				cs.Name = dc.Label
				cs.ID = *dc.ID
				//log.Printf("Data channel '%s'-'%d' open.\n", dc.Label, dc.ID)
			})

			dc.OnMessage(func(data datachannel.Payload) {
				switch p := data.(type) {
				case *datachannel.PayloadString:
					cs.ctrlC <- string(p.Data)
				case *datachannel.PayloadBinary:
					cs.mesgC <- p.Data
				default:
					//log.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), dc.Label)
				}
			})
		})

		err = wrc.SetRemoteDescription(*remote)
		if err != nil {
			panic(err)
		}

		answer, err := wrc.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		cs.offer = remote
		cs.answer = &answer
	} else {
		dc, err := wrc.CreateDataChannel(name, nil)
		if err != nil {
			panic(err)
		}

		dc.OnOpen(func() {
			//log.Printf("'%s'-'%d' open.\n", dc.Label, dc.ID)
			cs.ctrlC <- "OPEN"
		})

		dc.OnMessage(func(data datachannel.Payload) {
			switch p := data.(type) {
			case *datachannel.PayloadString:
				cs.ctrlC <- string(p.Data)
			case *datachannel.PayloadBinary:
				cs.mesgC <- p.Data
			default:
				log.Printf("Message '%s' from DataChannel '%s' no payload \n", p.PayloadType().String(), dc.Label)
			}
		})

		offer, err := wrc.CreateOffer(nil)
		if err != nil {
			panic(err)
		}
		cs.dc = dc
		cs.offer = &offer
	}

	return cs, nil
}

func (cs *CSession) Recv() (int, []byte) {
	select {
	case msg := <-cs.ctrlC:
		return CtrlMessage, []byte(msg)
	case msg := <-cs.mesgC:
		return DataMessage, msg
	}
}

func (cs *CSession) Send(typ int, data []byte) {
	if cs.dc == nil {
		log.Println("Cannot send message, datachannel invalid")
		return
	}

	<-cs.ratelimit

	if typ == CtrlMessage {
		err := cs.dc.Send(datachannel.PayloadString{Data: data})
		if err != nil {
			panic(err)
		}
	} else if typ == DataMessage {
		err := cs.dc.Send(datachannel.PayloadBinary{Data: data})
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("Try to send message of unknown type: %d", typ)
	}
}

func (cs *CSession) SetRemote(rdp webrtc.RTCSessionDescription) error {
	return cs.conn.SetRemoteDescription(rdp)
}

func (cs *CSession) GetOffer() webrtc.RTCSessionDescription {
	return *cs.offer
}

func (cs *CSession) GetAnswer() webrtc.RTCSessionDescription {
	return *cs.answer
}

func (cs *CSession) SetAnswer(rdp webrtc.RTCSessionDescription) error {
	cs.answer = &rdp
	return cs.conn.SetRemoteDescription(rdp)
}

func (cs *CSession) SetOffer(rdp webrtc.RTCSessionDescription) error {
	cs.offer = &rdp
	return cs.conn.SetRemoteDescription(rdp)
}

func (cs *CSession) Destroy() {
	if cs.conn != nil {
		cs.conn.Close()
	}
}
