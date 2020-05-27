package main

import (
	"fmt"
	"github.com/pkg/term"
	"github.com/pkg/term/termios"
	//"golang.org/x/crypto/ssh/terminal"

	"os"
	"os/signal"
	"syscall"
	"time"
)

type Term struct {
	name string
	fd   int
	orig syscall.Termios // original state of the terminal, see Open and Restore
}

func CBreakMode(t *Term) error {
	var a syscall.Termios
	if err := termios.Tcgetattr(uintptr(t.fd), (*syscall.Termios)(&a)); err != nil {
		return err
	}
	termios.Cfmakecbreak((*syscall.Termios)(&a))
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, (*syscall.Termios)(&a))
}

func main() {
	/*
		tos, _ := terminal.MakeRaw(int(os.Stdin.Fd()))

		defer terminal.Restore(int(os.Stdin.Fd()), tos)

	*/
	_, pts, err := termios.Pty()
	defer pts.Close()
	term, err := term.Open(pts.Name())
	if err != nil {
		panic(err)
	}
	fmt.Println(pts.Name())
	term.SetCbreak()
	defer term.Restore()
	/*
		name, err := termios.Ptsname(os.Stdin.Fd())
		if err != nil {
			panic(err)
		}
	*/

	sich := make(chan os.Signal)
	signal.Notify(sich, syscall.SIGKILL, syscall.SIGSTOP, syscall.SIGINT)

	go func() {
		for s := range sich {
			fmt.Println("Received signale")
			fmt.Printf("%#v\n", s)
			fmt.Printf("%s\n", s)
			os.Exit(0)
		}
	}()

	<-time.Tick(time.Second * 10000)
}

func init() {
}
