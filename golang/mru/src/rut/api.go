package rut

type API interface {
	GetAllInterfaces() (map[string]*Interface, error)
	GetInterfaceList() ([]*Interface, error)
	GetInterfaceByName(name string) (*Interface, error)

	GetAllVlans() (map[string]*Vlan, error)
	GetVlanList() ([]*Vlan, error)
	GetVlanByID(vid string) ([]*Vlan, error)

	GetStatistics() ([]*Statistics, error)
	GetInterfaceStatistics(name string) (*Statistics, error)
	GetCPUStatistics(name string) (*Statistics, error)

	GetAllRoutes() (map[string]*Route, error)
	GetRouteList() ([]*Route, error)
	GetRouteListByType(typ string) ([]*Route, error)
	GetRoute(prefix string) (*Route, error)

	GetAllARP() (map[string]*ARP, error)
	GetARPList() ([]*ARP, error)
	GetARP(prefix string) (*ARP, error)

	GetAllND() (map[string]*ND, error)
	GetNDList() ([]*ND, error)
	GetND(prefix string) (*ND, error)

	GetFDBList() ([]*FDB, error)

	GetAllOSPFs() (map[string]*OSPF, error)
	GetOSPFList() ([]*OSPF, error)
	GetOSPF(id string) (*OSPF, error)
	GetAllOSPFNeighbor() (map[string]*OSPFNeighbor, error)
	GetOSPFNeighborList() ([]*OSPFNeighbor, error)
	GetOSPFNeighbor(id string) (*OSPFNeighbor, error)
	GetOSPFLSDB() ([]*LSA, error)
	GetOSPFLSDBByType(typ string) ([]*LSA, error)
	GetOSPFLSA(typ, id string) (*LSA, error)
	GetOSPFRouteList() ([]*Route, error)
	GetAllOSPFRoutes() (map[string]*Route, error)
	GetOSPFRoute(prefix string) ([]*Route, error)
}
