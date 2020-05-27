package main

import (
	"fmt"
	"log"
	"n2x"
	"strings"
	"time"
)

func main() {
	n, err := n2x.New("10.71.20.231", "9001", n2x.DEFAULTSESSIONNAME)
	if err != nil {
		panic(err)
	}

	sess := n.Session

	err = n.ReservePort("103/1")
	if err != nil {
		panic(err)
	}

	err = n.ReservePort("103/2")
	if err != nil {
		panic(err)
	}

	for _, port := range sess.Ports {
		ips, _ := port.LegacyLinkGetSutIPAddresses()
		fmt.Printf("%q ", ips)
	}

	err = sess.StopTest()
	if err != nil {
		panic(err)
	}

	ports, err := sess.GetPorts()
	fmt.Printf("%q ", ports)
	//	time.Sleep(time.Duration(time.Second * 40))
	for i, port := range ports {
		err := port.RemoveAllTraffics()
		if err != nil {
			panic(err)
		}

		t, err := port.AddTraffic("10", fmt.Sprintf("TR%d", i))
		if err != nil {
			panic(err)
		}

		t.SetAverageLoad(n2x.PACKETS_PER_SEC, "100")
		fmt.Printf("%+v\n", t)

		_, err = t.ListStreamGroups()
		if err != nil {
			panic(err)
		}

		sgs, err := t.GetAllStreamGroups()
		if err != nil {
			panic(err)
		}

		err = sgs[0].SetIPv4UDP()
		if err != nil {
			panic(err)
		}
		err = sgs[1].SetIPv4TCP()
		if err != nil {
			panic(err)
		}

		err = sgs[2].SetIPv6TCP()
		if err != nil {
			panic(err)
		}
		err = sgs[3].SetIPv6UDP()
		if err != nil {
			panic(err)
		}

		for _, sg := range sgs {
			sg.Enable()

			protocols, err := sg.ListProtocolsInHeader()
			if err != nil {
				fmt.Println(err)
			}

			ps := strings.Split(protocols, " ")
			for _, p := range ps {
				line, _ := sg.ListProtocolFieldsInHeader(p)
				fields := strings.Split(line, " ")
				for _, f := range fields {
					sg.GetFieldFixedValue(p, f)
				}

				line, _ = sg.ListOptionalFields(p)
				fields = strings.Split(line, " ")
				for _, f := range fields {
					if f == "vlan_tag1" {
						sg.EnableOptionalField(p, f)
					}
					sg.SetFieldIncrementingValueRange("ipv6", "source_address", "20", "2001:db8:1000::", "1000", "22")
					sg.SetFieldIncrementingValueRange("ipv6", "destination_address", "30", "2001:db8:1000::", "1000", "22")
				}

				err = sg.SetSourceMacAddress("00:01:01:00:01:01", "24", "100", "1")
				if err != nil {
					panic(err)
				}
				err = sg.SetDestinationMacAddress("00:02:02:00:02:02", "10", "100", "1")
				if err != nil {
					panic(err)
				}
			}
		}

	}

	err = sess.StartTest()
	if err != nil {
		panic(err)
	}

	tick := time.Tick(time.Duration(time.Second * 10))
	for _ = range tick {
		rt, tt, rr, tr := ports[0].GetStatistics()
		fmt.Println(ports[0].Name, rt, tt, rr, tr)
		rt, tt, rr, tr = ports[1].GetStatistics()
		fmt.Println(ports[1].Name, rt, tt, rr, tr)
	}
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
