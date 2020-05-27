package main

import (
	"acase"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	var setting acase.Topology
	config, err := ioutil.ReadFile("./examples/atopology.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config, &setting)

	fmt.Printf("%+v", setting)
}
