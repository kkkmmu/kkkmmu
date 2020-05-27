package n2x

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"telnet"
	"util"
)

var DEFAULTSESSIONNAME = "ATSN2X"

var ResultR = regexp.MustCompile("{(?P<result>[[:alnum:][:space:]]*)}")
var BasicResultR = regexp.MustCompile("(?P<status>[[:alnum:]-]+)[[:space:]]+(?P<result>[[:alnum:][:space:]_-{}*\".:]*)")

var ErrorNoOpenSession = errors.New("There is no open session")

type N2X struct {
	IP       string
	Port     string
	Conn     *telnet.Session
	APIs     map[string][]string
	Session  *NSession
	Ports    map[string]*Port
	OSPFs    map[string]*OSPF /* RID to Handler Map */
	OSPF6s   map[string]*OSPF /* RID to Handler Map */
	Traffics map[string]*Traffic
	lock     *sync.Mutex
}

type APIs map[string][]string

func New(ip, port, name string) (*N2X, error) {
	n2x := &N2X{
		IP:       ip,
		Port:     port,
		APIs:     make(map[string][]string, 10),
		Ports:    make(map[string]*Port, 10),
		OSPFs:    make(map[string]*OSPF, 10),
		OSPF6s:   make(map[string]*OSPF, 10),
		Traffics: make(map[string]*Traffic, 10),
		lock:     &sync.Mutex{},
	}

	/* Create admin session. */
	sess, err := telnet.New3(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}

	n2x.Conn = sess

	os.Remove("n2x_log.txt")

	cmd := fmt.Sprintf("AgtSessionManager OpenSession RouterTester900 AGT_SESSION_ONLINE")
	res, err := n2x.Invoke(cmd)
	if err != nil {
		return nil, err
	}

	n2x.KillSessionByName(name)

	/* Connect to the real manangement session .*/
	nsess := &NSession{
		ID:    res,
		Ports: make(map[string]*Port, 10),
		lock:  &sync.Mutex{},
	}
	nport, err := n2x.GetSessionPort(nsess.ID)
	if err != nil {
		return nil, err
	}

	nsess.Port = nport

	nsess, err = n2x.ConnectToSession(nsess)
	if err != nil {
		return nil, err
	}

	err = n2x.SetSessionName(nsess, name)
	if err != nil {
		return nil, err
	}

	n2x.Session = nsess

	return n2x, nil
}

func (n *N2X) Lock() {
	n.lock.Lock()
}

func (n *N2X) Unlock() {
	n.lock.Unlock()
}

func (n *N2X) DeInit(name string) error {
	return n.KillSession(n.Session.ID)
}

func (n *N2X) KillAllSessions() error {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		if err == ErrorNoOpenSession {
			return nil
		}

		return fmt.Errorf("Cannot kill all session: %s", err.Error())
	}

	for _, sess := range sesss {
		err := n.KillSession(sess.ID)
		if err != nil {
			return fmt.Errorf("Cannot session %s(%d) : %s", sess.Label, sess.ID, err.Error())
		}
	}

	return nil
}

