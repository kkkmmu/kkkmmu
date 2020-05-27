package acase

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	//"strings"
	"dut"
	"time"
)

var NoTaskFile = errors.New("No Task file exist")
var NoConfigFile = errors.New("No Config file exist")
var NoTopologyFile = errors.New("No Topology file exist")

var DB CDB

type ACase struct {
	Module       string `yaml:"Module"`
	Name         string `yaml:"Name"`
	Topology     *Topology
	Config       *Config
	Tasks        []Task
	EnableTester bool `yaml:"EnableTester"`
	Enable       bool `yaml:"Enable"`
	Case         bool `yaml:"Case"`
	valid        bool
	Path         string
	Pass         string
}

func GetCase(name string) (*ACase, error) {
	return DB.GetCaseByName(name)
}

func GetAllCases() ([]*ACase, error) {
	return DB.GetAllCases()
}

func (acs *ACase) GetDUTByName(name string) (*dut.DUT, error) {
	return acs.Topology.GetDUTByName(name)
}

func (acs *ACase) GetTesterByName(name string) (*Tester, error) {
	return acs.Topology.GetTesterByName(name)
}

func (acs *ACase) PrintHeader() {
	fmt.Printf("=====================================================================================================\n")
	fmt.Printf("           Running Case %s:                                      \n", acs.Name)
	fmt.Printf("                    Path         :  %s                           \n", acs.Path)
	fmt.Printf("                    Config       :  %s                           \n", acs.Config.Path)
	fmt.Printf("                    Topology     :  %s                           \n", acs.Topology.Path)
	fmt.Printf("                    TesterEnabled:  %b                           \n", acs.EnableTester)
	fmt.Printf("                    Config's Ifps:                            \n")
	for n, itf := range acs.Config.Ifps {
		fmt.Printf("                     ifp %10s : Full %20s : Short: %10s : Media: %10s\n", n, itf.FullName, itf.ShortName, itf.Media)
	}

	fmt.Printf("                   Topology's Ifps:                            \n")
	for n, itf := range acs.Topology.Ifps {
		fmt.Printf("                     ifp %10s : Full %20s : Short: %10s : Media: %10s\n", n, itf.FullName, itf.ShortName, itf.Media)
	}
}

func (acs *ACase) PrintFooter() {
	fmt.Printf("                    Result       :  %s                            \n", acs.Pass)
	fmt.Printf("=====================================================================================================\n")
}

func (acs *ACase) Run() error {
	if acs.Case {
		if acs.Enable {
			err := acs.Topology.ConnectDuts(acs.Config)
			if err != nil {
				return fmt.Errorf("Cannot run case %s when connect dut with %s", err)
			}

			if acs.EnableTester {
				err := acs.InitTester()
				if err != nil {
					return fmt.Errorf("Cannot run case %s with %s", acs.Name, err)
				}

				//	defer acs.DeInitTester()
			}
			acs.PrintHeader()
			defer acs.PrintFooter()
			acs.Pass = "FAIELD"
			for _, t := range acs.Tasks {
				if !t.Enable {
					fmt.Printf("Task %s of case %s is disabled, skipped\n", t.Name, acs.Name)
					continue
				}
				err := acs.RunTask(&t)
				if err != nil {
					return fmt.Errorf("Task %s in case %s has failed with %s", t.Name, acs.Name, err)
				}
			}
		} else {
			acs.PrintHeader()
			acs.Pass = "Disabled"
			acs.PrintFooter()
		}
	} else {
		acses, err := DB.GetSubCasesByName(acs.Name)
		if err != nil {
			return fmt.Errorf("Cannot run case %s with %s", acs.Name, err)
		}

		for _, cs := range acses {
			err = cs.Run()
			if err != nil {
				return fmt.Errorf("Run subcase %s of case %s failed with %s", cs.Name, acs.Name, err)
			}
		}
	}

	acs.Pass = "PASS"

	return nil
}

func (acs *ACase) RunTask(t *Task) error {
	for _, r := range t.Routines {
		res, err := acs.RunRoutine(&r)
		if err != nil {
			return fmt.Errorf("Run Task %s failed %s with %s", t.Name, r.API, err)
		}

		if r.Type == "A" && res == false {
			return fmt.Errorf("Assert %s failed on task %s", r.API, t.Name)
		}
	}

	return nil
}

