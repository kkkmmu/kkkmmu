package main

import (
	"atscase"
	"fmt"
)

func main() {
	ats, err := atscase.Get("mlag")
	if err != nil {
		panic(err)
	}

	err = ats.Init()
	if err != nil {
		fmt.Println(err)
	}

	res := ats.Run()
	if res != true {
		fmt.Printf("Run %s failed\n", ats.Name)
	}
}