func (n *N2X) KillSession(id string) error {
	cmd := fmt.Sprintf("AgtSessionManager CloseSession %s", id)
	_, err := n.Invoke(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (n *N2X) KillSessionByName(name string) error {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		if err == ErrorNoOpenSession {
			return nil
		}
		return fmt.Errorf("Cannot kill all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.Label == name {
			err := n.KillSession(sess.ID)
			if err != nil {
				return fmt.Errorf("Cannot session %s(%d) : %s", sess.Label, sess.ID, err.Error())
			}
		}
	}

	return nil
}

func (n *N2X) NewSession(name string) (*NSession, error) {
	cmd := fmt.Sprintf("AgtSessionManager OpenSession RouterTester900 AGT_SESSION_ONLINE")
	res, err := n.Invoke(cmd)
	if err != nil {
		return nil, err
	}

	nsess := &NSession{ID: res, Ports: make(map[string]*Port, 10)}
	nport, err := n.GetSessionPort(nsess.ID)
	if err != nil {
		return nil, err
	}

	nsess.Port = nport

	nsess, err = n.ConnectToSession(nsess)
	if err != nil {
		return nil, err
	}

	err = n.SetSessionName(nsess, name)
	if err != nil {
		return nil, err
	}

	return nsess, nil
}

/*
func (n *N2X) NewSession(ports ...string) (*NSession, error) {
	cmd := fmt.Sprintf("AgtSessionManager OpenSession RouterTester900 AGT_SESSION_ONLINE")
	res, err := n.Invoke(cmd)
	if err != nil {
		return nil, err
	}

	nsess := &NSession{ID: res, Ports: make(map[string]*Port, 10)}
	nport, err := n.GetSessionPort(nsess.ID)
	if err != nil {
		return nil, err
	}

	nsess.Port = nport

	nsess, err = n.ConnectToSession(nsess)
	if err != nil {
		return nil, err
	}

	err = n.SetSessionName(nsess, DEFAULTSESSIONNAME)
	if err != nil {
		return nil, err
	}

	err = nsess.ReservePorts(ports...)
	if err != nil {
		return nil, err
	}

	return nsess, nil
}
*/

func (n *N2X) ReservePort(name string) error {
	port, err := n.Session.ReservePort(name)
	if err != nil {
		return fmt.Errorf("Cannot reserve n2x port %s with %s", name, err)
	}

	n.Ports[name] = port

	return nil
}

func (n *N2X) ReleasePort(name string) error {
	err := n.Session.ReleasePort(name)
	if err != nil {
		return fmt.Errorf("Cannot release n2x port %s with %s", err)
	}

	delete(n.Ports, name)

	return nil
}

func (n *N2X) SetPortMediaType(port string, media PortMediaType) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.SetMediaType(media)
}

func (n *N2X) GetPortMediaType(port string) (PortMediaType, error) {
	np, ok := n.Ports[port]
	if !ok {
		return MEDIA_UNKNOWN, fmt.Errorf("Port %s is not reserved", port)
	}

	return np.GetMediaType()
}

func (n *N2X) SetPortLegacyDUTIP(port, ip string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.LegacyLinkAddSutIPAddress(ip)
}

func (n *N2X) AddPortLegacyHost(port, vid, ip, masklen, mac string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.AddLegacyHost(vid, ip, masklen, mac)
}

func (n *N2X) AddPortLegacyHosts(port, vid, ip, masklen, mac, count string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.AddLegacyHosts(vid, ip, masklen, mac, count)
}

func (n *N2X) DelPortLegacyAllHosts(port string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.LegacyLinkReset()
}

func (n *N2X) SendPortAllArpRequests(port string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.SendAllArpRequests()
}

func (n *N2X) SetPortLegacyDUTIP6(port, ip string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.LegacyLinkAddSutIP6Address(ip)
}

func (n *N2X) AddPortLegacyHost6(port, vid, ip, masklen, mac string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.AddLegacyHost6(vid, ip, masklen, mac)
}

func (n *N2X) AddPortLegacyHosts6(port, vid, ip, masklen, mac, count string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.AddLegacyHosts6(vid, ip, masklen, mac, count)
}

func (n *N2X) SendPortAllNeighborSolicitations(port string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.SendAllNeighborSolicitations()
}

func (n *N2X) ConnectToSessionByID(id string) (*NSession, error) {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.ID == id {
			sess, err = n.ConnectToSession(sess)
			if err != nil {
				return nil, fmt.Errorf("Cannot connect to session %s(%d) : %s", sess.Label, sess.ID, err.Error())
			}

			err = sess.Sync()
			if err != nil {
				return nil, err
			}

			return sess, nil
		}
	}
	return nil, fmt.Errorf("Cannot find session with id: %s", id)
}

func (n *N2X) GetSessionByID(id string) (*NSession, error) {
	return n.ConnectToSessionByID(id)
}

func (n *N2X) ConnectToSessionByName(name string) (*NSession, error) {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.Label == name {
			sess, err = n.ConnectToSessionByID(sess.ID)
			if err != nil {
				return nil, fmt.Errorf("Cannot connect to session with: %s", err.Error())
			}

			return sess, nil
		}
	}

	return nil, fmt.Errorf("Cannot find session with name: %s", name)
}

func (n *N2X) GetSessionByName(name string) (*NSession, error) {
	return n.ConnectToSessionByName(name)
}

func (n *N2X) SetSessionName(sess *NSession, name string) error {
	cmd := fmt.Sprintf("AgtSessionManager SetSessionLabel %s %s", sess.ID, name)
	_, err := n.Invoke(cmd)
	if err != nil {
		return err
	}

	sess.Label = name

	return nil
}

