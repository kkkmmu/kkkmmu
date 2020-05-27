package acase

import (
	"emulator"
	"fmt"
	"ifp"
	"reflect"
	"strings"
)

type Tester struct {
	Name     string              `yaml:"Name"`
	Ifps     []*ifp.Ifp          `yaml:"Ifps"`
	IfpsMap  map[string]*ifp.Ifp `yaml:"IfpsMap"`
	Emulator emulator.EmulatorAPI
}

func (t *Tester) Verify(ifps map[string]*ifp.Ifp) error {
	if t.Ifps == nil {
		t.Ifps = make([]*ifp.Ifp, 0, 1)
	}

	if t.IfpsMap == nil {
		t.IfpsMap = make(map[string]*ifp.Ifp, 1)
	}

	for _, itf := range t.Ifps {
		if _, ok := t.IfpsMap[itf.Name]; ok {
			return fmt.Errorf("Duplicate interface %s exist on tester %s", itf.Name, t.Name)
		}
		t.IfpsMap[itf.Name] = itf
	}

	for _, itf := range ifps {
		if titf, ok := t.IfpsMap[itf.Name]; ok {
			titf.Media = itf.Media
			titf.Enable = itf.Enable
			titf.FullName = itf.FullName
			titf.ShortName = itf.ShortName
		}
	}

	return nil
}

func (t *Tester) GetAllInterface() []*ifp.Ifp {
	return t.Ifps
}

func (t *Tester) Init() error {
	err := t.Emulator.Init(t.Name)
	if err != nil {
		return fmt.Errorf("Init tester %s failed with %s", t.Name, err)
	}

	for _, itf := range t.Ifps {
		if itf.Enable {
			err = t.ReservePort(itf.Name)
			if err != nil {
				return fmt.Errorf("Init tester %s failed with %s", t.Name, err)
			}

			err = t.SetPortMediaType(itf.Name, itf.Media)
			if err != nil {
				return fmt.Errorf("Init tester %s failed with %s", t.Name, err)
			}
		}
	}

	return nil
}

func (t *Tester) DeInit() error {
	return t.Emulator.DeInit(t.Name)
}

func (t *Tester) ReservePort(port string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.ReservePort(itf.FullName)
}

func (t *Tester) ReleasePort(port string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.ReleasePort(itf.FullName)
}

func (t *Tester) SetPortMediaType(port string, media string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.SetPortMediaType(itf.FullName, media)
}

func (t *Tester) GetPortMediaType(port string) (string, error) {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return "", fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.GetPortMediaType(itf.FullName)
}

func (t *Tester) SetDutIpv4Address(port, ip string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.SetDutIpv4Address(itf.FullName, ip)
}

func (t *Tester) SetDutIpv6Address(port, ip string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.SetDutIpv6Address(itf.FullName, ip)
}

func (t *Tester) AddIpv4Host(port, vid, ip, masklen, mac string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.AddIpv4Host(itf.FullName, vid, ip, masklen, mac)
}

func (t *Tester) AddIpv6Host(port, vid, ip, masklen, mac string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.AddIpv6Host(itf.FullName, vid, ip, masklen, mac)
}

func (t *Tester) AddIpv4Hosts(port, vid, ip, masklen, mac, count string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.AddIpv4Hosts(itf.FullName, vid, ip, masklen, mac, count)
}

func (t *Tester) AddIpv6Hosts(port, vid, ip, masklen, mac, count string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.AddIpv6Hosts(itf.FullName, vid, ip, masklen, mac, count)
}

func (t *Tester) DelAllHosts(port string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.DelAllHosts(itf.FullName)
}

func (t *Tester) SendArpRequests(port string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.SendArpRequests(itf.FullName)
}

func (t *Tester) SendIpv6NeighborSolicitations(port string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Port %s does not exist on %s", port, t.Name)
	}
	return t.Emulator.SendIpv6NeighborSolicitations(itf.FullName)
}

func (t *Tester) AddL2Stream(name, sport, dport string) error {
	sitf, ok := t.IfpsMap[sport]
	if !ok {
		return fmt.Errorf("SPort %s does not exist on %s", sport, t.Name)
	}

	ditf, ok := t.IfpsMap[dport]
	if !ok {
		return fmt.Errorf("DPort %s does not exist on %s", dport, t.Name)
	}

	return t.Emulator.AddL2Stream(name, sitf.FullName, ditf.FullName)
}

func (t *Tester) AddIpv4Stream(name, sport, dport string) error {
	sitf, ok := t.IfpsMap[sport]
	if !ok {
		return fmt.Errorf("SPort %s does not exist on %s", sport, t.Name)
	}

	ditf, ok := t.IfpsMap[dport]
	if !ok {
		return fmt.Errorf("DPort %s does not exist on %s", dport, t.Name)
	}

	return t.Emulator.AddIpv4Stream(name, sitf.FullName, ditf.FullName)
}

func (t *Tester) AddIpv6Stream(name, sport, dport string) error {
	sitf, ok := t.IfpsMap[sport]
	if !ok {
		return fmt.Errorf("SPort %s does not exist on %s", sport, t.Name)
	}

	ditf, ok := t.IfpsMap[dport]
	if !ok {
		return fmt.Errorf("DPort %s does not exist on %s", dport, t.Name)
	}

	return t.Emulator.AddIpv6Stream(name, sitf.FullName, ditf.FullName)
}

func (t *Tester) SetStreamPPS(name, pps string) error {
	return t.Emulator.SetStreamPPS(name, pps)
}

func (t *Tester) SetStreamMPS(name, mps string) error {
	return t.Emulator.SetStreamMPS(name, mps)
}

