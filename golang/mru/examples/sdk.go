package main

import (
	"telnet"
	//"util"
	"fmt"
	//"unicode"
	"context"
	"github.com/fatih/color"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"util"
)

var count int = 0

func main() {

	os.Remove("all.txt")
	ctx := context.Background()
	c, err := telnet.New("admin", "Dasan123456", "10.71.20.169", "23")
	if err != nil {
		panic(err)
	}

	c.WriteLine("enable")
	n, err := c.ReadUntil("#")
	fmt.Println(string(n))

	c.WriteLine("q sh -l")
	n, err = c.ReadUntil("#")
	fmt.Println(string(n))
	//n, err = c.ReadUntilSkip([]string{"/"}, []string{"adH/", "iB/"})

	c.WriteLine("telnet 127.0.0.1 10000\n")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	for {
		RunLPMTest(ctx, c)
		count++
	}
}

func RunLPMTest(ctx context.Context, c *telnet.Session) {
	util.AppendToFile("all.txt", []byte(fmt.Sprintf("===============================%d=================================", count)))
	c.WriteLine("cint /etc/.config/cint.cint")
	n, err := c.ReadUntil("BCM.0>")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(n))

	c.WriteLine("d nocache chg L3_DEFIP_ALPM_IPV4")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	util.SaveToFile("L3_DEFIP_ALPM_IPV4", []byte(n))
	util.AppendToFile("all.txt", []byte(n))

	c.WriteLine("vlan clear")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	c.WriteLine("l3 intf clear")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	c.WriteLine("l3 l3table clear")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	c.WriteLine("l3 defip clear")
	n, err = c.ReadUntil("BCM.0>")
	fmt.Println(string(n))

	FindErrorEntry(ctx, c)
}

var RPrefix = regexp.MustCompile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}/[[:digit:]]{1,2}")
var RIP = regexp.MustCompile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")

var DR = regexp.MustCompile("DESTINATION=(?P<dest>[[:xdigit:]x]+),DATA=(?P<data>[[:xdigit:]x]+),")
var LK = regexp.MustCompile("LENGTH=(?P<length>[[:xdigit:]x]+),KEY=(?P<key>[[:xdigit:]x]+),")
var ID = regexp.MustCompile("ipipe0\\[(?P<idx>[[:xdigit:]]+)\\]")

func FindErrorEntry(ctx context.Context, c *telnet.Session) {
	var icount int
	var ecount int
	var ncount int
	invalid := make([]string, 0, 10)
	data, err := ioutil.ReadFile("L3_DEFIP_ALPM_IPV4")
	if err != nil {
		panic(err)
	}

	os.Remove("./" + "all_prefix.txt")
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		match := DR.FindStringSubmatch(line)
		if len(match) == 3 {
			icount++
			aidm := ID.FindStringSubmatch(line)
			alkm := LK.FindStringSubmatch(line)
			if len(alkm) == 3 {
				length, err := strconv.ParseInt(alkm[1], 0, 64)
				if err != nil {
					panic(err)
				}

				out := color.New(color.FgWhite)
				var p string
				if len(aidm) > 1 {
					p = out.Sprintf("%-18s : %18s/%02d %s\n", aidm[1], FixIPv4Address(alkm[2]), length, match[1])
				} else {
					p = out.Sprintf("%-18s : %18s/%02d %s\n", "ID Parse failed", FixIPv4Address(alkm[2]), length, match[1])
				}
				fmt.Printf(p)
				util.AppendToFile("./"+"all_prefix.txt", []byte(p))
			}

			if !strings.HasPrefix(match[1], "0x") {
				idm := ID.FindStringSubmatch(line)
				lkm := LK.FindStringSubmatch(line)
				if len(lkm) == 3 {
					length, err := strconv.ParseInt(lkm[1], 0, 64)
					if err != nil {
						panic(err)
					}
					out := color.New(color.FgGreen)
					p := out.Sprintf("%-18s : %18s/%02d %s\n", idm[1], FixIPv4Address(lkm[2]), length, match[1])

					fmt.Println("------------------------------------------------------------------------------------------------")
					fmt.Printf(p)
					util.AppendToFile("./"+"invalid_prefix.txt", []byte(p))
					ecount++
					os.Exit(0)
				}
			}
		} else {
			invalid = append(invalid, line)
		}
	}

	out := color.New(color.FgRed)
	out.Printf("\n\nTotal entry count: %d, not parsed count: %d, parsed count: %d, error entry count: %d, \n", icount+ncount, ncount, icount, ecount)
}

func FixIPv4Address(s string) net.IP {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len(s) == 0 {
		s = "00000000"
	} else if len(s) == 1 {
		s = "0000000" + s
	} else if len(s) == 2 {
		s = "000000" + s
	} else if len(s) == 3 {
		s = "00000" + s
	} else if len(s) == 4 {
		s = "0000" + s
	} else if len(s) == 5 {
		s = "000" + s
	} else if len(s) == 6 {
		s = "00" + s
	} else if len(s) == 7 {
		s = "0" + s
	}

	f1, _ := strconv.ParseInt("0x"+s[:2], 0, 32)
	f2, _ := strconv.ParseInt("0x"+s[2:4], 0, 32)
	f3, _ := strconv.ParseInt("0x"+s[4:6], 0, 32)
	f4, _ := strconv.ParseInt("0x"+s[6:8], 0, 32)

	return net.IPv4(byte(f1), byte(f2), byte(f3), byte(f4))
}

func GetLinesContainPrefix(rdata string) []string {
	rs := make([]string, 0, 10)
	lines := strings.Split(rdata, "\n")
	for _, line := range lines {
		if strings.Contains(line, "cmd") ||
			strings.Contains(line, "show") ||
			strings.Contains(line, "ent") ||
			strings.Contains(line, "DUT") ||
			strings.Contains(line, "DUT") ||
			strings.Contains(line, "FIB") ||
			strings.Contains(line, "VRF") ||
			strings.Contains(line, "done") {
			continue
		}
		match := RIP.FindStringSubmatch(line)
		if len(match) > 0 {
			rs = append(rs, line)
		}
	}

	if len(rs) > 0 {
		return rs
	}

	return nil
}
