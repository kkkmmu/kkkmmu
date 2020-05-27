package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"rut"
	"strconv"
	"strings"
	"util"
)

var RPrefix = regexp.MustCompile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}/[[:digit:]]{1,2}")
var RIP = regexp.MustCompile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")

var IP = flag.String("ip", "10.42.23.21", "IP address of the remote device")
var Host = flag.String("hostname", "DUT01", "Host name of the remote device")
var User = flag.String("username", "admin", "Username of the remote device")
var Password = flag.String("password", "krcho", "Passwrod of the remote device")

func main() {
	flag.Parse()
	dev, err := rut.New(&rut.RUT{
		Name:     *Host,
		Device:   "V5",
		IP:       *IP,
		Port:     "23",
		Username: *User,
		Hostname: *Host,
		Password: *Password,
	})

	if err != nil {
		panic(err)
	}

	dev.Init()

	ctx := context.Background()

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "show ip route",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fmt.Println(string(data))
	util.SaveToFile("rib.txt", []byte(data))
	ParsePrefix("rib", data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "show hsl prefix-table",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//fmt.Println(string(data))
	util.SaveToFile("hsl.txt", []byte(data))
	ParsePrefix("hsl", data)

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "@ent bcm l3 defip show",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ParsePrefix("chip", data)
	util.SaveToFile("chip.txt", []byte(data))
	//fmt.Println(string(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "@ent bcm d chg ing_l3_next_hop",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	util.SaveToFile("ing_l3_next_hop.txt", []byte(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "@ent bcm d chg egr_l3_next_hop",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	util.SaveToFile("egr_l3_next_hop.txt", []byte(data))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "@ent bcm d chg L3_DEFIP_ALPM_IPV4",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	util.SaveToFile("L3_DEFIP_ALPM_IPV4", []byte(data))

	FindErrorEntry(ctx, dev)
}

func ParsePrefix(name, data string) {
	os.Remove("./" + *Host + name + "_prefix.txt")
	ps := RPrefix.FindAllStringSubmatch(data, -1)
	for _, p := range ps {
		util.AppendToFile("./"+*Host+name+"_prefix.txt", []byte(p[0]+"\n"))
	}
}

var DR = regexp.MustCompile("DESTINATION=(?P<dest>[[:xdigit:]x]+),DATA=(?P<data>[[:xdigit:]x]+),")
var LK = regexp.MustCompile("LENGTH=(?P<length>[[:xdigit:]x]+),KEY=(?P<key>[[:xdigit:]x]+),")
var ID = regexp.MustCompile("ipipe0\\[(?P<idx>[[:xdigit:]]+)\\]")

func FindErrorEntry(ctx context.Context, dev *rut.RUT) {
	var icount int
	var ecount int
	var ncount int
	invalid := make([]string, 0, 10)
	data, err := ioutil.ReadFile("L3_DEFIP_ALPM_IPV4")
	if err != nil {
		panic(err)
	}

	os.Remove("./" + "invalid_prefix.txt")
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

				var p string
				if len(aidm) > 1 {
					p = fmt.Sprintf("%-18s : %18s/%02d %s\n", aidm[1], FixIPv4Address(alkm[2]), length, match[1])
				} else {
					p = fmt.Sprintf("%-18s : %18s/%02d %s\n", "ID Parse failed", FixIPv4Address(alkm[2]), length, match[1])
				}
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
					ip := FixIPv4Address(lkm[2])
					out := color.New(color.FgGreen)
					p := out.Sprintf("%-18s : %18s/%02d %s\n", idm[1], FixIPv4Address(lkm[2]), length, match[1])

					fmt.Println("------------------------------------------------------------------------------------------------")
					fmt.Printf(p)
					util.AppendToFile("./"+"invalid_prefix.txt", []byte(p))

					GetChipHslCache(ctx, dev, ip.String(), int(length))
					ecount++
				}

			}
		} else {
			invalid = append(invalid, line)
		}
	}

	//fmt.Println("Invalid lines:")
	os.Remove("./" + "not_parsed_lines.txt")
	for _, in := range invalid {
		if !strings.Contains(in, "DUT") &&
			!strings.Contains(in, "ent") &&
			!strings.Contains(in, "done") &&
			!strings.Contains(in, "cmd") {
			ncount++
			util.AppendToFile("./"+"not_parsed_lines.txt", []byte(in+"\n"))
			//fmt.Println(in)
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

func GetChipHslCache(ctx context.Context, dev *rut.RUT, ip string, plen int) {
	rdata, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  fmt.Sprintf("@ent bcm l3 alpm find IP4=%s LENGTH=%d", ip, plen),
	})

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Println("  ALPM entry: ")
	ls := GetLinesContainPrefix(string(rdata))
	if len(ls) > 0 {
		for _, l := range ls {
			fmt.Printf("        %s\n", l)
		}
	}

	fmt.Println("  Hsl entry: ")
	rdata, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  fmt.Sprintf("show hsl prefix-table ipv4 %s/%d", ip, plen),
	})

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	ls = GetLinesContainPrefix(string(rdata))
	if len(ls) > 0 {
		for _, l := range ls {
			fmt.Printf("        %s\n", l)
		}
	}
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
