package device

import ()

type Switch interface {
	IsVlanExist(dev *Device, vid string) (error, string)

	SetInterfaceUp(dev *Device, id string) (error, string)
	SetInterfaceDown(dev *Device, id string) (error, string)
	IsInterfaceUp(dev *Device, id string) (error, string)
	IsInterfaceRunning(dev *Device, id string) (error, string)

	SetInterfaceSpeed(dev *Device, ifname, speed string) (error, string)
	UnSetInterfaceSpeed(dev *Device, ifname, speed string) (error, string)
	SetInterfaceMtu(dev *Device, ifname, mtu string) (error, string)
	UnSetInterfaceMtu(dev *Device, ifname, mtu string) (error, string)
	SetInterfaceTypeL3(dev *Device, ifname string) (error, string)
	SetInterfaceTypeL2(dev *Device, ifname string) (error, string)
	AddInterfaceToVlan(dev *Device, ifname, vid, tagged string) (error, string)
	DelInterfaceFromVlan(dev *Device, ifname, vid, tagged string) (error, string)
	SetInterfaceIPAddress(dev *Device, ifname, ip string) (error, string)
	SetInterfaceSecondaryIPAddress(dev *Device, ifname, ip string) (error, string)
	UnSetInterfaceIPAddress(dev *Device, ifname, ip string) (error, string)
	UnSetInterfaceSecondaryIPAddress(dev *Device, ifname, ip string) (error, string)

	/* Statistics */
	GetInterfaceStatisticsAvg5Sec(dev *Device, id string) (error, string)
	GetInterfaceStatisticsAvg1Min(dev *Device, id string) (error, string)
	GetInterfaceStatisticsAvg10Min(dev *Device, id string) (error, string)
	GetInterfaceStatisticsInUcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsInMcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsInBcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsInDiscardPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsInErrorPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsOutUcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsOutMcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsOutBcastPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsOutDiscardPkts(dev *Device, id string) (error, string)
	GetInterfaceStatisticsOutErrorPkts(dev *Device, id string) (error, string)

	/* VRF */
	AddVRF(dev *Device, vrf string) (error, string)
	DelVRF(dev *Device, vrf string) (error, string)
	AddInterfaceToVRF(dev *Device, ifname, vrf string) (error, string)
	DelInterfaceFromVRF(dev *Device, ifname, vrf string) (error, string)

	/* MLAG */
	AddMcecIntraDomainLink(dev *Device, ifname string) (error, string)
	DelMcecIntraDomainLink(dev *Device, ifname string) (error, string)
	AddMcecDomainDataLink(dev *Device, ifname string) (error, string)
	DelMcecDomainDataLink(dev *Device, ifname string) (error, string)
	IsMlagDomainUp(dev *Device) (error, string)
	IsMlagInSync(dev *Device, mlag_id string) (error, string)
	AddMlag(dev *Device, lacp_id, mlag_id string) (error, string)
	DelMlag(dev *Device, lacp_id, mlag_id string) (error, string)
	DelMlagAll(dev *Device) (error, string)
	IsMlagExist(dev *Device, id string) (error, string)
	IsInterfaceMlagIntraDomainLink(dev *Device, id string) (error, string)
	IsInterfaceMlagDomainDataLink(dev *Device, id string) (error, string)
	IsMlagDomainInSync(dev *Device) (error, string)
	IsInterfaceMlagLocalMember(dev *Device, id, mlag_id string) (error, string)
	IsInterfaceMlagRemoteMember(dev *Device, id, mlag_id string) (error, string)

	/* LACP */
	AddLacpInterface(dev *Device, id string) (error, string)
	DelLacpInterface(dev *Device, id string) (error, string)
	DelLacpInterfaceAll(dev *Device) (error, string)
	AddLacpMember(dev *Device, ifname, id string) (error, string)
	DelLacpMember(dev *Device, ifname, id string) (error, string)
	IsLacpUp(dev *Device, id string) (error, string)
	IsLacpMemberInSync(dev *Device, id, ifname string) (error, string)
	IsLacpInterfaceExist(dev *Device, id string) (error, string)
	IsInterfaceLacpMember(dev *Device, ifname, id string) (error, string)

	SetLacpTrafficDistMode(dev *Device, id, mode string) (error, string)
	UnSetLacpTrafficDistMode(dev *Device, id, mode string) (error, string)
	AddRoute(dev *Device, prefix, masklen, nexthop, oif string) (error, string)
	DelRoute(dev *Device, prefix, masklen, nexthop, oif string) (error, string)
	DelAllRoute(dev *Device) (error, string)

	/* Vlan */
	AddVlan(dev *Device, vid string) (error, string)
	DelVlan(dev *Device, vid string) (error, string)
	IsInterfaceVlanMember(dev *Device, id, vid string) (error, string)
	GetInterfaceIndexByName(dev *Device, ifname string) (error, string)

	/* OSPF */
	AddOspfInstance(dev *Device, id string) (error, string)
	DelOspfInstance(dev *Device, id string) (error, string)
	AddOspfNetwork(dev *Device, id, network, area string) (error, string)
	DelOspfNetwork(dev *Device, id, network, area string) (error, string)
	SetOspfInstanceRid(dev *Device, id, rid string) (error, string)
	UnsetOspfInstanceRid(dev *Device, id, rid string) (error, string)
	IsOspfNeighorUp(dev *Device, id, neigh string) (error, string)
	IsOspfNeighorState(dev *Device, id, neigh, state string) (error, string)
	IsOspfInterfaceState(dev *Device, id, ifname, state string) (error, string)
}
