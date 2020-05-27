package main

import (
	"acase"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	var setting []acase.Task
	config, err := ioutil.ReadFile("./examples/ascript.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config, &setting)

	fmt.Printf("%+v", setting)
}
