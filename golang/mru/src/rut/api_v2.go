package rut

type V2 struct {
}

func (v2 V2) GetAllInterfaces() (map[string]*Interface, error) {
	return nil, nil
}

func (v2 V2) GetInterfaceList() ([]*Interface, error) {
	return nil, nil
}

func (v2 V2) GetInterfaceByName(name string) (*Interface, error) {
	return nil, nil
}

func (v2 V2) GetAllVlans() (map[string]*Vlan, error) {
	return nil, nil
}

func (v2 V2) GetVlanList() ([]*Vlan, error) {
	return nil, nil
}

func (v2 V2) GetVlanByID(vid string) ([]*Vlan, error) {
	return nil, nil
}

func (v2 V2) GetStatistics() ([]*Statistics, error) {
	return nil, nil
}

func (v2 V2) GetInterfaceStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v2 V2) GetCPUStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v2 V2) GetAllRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v2 V2) GetRouteList() ([]*Route, error) {
	return nil, nil
}

func (v2 V2) GetRouteListByType(typ string) ([]*Route, error) {
	return nil, nil
}

func (v2 V2) GetRoute(prefix string) (*Route, error) {
	return nil, nil
}

func (v2 V2) GetAllARP() (map[string]*ARP, error) {
	return nil, nil
}

func (v2 V2) GetARPList() ([]*ARP, error) {
	return nil, nil
}

func (v2 V2) GetARP(prefix string) (*ARP, error) {
	return nil, nil
}

func (v2 V2) GetAllND() (map[string]*ND, error) {
	return nil, nil
}

func (v2 V2) GetNDList() ([]*ND, error) {
	return nil, nil
}

func (v2 V2) GetND(prefix string) (*ND, error) {
	return nil, nil
}

func (v2 V2) GetFDBList() ([]*FDB, error) {
	return nil, nil
}

func (v2 V2) GetAllOSPFs() (map[string]*OSPF, error) {
	return nil, nil
}

func (v2 V2) GetOSPFList() ([]*OSPF, error) {
	return nil, nil
}

func (v2 V2) GetOSPF(id string) (*OSPF, error) {
	return nil, nil
}

func (v2 V2) GetAllOSPFNeighbor() (map[string]*OSPFNeighbor, error) {
	return nil, nil
}

func (v2 V2) GetOSPFNeighborList() ([]*OSPFNeighbor, error) {
	return nil, nil
}

func (v2 V2) GetOSPFNeighbor(id string) (*OSPFNeighbor, error) {
	return nil, nil
}

func (v2 V2) GetOSPFLSDB() ([]*LSA, error) {
	return nil, nil
}

func (v2 V2) GetOSPFLSDBByType(typ string) ([]*LSA, error) {
	return nil, nil
}

func (v2 V2) GetOSPFLSA(typ, id string) (*LSA, error) {
	return nil, nil
}

func (v2 V2) GetOSPFRouteList() ([]*Route, error) {
	return nil, nil
}

func (v2 V2) GetAllOSPFRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v2 V2) GetOSPFRoute(prefix string) ([]*Route, error) {
	return nil, nil
}