func (n *N2X) GetAllOpenSessions() ([]*NSession, error) {
	res, err := n.Invoke("AgtSessionManager ListOpenSessions")
	if err != nil {
		return nil, err
	}

	var sessions = make([]*NSession, 0, 10)
	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		sess := strings.Split(matches[1], " ")
		for _, s := range sess {
			if strings.TrimSpace(s) == "" {
				continue
			}
			nsess := &NSession{ID: s, Ports: make(map[string]*Port, 10)}
			sessions = append(sessions, nsess)
		}
	}

	if len(sessions) == 0 {
		return nil, ErrorNoOpenSession
	}

	for _, sess := range sessions {
		port, err := n.GetSessionPort(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session port of %s", sess.ID)
		}

		label, err := n.GetSessionLabel(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session label of %s", sess.ID)
		}

		pid, err := n.GetSessionPid(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session Pid of %s", sess.ID)
		}

		sess.Port = port
		sess.Label = strings.Trim(label, "\"")
		sess.Pid = pid
	}

	return sessions, nil
}

func (n *N2X) ConnectToSession(sess *NSession) (*NSession, error) {
	addr := fmt.Sprintf("%s:%s", n.IP, sess.Port)
	conn, err := telnet.New3(addr)
	if err != nil {
		return nil, err
	}

	sess.Conn = conn

	return sess, nil
}

func (n *N2X) GetSessionPort(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionPort %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) GetSessionLabel(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionLabel %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) GetSessionPid(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionPid %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) Invoke(cmds ...string) (string, error) {
	n.Lock()
	defer n.Unlock()

	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range cmds {
		cmd += fmt.Sprintf(" %s", p)
	}

	util.AppendToFile("n2x_log.txt", []byte(cmd+"\n"))

	_, err := n.Conn.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := n.GetCommandResult()
	if err != nil {
		return "", fmt.Errorf("Cannot get result of: %s with error: %s", cmd, err.Error())
	}

	util.AppendToFile("n2x_log.txt", []byte("Result: "+line))
	res := BasicResultR.FindStringSubmatch(line)
	if len(res) != 3 {
		return "", fmt.Errorf("Run %s with invalid result: %s", cmd, line)
	}

	if res[1] == "0" {
		return strings.TrimSpace(res[2]), nil
	}

	return "", fmt.Errorf("Run %s with result: (%s: %s)", cmd, res[1], res[2])

}

func (n *N2X) GetAllMethods() error {
	res, err := n.Invoke("AgtHelp ListObjects")
	if err != nil {
		return err
	}

	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		res = matches[1]
	}

	objects := strings.Split(res, " ")
	for _, obj := range objects {
		if _, ok := n.APIs[obj]; !ok {
			n.APIs[obj] = nil
		} else {
			fmt.Printf("Duplicate object: %s found!\n", obj)
		}
	}

	for obj, _ := range n.APIs {
		cmd := fmt.Sprintf("AgtHelp ListMethods %s", obj)
		res, err := n.Invoke(cmd)
		if err != nil {
			n.APIs[obj] = nil
			return err
		}

		matches := ResultR.FindStringSubmatch(res)
		if len(matches) == 2 {
			res = matches[1]
		}

		fields := strings.Split(res, " ")
		methods := make([]string, 0, 10)
		for _, field := range fields {
			methods = append(methods, field)
		}

		n.APIs[obj] = methods
	}

	return nil
}

