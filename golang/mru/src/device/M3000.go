package device

import (
	"command"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

var M3000 M3000_API

type M3000_API struct {
	Switch
}

var PoR = regexp.MustCompile("(?P<po>po[0-9]+)")
var MlagPortalR = regexp.MustCompile("(?P<mp>MLAG-[0-9]+)")
var MlagPortalLacpR = regexp.MustCompile("(?P<mp>MLAG-[0-9]+) Mapped Aggregator : po(?P<po>[0-9]+)")

/*
po20 100020
% Admin Key: 0020 - Oper Key 0020
%   Link: eth0/20 (5020) sync: 0
*/
var LacpSummaryPoR = regexp.MustCompile("po(?P<id>[0-9]+)[[:space:]]+[0-9]+")
var LacpSummaryLinkR = regexp.MustCompile("Link:[[:space:]]+(?P<link>[[:alnum:]/]+)[[:space:]]+")

func (m M3000_API) IsVlanExist(dev *Device, vid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show vlan brief"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show ip interface brief"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface vlan 1." + vid})

	fmt.Println("Call IsVlanExist on M3000_API")

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happend"), data
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, "Interface vlan1."+vid+" is") &&
		!strings.Contains(data, vid+" Vlan") &&
		!strings.Contains(data, "vlan1."+vid+" up") &&
		!strings.Contains(data, "vlan1."+vid+" down") {
		return fmt.Errorf("[ERROR] vlan exist check failed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) SetInterfaceUp(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no shutdown"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "exit"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happend"), data
	}
	return nil, data
}

func (m M3000_API) SetInterfaceDown(dev *Device, ifname string) (error, string) {
	fmt.Println("Call SetInterfaceDown with ", ifname, "on M3000_API")
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "shutdown"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "exit"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happend"), data
	}
	return nil, data
}

func (m M3000_API) IsInterfaceUp(dev *Device, ifname string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface " + ifname})

	fmt.Println("Call IsInterfaceUp with ", ifname)
	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happend: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "<UP,") {
		return nil, "true"
	}

	return fmt.Errorf("Interface up check failed with: %s", data), "false"
}

func (m M3000_API) IsInterfaceRunning(dev *Device, ifname string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface " + ifname})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happend: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "protocol is up") &&
		strings.Contains(data, "RUNNING") {
		return nil, "true"
	}

	return fmt.Errorf("[ERROR] Interaface running check failed with: %s", data), "false"
}

func (m M3000_API) GetInterfaceStatisticsAvg5Sec(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsAvg1Min(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsAvg10Min(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsInUcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsInMcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsInBcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsInDiscardPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsInErrorPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsOutUcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsOutMcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsOutBcastPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsOutDiscardPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) GetInterfaceStatisticsOutErrorPkts(dev *Device, ifname string) (error, string) {
	return nil, ""
}

func (m M3000_API) SetInterfaceSpeed(dev *Device, ifname, speed string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "speed " + speed})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) UnSetInterfaceSpeed(dev *Device, ifname, speed string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no speed "})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) SetInterfaceMtu(dev *Device, ifname, mtu string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "mtu " + mtu})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) UnSetInterfaceMtu(dev *Device, ifname, mtu string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no mtu "})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) SetInterfaceTypeL3(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no switchport"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) SetInterfaceTypeL2(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "bridge-group 1"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddInterfaceToVlan(dev *Device, ifname, vid, tagged string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	if tagged == "true" {
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport mode trunk"})
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport trunk allowed vlan add " + vid})
	} else {
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport mode access"})
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport access vlan " + vid})
	}
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelInterfaceFromVlan(dev *Device, ifname, vid, tagged string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	if tagged == "true" {
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport mode trunk"})
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport trunk allowed vlan remove " + vid})
	} else {
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "switchport mode access"})
		cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no switchport access vlan"})
	}
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) SetInterfaceIPAddress(dev *Device, ifname, ip string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "ip address " + ip})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) SetInterfaceSecondaryIPAddress(dev *Device, ifname, ip string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "ip address " + ip + " secondary"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) UnSetInterfaceIPAddress(dev *Device, ifname, ip string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no ip address " + ip})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) UnSetInterfaceSecondaryIPAddress(dev *Device, ifname, ip string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no ip address " + ip + " secondary"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddInterfaceToVRF(dev *Device, ifname, vrf string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "ip vrf forwarding " + vrf})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelInterfaceFromVRF(dev *Device, ifname, vrf string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no ip vrf forwarding " + vrf})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddVRF(dev *Device, vrf string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "ip vrf " + vrf})
	cmds = append(cmds, &command.Command{Mode: "config-vrf", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelVRF(dev *Device, vrf string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "no ip vrf " + vrf})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddMcecIntraDomainLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "mcec domain configuration"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "intra-domain-link " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelMcecIntraDomainLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "mcec domain configuration"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "no intra-domain-link"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddMcecDomainDataLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "mcec domain configuration"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "domain-data-link " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelMcecDomainDataLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "mcec domain configuration"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "no domain-data-link"})
	cmds = append(cmds, &command.Command{Mode: "config-mcec", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) IsMlagDomainUp(dev *Device) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show mlag domain details"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "Domain Sync : IN_SYNC") &&
		strings.Contains(data, "Neigh Domain Sync : IN_SYNC") &&
		strings.Contains(data, "Domain Adjacency : UP") {
		return nil, data
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) IsMlagInSync(dev *Device, mlag_id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show lagd mlag " + mlag_id})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "Mlag Sync : IN_SYNC") {
		return nil, data
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) AddMlag(dev *Device, lacp_id, mlag_id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface po" + lacp_id})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "mlag " + mlag_id})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelMlag(dev *Device, lacp_id, mlag_id string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface po" + lacp_id})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no mlag"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelMlagAll(dev *Device) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show mlag domain details"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	matches := MlagPortalLacpR.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		err, _ := m.DelMlag(dev, match[2], match[1])
		if err != nil {
			return fmt.Errorf("error happened when delete mlag all with %s", err), "false"
		}
	}

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) IsMlagExist(dev *Device, id string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show running-config"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show mlag domain details"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show lagd mlag " + id})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	matches := MlagPortalR.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		if match[0] == "MLAG-"+id {
			return nil, "true"
		}
	}

	return fmt.Errorf("mlag exist check failed with %s", data, "false"), "false"
}

