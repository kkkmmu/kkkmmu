package dut

import (
	"defaults"
	"device"
	"fmt"
	"ifp"
	"reflect"
	"strings"
)

/* DUT is controlled by API, not CLI.
   If you want to use CLI, use Device of this DUT.
*/
type DUT struct {
	Name   string     `yaml:"Name"`
	Ifps   []*ifp.Ifp `yaml:"Ifps"`
	IfpMap map[string]*ifp.Ifp
	Vlans  map[string]string
	dev    *device.Device
}

func New(name string) (*DUT, error) {
	return &DUT{
		Name:   name,
		Ifps:   make([]*ifp.Ifp, 0, 2),
		IfpMap: make(map[string]*ifp.Ifp, 2),
	}, nil
}

func (d *DUT) Init() error {
	if d.IfpMap == nil {
		d.IfpMap = make(map[string]*ifp.Ifp, 1)
	}

	if d.Ifps == nil {
		d.Ifps = make([]*ifp.Ifp, 0, 1)
	}

	if d.Vlans == nil {
		d.Vlans = make(map[string]string, 1)
	}

	for _, itf := range d.Ifps {
		d.IfpMap[itf.Name] = itf
	}

	return nil
}

func (d *DUT) GetAllInterface() []*ifp.Ifp {
	return d.Ifps
}

func (d *DUT) GetCurrentMode() (string, error) {
	return d.dev.GetCurrentMode()
}

func (d *DUT) GetDevice() (*device.Device, error) {
	if d.dev == nil {
		return nil, fmt.Errorf("DUT %s has not been assigned for any real device", d.Name)
	}

	return d.dev, nil
}

func (d *DUT) Connect() error {
	if d.dev == nil {
		return fmt.Errorf("DUTs %s has no real device assigned", d.Name)
	}

	err := d.dev.Init()
	if err != nil {
		return fmt.Errorf("Connect to DUT %s failed with: %s", d.Name, err)
	}

	return nil
}

func (d *DUT) Verify() error {
	d.IfpMap = nil
	d.IfpMap = make(map[string]*ifp.Ifp, len(d.Ifps))
	for _, itf := range d.Ifps {
		if _, ok := d.IfpMap[itf.Name]; ok {
			return fmt.Errorf("Duplicate ifp %s already exist on %s", itf.Name, d.Name)
		}
		d.IfpMap[itf.Name] = itf
	}

	return nil
}

func (d *DUT) Call(function string, args ...string) {
	rargs := make([]reflect.Value, len(args))
	for i, _ := range args {
		rargs[i] = reflect.ValueOf(args[i])
	}

	fmt.Printf("%s CALL %s with (%d)%s\n", d.Name, function, len(rargs), rargs)
	result := reflect.ValueOf(d).MethodByName(function).Call(rargs)
	if result[0].Interface() != nil {
		fmt.Printf("Run function %s on %s with error %s\n", function, d.dev.Name, result[0].Interface())
	}
}

func (d *DUT) Assert(function string, expected string, args ...string) bool {
	rargs := make([]reflect.Value, len(args))
	for i, _ := range args {
		rargs[i] = reflect.ValueOf(args[i])
	}

	result := reflect.ValueOf(d).MethodByName(function).Call(rargs)
	expected = strings.Trim(expected, `"`)
	/*
		if result[0].Interface() == nil && result[1].Interface() == expected {
			return true
		}
	*/

	fmt.Printf("%s ASSERT %s with %s expected: %s resut: %s\n", d.Name, function, args, expected, result[1].Interface())
	if result[1].Interface() == expected {
		return true
	}

	return false
}

func (d *DUT) AddInterface(nifp *ifp.Ifp) (string, string) {
	if _, ok := d.IfpMap[nifp.Name]; ok {
		return "false", fmt.Sprintf("Interface %s already exist on DUT: %s", nifp.Name, d.Name)
	}

	d.IfpMap[nifp.Name] = nifp

	return "true", ""
}

func (d *DUT) DelInterface(nifp *ifp.Ifp) (string, string) {
	if _, ok := d.IfpMap[nifp.Name]; !ok {
		return "false", fmt.Sprintf("Interface %s not exist on DUT: %s", nifp.Name, d.Name)
	}

	delete(d.IfpMap, nifp.Name)

	return "true", ""
}

