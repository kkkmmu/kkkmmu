package main

import (
	"flag"
	//	"fmt"
	"github.com/kr/pty"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var RIP = regexp.MustCompile("\\((?P<ip>[[:xdigit:]]{1,8})/")

func ParseLog() {
	data, err := ioutil.ReadFile("log.txt")
	if err != nil {
		panic(err)
	}

	os.Remove("plog.txt")
	sdata := strings.Replace(string(data), "\000", "", -1)
	sdata = strings.Replace(sdata, "\r", "\n", -1)
	lines := strings.Split(string(sdata), "\n")
	for _, line := range lines {
		if strings.Contains(line, "HSL") && !strings.Contains(line, "Event") {
			line = strings.Replace(line, "PRAS[1998]: HSL-ERR In function ", "", -1)
			line = strings.Replace(line, "hsl_bcm_prefix_delete ", "", -1)
			line = strings.Replace(line, "hsl_bcm_prefix_add ", "", -1)
			line = strings.Replace(line, " DEFIP entry ", "", -1)
			matches := RIP.FindAllStringSubmatch(line, -1)
			if len(matches) == 0 {
				continue
			}

			for i := 0; i < len(matches); i++ {
				ip := FixIPv4Address(matches[i][1])
				line = strings.Replace(line, matches[i][1], ip.String(), -1)
			}
			//fmt.Printf("%+v, %d\n", matches, len(matches))
			AppendToFile("plog.txt", []byte(line+"\n"))
		}
	}
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

func main() {
	Cmd := flag.String("cmd", "telnet 10.42.13.13", "command to execute.")
	Parse := flag.Bool("parse", false, "parser log.")

	flag.Parse()
	if *Parse {
		ParseLog()
		return
	}

	os.Remove("log.txt")

	cmds := strings.Split(*Cmd, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var err error
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		if _, err = io.Copy(ptmx, os.Stdin); err != nil {
			log.Println(err)
		}
	}()
	ptmx.Write([]byte("admin\n"))
	time.Sleep(time.Second)
	ptmx.Write([]byte("test1234A\n"))
	time.Sleep(time.Second)
	ptmx.Write([]byte("enable\n"))
	time.Sleep(time.Second)
	ptmx.Write([]byte("terminal monitor\n"))
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

		AppendToFile("log.txt", buf)

		/* Send to local console. */
		if _, err = os.Stdout.Write(buf[0:n]); err != nil {
			log.Println(err)
			panic(err)
		}
	}
}

func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}
