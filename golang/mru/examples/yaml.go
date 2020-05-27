package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//Nginx nginx  配置
type Nginx struct {
	Port    int    `yaml:"Port"`
	LogPath string `yaml:"LogPath"`
	Path    string `yaml:"Path"`
}

//Config   系统配置配置
type Config struct {
	Name      string `yaml:"SiteName"`
	Addr      string `yaml:"SiteAddr"`
	HTTPS     bool   `yaml:"Https"`
	SiteNginx Nginx  `yaml:"Nginx"`
	DUTS      []DUT  `yaml:"DUTS"`
	Test      string `yaml:"Test"`
	Tasks     []Task `yaml:"Tasks"`
}

type Topology struct {
	DUTS []DUT
}

type Routine struct {
	Name string `yaml:"Name"`
	ID   string `yaml:"ID"`
}

type Task struct {
	Name     string    `yaml:"Name"`
	Routines []Routine `yaml:"Routines"`
}

type DUT struct {
	Name string `yaml:"Name"`
	Ifps []Ifp  `yaml:"Ifps"`
}

type Ifp struct {
	Name string `yaml:"Name"`
}

func main() {

	var setting Config
	config, err := ioutil.ReadFile("./examples/yaml.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config, &setting)

	fmt.Println(setting.Name)
	fmt.Println(setting.Addr)
	fmt.Println(setting.HTTPS)
	fmt.Println(setting.SiteNginx.Port)
	fmt.Println(setting.SiteNginx.LogPath)
	fmt.Println(setting.SiteNginx.Path)
	fmt.Println(setting.Test)
	fmt.Printf("%+q\n", setting.DUTS)
	fmt.Printf("%+v\n", setting.Tasks)
}
