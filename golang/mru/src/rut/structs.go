package rut

type Interface struct {
	Name            string
	VRF             string
	Type            string
	AdminStatus     string
	OperationStatus string
	Index           string
	MTU             string
	Bandwidth       string
	Mode            string
	Hardware        string
	HardwareAddress string
	Loopback        bool
	Multicast       bool
	Broadcast       bool
	IP              []string
	IP6             []string
	RxPPS           string
	RxBPS           string
	TxPPS           string
	TxBPS           string
	RxUcastPacket   uint64
	RxMcastPacket   uint64
	RxBcastPacket   uint64
	RxTotalPacket   uint64
	RxTotalByte     uint64
	TxUcastPacket   uint64
	TxMcastPacket   uint64
	TxBcastPacket   uint64
	TxTotalPacket   uint64
	TxTotalByte     uint64
	RxDiscard       uint64
	RxErrors        uint64
	TxDiscard       uint64
	TxErrors        uint64
}

type Vlan struct {
}

type ARP struct {
}

type Statistics struct {
}

type Route struct {
}

type ND struct {
}

type FDB struct {
}

type OSPF struct {
}

type OSPFNeighbor struct {
}

type LSA struct {
}
