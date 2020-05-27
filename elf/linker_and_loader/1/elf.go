package main

import (
	"debug/elf"
	"fmt"
)

func main() {
	file, err := elf.Open("m")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", file.FileHeader)

	for _, prog := range file.Progs {
		fmt.Printf("%+v\n", prog)
	}

	for _, section := range file.Sections {
		fmt.Printf("%+v\n", section)
	}
}
