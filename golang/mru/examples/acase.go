package main

import (
	"acase"
	//	"fmt"
)

func main() {
	nac, err := acase.GetCase("MLAG_SINGLE_DU")
	//nac, err := acase.GetCase("OSPF")
	if err != nil {
		panic(err)
	}

	//	fmt.Println(dut.GetCurrentMode())

	err = nac.Run()
	if err != nil {
		panic(err)
	}

}