func (m M3000_API) AddLacpInterface(dev *Device, id string) (error, string) {
	return nil, ""
}

func (m M3000_API) DelLacpInterface(dev *Device, id string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show etherchannel summary"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	lacps := strings.Split(string(data), "Aggregator")
	for _, lacp := range lacps {
		lacp = SpaceRegex.ReplaceAllString(lacp, " ")
		lacp = strings.Replace(lacp, "%", "", -1)
		pos := LacpSummaryPoR.FindStringSubmatch(lacp)
		if len(pos) != 2 {
			continue
		}

		if pos[1] != id {
			continue
		}

		links := LacpSummaryLinkR.FindAllStringSubmatch(lacp, -1)
		for _, link := range links {
			err, _ := m.DelLacpMember(dev, link[1], id)
			if err != nil {
				return fmt.Errorf("Cannot delete lacp %s with %s", id, err), ""
			}
		}
	}

	return nil, data
}

func (m M3000_API) DelLacpInterfaceAll(dev *Device) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show etherchannel summary"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	lacps := strings.Split(string(data), "Aggregator")
	for _, lacp := range lacps {
		lacp = SpaceRegex.ReplaceAllString(lacp, " ")
		lacp = strings.Replace(lacp, "%", "", -1)
		pos := LacpSummaryPoR.FindStringSubmatch(lacp)
		if len(pos) != 2 {
			continue
		}

		links := LacpSummaryLinkR.FindAllStringSubmatch(lacp, -1)
		for _, link := range links {
			err, _ := m.DelLacpMember(dev, link[1], pos[1])
			if err != nil {
				return fmt.Errorf("Cannot delete all lacp interface with %s", err), ""
			}
		}
	}

	return nil, data
}

func (m M3000_API) AddLacpMember(dev *Device, ifname, id string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "channel-group " + id + " mode active"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelLacpMember(dev *Device, ifname, id string) (error, string) {
	defer dev.GoInitMode()
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no channel-group"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) IsLacpUp(dev *Device, id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface po" + id})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "Interface po"+id+" is up") ||
		strings.Contains(data, "protocol is up") ||
		strings.Contains(data, "UP") ||
		strings.Contains(data, "RUNNING") {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, data
}

func (m M3000_API) IsLacpMemberInSync(dev *Device, id, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show etherchannel detail"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	matches := LacpMemberRegex.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		if match[1] == ifname {
			return nil, "true"
		}
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) SetLacpTrafficDistMode(dev *Device, id, mode string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface po" + id})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "port-channel load-balance " + mode})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) UnSetLacpTrafficDistMode(dev *Device, id, mode string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface po" + id})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no port-channel load-balance"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) AddRoute(dev *Device, prefix, masklen, nexthop, oif string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "ip route " + prefix + "/" + masklen + " " + nexthop + " " + oif})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelRoute(dev *Device, prefix, masklen, nexthop, oif string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "no ip route " + prefix + "/" + masklen + " " + nexthop + " " + oif})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) DelAllRoute(dev *Device) (error, string) {
	err, routes := m.GetAllStaticRoute(dev)
	if err != nil {
		return fmt.Errorf("Cannot delete all route with %s", err), "false"
	}

	if len(routes) == 0 {
		return nil, "true"
	}

	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	for _, route := range routes {
		cmds = append(cmds, &command.Command{Mode: "config", CMD: "no " + route})
	}
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) GetAllStaticRoute(dev *Device) (error, []string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show running-config ip route"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	lines := strings.Split(string(data), "\n")

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), nil
	}

	routes := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "ip route") && !strings.Contains(line, "default") && !strings.Contains(line, "0.0.0.0/0") {
			routes = append(routes, line)
		}
	}
	return nil, routes
}

