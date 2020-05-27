package atscase

import (
	"defaults"
	"device"
	"dut"
	"emulator"
	"encoding/json"
	"fmt"
	"ifp"
	"io/ioutil"
	"its"
	"link"
	"routine"
	"strings"
	"topology"
	"util"
)

type ATSCase struct {
	Name        string
	Topology    *topology.Topology
	prtasks     []*routine.Routine
	tasks       []*routine.Routine
	potasks     []*routine.Routine
	Result      bool
	lastTask    *routine.Routine
	initialized bool
}

func (atc *ATSCase) Run() bool {
	if !atc.initialized {
		fmt.Println("Cannot run test case: ", atc.Name, " not initialized!")
		return false
	}
	fmt.Println("Running test case: ", atc.Name)
	fmt.Println("Start Check Pre-Conditions")

	for _, t := range atc.prtasks {
		if t.Do() != true {
			atc.lastTask = t
			atc.Result = false
			return false
		}
	}

	fmt.Println("Running Test Case")
	for _, t := range atc.tasks {
		if t.Do() != true {
			atc.lastTask = t
			atc.Result = false
			return false
		}
	}

	fmt.Println("Start Check Post-Conditions")
	for _, t := range atc.potasks {
		if t.Do() != true {
			atc.lastTask = t
			atc.Result = false
			return false
		}
	}

	atc.Result = true

	return true
}

func (atc *ATSCase) isTopologyFileValid(data string) (bool, error) {
	if !strings.HasPrefix(data, defaults.TopologyFileMark) {
		return false, fmt.Errorf("Invalid topology file", atc.Name)
	}

	devs := defaults.Parsers[defaults.DUTMark].FindStringSubmatch(data)
	if len(devs) == 0 {
		return false, fmt.Errorf("Configure file should include the device section")
	}

	ifps := defaults.Parsers[defaults.LinkMark].FindStringSubmatch(data)
	if len(ifps) == 0 {
		return false, fmt.Errorf("Configure file should include the interface section")
	}

	return true, nil
}

func (atc *ATSCase) BuildTopology() error {

	if ok, _ := util.PathExists(defaults.GetTopologyFile(atc.Name)); !ok {
		return fmt.Errorf("Topology file %s for %s does not exist", defaults.GetTopologyFile(atc.Name), atc.Name)
	}

	data, err := ioutil.ReadFile(defaults.GetTopologyFile(atc.Name))
	if err != nil {
		return fmt.Errorf("Read topology file for %s failed with %s", atc.Name, err)
	}

	duts, err := atc.parseByMark(string(data), defaults.DUTMark)
	if err != nil {
		return fmt.Errorf("Parse topology for %s failed with %s", atc.Name, err)
	}

	for _, d := range duts {
		var newdut dut.DUT
		err := json.Unmarshal([]byte(d), &newdut)
		if err != nil {
			return fmt.Errorf("Cannot parse DUT %s with %s", d, err)
		}

		fmt.Printf("%+v\n", newdut)
		atc.Topology.AddDUT(&newdut)
	}

	testers, err := atc.parseByMark(string(data), defaults.TesterMark)
	if err != nil {
		return fmt.Errorf("Parse topology for %s failed with %s", atc.Name, err)
	}

	for _, t := range testers {
		var newits its.ITS
		err := json.Unmarshal([]byte(t), &newits)
		if err != nil {
			return fmt.Errorf("Cannot parse Tester %s with %s", t, err)
		}

		fmt.Printf("%+v\n", newits)
		atc.Topology.AddTester(&newits)
	}

	links, err := atc.parseByMark(string(data), defaults.LinkMark)
	if err != nil {
		return fmt.Errorf("Parse topology for %s failed with %s", atc.Name, err)
	}

	for _, l := range links {
		var newlink link.Link
		err := json.Unmarshal([]byte(l), &newlink)
		if err != nil {
			return fmt.Errorf("Cannot parse link %s with %s", l, err)
		}

		fmt.Printf("%+v\n", newlink)
		atc.Topology.AddLink(&newlink)
	}

	err = atc.Topology.Verify()
	if err != nil {
		return fmt.Errorf("Topology verfiy failed with: %s", err)
	}

	return nil
}