func (n *N2X) GetCommandResult() (string, error) {
	var line []byte

	b, err := n.Conn.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = n.Conn.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else if b == 'm' {
		line, err = n.Conn.ReadUntil("brace")
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := n.Conn.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func (n *N2X) DeletePortAllOSPFs(port string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.DeleteAllOSPFs()
}

func (n *N2X) AddPortOSPF(port, area, rid, drid string) (*OSPF, error) {
	np, ok := n.Ports[port]
	if !ok {
		return nil, fmt.Errorf("Port %s is not reserved", port)
	}

	_, ok = np.OSPFs[rid]
	if ok {
		return nil, fmt.Errorf("Same OSPF instance %s already exist on port %s", rid, port)
	}

	oi, err := np.AddOSPF(area, rid, drid, rid)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf %s on port %s with %s", rid, port, err)
	}

	n.OSPFs[rid] = oi

	return oi, nil
}

func (n *N2X) AddOSPFExternalRoute(rid, prefix, plen, count, step string) error {
	oi, ok := n.OSPFs[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf instance %s", rid)
	}

	return oi.AddExternalRoute(prefix, plen, count, step)
}

func (n *N2X) DelOSPFExternalRoute(rid, prefix string) error {
	oi, ok := n.OSPFs[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf instance %s", rid)
	}

	return oi.DeleteExternalRoute(prefix)
}

func (n *N2X) AdvertiseOSPFExternalRoute(rid, prefix string) error {
	oi, ok := n.OSPFs[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf instance %s", rid)
	}

	return oi.AdvertiseExternalRoute(prefix)
}

func (n *N2X) WithdrawOSPFExternalRoute(rid, prefix string) error {
	oi, ok := n.OSPFs[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf instance %s", rid)
	}

	return oi.WithdrawExternalRoute(prefix)
}

func (n *N2X) DeletePortAllOSPF6s(port string) error {
	np, ok := n.Ports[port]
	if !ok {
		return fmt.Errorf("Port %s is not reserved", port)
	}

	return np.DeleteAllOSPF6s()
}

func (n *N2X) AddPortOSPF6(port, area, rid, drid string) (*OSPF, error) {
	np, ok := n.Ports[port]
	if !ok {
		return nil, fmt.Errorf("Port %s is not reserved", port)
	}

	_, ok = np.OSPF6s[rid]
	if ok {
		return nil, fmt.Errorf("Same OSPF instance %s already exist on port %s", rid, port)
	}

	oi, err := np.AddOSPF6(area, rid, drid, rid)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf %s on port %s with %s", rid, port, err)
	}

	n.OSPF6s[rid] = oi

	return oi, nil
}

func (n *N2X) AddOSPF6ExternalRoute(rid, prefix, plen, count, step string) error {
	oi, ok := n.OSPF6s[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf6 instance %s", rid)
	}

	return oi.AddExternalRoute6(prefix, plen, count, step)
}

func (n *N2X) DelOSPF6ExternalRoute(rid, prefix string) error {
	oi, ok := n.OSPF6s[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf6 instance %s", rid)
	}

	return oi.DeleteExternalRoute6(prefix)
}

func (n *N2X) AdvertiseOSPF6ExternalRoute(rid, prefix string) error {
	oi, ok := n.OSPF6s[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf6 instance %s", rid)
	}

	return oi.AdvertiseExternalRoute6(prefix)
}

func (n *N2X) WithdrawOSPF6ExternalRoute6(rid, prefix string) error {
	oi, ok := n.OSPF6s[rid]
	if !ok {
		return fmt.Errorf("Cannot find ospf6 instance %s", rid)
	}

	return oi.WithdrawExternalRoute6(prefix)
}

func (n *N2X) Reset() error {
	return n.Session.Reset()
}

func (n *N2X) StartRoutingEngine() error {
	return n.Session.StartRoutingEngine()
}

func (n *N2X) StopRoutingEngine() error {
	return n.Session.StopRoutingEngine()
}

func (n *N2X) StartTraffic() error {
	return n.Session.StartTest()
}

func (n *N2X) StopTraffic() error {
	return n.Session.StopTest()
}

func (n *N2X) AddPortStreams(name, sport, dport string) (*Traffic, error) {
	dnp, ok := n.Ports[dport]
	if !ok {
		return nil, fmt.Errorf("port %s has not been reserverd", dport)
	}

	snp, ok := n.Ports[sport]
	if !ok {
		return nil, fmt.Errorf("port %s has not been reserverd", sport)
	}

	_, ok = n.Traffics[name]
	if ok {
		return nil, fmt.Errorf("Same traffic %s already exist")
	}

	tr, err := snp.AddStreams(name, dnp)
	if err != nil {
		return nil, fmt.Errorf("Failed to and stream %s on port %s with %s", name, snp.Name, err)
	}

	err = tr.DefaultStreamGroup.SetIPv4TCP()
	if err != nil {
		return nil, fmt.Errorf("Cannot set steam %s type type IPv4 with %s", name, err)
	}

	tr.StreamType = "ipv4"
	n.Traffics[name] = tr

	return tr, nil
}

