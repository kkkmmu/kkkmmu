package acase

import (
	"dut"
	"emulator"
	"fmt"
	"gopkg.in/yaml.v2"
	"ifp"
)

type Topology struct {
	//DUTS []dut.DUT `yaml:"DUTS"`
	DUTS        map[string]*dut.DUT `yaml:"DUTS"`
	Links       []Link              `yaml:"Links"`
	Tester      *Tester             `yaml:"Tester"`
	Path        string              `yaml:"Path"`
	Ifps        map[string]*ifp.Ifp `yaml:"Ifps"`
	Initialized bool
}

func (tp *Topology) Init(config *Config) error {
	/*8
	if tp.Initialized == true {
		return nil
	}
	*/

	if tp.Ifps == nil {
		tp.Ifps = make(map[string]*ifp.Ifp, 1)
	}

	if len(tp.DUTS) != len(config.Devices) {
		return fmt.Errorf("Topology %s need %d device, config %s include %s", tp.Path, len(tp.DUTS), config.Path, len(config.Devices))
	}

	/* DUT Basic structure set up. */
	for _, d := range tp.DUTS {
		err := d.Init()
		if err != nil {
			return fmt.Errorf("Topolgy %s init failed with: %s", tp.Path, err)
		}

	}

	/* Tester Basic structure setup */
	err := tp.Tester.Verify(config.Ifps)
	if err != nil {
		return fmt.Errorf("Tester %s init failed with: %s", tp.Tester.Name, err)
	}

	for _, d := range tp.DUTS {
		itfs := d.GetAllInterface()

		for _, itf := range itfs {
			if _, ok := tp.Ifps[itf.Name]; ok {
				return fmt.Errorf("Duplicate interface %s exist in toplogy %s", itf.Name, tp.Path)
			}
			tp.Ifps[itf.Name] = itf
		}
	}

	itfs := tp.Tester.GetAllInterface()
	for _, itf := range itfs {
		if _, ok := tp.Ifps[itf.Name]; ok {
			return fmt.Errorf("Duplicate interface %s exist in toplogy %s", itf.Name, tp.Path)
		}
		tp.Ifps[itf.Name] = itf
	}

	for _, itf1 := range config.Ifps {
		itf2, ok := tp.Ifps[itf1.Name]
		if !ok {
			return fmt.Errorf("There is no interface %s in toplology %s", itf1.Name, tp.Path)
		}

		itf2.FullName = itf1.FullName
		itf2.ShortName = itf1.ShortName
	}

	/* SetUp Tester */
	if tp.Tester.Name != config.Emulator.Name {
		return fmt.Errorf("Topology init failed there is no %s in topology: %s", config.Emulator.Name, tp.Path)
	}

	if config.Emulator.Device != emulator.N2X.Device {
		return fmt.Errorf("Unspported emulator %s, only support %s", config.Emulator.Device, emulator.N2X.Device)
	}

	tp.Tester.Emulator = &emulator.EN2X{
		Name:   config.Emulator.Name,
		Device: config.Emulator.Device,
		IP:     config.Emulator.IP,
		Port:   config.Emulator.Port,
	}

	/* This should be run when start runing a case.
	err = tp.Tester.Init()
	if err != nil {
		return fmt.Errorf("Cannot init test %s with %s", tp.Tester.Name, err)
	}
	*/

	tp.Initialized = true

	return nil
}

func (tp *Topology) ConnectDuts(config *Config) error {
	for dn, dt := range tp.DUTS {
		dev, ok := config.Devices[dn]
		if !ok {
			return fmt.Errorf("DUT %s is not exist in config %s", dn, config.Path)
		}

		err := dt.SetDevice(dev)
		if err != nil {
			return fmt.Errorf("Topology %s init failed with %s", tp.Path, err)
		}

		err = dt.Connect()
		if err != nil {
			return fmt.Errorf("Topoloy %s init failed when connect to dut with %s", tp.Path, err)
		}
	}

	return nil
}

func (tp *Topology) Copy() *Topology {
	data, err := yaml.Marshal(tp)
	if err != nil {
		panic(err)
	}

	var ntp Topology

	err = yaml.Unmarshal(data, &ntp)
	if err != nil {
		panic(err)
	}

	return &ntp
}

func (tp *Topology) InitTester() error {
	return tp.Tester.Init()
}

func (tp *Topology) DeInitTester() error {
	return tp.Tester.DeInit()
}

func (tp *Topology) GetDUTByName(name string) (*dut.DUT, error) {
	if dt, ok := tp.DUTS[name]; ok {
		return dt, nil
	}

	return nil, fmt.Errorf("There is no dut name %s in topology %s", name, tp.Path)
}

func (tp *Topology) GetTesterByName(name string) (*Tester, error) {
	if tp.Tester.Name == name {
		return tp.Tester, nil
	}

	return nil, fmt.Errorf("There is no tester name %s in topology %s", name, tp.Path)
}

func (tp *Topology) GetInterfaceByName(name string) (*ifp.Ifp, error) {
	if itf, ok := tp.Ifps[name]; ok {
		return itf, nil
	}

	return nil, fmt.Errorf("There is no interface name %s in topology %s", name, tp.Path)
}
