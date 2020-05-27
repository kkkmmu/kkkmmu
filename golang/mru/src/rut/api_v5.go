package rut

type V5 struct {
}

func (v5 V5) GetAllInterfaces() (map[string]*Interface, error) {
	return nil, nil
}

func (v5 V5) GetInterfaceList() ([]*Interface, error) {
	return nil, nil
}

func (v5 V5) GetInterfaceByName(name string) (*Interface, error) {
	return nil, nil
}

func (v5 V5) GetAllVlans() (map[string]*Vlan, error) {
	return nil, nil
}

func (v5 V5) GetVlanList() ([]*Vlan, error) {
	return nil, nil
}

func (v5 V5) GetVlanByID(vid string) ([]*Vlan, error) {
	return nil, nil
}

func (v5 V5) GetStatistics() ([]*Statistics, error) {
	return nil, nil
}

func (v5 V5) GetInterfaceStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v5 V5) GetCPUStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v5 V5) GetAllRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v5 V5) GetRouteList() ([]*Route, error) {
	return nil, nil
}

func (v5 V5) GetRouteListByType(typ string) ([]*Route, error) {
	return nil, nil
}

func (v5 V5) GetRoute(prefix string) (*Route, error) {
	return nil, nil
}

func (v5 V5) GetAllARP() (map[string]*ARP, error) {
	return nil, nil
}

func (v5 V5) GetARPList() ([]*ARP, error) {
	return nil, nil
}

func (v5 V5) GetARP(prefix string) (*ARP, error) {
	return nil, nil
}

func (v5 V5) GetAllND() (map[string]*ND, error) {
	return nil, nil
}

func (v5 V5) GetNDList() ([]*ND, error) {
	return nil, nil
}

func (v5 V5) GetND(prefix string) (*ND, error) {
	return nil, nil
}

func (v5 V5) GetFDBList() ([]*FDB, error) {
	return nil, nil
}

func (v5 V5) GetAllOSPFs() (map[string]*OSPF, error) {
	return nil, nil
}

func (v5 V5) GetOSPFList() ([]*OSPF, error) {
	return nil, nil
}

func (v5 V5) GetOSPF(id string) (*OSPF, error) {
	return nil, nil
}

func (v5 V5) GetAllOSPFNeighbor() (map[string]*OSPFNeighbor, error) {
	return nil, nil
}

func (v5 V5) GetOSPFNeighborList() ([]*OSPFNeighbor, error) {
	return nil, nil
}

func (v5 V5) GetOSPFNeighbor(id string) (*OSPFNeighbor, error) {
	return nil, nil
}

func (v5 V5) GetOSPFLSDB() ([]*LSA, error) {
	return nil, nil
}

func (v5 V5) GetOSPFLSDBByType(typ string) ([]*LSA, error) {
	return nil, nil
}

func (v5 V5) GetOSPFLSA(typ, id string) (*LSA, error) {
	return nil, nil
}

func (v5 V5) GetOSPFRouteList() ([]*Route, error) {
	return nil, nil
}

func (v5 V5) GetAllOSPFRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v5 V5) GetOSPFRoute(prefix string) ([]*Route, error) {
	return nil, nil
}