func (n *N2X) AddPortStreams6(name, sport, dport string) (*Traffic, error) {
	dnp, ok := n.Ports[dport]
	if !ok {
		return nil, fmt.Errorf("port %s has not been reserverd", dport)
	}

	snp, ok := n.Ports[sport]
	if !ok {
		return nil, fmt.Errorf("port %s has not been reserverd", sport)
	}

	_, ok = n.Traffics[name]
	if ok {
		return nil, fmt.Errorf("Same traffic %s already exist")
	}

	tr, err := snp.AddStreams(name, dnp)
	if err != nil {
		return nil, fmt.Errorf("Failed to and stream %s on port %s with %s", name, snp.Name, err)
	}

	err = tr.DefaultStreamGroup.SetIPv6TCP()
	if err != nil {
		return nil, fmt.Errorf("Cannot set steam %s type type IPv6 with %s", name, err)
	}

	tr.StreamType = "ipv6"

	n.Traffics[name] = tr

	return tr, nil
}

func (n *N2X) SetStreamsPPS(name, pps string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	err := tr.SetAverageLoad(PACKETS_PER_SEC, pps)
	if err != nil {
		return fmt.Errorf("Failed to set pps %s to stream %s with", pps, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsMPS(name, mps string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	err := tr.SetAverageLoad(MBITS_PER_SEC, mps)
	if err != nil {
		return fmt.Errorf("Failed to set mps %s to stream %s with", mps, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsVLAN(name, vid, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}
	return tr.SetStreamsVLAN(vid, count)
}

func (n *N2X) SetStreamsSrcMAC(name, mac, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	return tr.SetStreamsSrcMAC(mac, count)
}

func (n *N2X) SetStreamsDstMAC(name, mac, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	return tr.SetStreamsDstMAC(mac, count)
}

func (n *N2X) SetStreamsSrcIP(name, ip, plen, step, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv4" {
		return fmt.Errorf("Traffic %s is not ipv4 stream", name)
	}

	err := tr.SetStreamsSrcIP(ip, plen, step, count)
	if err != nil {
		return fmt.Errorf("Failed to set sip %s to stream %s with", ip, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsDstIP(name, ip, plen, step, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv4" {
		return fmt.Errorf("Traffic %s is not ipv4 stream", name)
	}

	err := tr.SetStreamsDstIP(ip, plen, step, count)
	if err != nil {
		return fmt.Errorf("Failed to set dip %s to stream %s with", ip, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsSrcIP6(name, ip, plen, step, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv6" {
		return fmt.Errorf("Traffic %s is not ipv6 stream", name)
	}

	err := tr.SetStreamsSrcIP6(ip, plen, step, count)
	if err != nil {
		return fmt.Errorf("Failed to set sip %s to stream %s with", ip, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsDstIP6(name, ip, plen, step, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv6" {
		return fmt.Errorf("Traffic %s is not ipv6 stream", name)
	}

	err := tr.SetStreamsDstIP6(ip, plen, step, count)
	if err != nil {
		return fmt.Errorf("Failed to set dip %s to stream %s with", ip, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsIPProtocol(name, proto, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv4" {
		return fmt.Errorf("Traffic %s is not ipv4 stream", name)
	}

	err := tr.SetStreamsIPProtocol(proto, count)
	if err != nil {
		return fmt.Errorf("Failed to set proto %s to stream %s with", proto, name, err)
	}

	return nil
}

func (n *N2X) SetStreamsIPv6NextHeader(name, nh, count string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	if tr.StreamType != "ipv6" {
		return fmt.Errorf("Traffic %s is not ipv6 stream", name)
	}

	err := tr.SetStreamsIPv6NextHeader(nh, count)
	if err != nil {
		return fmt.Errorf("Failed to set nh %s to stream %s with", nh, name, err)
	}

	return nil
}

func (n *N2X) DisableStreams(name string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	return tr.Disable()
}

func (n *N2X) EnableStreams(name string) error {
	tr, ok := n.Traffics[name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist", name)
	}

	return tr.Enable()
}

func (n *N2X) GetPortStatistics(port string) (uint64, uint64, uint64, error) {
	np, ok := n.Ports[port]
	if !ok {
		return 0, 0, 0, fmt.Errorf("Port %s is not reserved", port)
	}

	return np.GetStatistics()
}

func (n *N2X) GetStreamStatistics(name string) (uint64, uint64, uint64, error) {
	tr, ok := n.Traffics[name]
	if !ok {
		return 0, 0, 0, fmt.Errorf("cannot get traffic statistics Traffic %s does not exist", name)
	}

	return tr.GetStatistics()
}

func init() {

}
