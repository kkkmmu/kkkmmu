package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"rut"
	"strings"
)

var (
	Space  = regexp.MustCompile(`\s+`)
	Mutex  = regexp.MustCompile("sal_mutex_take \\(m=(?P<address>0x[[:xdigit:]]+),")
	LockR  = regexp.MustCompile(" = (?P<value>[[:alnum:]]+)")
	Params = regexp.MustCompile(`\(.*\)[[:space:]]+at`)
	From   = regexp.MustCompile(` from.*$`)
	String = regexp.MustCompile(`\".*\"`)
)

//$4 = {mutex = {__data = {__lock = 2, __count = 1, __owner = 1035, __kind = 1, __nusers = 1, {__spins = 0, __list = {__next = 0x0}}}, __size = "\002\000\000\000\001\000\000\000\v\004\000\000\001\000\000\000\001\000\000\000\000\000\000", __align = 2}, desc = 0xbefff8a8 ""}

var (
	IP       = flag.String("ip", "10.71.20.156", "IP address of the remote device")
	Server   = flag.String("server", "10.71.1.3", "IP address of file server")
	Port     = flag.String("port", "23", "Port to connect")
	Host     = flag.String("hostname", "SWITCH", "Host name of the remote device")
	Protocol = flag.String("prot", "telnet", "Login protocol(ssh/telnet)")
	User     = flag.String("u", "admin", "Username of the remote device")
	Password = flag.String("p", "Dasan123456", "Passwrod of the remote device")

	Process = flag.String("process", "hsl", "process name to debug")
	Bin     = flag.String("bin", "hsl", "binary file name path")

	CTX = context.Background()
)

type Frame struct {
	ID       string
	Ptr      string
	Function string
	Position string
}

type Lock struct {
	Name   string
	Ptr    string
	Owner  *Thread
	Lock   string
	Count  string
	Kind   string
	Nusers string
	Desc   string
}

func (l *Lock) String() string {
	//return fmt.Sprintf("[ Name: %10s:%20s, Ptr: %10s, Owner: %10s ]", l.Name, l.Desc, l.Ptr, l.Owner.Name)
	return fmt.Sprintf("[ Name: %20s, Owner: %10s ]", l.Desc, l.Owner.Name)
}

type Thread struct {
	Name        string
	ID          string
	PID         string
	SP          string
	Frame       string
	OwnedLock   map[string]*Lock
	WaitForLock map[string]*Lock
	Running     bool
	Main        bool
	WaitFor     *Thread
	Frames      []*Frame
}

func (t *Thread) String() string {
	var Out *color.Color
	if len(t.OwnedLock) != 0 && len(t.WaitForLock) != 0 {
		Out = color.New(color.FgRed)
	} else {
		Out = color.New(color.FgWhite)
	}

	var str string
	str = Out.Sprintf("{%16s %6s %6s  %10s", t.Name, t.ID, t.PID, t.SP)
	if len(t.OwnedLock) != 0 {
		str += Out.Sprintf("\n                            Hold (%s)\n", t.OwnedLockInfo())
	}

	if len(t.WaitForLock) != 0 {
		if len(t.OwnedLock) == 0 {
			str += Out.Sprintf("                                      \n")
		}
		str += Out.Sprintf("                            WaitFor %10s on (%20s)\n", t.WaitFor.Name, t.WaitForLockInfo())
	}

	str += Out.Sprintf("}")

	return str
}

func (t *Thread) LockInfo() string {
	var Out *color.Color
	if len(t.OwnedLock) != 0 || len(t.WaitForLock) != 0 {
		Out = color.New(color.FgRed)
	} else {
		Out = color.New(color.FgGreen)
	}

	var owned string
	var waited string

	if len(t.OwnedLock) != 0 {
		for _, ol := range t.OwnedLock {
			owned += Out.Sprintf(" (%10s:%10s)", ol.Name, ol.Desc)
		}
	}
	if len(t.WaitForLock) != 0 {
		for _, wl := range t.WaitForLock {
			waited += Out.Sprintf(" (%10s:%10s)", wl.Name, wl.Desc)
		}
	}
	return Out.Sprintf("{%16s %6s %6s  %10s Hold (%20s) WaitFor(%22s)", t.Name, t.ID, t.PID, t.SP, owned, waited)
}

func (t *Thread) OwnedLockInfo() string {

	var owned string

	if len(t.OwnedLock) != 0 {
		for _, ol := range t.OwnedLock {
			owned += fmt.Sprintf("(%10s:%10s)", ol.Name, ol.Desc)
		}
	}

	return owned
}

func (t *Thread) WaitForLockInfo() string {
	var waited string

	if len(t.WaitForLock) != 0 {
		for _, wl := range t.WaitForLock {
			waited += fmt.Sprintf("(%10s:%10s)", wl.Name, wl.Desc)
		}
	}

	return waited
}

func (t *Thread) FramesInfo() string {
	var frames string

	if len(t.Frames) != 0 {
		for _, f := range t.Frames {
			frames += fmt.Sprintf("                           %2s %-30s %-50s\n", f.ID, f.Function, f.Position)
		}
	}

	return frames
}