func (d *DUT) AddLink(lifp, rifp *ifp.Ifp, rd *DUT) string {
	/*
		if _, ok := d.Links[lifp]; ok {
			return fmt.Errorf("%s of %s already connnect to %s of %s", lifp.Name, d.Name, ifp.LinkTo.Name, d.Links[lifp].Name)
		}

		d.Links[lifp] = rd
	*/

	return ""
}

func (d *DUT) DelLink(lifp, rifp *ifp.Ifp, rd *DUT) string {
	/*
		if _, ok := d.Links[lifp]; !ok {
			return fmt.Errorf("There is not link on %s of %s", lifp.Name, d.Name)
		}

		delete(d.Links, lifp)

		lifp.LinkTo = nil
	*/

	return ""
}

func (d *DUT) SetDevice(dev *device.Device) error {
	de, ok := defaults.Devices[dev.Device]
	if !ok {
		return fmt.Errorf("Device %s is not supported currently", dev.Device)
	}

	dev.APIs = de
	d.dev = dev

	for _, itf := range d.Ifps {
		itf.Dev = dev
	}

	return nil
}

func (d *DUT) SetInterfaceUp(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {

		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.SetInterfaceUp(d.dev, itf.FullName)
}

func (d *DUT) SetVlanInterfaceUp(vid string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.SetInterfaceUp(d.dev, name)
}

func (d *DUT) SetInterfaceDown(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {

		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.SetInterfaceDown(d.dev, itf.FullName)
}

func (d *DUT) SetVlanInterfaceDown(vid string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.SetInterfaceDown(d.dev, name)
}

func (d *DUT) IsInterfaceUp(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {

		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceUp(d.dev, itf.FullName)
}

func (d *DUT) IsVlanInterfaceUp(vid string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.IsInterfaceUp(d.dev, name)
}

func (d *DUT) IsInterfaceRunning(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"

	}
	return d.dev.APIs.IsInterfaceRunning(d.dev, itf.FullName)
}

func (d *DUT) IsVlanInterfaceRunning(vid string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.IsInterfaceRunning(d.dev, name)
}

func (d *DUT) GetInterfaceStatisticsAvg5Sec(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsAvg1Min(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsAvg10Min(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsInUcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsInMcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsInBcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsInDiscardPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsInErrorPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsOutUcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsOutMcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsOutBcastPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsOutDiscardPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) GetInterfaceStatisticsOutErrorPkts(ifname string) (error, string) {
	return nil, "false"
}

func (d *DUT) SetInterfaceSpeed(ifname, speed string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceSpeed(d.dev, itf.FullName, speed)
}

func (d *DUT) UnSetInterfaceSpeed(ifname, speed string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.UnSetInterfaceSpeed(d.dev, itf.FullName, speed)
}

func (d *DUT) SetInterfaceMtu(ifname, mtu string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceMtu(d.dev, itf.FullName, mtu)
}

func (d *DUT) UnSetInterfaceMtu(ifname, mtu string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.UnSetInterfaceMtu(d.dev, itf.FullName, mtu)
}

func (d *DUT) SetVlanInterfaceMtu(vid, mtu string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.SetInterfaceMtu(d.dev, name, mtu)
}

func (d *DUT) UnVlanSetInterfaceMtu(vid, mtu string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.UnSetInterfaceMtu(d.dev, name, mtu)
}

func (d *DUT) SetInterfaceTypeL3(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceTypeL3(d.dev, itf.FullName)
}

func (d *DUT) SetInterfaceTypeL2(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceTypeL2(d.dev, itf.FullName)
}

func (d *DUT) AddInterfaceToVlan(ifname, vid, tagged string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.AddInterfaceToVlan(d.dev, itf.FullName, vid, tagged)
}

func (d *DUT) DelInterfaceFromVlan(ifname, vid, tagged string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.DelInterfaceFromVlan(d.dev, itf.FullName, vid, tagged)
}

func (d *DUT) SetInterfaceIPAddress(ifname, ip string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceIPAddress(d.dev, itf.FullName, ip)
}

func (d *DUT) SetVlanInterfaceIPAddress(vid, ip string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.SetInterfaceIPAddress(d.dev, name, ip)
}

func (d *DUT) SetInterfaceSecondaryIPAddress(ifname, ip string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.SetInterfaceSecondaryIPAddress(d.dev, itf.FullName, ip)
}

func (d *DUT) SetVlanInterfaceSecondaryIPAddress(vid, ip string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.SetInterfaceSecondaryIPAddress(d.dev, name, ip)
}

func (d *DUT) UnSetInterfaceIPAddress(ifname, ip string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.UnSetInterfaceIPAddress(d.dev, itf.FullName, ip)
}

func (d *DUT) UnSetVlanInterfaceIPAddress(vid, ip string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.UnSetInterfaceIPAddress(d.dev, name, ip)
}

func (d *DUT) UnSetInterfaceSecondaryIPAddress(ifname, ip string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.UnSetInterfaceSecondaryIPAddress(d.dev, itf.FullName, ip)
}

func (d *DUT) UnSetVlanInterfaceSecondaryIPAddress(vid, ip string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}

	return d.dev.APIs.UnSetInterfaceSecondaryIPAddress(d.dev, name, ip)
}

func (d *DUT) AddInterfaceToVRF(ifname, vrf string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.AddInterfaceToVRF(d.dev, itf.FullName, vrf)
}

func (d *DUT) AddVlanInterfaceToVRF(vid, vrf string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.AddInterfaceToVRF(d.dev, name, vrf)
}

func (d *DUT) DelInterfaceFromVRF(ifname, vrf string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}

	return d.dev.APIs.DelInterfaceFromVRF(d.dev, itf.FullName, vrf)
}

func (d *DUT) DelVlanInterfaceFromVRF(vid, vrf string) (error, string) {
	name, ok := d.Vlans[vid]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", vid), "false"
	}
	return d.dev.APIs.DelInterfaceFromVRF(d.dev, name, vrf)
}

func (d *DUT) AddVRF(vrf string) (error, string) {
	return d.dev.APIs.AddVRF(d.dev, vrf)
}

func (d *DUT) DelVRF(vrf string) (error, string) {
	return d.dev.APIs.DelVRF(d.dev, vrf)
}

func (d *DUT) AddMcecIntraDomainLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.AddMcecIntraDomainLink(d.dev, itf.ShortName)
}

func (d *DUT) DelMcecIntraDomainLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.DelMcecIntraDomainLink(d.dev, itf.ShortName)
}

func (d *DUT) AddMcecDomainDataLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.AddMcecDomainDataLink(d.dev, itf.ShortName)
}

func (d *DUT) DelMcecDomainDataLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.DelMcecDomainDataLink(d.dev, itf.ShortName)
}

func (d *DUT) IsMlagDomainUp() (error, string) {
	return d.dev.APIs.IsMlagDomainUp(d.dev)
}

func (d *DUT) IsMlagInSync(mlag_id string) (error, string) {
	return d.dev.APIs.IsMlagInSync(d.dev, mlag_id)
}

func (d *DUT) AddMlag(lacp_id, mlag_id string) (error, string) {
	return d.dev.APIs.AddMlag(d.dev, lacp_id, mlag_id)
}

func (d *DUT) DelMlag(lacp_id, mlag_id string) (error, string) {
	return d.dev.APIs.DelMlag(d.dev, lacp_id, mlag_id)
}

func (d *DUT) DelMlagAll() (error, string) {
	return d.dev.APIs.DelMlagAll(d.dev)
}

func (d *DUT) IsMlagExist(mlag_id string) (error, string) {
	return d.dev.APIs.IsMlagExist(d.dev, mlag_id)
}

func (d *DUT) AddLacpInterface(id string) (error, string) {
	return d.dev.APIs.AddLacpInterface(d.dev, id)
}

func (d *DUT) DelLacpInterface(id string) (error, string) {
	return d.dev.APIs.DelLacpInterface(d.dev, id)
}

func (d *DUT) DelLacpInterfaceAll() (error, string) {
	return d.dev.APIs.DelLacpInterfaceAll(d.dev)
}

func (d *DUT) AddLacpMember(ifname, id string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.AddLacpMember(d.dev, itf.ShortName, id)
}

func (d *DUT) DelLacpMember(ifname, id string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.DelLacpMember(d.dev, itf.ShortName, id)
}

func (d *DUT) IsLacpUp(id string) (error, string) {
	return d.dev.APIs.IsLacpUp(d.dev, id)
}

func (d *DUT) IsLacpMemberInSync(id, ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsLacpMemberInSync(d.dev, id, itf.ShortName)
}

func (d *DUT) SetLacpTrafficDistMode(id, mode string) (error, string) {
	return d.dev.APIs.SetLacpTrafficDistMode(d.dev, id, mode)
}

func (d *DUT) UnSetLacpTrafficDistMode(id, mode string) (error, string) {
	return d.dev.APIs.UnSetLacpTrafficDistMode(d.dev, id, mode)
}

