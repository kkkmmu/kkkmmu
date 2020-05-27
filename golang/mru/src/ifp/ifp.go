package ifp

import (
	"device"
	"emulator"
)

const (
	ATS_IF_TYPE_L2 = iota
	ATS_IF_TYPE_L3
)

type Ifp struct {
	Type      int                `yaml:"Type"`
	Index     int                `yaml:"Index"`
	Name      string             `yaml:"Name"`
	FullName  string             `yaml:"FullName"`
	Media     string             `yaml:"Media"`
	ShortName string             `yaml:"ShortName"`
	RealName  string             `yaml:"RealName"`
	IsUp      bool               `yaml:"IsUp"`
	IsRunning bool               `yaml:"IsRunning"`
	Enable    bool               `yaml:"Enable"`
	LinkTo    *Ifp               `yaml:"LinkTo"`
	Dev       *device.Device     `yaml:"Dev"`
	Emu       *emulator.Emulator `yaml:"Emu"`
}
