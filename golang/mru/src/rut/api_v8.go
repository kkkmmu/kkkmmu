package rut

type V8 struct {
}

func (v8 V8) GetAllInterfaces() (map[string]*Interface, error) {
	return nil, nil
}

func (v8 V8) GetInterfaceList() ([]*Interface, error) {
	return nil, nil
}

func (v8 V8) GetInterfaceByName(name string) (*Interface, error) {
	return nil, nil
}

func (v8 V8) GetAllVlans() (map[string]*Vlan, error) {
	return nil, nil
}

func (v8 V8) GetVlanList() ([]*Vlan, error) {
	return nil, nil
}

func (v8 V8) GetVlanByID(vid string) ([]*Vlan, error) {
	return nil, nil
}

func (v8 V8) GetStatistics() ([]*Statistics, error) {
	return nil, nil
}

func (v8 V8) GetInterfaceStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v8 V8) GetCPUStatistics(name string) (*Statistics, error) {
	return nil, nil
}

func (v8 V8) GetAllRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v8 V8) GetRouteList() ([]*Route, error) {
	return nil, nil
}

func (v8 V8) GetRouteListByType(typ string) ([]*Route, error) {
	return nil, nil
}

func (v8 V8) GetRoute(prefix string) (*Route, error) {
	return nil, nil
}

func (v8 V8) GetAllARP() (map[string]*ARP, error) {
	return nil, nil
}

func (v8 V8) GetARPList() ([]*ARP, error) {
	return nil, nil
}

func (v8 V8) GetARP(prefix string) (*ARP, error) {
	return nil, nil
}

func (v8 V8) GetAllND() (map[string]*ND, error) {
	return nil, nil
}

func (v8 V8) GetNDList() ([]*ND, error) {
	return nil, nil
}

func (v8 V8) GetND(prefix string) (*ND, error) {
	return nil, nil
}

func (v8 V8) GetFDBList() ([]*FDB, error) {
	return nil, nil
}

func (v8 V8) GetAllOSPFs() (map[string]*OSPF, error) {
	return nil, nil
}

func (v8 V8) GetOSPFList() ([]*OSPF, error) {
	return nil, nil
}

func (v8 V8) GetOSPF(id string) (*OSPF, error) {
	return nil, nil
}

func (v8 V8) GetAllOSPFNeighbor() (map[string]*OSPFNeighbor, error) {
	return nil, nil
}

func (v8 V8) GetOSPFNeighborList() ([]*OSPFNeighbor, error) {
	return nil, nil
}

func (v8 V8) GetOSPFNeighbor(id string) (*OSPFNeighbor, error) {
	return nil, nil
}

func (v8 V8) GetOSPFLSDB() ([]*LSA, error) {
	return nil, nil
}

func (v8 V8) GetOSPFLSDBByType(typ string) ([]*LSA, error) {
	return nil, nil
}

func (v8 V8) GetOSPFLSA(typ, id string) (*LSA, error) {
	return nil, nil
}

func (v8 V8) GetOSPFRouteList() ([]*Route, error) {
	return nil, nil
}

func (v8 V8) GetAllOSPFRoutes() (map[string]*Route, error) {
	return nil, nil
}

func (v8 V8) GetOSPFRoute(prefix string) ([]*Route, error) {
	return nil, nil
}
