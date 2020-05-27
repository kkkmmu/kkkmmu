package main

import (
	"fmt"
	"log"
	"n2x"
	"strings"
	//	"time"
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

		fmt.Printf("%+v\n", sgs)

		for _, sg := range sgs {
			sg.Disable()
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

		err = sgs[4].SetIPv6ND()
		if err != nil {
			panic(err)
		}
		err = sgs[5].SetIPv4ARP()
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
				sg.SetIPv4Tos("1", "2", "1")
				sg.SetIPv6DSCP("1", "20", "1")
				sg.SetIPv4DSCP("1", "20", "1")
				line, _ = sg.ListProtocolFieldsInHeader(p)
				err = sg.SetVlan("100", "1", "1")
				if err != nil {
					panic(err)
				}

				err = sg.SetCos("4", "1", "1")
				if err != nil {
					panic(err)
				}

				err = sg.SetCos("4", "1", "1")
				if err != nil {
					panic(err)
				}

				err = sg.UnsetVlan()
				if err != nil {
					panic(err)
				}

				sg.SetIPv4Protocol("4", "100", "1")
			}
		}

	}

	err = sess.StartTest()
	if err != nil {
		panic(err)
	}
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