func (m M3000_API) AddVlan(dev *Device, vid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "vlan database"})
	cmds = append(cmds, &command.Command{Mode: "config-vlan", CMD: "vlan " + vid + " bridge 1 state enable"})
	cmds = append(cmds, &command.Command{Mode: "config-vlan", CMD: "exit"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, fmt.Sprintf("vlan1.%s", vid)
}
func (m M3000_API) DelVlan(dev *Device, vid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "interface vlan1." + vid})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "shutdown"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "no ip address"})
	cmds = append(cmds, &command.Command{Mode: "config-if", CMD: "exit"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "no interface vlan1." + vid})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "vlan database"})
	cmds = append(cmds, &command.Command{Mode: "config-vlan", CMD: "no vlan " + vid + " bridge 1"})
	cmds = append(cmds, &command.Command{Mode: "config-vlan", CMD: "exit"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "exit"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}
	return nil, data
}

func (m M3000_API) IsInterfaceVlanMember(dev *Device, ifname, vid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show vlan brief"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show vlan " + vid})

	fmt.Println("Call IsVlanExist on M3000_API")

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, ifname) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsLacpInterfaceExist(dev *Device, id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface po" + id})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, "Interface po"+id+" is") {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsInterfaceLacpMember(dev *Device, ifname, id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show running-config interface " + ifname})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show etherchannel detail"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, "channel-group "+id+" mode ") ||
		!strings.Contains(data, "Link: "+ifname) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsInterfaceMlagIntraDomainLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show running-config mcec"})
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show mlag domain details"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, "intra-domain-link "+ifname) ||
		!strings.Contains(data, "Intra Domain Interface : "+ifname) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsInterfaceMlagDomainDataLink(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show running-config mcec"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if !strings.Contains(data, "domain-data-link "+ifname) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsMlagDomainInSync(dev *Device) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show mlag domain details"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, "Domain Sync : IN_SYNC") &&
		strings.Contains(data, "Neigh Domain Sync : IN_SYNC") &&
		strings.Contains(data, "Domain Adjacency : UP") {
		return nil, "true"
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) IsInterfaceMlagLocalMember(dev *Device, ifname, mlag_id string) (error, string) {
	err, index := m.GetInterfaceIndexByName(dev, ifname)
	if err != nil {
		return fmt.Errorf("Cannot get interface %s's index with %s", ifname, err), ""
	}

	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show lagd mlag " + mlag_id + " | include Home member List"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, index) {
		return nil, "true"
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) IsInterfaceMlagRemoteMember(dev *Device, ifname, mlag_id string) (error, string) {
	err, index := m.GetInterfaceIndexByName(dev, ifname)
	if err != nil {
		return fmt.Errorf("Cannot get interface %s's index with %s", ifname, err), ""
	}

	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show lagd mlag " + mlag_id + " | include Remote member List"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	if strings.Contains(data, index) {
		return nil, "true"
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) GetInterfaceIndexByName(dev *Device, ifname string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show interface " + ifname + " | include index"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	data = SpaceRegex.ReplaceAllString(data, " ")

	matches := InterfaceIndexRegex.FindStringSubmatch(data)
	if len(matches) == 2 {
		return nil, matches[1]
	}

	return fmt.Errorf("error happed: %s", data), "false"
}

func (m M3000_API) AddOspfInstance(dev *Device, id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) DelOspfInstance(dev *Device, id string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "no router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) AddOspfNetwork(dev *Device, id, network, area string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "network " + network + " area " + area})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) DelOspfNetwork(dev *Device, id, network, area string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "no network " + network + " area " + area})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) SetOspfInstanceRid(dev *Device, id, rid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "router-id " + rid})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) UnsetOspfInstanceRid(dev *Device, id, rid string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "configure terminal"})
	cmds = append(cmds, &command.Command{Mode: "config", CMD: "router ospf " + id})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "no router-id " + rid})
	cmds = append(cmds, &command.Command{Mode: "config-router", CMD: "end"})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	return nil, "true"
}

func (m M3000_API) IsOspfNeighorUp(dev *Device, id, neigh string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show ip ospf neighbor " + neigh})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	if strings.Contains(string(data), "Full") {
		return nil, "true"
	}

	return nil, "false"
}

func (m M3000_API) IsOspfNeighorState(dev *Device, id, neigh, state string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show ip ospf neighbor " + neigh})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	if strings.Contains(string(data), state) {
		return nil, "true"
	}

	return nil, "false"
}

func (m M3000_API) IsInterfaceState(dev *Device, id, ifname, state string) (error, string) {
	cmds := make([]*command.Command, 0, 1)
	cmds = append(cmds, &command.Command{Mode: "normal", CMD: "show ip ospf interface " + ifname})

	var data string
	wg := sync.WaitGroup{}
	for _, c := range cmds {
		wg.Add(1)
		res, err := dev.RunCommand(CTX, c)
		if err != nil {
			data += res
			data += fmt.Sprintf("Run Command: %s failed with: %s", c.CMD, err.Error())
		}
		wg.Done()
		data += res
	}
	wg.Wait()

	if IsErrorHappened(data) {
		return fmt.Errorf("error happed: %s", data), "false"
	}

	if strings.Contains(string(data), state) {
		return nil, "true"
	}

	return nil, "false"
}