func (acs *ACase) RunRoutine(r *Routine) (bool, error) {
	var dut *dut.DUT
	var tster *Tester

	dt, err := acs.GetDUTByName(r.Device)
	if err != nil {
		tst, err := acs.GetTesterByName(r.Device)
		if err != nil {
			return false, fmt.Errorf("Cannot find DUT/Tester %s in topology %s", r.Device, acs.Topology.Path)
		} else {
			tster = tst
		}
	} else {
		dut = dt
	}

	if r.Delay > 0 {
		delay := time.Tick(time.Second * time.Duration(r.Delay))
		<-delay
	}

	if dut != nil {
		if r.Type == "A" {
			return dut.Assert(r.API, r.Expect, r.Params...), nil
		} else if r.Type == "R" {
			dut.Call(r.API, r.Params...)
		} else {
			return false, fmt.Errorf("Unkown routine type: %s on DUT %s", r.Type, r.Device)
		}
	} else if tster != nil {
		if r.Type == "A" {
			return tster.Assert(r.API, r.Expect, r.Params...), nil
		} else if r.Type == "R" {
			tster.Call(r.API, r.Params...)
		} else {
			return false, fmt.Errorf("Unkown routine type: %s on DUT %s", r.Type, r.Device)
		}
	}

	return true, nil
}

func (acs *ACase) IsValid() bool {
	if !DB.IsCaseExist(acs.Name) {
		return false
	}

	return true
}

func (acs *ACase) Init() error {
	if !acs.Enable {
		return nil
	}

	if acs.Case {
		data, err := ioutil.ReadFile(acs.Path + "/tasks.yaml")
		if err != nil {
			return fmt.Errorf("Read task file failed with: %s", err)
		}

		err = yaml.Unmarshal(data, &acs.Tasks)
		if err != nil {
			return err
		}

		if acs.Topology != nil && acs.Config != nil {
			err := acs.Topology.Init(acs.Config)
			if err != nil {
				return fmt.Errorf("Topology %s init failed with config %s with: %s", acs.Topology.Path, acs.Config.Path, err)
			}
		} else if acs.Topology != nil && acs.Config == nil {
			return fmt.Errorf("Case %s has no config and topology exist", acs.Path)
		} else if acs.Topology != nil {
			return fmt.Errorf("Case %s has no topology exist", acs.Path)
		} else if acs.Config == nil {
			return fmt.Errorf("Case %s has no config exist", acs.Path)
		}
	}

	return nil
}

func (ac *ACase) Verify() error {
	return nil
}

func (ac *ACase) InitTester() error {
	return ac.Topology.InitTester()
}

func (ac *ACase) DeInitTester() error {
	return ac.Topology.DeInitTester()
}

func GetCaseFromFile(fname string) ([]*ACase, error) {
	CS := make([]*ACase, 0, 10)
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("Cannot get cases from file %s with %s", fname, err)
	}

	err = yaml.Unmarshal(data, &CS)
	if err != nil {
		return nil, fmt.Errorf("Cannot get cases from file %s with %s", fname, err)
	}

	for _, c := range CS {
		c.Path = filepath.Dir(fname) + "/" + c.Name
	}

	return CS, nil
}

func GetConfigFromFile(fname string) (*Config, error) {
	var con Config
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("Cannot get config from file %s with %s", fname, err)
	}

	err = yaml.Unmarshal(data, &con)
	if err != nil {
		return nil, fmt.Errorf("Cannot get config from file %s with %s", fname, err)
	}

	con.Path = filepath.Dir(fname)

	return &con, nil
}

func GetTopologyFromFile(fname string) (*Topology, error) {
	var tp Topology
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("Cannot get topology from file %s with %s", fname, err)
	}

	err = yaml.Unmarshal(data, &tp)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse topology from file %s with %s", fname, err)
	}

	tp.Path = filepath.Dir(fname)
	return &tp, nil
}

func init() {
	var files []string
	filepath.Walk("asset/cases", func(path string, info os.FileInfo, err error) error {
		if info.Name() == "all.yaml" {
			files = append(files, path)
		}
		return nil
	})

	var cases = make([]*ACase, 0, 10)
	for _, file := range files {
		css, err := GetCaseFromFile(file)
		if err != nil {
			panic(err)
		}
		cases = append(cases, css...)
	}

	/* Build the case DB */
	for _, c := range cases {
		DB.AddCase(c)
	}

	var tfiles []string
	filepath.Walk("asset/cases", func(path string, info os.FileInfo, err error) error {
		if info.Name() == "topology.yaml" {
			tfiles = append(tfiles, path)
		}
		return nil
	})

	/* Init Topology configure for each case */
	tps := make([]*Topology, 0, len(tfiles))
	for _, file := range tfiles {
		tp, err := GetTopologyFromFile(file)
		if err != nil {
			panic(err)
		}

		tps = append(tps, tp)
	}

	for _, t := range tps {
		DB.AddTopology(t)
	}

	/* Init Configuration for each case */
	var cfiles []string
	filepath.Walk("asset/cases", func(path string, info os.FileInfo, err error) error {
		if info.Name() == "config.yaml" {
			cfiles = append(cfiles, path)
		}
		return nil
	})

	cs := make([]*Config, 0, len(cfiles))
	for _, file := range cfiles {
		con, err := GetConfigFromFile(file)
		if err != nil {
			panic(err)
		}
		cs = append(cs, con)
	}

	for _, c := range cs {
		DB.AddConfig(c)
	}

	err := DB.Init()
	if err != nil {
		panic(err)
	}
}
