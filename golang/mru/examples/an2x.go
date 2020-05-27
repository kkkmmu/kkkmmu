package main

import (
	"fmt"
	"n2x"
	"time"
)

func main() {
	nn, err := n2x.New("10.71.20.231", "9001", "LIWEI_TEST")
	if err != nil {
		panic(err)
	}

	err = nn.ReservePort("103/1")
	if err != nil {
		panic(err)
	}

	err = nn.SetPortMediaType("103/1", n2x.MEDIA_SFP)
	if err != nil {
		panic(err)
	}

	err = nn.ReservePort("103/2")
	if err != nil {
		panic(err)
	}

	err = nn.SetPortMediaType("103/2", n2x.MEDIA_SFP)
	if err != nil {
		panic(err)
	}

	err = nn.SetPortLegacyDUTIP("103/1", "1.1.1.2")
	if err != nil {
		panic(err)
	}

	err = nn.AddPortLegacyHost("103/1", "0", "1.1.1.4", "24", "00:01:00:00:00:01")
	if err != nil {
		panic(err)
	}

	err = nn.SetPortLegacyDUTIP6("103/1", "2001:0db8:1030::254")
	if err != nil {
		panic(err)
	}

	err = nn.AddPortLegacyHost6("103/1", "0", "2001:0db8:1030::1", "64", "00:01:00:00:00:01")
	if err != nil {
		panic(err)
	}

	err = nn.SetPortLegacyDUTIP("103/2", "2.1.1.2")
	if err != nil {
		panic(err)
	}

	err = nn.AddPortLegacyHosts("103/2", "0", "2.1.1.4", "24", "00:01:00:00:00:01", "100")
	if err != nil {
		panic(err)
	}

	err = nn.SetPortLegacyDUTIP6("103/2", "2001:0db8:1031::254")
	if err != nil {
		panic(err)
	}

	err = nn.AddPortLegacyHosts6("103/2", "0", "2001:0db8:1031::128", "64", "00:01:00:00:00:01", "100")
	if err != nil {
		panic(err)
	}

	_, err = nn.AddPortOSPF("103/1", "0.0.0.0", "1.1.1.1", "1.1.1.2")
	if err != nil {
		panic(err)
	}

	_, err = nn.AddPortOSPF6("103/1", "0.0.0.0", "1.2.1.1", "1.2.1.2")
	if err != nil {
		panic(err)
	}

	oi5, err := nn.AddPortOSPF("103/2", "0.0.0.0", "2.1.1.1", "2.1.1.2")
	if err != nil {
		panic(err)
	}
	err = oi5.AddExternalRoute("4.0.0.0", "24", "10000", "1")
	if err != nil {
		panic(err)
	}

	oi6, err := nn.AddPortOSPF6("103/2", "0.0.0.0", "2.2.1.1", "2.2.1.2")
	if err != nil {
		panic(err)
	}

	err = oi6.AddExternalRoute6("4000::", "64", "10000", "1")
	if err != nil {
		panic(err)
	}

	_, err = nn.AddPortStreams("1031", "103/1", "103/2")
	if err != nil {
		panic(err)
	}

	/*
		err = nn.SetStreamsVLAN("1031", "1234", "10")
		if err != nil {
			panic(err)
		}
	*/

	err = nn.SetStreamsPPS("1031", "10000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsSrcIP("1031", "11.1.1.1", "24", "1", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsDstIP("1031", "111.1.1.1", "24", "1", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsSrcMAC("1031", "12:01:00:01:00:01", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsDstMAC("1031", "11:01:00:01:00:01", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsIPProtocol("1031", "3", "10")
	if err != nil {
		panic(err)
	}

	_, err = nn.AddPortStreams6("1032", "103/2", "103/1")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsMPS("1032", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsSrcIP6("1032", "2001::1", "64", "1", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsDstIP6("1032", "2002::1", "64", "1", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsDstMAC("1032", "12:01:00:01:00:01", "1000")
	if err != nil {
		panic(err)
	}
	err = nn.SetStreamsIPv6NextHeader("1032", "111", "10")
	if err != nil {
		panic(err)
	}

	_, err = nn.AddPortStreams6("10322", "103/2", "103/1")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsMPS("10322", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsSrcIP6("10322", "2001::1", "64", "1", "1000")
	if err != nil {
		panic(err)
	}
	err = nn.SetStreamsDstMAC("10322", "22:01:00:01:00:01", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsDstIP6("10322", "2002::1", "64", "1", "1000")
	if err != nil {
		panic(err)
	}

	err = nn.SetStreamsIPv6NextHeader("10322", "111", "10")
	if err != nil {
		panic(err)
	}

	err = nn.StartRoutingEngine()
	if err != nil {
		panic(err)
	}

	err = nn.StartTraffic()
	if err != nil {
		panic(err)
	}

	for {
		<-time.Tick(time.Duration(time.Second * 10))

		tx, rx, err := nn.GetPortStatistics("103/1")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("p"+"103/1", tx, rx)

		tx, rx, err = nn.GetPortStatistics("103/2")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("p"+"103/2", tx, rx)

		tx, rx, err = nn.GetStreamStatistics("1031")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("stream"+"1031", tx, rx)

		tx, rx, err = nn.GetStreamStatistics("1032")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("stream"+"1032", tx, rx)

		tx, rx, err = nn.GetStreamStatistics("10322")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("stream"+"10322", tx, rx)
	}

	/*
		err = nn.DeletePortAllOSPFs("103/1")
		if err != nil {
			panic(err)
		}
		err = nn.DeletePortAllOSPF6s("103/1")
		if err != nil {
			panic(err)
		}
		err = nn.DeletePortAllOSPFs("103/2")
		if err != nil {
			panic(err)
		}
		err = nn.DeletePortAllOSPF6s("103/2")
		if err != nil {
			panic(err)
		}

			err = nn.DelPortLegacyAllHosts("103/1")
			if err != nil {
				panic(err)
			}

			err = nn.DelPortLegacyAllHosts("103/2")
			if err != nil {
				panic(err)
			}

				err = nn.ReleasePort("103/1")
				if err != nil {
					panic(err)
				}

				err = nn.ReleasePort("103/2")
				if err != nil {
					panic(err)
				}

				err = nn.KillSessionByName("LIWEI_TEST")
				if err != nil {
					panic(err)
				}
	*/
}