func (d *DUT) AddRoute(prefix, masklen, nexthop, oif string) (error, string) {
	itf, ok := d.IfpMap[oif]
	if oif != "" && !ok {
		return fmt.Errorf("Interface %s does not exist", oif), "false"
	}
	return d.dev.APIs.AddRoute(d.dev, prefix, masklen, nexthop, itf.ShortName)
}

func (d *DUT) DelRoute(prefix, masklen, nexthop, oif string) (error, string) {
	itf, ok := d.IfpMap[oif]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", oif), "false"
	}
	return d.dev.APIs.DelRoute(d.dev, prefix, masklen, nexthop, itf.ShortName)
}

func (d *DUT) DelAllRoute() (error, string) {
	return d.dev.APIs.DelAllRoute(d.dev)
}

func (d *DUT) AddVlan(vid string) (error, string) {
	err, name := d.dev.APIs.AddVlan(d.dev, vid)
	if err != nil {
		return fmt.Errorf("Create vlan %s failed with %s", vid, err), ""
	}

	d.Vlans[vid] = name

	return nil, name
}

func (d *DUT) DelVlan(vid string) (error, string) {
	_, ok := d.Vlans[vid]
	if ok {
		delete(d.Vlans, vid)
	}
	return d.dev.APIs.DelVlan(d.dev, vid)
}

func (d *DUT) IsVlanExist(vid string) (error, string) {
	return d.dev.APIs.IsVlanExist(d.dev, vid)
}

func (d *DUT) IsInterfaceVlanMember(ifname, vid string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceVlanMember(d.dev, itf.ShortName, vid)
}

func (d *DUT) IsLacpInterfaceExist(id string) (error, string) {
	return d.dev.APIs.IsLacpInterfaceExist(d.dev, id)
}

func (d *DUT) IsInterfaceLacpMember(ifname, id string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceLacpMember(d.dev, itf.ShortName, id)
}

func (d *DUT) IsInterfaceMlagIntraDomainLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceMlagIntraDomainLink(d.dev, itf.ShortName)
}

func (d *DUT) IsInterfaceMlagDomainDataLink(ifname string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceMlagDomainDataLink(d.dev, itf.ShortName)
}

func (d *DUT) IsMlagDomainInSync() (error, string) {
	return d.dev.APIs.IsMlagDomainInSync(d.dev)
}

func (d *DUT) IsInterfaceMlagLocalMember(ifname, mlag_id string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceMlagLocalMember(d.dev, itf.FullName, mlag_id)
}

func (d *DUT) IsInterfaceMlagRemoteMember(ifname, mlag_id string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {
		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsInterfaceMlagRemoteMember(d.dev, itf.FullName, mlag_id)
}

func (d *DUT) AddOspfInstance(id string) (error, string) {
	return d.dev.APIs.AddOspfInstance(d.dev, id)
}

func (d *DUT) DelOspfInstance(id string) (error, string) {
	return d.dev.APIs.DelOspfInstance(d.dev, id)
}

func (d *DUT) AddOspfNetwork(id, network, area string) (error, string) {
	return d.dev.APIs.AddOspfNetwork(d.dev, id, network, area)
}

func (d *DUT) DelOspfNetwork(id, network, area string) (error, string) {
	return d.dev.APIs.DelOspfNetwork(d.dev, id, network, area)
}

func (d *DUT) SetOspfInstanceRid(id, rid string) (error, string) {
	return d.dev.APIs.SetOspfInstanceRid(d.dev, id, rid)
}

func (d *DUT) UnsetOspfInstanceRid(id, rid string) (error, string) {
	return d.dev.APIs.UnsetOspfInstanceRid(d.dev, id, rid)
}

func (d *DUT) IsOspfNeighorUp(id, neigh string) (error, string) {
	return d.dev.APIs.IsOspfNeighorUp(d.dev, id, neigh)
}

func (d *DUT) IsOspfNeighorState(id, neigh, state string) (error, string) {
	return d.dev.APIs.IsOspfNeighorState(d.dev, id, neigh, state)
}

func (d *DUT) IsOspfInterfaceState(id, ifname, state string) (error, string) {
	itf, ok := d.IfpMap[ifname]
	if !ok {

		return fmt.Errorf("Interface %s does not exist", ifname), "false"
	}
	return d.dev.APIs.IsOspfInterfaceState(d.dev, id, itf.FullName, state)
}
