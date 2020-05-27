package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ()

func main() {
	fmt.Println(os.Getwd())
	fmt.Println(filepath.Abs("."))
	//fmt.Println(os.Chdir("/home/kkkmmu"))
	top, err := GetTopDir()
	if err != nil {
		fmt.Println("Cannot find current project with ", err)
		os.Exit(-1)
	}
	os.Chdir(top)
}

func RunCommand(cmdName string, arg ...string) {
	cmd := exec.Command(cmdName, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func GetTopDir() (string, error) {
	cdr, err := os.Getwd()
	if err != nil || !strings.Contains(cdr, "kkkmmu") {
		return "", errors.New("Project must under work dir")
	}

	project, err := GetProject(cdr)
	if err != nil {
		return "", fmt.Errorf("Get top failed with %s", err)
	}

	return cdr[:strings.Index(cdr, project)+len(project)], nil
}

func GetProject(dir string) (string, error) {
	dirs := strings.Split(dir, "/")

	for i, dir := range dirs {
		if dir == "kkkmmu" && len(dirs) >= i+1 {
			return dirs[i+1], nil
		}
	}

	return "", errors.New("No project exist")
}
