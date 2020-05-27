package main

import (
	"sftp"
)

func main() {
	cli, err := sftp.NewClient("10.55.2.202", "22", "liwei", "Lee123!@#")
	if err != nil {
		panic(err)
	}

	cli.Put("switch.go", "dump.go")
	//cli.Get("hsl.unstrip", "hsl.unstrip")
}