func (atc *ATSCase) isScriptFileFile(data string) bool {
	if !strings.HasPrefix(data, defaults.ScriptFileMark) {
		return false
	}

	return true
}

func (atc *ATSCase) BuildTestTasks(name string) error {
	if ok, _ := util.PathExists(defaults.GetScriptFile(atc.Name)); !ok {
		return fmt.Errorf("Test script file %s for %s exist", defaults.GetConfigureFile(atc.Name), atc.Name)
	}

	data, err := ioutil.ReadFile(defaults.GetScriptFile(atc.Name))
	if err != nil {
		return fmt.Errorf("Read Test script file %s failed with %s", defaults.GetConfigureFile(atc.Name), err)
	}

	lines := strings.Split(string(data), "\n")
	lines = lines[1:]
	fnp := defaults.Parsers[defaults.FunctionMark]
	asp := defaults.Parsers[defaults.AssertionMark]
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line != "" {
			if strings.HasPrefix(line, defaults.CommentMark) {
				continue
			}
			var nr routine.Routine

			matches := asp.FindStringSubmatch(line)
			if len(matches) == 5 {
				nr.Function = matches[2]
				params := strings.Split(matches[3], ",")
				nr.Paramerters = make([]string, 0, 2)
				for _, p := range params {
					nr.Paramerters = append(nr.Paramerters, strings.TrimSpace(p))
				}

				nr.Dut = atc.Topology.DUTs[matches[1]]
				nr.Expected = matches[4]
				nr.IsAssert = true

				fmt.Println(&nr)
				atc.tasks = append(atc.tasks, &nr)
				continue
			}

			/* Due to match issue, this must be after the previous one.*/
			matches = fnp.FindStringSubmatch(line)
			if len(matches) == 4 {
				nr.Function = matches[2]
				params := strings.Split(matches[3], ",")
				nr.Paramerters = make([]string, 0, 2)
				for _, p := range params {
					nr.Paramerters = append(nr.Paramerters, strings.TrimSpace(p))
				}

				nr.Dut = atc.Topology.DUTs[matches[1]]
				nr.IsAssert = false
				atc.tasks = append(atc.tasks, &nr)
				fmt.Println(&nr)
				continue
			}

			return fmt.Errorf("Invalid instruction %s in script", line)
		}
	}
	return nil
}

func (atc *ATSCase) isConfigureFileValid(data string) (bool, error) {
	if !strings.HasPrefix(data, defaults.ConfigureFileMark) {
		return false, fmt.Errorf("Invalid configure file", atc.Name)
	}

	devs := defaults.Parsers[defaults.DevMark].FindStringSubmatch(data)
	if len(devs) == 0 {
		return false, fmt.Errorf("Configure file should include the device section")
	}

	ifps := defaults.Parsers[defaults.IfpMark].FindStringSubmatch(data)
	if len(ifps) == 0 {
		return false, fmt.Errorf("Configure file should include the interface section")
	}

	emulators := defaults.Parsers[defaults.EmulatorMark].FindStringSubmatch(data)
	if len(emulators) == 0 {
		return false, fmt.Errorf("Configure file should include the emulator section")
	}
	/*
		others := defaults.Parsers[defaults.OtherMark].FindStringSubmatch(data)
		if len(others) == 0 {
			return false, fmt.Errorf("Configure file should include the other section")
		}
	*/

	return true, nil
}

func (atc *ATSCase) parseByMark(data, mark string) ([]string, error) {
	parser, ok := defaults.Parsers[mark]
	if !ok {
		return nil, fmt.Errorf("Unknown mark %s cannot be parsed")
	}

	matches := parser.FindStringSubmatch(data)
	if len(matches) == 0 {
		return nil, fmt.Errorf("There is no mark %s exist", mark)
	}

	lines := strings.Split(strings.TrimSpace(matches[1]), "\n")

	res := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}

	return res, nil
}

