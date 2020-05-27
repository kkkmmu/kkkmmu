package rut

type XP struct {
}

func (xp XP) GetAllInterfaces() (map[string]*Interface, error) {
	return nil, nil
}

func (xp XP) GetInterfaceList() ([]*Interface, error) {
	return nil, nil
}

func (xp XP) GetInterfaceByName(name string) (*Interface, error) {
	return nil, nil
}

func (xp XP) GetAllVlans() (map[string]*Vlan, error) {
	return nil, nil
}

func (xp XP) GetVlanList() ([]*Vlan, error) {
	return nil, nil
}

func (xp XP) GetVlanByID(vid string) ([]*Vlan, error) {
	return nil, nil
}

func (xp XP) GetStatistics() ([]*Statistics, error) {
	return nil, nil
}

func (xp XP) GetInterfaceStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (xp XP) GetCPUStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (xp XP) GetAllRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (xp XP) GetRouteList() ([]*Route, error) {
	return nil, nil
}

func (xp XP) GetRouteListByType(typ string) ([]*Route, error) {
	return nil, nil
}

func (xp XP) GetRoute(prefix string) (*Route, error) {
	return nil, nil
}

func (xp XP) GetAllARP() (map[string]*ARP, error) {
	return nil, nil
}

func (xp XP) GetARPList() ([]*ARP, error) {
	return nil, nil
}

func (xp XP) GetARP(prefix string) (*ARP, error) {
	return nil, nil
}

func (xp XP) GetAllND() (map[string]*ND, error) {
	return nil, nil
}

func (xp XP) GetNDList() ([]*ND, error) {
	return nil, nil
}

func (xp XP) GetND(prefix string) (*ND, error) {
	return nil, nil
}

func (xp XP) GetFDBList() ([]*FDB, error) {
	return nil, nil
}

func (xp XP) GetAllOSPFs() (map[string]*OSPF, error) {
	return nil, nil
}

func (xp XP) GetOSPFList() ([]*OSPF, error) {
	return nil, nil
}

func (xp XP) GetOSPF(id string) (*OSPF, error) {
	return nil, nil
}

func (xp XP) GetAllOSPFNeighbor() (map[string]*OSPFNeighbor, error) {
	return nil, nil
}

func (xp XP) GetOSPFNeighborList() ([]*OSPFNeighbor, error) {
	return nil, nil
}

func (xp XP) GetOSPFNeighbor(id string) (*OSPFNeighbor, error) {
	return nil, nil
}

func (xp XP) GetOSPFLSDB() ([]*LSA, error) {
	return nil, nil
}

func (xp XP) GetOSPFLSDBByType(typ string) ([]*LSA, error) {
	return nil, nil
}

func (xp XP) GetOSPFLSA(typ, id string) (*LSA, error) {
	return nil, nil
}

func (xp XP) GetOSPFRouteList() ([]*Route, error) {
	return nil, nil
}

func (xp XP) GetAllOSPFRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (xp XP) GetOSPFRoute(prefix string) ([]*Route, error) {
	return nil, nil
}
