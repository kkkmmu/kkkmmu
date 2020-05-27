package main

import (
	"logd"
)

func main() {
	ld := logd.New("10.71.1.35")
	ld.Start()
}
