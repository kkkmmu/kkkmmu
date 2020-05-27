package configuration

import (
	//"log"
	"fileserver"
)

type Configuration struct {
	Device         string
	Protocol       string
	IP             string
	Port           string
	Username       string
	Password       string
	EnablePrompt   string
	LoginPrompt    string
	PasswordPrompt string
	BasePrompt     string
	Prompt         string
	ModeDB         map[string]string
	Name           string
	Hostname       string
	PromptHostname string
	SessionID      string
	SFU            string
	ActiveSFU      string
}

var DefaultFileServer = map[string]*fileserver.FileServer{
	"TSLF": &fileserver.FileServer{
		Protocol: "ftp",
		Name:     "TSL",
		IP:       "10.71.1.3",
		Port:     "23",
		Username: "tsl",
		Password: "tsl",
	},

	"TSLS": &fileserver.FileServer{
		Protocol: "ssh",
		Name:     "TSL",
		IP:       "10.71.1.3",
		Port:     "23",
		Username: "tsl",
		Password: "tsl",
	},

	"APPLE": &fileserver.FileServer{
		Protocol: "ssh",
		Name:     "APPLE",
		IP:       "10.55.2.65",
		Port:     "22",
		Username: "tsl",
		Password: "tsl",
	},
}

/*

ModeDB: map[string]string{
		"login":         "login",
		"password":      "Passowrd:",
		"enable":        "SWITCH>",
		"normal":        "SWITCH[A]#",
		"config":        "(config)",
		"config-vlan":   "(config-vlan)",
		"config-if":     "(config-if[",
		"config-dhcp":   "(config-dhcp[",
		"config-router": "(config-router)",
		"bridge":        "(bridge)",
		"shell":         "*SWITCH",
		"bcmshell":      "BCM.0>",
	},
*/

var Modes = []string{
	"login",
	"password",
	"enable",
	"normal",
	"config",
	"config-vlan",
	"config-if",
	"config-vrf",
	"config-dhcp",
	"config-mcec",
	"config-router",
	"config-route-map",
	"config-flow",
	"config-policer",
	"config-policy",
	"config-cmap-qos",
	"config-pmap-qos",
	"config-pmap-c-qos",
	"bridge",
	"shell",
	"bcmshell",
	"rtkshell",
	"gdb",
	"cint",
}

var DefaultHostName = "SWITCH"
var DefaultBasePrompt = "SWITCH"
var DefaultEnablePrompt = ">"
var DefaultLoginPrompt = "login"
var DefaultPasswordPrompt = "Password"
var DefaultSFU = "B"
var PromptEnd = "#"
var DefaultProtocol = "telnet"
var DefaultTelnetPort = "23"
var DefaultSshPort = "22"

func BuildModeDBFromHostNameAndBasePrompt(host, base string) map[string]string {
	if host == "" {
		host = DefaultHostName
	}
	if base == "" {
		base = DefaultBasePrompt
	}

	db := make(map[string]string, len(Modes))
	for _, m := range Modes {
		if m == "login" {
			db[m] = "login"
		} else if m == "enable" {
			db[m] = base + ">"
		} else if m == "passowrd" {
			db[m] = "Password:"
		} else if m == "shell" {
			db[m] = "*" + host + "#"
		} else if m == "bcmshell" {
			db[m] = "BCM.0>"
		} else if m == "normal" {
			db[m] = base + "#"
		} else if m == "config" {
			db[m] = base + "(config)"
		} else if m == "config-vlan" {
			db[m] = base + "(config-vlan"
		} else if m == "config-if" {
			db[m] = base + "(config-if"
		} else if m == "config-dhcp" {
			db[m] = base + "(config-dhcp"
		} else if m == "config-mcec" {
			db[m] = base + "(config-mcec"
		} else if m == "config-router" {
			db[m] = base + "(config-router"
		} else if m == "config-flow" {
			db[m] = base + "(config-flow"
		} else if m == "config-policer" {
			db[m] = base + "(config-policer"
		} else if m == "config-cmap-qos" {
			db[m] = base + "(config-cmap-qos"
		} else if m == "config-pmap-qos" {
			db[m] = base + "(config-pmap-qos"
		} else if m == "config-pmap-c-qos" {
			db[m] = base + "(config-pmap-c-qos"
		} else if m == "config-policy" {
			db[m] = base + "(config-policy"
		} else if m == "config-vrf" {
			db[m] = base + "(config-vrf"
		} else if m == "bridge" {
			db[m] = base + "(bridge)"
		} else if m == "config-route-map" {
			db[m] = base + "(config-route-map)"
		} else if m == "rtkshell" {
			db[m] = "RTK.0>"
		} else if m == "gdb" {
			db[m] = "(gdb)"
		} else if m == "cint" {
			db[m] = "cint>"
		}
	}

	//log.Printf("Mode DB: %q", db)
	return db
}
