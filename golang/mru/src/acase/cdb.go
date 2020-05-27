package acase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CDB struct {
	Cases      []*ACase
	CasesMap   map[string]*ACase
	Topologies map[string]*Topology
	Configs    map[string]*Config
	WG         sync.WaitGroup
}

func (cdb *CDB) GetAllCases() ([]*ACase, error) {
	acses := make([]*ACase, 0, 10)
	for _, ac := range cdb.Cases {
		if ac.Case == false {
			acs, err := cdb.GetSubCasesByName(ac.Name)
			if err != nil {
				return nil, fmt.Errorf("Failed to get subcase for %s with %s", ac.Name, err)
			}

			acses = append(acses, acs...)
		}
		acses = append(acses, ac)
	}

	if len(acses) == 0 {
		return nil, fmt.Errorf("The case DB is empty")
	}

	return acses, nil
}

func (cdb *CDB) GetCaseByName(name string) (*ACase, error) {

	for k, c := range cdb.CasesMap {
		if filepath.Base(k) == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("case %s does not exist", name)
}

func (cdb *CDB) AddCase(nac *ACase) error {

	if cdb.Cases == nil {
		cdb.Cases = make([]*ACase, 0, 1)
	}

	if cdb.CasesMap == nil {
		cdb.CasesMap = make(map[string]*ACase, 1)
	}

	if _, ok := cdb.CasesMap[nac.Path]; ok {
		return fmt.Errorf("Case %s alread exist", nac.Path)
	}

	cdb.Cases = append(cdb.Cases, nac)
	cdb.CasesMap[nac.Path] = nac

	return nil
}

func (cdb *CDB) AddTopology(tp *Topology) error {
	if cdb.Topologies == nil {
		cdb.Topologies = make(map[string]*Topology, 1)
	}

	if _, ok := cdb.Topologies[tp.Path]; ok {
		return fmt.Errorf("Topology %s alread exist", tp.Path)
	}

	cdb.Topologies[tp.Path] = tp

	filepath.Walk(tp.Path, func(path string, info os.FileInfo, err error) error {
		cdb.WG.Add(1)
		defer cdb.WG.Done()
		if c, ok := cdb.CasesMap[path]; ok {
			if c.Topology != nil {
				if strings.HasPrefix(tp.Path, c.Topology.Path) {
					c.Topology = tp.Copy()
				}
				return nil
			}
			c.Topology = tp.Copy()
		}
		return nil
	})

	cdb.WG.Wait()

	return nil
}

func (cdb *CDB) AddConfig(conf *Config) error {
	if cdb.Configs == nil {
		cdb.Configs = make(map[string]*Config, 1)
	}

	if _, ok := cdb.Configs[conf.Path]; ok {
		return fmt.Errorf("Config %s alread exist", conf.Path)
	}

	cdb.Configs[conf.Path] = conf

	filepath.Walk(conf.Path, func(path string, info os.FileInfo, err error) error {
		cdb.WG.Add(1)
		defer cdb.WG.Done()
		if c, ok := cdb.CasesMap[path]; ok {
			if c.Config != nil {
				if strings.HasPrefix(conf.Path, c.Config.Path) {
					c.Config = conf
				}
				return nil
			}
			c.Config = conf
		}
		return nil
	})
	cdb.WG.Wait()

	return nil
}

func (cdb *CDB) IsCaseExist(name string) bool {
	for k, _ := range cdb.CasesMap {
		if strings.Contains(k, name) {
			return true
		}
	}

	return false
}

func (cdb *CDB) Init() error {
	for _, c := range cdb.CasesMap {
		cdb.WG.Add(1)
		err := c.Init()
		if err != nil {
			return fmt.Errorf("Case %s init failed with: %s", c.Name, err)
		}
		cdb.WG.Done()
	}

	cdb.WG.Wait()

	return nil
}

func (cdb *CDB) GetSubCasesByName(name string) ([]*ACase, error) {
	nac, err := cdb.GetCaseByName(name)
	if err != nil {
		return nil, fmt.Errorf("Case %s does not exist", name)
	}

	acses := make([]*ACase, 0, 10)
	acs, ok := cdb.CasesMap[nac.Path]
	if !ok {
		return nil, fmt.Errorf("Case %s does not exist", name)
	}

	for _, ac := range cdb.Cases {
		if ac.Module == acs.Name {
			acses = append(acses, ac)
		}
	}

	if len(acses) == 0 {
		return nil, fmt.Errorf("There is no subcase for %s", acs.Name)
	}

	return acses, nil
}
