package main

import (
	"acase"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	var setting acase.Config
	config, err := ioutil.ReadFile("./examples/aconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config, &setting)

	fmt.Printf("%+q", setting)
}