func (t *Tester) SetStreamSrcMac(name, mac, count string) error {
	return t.Emulator.SetStreamSrcMac(name, mac, count)
}

func (t *Tester) SetStreamDstMac(name, mac, count string) error {
	return t.Emulator.SetStreamDstMac(name, mac, count)
}

func (t *Tester) SetStreamSrcIp(name, ip, plen, step, count string) error {
	return t.Emulator.SetStreamSrcIp(name, ip, plen, step, count)
}

func (t *Tester) SetStreamDstIp(name, ip, plen, step, count string) error {
	return t.Emulator.SetStreamDstIp(name, ip, plen, step, count)
}

func (t *Tester) SetStreamSrcIpv6(name, ip, plen, step, count string) error {
	return t.Emulator.SetStreamSrcIpv6(name, ip, plen, step, count)
}

func (t *Tester) SetStreamVlan(name, vid, count string) error {
	return t.Emulator.SetStreamVlan(name, vid, count)
}

func (t *Tester) SetStreamIpProtocol(name, proto, count string) error {
	return t.Emulator.SetStreamIpProtocol(name, proto, count)
}

func (t *Tester) SetStreamIpv6NextHeader(name, nh, count string) error {
	return t.Emulator.SetStreamIpv6NextHeader(name, nh, count)
}

func (t *Tester) SetStreamDstIpv6(name, ip, plen, step, count string) error {
	return t.Emulator.SetStreamDstIpv6(name, ip, plen, step, count)
}

func (t *Tester) IsPortTrafficLostOccured(sport, dport string) (error, string) {
	sitf, ok := t.IfpsMap[sport]
	if !ok {
		return fmt.Errorf("SPort %s does not exist on %s", sport, t.Name), "true"
	}

	ditf, ok := t.IfpsMap[dport]
	if !ok {
		return fmt.Errorf("DPort %s does not exist on %s", dport, t.Name), "true"
	}
	return t.Emulator.IsPortTrafficLostOccured(sitf.FullName, ditf.FullName)
}

func (t *Tester) IsStreamTrafficLostOccured(stream string) (error, string) {
	return t.Emulator.IsStreamTrafficLostOccured(stream)
}

func (t *Tester) StartTraffic() error {
	return t.Emulator.StartTraffic()
}

func (t *Tester) StopTraffic() error {
	return t.Emulator.StopTraffic()
}

func (t *Tester) StartRouting() error {
	return t.Emulator.StartRouting()
}

func (t *Tester) StopRouting() error {
	return t.Emulator.StopRouting()
}

func (t *Tester) AddOspf(port, area, rid, drid string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Add ospf failed with %s does not exist", port)
	}

	return t.Emulator.AddOspf(itf.FullName, area, rid, drid)
}

func (t *Tester) AddOspfExternalRoute(rid, prefix, plen, count, step string) error {
	return t.Emulator.AddOspfExternalRoute(rid, prefix, plen, count, step)
}

func (t *Tester) DelOspfExternalRoute(rid, prefix string) error {
	return t.Emulator.DelOspfExternalRoute(rid, prefix)
}

func (t *Tester) AdvertiseOspfExternalRoute(rid, prefix string) error {
	return t.Emulator.AdvertiseOspfExternalRoute(rid, prefix)
}

func (t *Tester) WithdrawOspfExternalRoute(rid, prefix string) error {
	return t.Emulator.WithdrawOspfExternalRoute(rid, prefix)
}

func (t *Tester) AddOspf6(port, area, rid, drid string) error {
	itf, ok := t.IfpsMap[port]
	if !ok {
		return fmt.Errorf("Add ospf6 failed with %s does not exist", port)
	}
	return t.Emulator.AddOspf6(itf.FullName, area, rid, drid)
}

func (t *Tester) AddOspf6ExternalRoute(rid, prefix, plen, count, step string) error {
	return t.Emulator.AddOspf6ExternalRoute(rid, prefix, plen, count, step)
}

func (t *Tester) DelOspf6ExternalRoute(rid, prefix string) error {
	return t.Emulator.DelOspf6ExternalRoute(rid, prefix)
}

func (t *Tester) AdvertiseOspf6ExternalRoute(rid, prefix string) error {
	return t.Emulator.AdvertiseOspf6ExternalRoute(rid, prefix)
}

func (t *Tester) WithdrawOspf6ExternalRoute(rid, prefix string) error {
	return t.Emulator.WithdrawOspf6ExternalRoute(rid, prefix)
}

func (t *Tester) Call(function string, args ...string) {
	rargs := make([]reflect.Value, len(args))
	for i, _ := range args {
		rargs[i] = reflect.ValueOf(args[i])
	}

	fmt.Printf("%s CALL %s with (%d)%s\n", t.Name, function, len(rargs), rargs)
	result := reflect.ValueOf(t).MethodByName(function).Call(rargs)
	if result[0].Interface() != nil {
		fmt.Printf("Run function %s on %s with error %s\n", function, t.Emulator.GetName(), result[0].Interface())
	}
}

func (t *Tester) Assert(function string, expected string, args ...string) bool {
	rargs := make([]reflect.Value, len(args))
	for i, _ := range args {
		rargs[i] = reflect.ValueOf(args[i])
	}

	result := reflect.ValueOf(t).MethodByName(function).Call(rargs)
	expected = strings.Trim(expected, `"`)
	/*
		if result[0].Interface() == nil && result[1].Interface() == expected {
			return true
		}
	*/

	fmt.Printf("%s ASSERT %s with %s expected: %s resut: %s err: %s\n", t.Name, function, args, expected, result[1].Interface(), result[0].Interface())
	if result[1].Interface() == expected {
		return true
	}

	return false
}