/* Initialize test case by real environment. */
func (atc *ATSCase) Init() error {
	if ok, _ := util.PathExists(defaults.GetConfigureFile(atc.Name)); !ok {
		return fmt.Errorf("Cannot init case %s, there is no configure file %s exist", atc.Name, defaults.GetConfigureFile(atc.Name))
	}

	data, err := ioutil.ReadFile(defaults.GetConfigureFile(atc.Name))
	if err != nil {
		return fmt.Errorf("Cannot init case %s when read %s with %s", atc.Name, defaults.GetConfigureFile(atc.Name), err)
	}

	if ok, err := atc.isConfigureFileValid(string(data)); !ok {
		return fmt.Errorf("Cannot init case %s with %s", atc.Name, err)
	}

	devs, err := atc.parseByMark(string(data), defaults.DevMark)
	if err != nil {
		return fmt.Errorf("Cannot init case %s with %s", err)
	}

	if len(devs) != len(atc.Topology.DUTs) {
		return fmt.Errorf("Case %s need %d device but current device count is %d", atc.Name, len(atc.Topology.DUTs), len(devs))
	}

	for _, dev := range devs {
		var newdev device.Device
		err := json.Unmarshal([]byte(dev), &newdev)
		if err != nil {
			return fmt.Errorf("Cannot parse device %s with %s", dev, err)
		}

		atc.Topology.DUTs[newdev.Name].SetDevice(&newdev)
	}

	emus, err := atc.parseByMark(string(data), defaults.EmulatorMark)
	if err != nil {
		return fmt.Errorf("Cannot init case %s with %s", err)
	}

	if len(emus) != len(atc.Topology.ITSs) {
		return fmt.Errorf("Case %s need %d testers but current tester count is %d", atc.Name, len(atc.Topology.ITSs), len(emus))
	}

	for _, emu := range emus {
		var newemu emulator.Emulator
		err := json.Unmarshal([]byte(emu), &newemu)
		if err != nil {
			return fmt.Errorf("Cannot parse emulator %s with %s", emu, err)
		}

		it, ok := atc.Topology.ITSs[newemu.Name]
		if !ok {
			return fmt.Errorf("There is no test %s in this topology", newemu.Name)
		}
		it.SetEmulator(&newemu)
	}

	itfs, err := atc.parseByMark(string(data), defaults.IfpMark)
	if err != nil {
		return fmt.Errorf("Cannot init case %s with %s", err)
	}

	if len(itfs) != atc.Topology.GetInterfaceCount() {
		return fmt.Errorf("Case %s need %d interfaces but current interface count is %d", atc.Name, atc.Topology.GetInterfaceCount(), len(itfs))
	}

	for _, itf := range itfs {
		var newitf ifp.Ifp
		err := json.Unmarshal([]byte(itf), &newitf)
		if err != nil {
			return fmt.Errorf("Cannot parse interface %s with %s", itf, err)
		}

		i := atc.Topology.GetInteraceByName(newitf.Name)
		if i == nil {
			return fmt.Errorf("Cannot find interface %s in current topology", newitf.Name)
		}

		i.FullName = newitf.FullName
		i.ShortName = newitf.ShortName
	}

	err = atc.Topology.Init()
	if err != nil {
		return fmt.Errorf("Test topology initialization failed with: %s", err)
	}

	atc.initialized = true

	return nil
}

func (atc *ATSCase) New(name string) (*ATSCase, error) {
	return nil, nil
}

func Get(name string) (*ATSCase, error) {
	if ok, _ := util.DirExists("asset/cases/" + name); !ok {
		return nil, fmt.Errorf("This is no test case %s exist", name)
	}

	ats := &ATSCase{
		Name:     name,
		prtasks:  make([]*routine.Routine, 0, 1),
		tasks:    make([]*routine.Routine, 0, 1),
		potasks:  make([]*routine.Routine, 0, 1),
		Topology: topology.New(),
	}

	err := ats.BuildTopology()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse topology for %s with %s", name, err)
	}

	err = ats.BuildTestTasks(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse test tasks for %s with %s", name, err)
	}

	return ats, nil
}