func main() {
	flag.Parse()
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		Protocol: *Protocol,
		IP:       *IP,
		Port:     *Port,
		Username: *User,
		Hostname: *Host,
		Password: *Password,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Start connect to %s!\n", *IP)
	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Connect to %s sucessfull!\n", *IP)

	ctx := context.Background()

	if dev.CurrentMode() == "normal" {
		_, err = dev.RunCommand(ctx, &command.Command{
			Mode: "normal",
			CMD:  "q sh -l",
		})

		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Start download %s from %s!\n", *Bin, *Server)
	if *Server == "10.71.1.3" {
		dev.Get("TSLS", *Bin)
	} else if *Server == "10.55.2.65" {
		dev.Get("APPLES", *Bin)
	} else {
		fmt.Println("Unkown server %s!\n", *Server)
		return
	}
	fmt.Printf("Download %s from %s finished!\n", *Bin, *Server)

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "which gdb",
	})

	if !strings.Contains(string(data), "bin") {
		fmt.Println("Cannot find gdb on this device ", *IP)
		return
	}

	fmt.Printf("Start analysis process %s with binary %s!\n", *Process, *Bin)
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  `ps -T -p $(pidof ` + *Process + `)`,
	})
	if err != nil {
		panic(err)
	}
	threads := GetAllThreadFromPS(data)
	fmt.Printf("Process %s has %d threads!\n", *Process, len(threads))

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  `gdb -p $(pidof ` + *Process + `) ` + *Bin,
	})
	if err != nil {
		panic(err)
	}

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "gdb",
		CMD:  "info threads",
	})
	if err != nil {
		panic(err)
	}

	UpdateThreadInfoByGDB(threads, data)

	for _, thread := range threads {
		if strings.Contains(thread.SP, "lll_lock_wait") {
			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "gdb",
				CMD:  "thread apply " + thread.ID + " bt",
			})
			if err != nil {
				panic(err)
			}
			thread.Frames = ParseStackTrace(data)
			match := Mutex.FindStringSubmatch(string(data))
			if len(match) == 2 {
				data, err = dev.RunCommand(ctx, &command.Command{
					Mode: "gdb",
					CMD:  "thread apply " + thread.ID + " print *(recursive_mutex_t *)" + match[1],
				})
				if err != nil {
					panic(err)
				}
				lmatch := LockR.FindAllStringSubmatch(string(data), -1)
				if len(lmatch) > 0 {
					data, err = dev.RunCommand(ctx, &command.Command{
						Mode: "gdb",
						CMD:  "x /s " + lmatch[8][1],
					})
					if err != nil {
						panic(err)
					}
					var desc string
					lines := strings.Split(data, "\n")
					for _, line := range lines {
						if strings.Contains(line, lmatch[8][1]) {
							desc = strings.Trim(line, lmatch[8][1]+":")
							desc = strings.TrimSpace(desc)
						}
					}
					nLock := &Lock{
						Name:   lmatch[8][1],
						Ptr:    match[1],
						Owner:  threads[lmatch[2][1]],
						Lock:   lmatch[0][1],
						Count:  lmatch[1][1],
						Kind:   lmatch[3][1],
						Nusers: lmatch[4][1],
						Desc:   desc,
					}
					thread.WaitForLock[nLock.Ptr] = nLock
					threads[lmatch[2][1]].OwnedLock[nLock.Ptr] = nLock
					thread.WaitFor = threads[lmatch[2][1]]
				}
			}
		}
	}

	fmt.Println("ThreadInfo")
	for _, thread := range threads {
		if len(thread.OwnedLock) != 0 && len(thread.WaitForLock) != 0 {
			continue
		}
		fmt.Println(thread)
	}

	for _, thread := range threads {
		if len(thread.OwnedLock) != 0 && len(thread.WaitForLock) != 0 {
			fmt.Println(thread)
		}
	}

	fmt.Println("BackTrace")
	for _, thread := range threads {
		if len(thread.OwnedLock) != 0 && len(thread.WaitForLock) != 0 {
			fmt.Println(thread)
			fmt.Println(thread.FramesInfo())
		}
	}
}

func GetAllThreadFromPS(data string) map[string]*Thread {
	threads := make(map[string]*Thread, 10)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = Space.ReplaceAllString(line, " ")
		if strings.Contains(line, "?") {
			main := false
			fields := strings.Split(line, " ")

			if fields[0] == fields[1] {
				main = true
			}

			threads[fields[1]] = &Thread{
				PID:         fields[1],
				Name:        fields[4],
				Main:        main,
				OwnedLock:   make(map[string]*Lock, 1),
				WaitForLock: make(map[string]*Lock, 1),
				Frames:      make([]*Frame, 0, 1),
			}
		}
	}

	return threads
}

func UpdateThreadInfoByGDB(db map[string]*Thread, data string) map[string]*Thread {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = String.ReplaceAllString(line, "")
		line = strings.Replace(line, "(", "", 1)
		line = strings.Replace(line, ")", "", 1)
		line = strings.Replace(line, "*", "", 1)
		line = strings.Replace(line, `"`, "", -1)
		line = strings.Replace(line, " at", "", -1)
		line = strings.Replace(line, " from", "", -1)
		line = Space.ReplaceAllString(line, " ")
		line = strings.TrimSpace(line)
		if strings.Contains(line, "LWP") {
			fields := strings.Split(line, " ")

			thread, ok := db[fields[4]]
			if !ok {
				panic("Cannot find thread " + fields[5])
			}

			thread.ID = fields[0]
			thread.SP = fields[7]
		}
	}

	return db
}

func ParseStackTrace(data string) []*Frame {
	frames := make([]*Frame, 0, 10)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		match := Params.FindStringSubmatch(line)
		if len(match) > 0 {
			line = Params.ReplaceAllString(line, "at")
			line = Space.ReplaceAllString(line, " ")
			line = From.ReplaceAllString(line, "")
			line = strings.Replace(line, "?", "", -1)
			line = strings.Replace(line, "#", "", -1)
			line = strings.TrimSpace(line)
			fields := strings.Split(line, " ")
			if len(fields) == 6 {
				frames = append(frames, &Frame{
					ID:       fields[0],
					Ptr:      fields[1],
					Function: fields[3],
					Position: fields[5],
				})
			}
		}
	}

	return frames
}

func GetWaitForPosition(data string) string {
	return ""
}
