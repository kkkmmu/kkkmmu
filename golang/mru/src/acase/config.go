package acase

import (
	"device"
	//"dut"
	"emulator"
	"ifp"
)

type Config struct {
	Devices  map[string]*device.Device `yaml:"Devices"`
	Emulator emulator.Emulator         `yaml:"Emulator"`
	Ifps     map[string]*ifp.Ifp       `yaml:"Ifps"`
	Path     string
}
