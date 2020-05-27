package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	runCommand("ssh", "10.71.1.3", "-l", "tsl")
}

func runCommand(cmdName string, arg ...string) {
	cmd := exec.Command(cmdName, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
<<<<<<< HEAD
=======
		fmt.Printf("Failed to start Ruby. %s\n", err.Error())
>>>>>>> 384c60945f80088ff2a569a6bbe0f7e34886563f
		os.Exit(1)
	}
}
