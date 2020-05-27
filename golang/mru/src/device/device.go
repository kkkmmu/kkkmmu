package device

import (
	"cli"
	"command"
	"configuration"
	"context"
	"errors"
	"fileserver"
	"fmt"
	"log"
	"logger"
	"path/filepath"
	"regexp"
	"result"
	"script"
	"strconv"
	"strings"
	"time"
)

var (
	CTX                 = context.Background()
	SpaceRegex          = regexp.MustCompile(`\s+`)
	LacpMemberRegex     = regexp.MustCompile(`Link: (?P<ifname>[[:alnum:]/]+) \([[:alnum:]]+\) sync: 1`)
	InterfaceIndexRegex = regexp.MustCompile(`index (?P<ifindex>[[:alnum:]]+) metric`)
)

//Device should be and interface
type Device struct {
	Name        string `yaml:"Name"`
	cli         *cli.Cli
	Device      string `yaml:"Device"`
	Type        string `yaml:"Type"` // PB: Pica Box, CH: Chassic
	Username    string `yaml:"Username"`
	Password    string `yaml:"Password"`
	Protocol    string `yaml:"Protocol"`
	IP          string `yaml:"IP"`
	Port        string `yaml:"Port"`
	BasePrompt  string `yaml:"BasePrompt"`
	Hostname    string `yaml:"Hostname"` //hostName
	SessionID   string `yaml:"SessionID"`
	SFU         string `yaml:"SFU"`
	FileServers map[string]*fileserver.FileServer
	APIs        Switch
}

type Config struct {
	Index      int    `json:"index"`
	Device     string `json:"device"`
	Protocol   string `json:"protocol"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"passowrd"`
	BasePrompt string `json:"baseprompt"`
	Hostname   string `json:"hostname"`
	SessionID  string `json:"-"`
	SFU        string `json:"sfu"`
}

type DB struct {
	DB map[string]*Device `json:"-"`
}

func NewDB() *DB {
	return &DB{
		DB: make(map[string]*Device, 1),
	}
}

func (db *DB) AddDUT(name string, dut *Device) {
	db.DB["DUT"+name] = dut
}

func (db *DB) GetDeviceByName(name string) *Device {
	if r, ok := db.DB[name]; ok {
		return r
	}
	return nil
}

func IsErrorHappened(in string) bool {
	if strings.Contains(in, "Invalid") || strings.Contains(in, "invalid") || strings.Contains(in, "INVALID") ||
		strings.Contains(in, "Error") || strings.Contains(in, "error") || strings.Contains(in, "ERROR") ||
		strings.Contains(in, "locked by other VTY") ||
		strings.Contains(in, "received SIGSEGV") || strings.Contains(in, "Call backtrace") {
		if !strings.Contains(in, "input errors") && !strings.Contains(in, "crc-error") &&
			!strings.Contains(in, "FCS error") && !strings.Contains(in, "UndersizeErrors") &&
			!strings.Contains(in, " OverSizeErrors") { //show interface packet statistics
			return true
		}
	}

	return false
}

func buildDefaultConfiguration(r *Device) *configuration.Configuration {
	//log.Println(r.Hostname, r.Device)
	var conf configuration.Configuration
	conf.Name = r.Name
	conf.Username = r.Username
	conf.Password = r.Password
	conf.Device = r.Device
	conf.Hostname = r.Hostname
	conf.BasePrompt = r.BasePrompt
	conf.IP = r.IP

	if r.Protocol == "" {
		conf.Protocol = "telnet"
	} else {
		conf.Protocol = r.Protocol
	}

	if r.Port != "" {
		conf.Port = r.Port
	} else if conf.Protocol == "telnet" {
		conf.Port = configuration.DefaultTelnetPort
	} else if conf.Protocol == "ssh" {
		conf.Port = configuration.DefaultSshPort
	}

	conf.EnablePrompt = configuration.DefaultEnablePrompt
	conf.LoginPrompt = configuration.DefaultLoginPrompt
	conf.PasswordPrompt = configuration.DefaultPasswordPrompt
	conf.Prompt = configuration.PromptEnd
	conf.ModeDB = configuration.BuildModeDBFromHostNameAndBasePrompt(r.Hostname, r.BasePrompt)
	conf.SessionID = r.SessionID

	if r.SFU == "" {
		conf.SFU = configuration.DefaultSFU
	} else {
		conf.SFU = r.SFU
	}

	//log.Printf("%#v", conf)
	return &conf
}

func New(r *Device) (*Device, error) {
	return r, nil
}

func (d *Device) Init() error {
	if d.Type == "CH" {
		d.BasePrompt = d.Hostname + "[" + d.SFU + "]"
	} else {
		d.BasePrompt = d.Hostname
	}

	d.FileServers = configuration.DefaultFileServer
	conf := buildDefaultConfiguration(d)
	c, err := cli.NewCli(conf)
	if err != nil {
		return errors.New("Cannot create CLI instance: " + err.Error())
	}

	d.cli = c

	name, err := c.GetHostname()
	if name != "" && name != d.Hostname {
		d.Hostname = name
		d.UpdateConfiguration()
	}

	err = c.Init()
	if err != nil {
		return errors.New("Cannot init Device with: " + err.Error())
	}

	return nil
}

func (d *Device) UpdateConfiguration() {
	//log.Println(r.Hostname, r.Device)
	var conf configuration.Configuration
	conf.Name = d.Name
	conf.Username = d.Username
	conf.Password = d.Password
	conf.Device = d.Device
	conf.Hostname = d.Hostname
	conf.BasePrompt = d.BasePrompt
	conf.IP = d.IP
	conf.Port = d.Port
	if d.Protocol == "" {
		conf.Protocol = "telnet"
	} else {
		conf.Protocol = d.Protocol
	}

	conf.EnablePrompt = configuration.DefaultEnablePrompt
	conf.LoginPrompt = configuration.DefaultLoginPrompt
	conf.PasswordPrompt = configuration.DefaultPasswordPrompt
	conf.Prompt = configuration.PromptEnd
	conf.ModeDB = configuration.BuildModeDBFromHostNameAndBasePrompt(d.Hostname, d.BasePrompt)
	conf.SessionID = d.SessionID

	if d.SFU == "" {
		conf.SFU = configuration.DefaultSFU
	} else {
		conf.SFU = d.SFU
	}

	//log.Printf("%#v", conf)
	d.cli.UpdateConfig(&conf)
}

func (d *Device) GoInitMode() {
	d.cli.GoNormalMode()
}

func (d *Device) SetModeDB(db map[string]string) {
	d.cli.SetModeDB(db)
}

func (d *Device) CurrentMode() string {
	return d.cli.CurrentMode()
}

func (d *Device) GetCurrentMode() (string, error) {
	return d.cli.GetCurrentMode()
}

func (d *Device) RunCommand(ctx context.Context, cmd *command.Command) (string, error) {
	logger.Push(ctx, fmt.Sprintf("Run Command: %-40s cmode:%15s mode: %15s on %20s\n", cmd.CMD, cmd.Mode, d.cli.CurrentMode(), d.IP))
	return d.runCommand(cmd)
}

func (d *Device) RunCommands(ctx context.Context, cmds []*command.Command) (string, error) {
	for _, c := range cmds {
		data, err := d.RunCommand(ctx, c)
		if err != nil {
			return data, err
		}
	}

	return "", nil
}

func (d *Device) WriteLine(line string) (int, error) {
	return d.cli.WriteLine(line)
}

func (d *Device) Expect(delims ...string) ([]byte, error) {
	return d.cli.Expect(delims...)
}

func (d Device) runCommand(cmd *command.Command) (string, error) {
	if cmd.Delay != 0 {
		<-time.After(time.Second * time.Duration(cmd.Delay))
	}
	data, err := d.cli.RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *Device) RunScript(sc *script.Script) <-chan result.Result {
	log.Printf("Start Runing Script: %v", sc)
	res := make(chan result.Result)
	go func(chan<- result.Result) {
		for i := 0; i < sc.Count; i++ {
			for _, c := range sc.Commands {
				<-time.After(time.Second * time.Duration(c.Delay))
				log.Printf("Run command: %v", c)
				data, err := r.cli.RunCommand(c)
				res <- result.Result{
					Command: c.CMD,
					Result:  string(data),
					Err:     err,
				}
			}
			<-time.After(time.Second * time.Duration(sc.Timer))
		}
		close(res)
	}(res)

	return res
}

func (d *Device) CreateVlan(id int) error {
	return nil

}

func (d *Device) DestroyVlan(id int) error {

	return nil
}

func (d *Device) DestroyAllVlan() error {

	return nil
}

func (d *Device) CreateVlanInterface(id int, ip string) error {

	return nil
}

func (d *Device) DestroyVlanInterface(id int) error {

	return nil
}

func (d *Device) AddIPAddress(ifname, ip string) error {

	return nil
}

func (d *Device) DelIPAddress(ifname, ip string) error {

	return nil
}

func (d *Device) AddSecondaryIPAddress(ifname, ip string) error {

	return nil
}

func (d *Device) DelSecondaryIPAddress(ifname, ip string) error {

	return nil
}

func (d *Device) DelAllIPAddress(ifname string) error {

	return nil
}

func (d *Device) CreateOSPFInstance(id, tag string) error {

	return nil
}

func (d *Device) DestroyOSPFInstance(id string) error {

	return nil
}

func isValidDeviceConfig(c *Config) bool {
	if c.Index < 0 ||
		c.Device == "" ||
		c.IP == "" ||
		c.Port == "" ||
		c.Username == "" {
		return false
	}

	return true
}

func (d *Device) IsAlive(ctx context.Context) bool {
	/*
		msg, err := d.cli.GoNormalMode()
		if err != nil {
			log.Println(err, msg)
			return false
		}
	*/

	res, err := d.RunCommand(ctx, &command.Command{
		Mode: d.cli.CurrentMode(),
		CMD:  "show running-config",
	})

	if err != nil {
		log.Println(err, res)
		return false
	}

	if strings.Contains(res, d.Hostname) {
		return true
	}
	return false
}

func (d Device) GoShellMode() ([]byte, error) {
	return d.cli.GoShellMode()
}

func (d *Device) FTPPut(local, ip, user, pass, dir string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	ldir := filepath.Dir(local)
	file := filepath.Base(local)

	_, err := d.WriteLine("ftp " + ip)
	if err != nil {
		return err
	}

	_, err = d.Expect("):")
	//log.Println(string(data))
	_, err = d.WriteLine(user)
	if err != nil {
		return err
	}
	_, err = d.Expect("Password:")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("lcd " + ldir)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("cd " + dir)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("put " + file)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}

	_, err = d.WriteLine("exit")
	if err != nil {
		return err
	}

	_, err = d.Expect("#")
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) FTPGet(local, ip, user, pass, dir, file string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	ldir := filepath.Dir(local)

	_, err := d.WriteLine("ftp " + ip)
	if err != nil {
		return err
	}

	_, err = d.Expect("):")
	//log.Println(string(data))
	_, err = d.WriteLine(user)
	if err != nil {
		return err
	}
	_, err = d.Expect("Password:")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("lcd " + ldir)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("cd " + dir)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("get " + file)
	if err != nil {
		return err
	}

	_, err = d.Expect("ftps>")
	if err != nil {
		return err
	}

	_, err = d.WriteLine("exit")
	if err != nil {
		return err
	}

	_, err = d.Expect("#")
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) TCPDUMP(intf, filter, file, count string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("tcpdump should working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}

	if file == "" || !strings.HasSuffix(file, ".pcap") {
		return fmt.Errorf("Invalid file name: ", file)
	}

	dump := "tcpdump "
	if intf != "" {
		dump += " -i " + intf
	}

	if count != "" {
		dump += " -c " + count
	}

	if filter != "" {
		dump += " " + filter
	}

	dump += " -w " + file

	if count != "" {
		_, err := d.WriteLine(dump)
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	} else {
		_, err := d.WriteLine(dump + "&")
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 10)
		_, err = d.WriteLine("killall tcpdump")
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Device) Put(server, local string) error {
	s, ok := d.FileServers[server]
	if !ok {
		return fmt.Errorf("remote server %s is not exist", server)
	}

	if s.Protocol == "ssh" {
		if d.cli.CurrentMode() != "shell" {
			_, err := d.GoShellMode()
			if err != nil {
				return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
			}
		}

		if !filepath.IsAbs(local) {
			_, err := d.WriteLine("scp " + local + " " + s.Username + "@" + s.IP + ":" + local)
			if err != nil {
				return err
			}
		} else {
			_, err := d.WriteLine("scp " + local + " " + s.Username + "@" + s.IP + ":" + filepath.Base(local))
			if err != nil {
				return err
			}
		}

		known, err := d.Expect("(yes/no)?", "password:")
		if strings.Contains(string(known), "(yes/no)?") {
			d.WriteLine("yes")
			_, err := d.Expect("password:")
			if err != nil {
				return err
			}
		}

		_, err = d.WriteLine(s.Password)
		if err != nil {
			return err
		}

		data, err := d.Expect("#")
		if err != nil {
			return err
		}

		if !strings.Contains(string(data), "100%") {
			return fmt.Errorf("Upload file with error: %s", string(data))
		}
	} else if s.Protocol == "ftp" {
		if d.cli.CurrentMode() != "shell" {
			_, err := d.GoShellMode()
			if err != nil {
				return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
			}
		}

		_, err := d.WriteLine("ftp " + s.IP)
		if err != nil {
			return err
		}

		_, err = d.Expect("):")
		_, err = d.WriteLine(s.Username)
		if err != nil {
			return err
		}
		_, err = d.Expect("Password:")
		if err != nil {
			return err
		}

		_, err = d.WriteLine(s.Password)
		if err != nil {
			return err
		}

		_, err = d.Expect("ftps>")
		if err != nil {
			return err
		}

		if !filepath.IsAbs(local) {
			_, err = d.WriteLine("put " + local)
			if err != nil {
				return err
			}
		} else {
			_, err = d.WriteLine("put " + filepath.Base(local))
			if err != nil {
				return err
			}
		}

		_, err = d.Expect("ftps>")
		if err != nil {
			return err
		}

		_, err = d.WriteLine("exit")
		if err != nil {
			return err
		}

		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown proto %s for server %s", s.Protocol, s.Name)
	}

	return nil
}

func (d *Device) Get(server, remote string) error {
	s, ok := d.FileServers[server]
	if !ok {
		return fmt.Errorf("remote server %s is not exist", server)
	}

	if s.Protocol == "ssh" {
		if d.cli.CurrentMode() != "shell" {
			_, err := d.GoShellMode()
			if err != nil {
				return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
			}
		}

		if !filepath.IsAbs(remote) {
			_, err := d.WriteLine("scp " + s.Username + "@" + s.IP + ":" + remote + " " + remote)
			if err != nil {
				return err
			}
		} else {
			_, err := d.WriteLine("scp " + s.Username + "@" + s.IP + ":" + remote + " " + filepath.Base(remote))
			if err != nil {
				return err
			}
		}

		known, err := d.Expect("(yes/no)?", "password:")
		if strings.Contains(string(known), "(yes/no)?") {
			d.WriteLine("yes")
			_, err := d.Expect("password:")
			if err != nil {
				return err
			}
		}

		_, err = d.WriteLine(s.Password)
		if err != nil {
			return err
		}

		data, err := d.Expect("#")
		if err != nil {
			return err
		}

		if !strings.Contains(string(data), "100%") {
			return fmt.Errorf("Upload file with error: %s", string(data))
		}
	} else if s.Protocol == "ftp" {
		if d.cli.CurrentMode() != "shell" {
			_, err := d.GoShellMode()
			if err != nil {
				return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
			}
		}

		_, err := d.WriteLine("ftp " + s.IP)
		if err != nil {
			return err
		}

		_, err = d.Expect("):")
		//log.Println(string(data))
		_, err = d.WriteLine(s.Username)
		if err != nil {
			return err
		}
		_, err = d.Expect("Password:")
		if err != nil {
			return err
		}
		//log.Println(string(data))

		_, err = d.WriteLine(s.Password)
		if err != nil {
			return err
		}

		_, err = d.Expect("ftps>")
		if err != nil {
			return err
		}
		//log.Println(string(data))

		if !filepath.IsAbs(remote) {
			_, err = d.WriteLine("get " + remote)
			if err != nil {
				return err
			}
		} else {
			_, err = d.WriteLine("get " + filepath.Base(remote))
			if err != nil {
				return err
			}
		}

		_, err = d.Expect("ftps>")
		if err != nil {
			return err
		}

		_, err = d.WriteLine("exit")
		if err != nil {
			return err
		}

		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown proto %s for server %s", s.Protocol, s.Name)
	}

	return nil
}

func (d *Device) SCPPut(local, ip, user, pass, dir string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	file := filepath.Base(local)

	_, err := d.WriteLine("scp " + local + " " + user + "@" + ip + ":" + dir + "/" + file)
	if err != nil {
		return err
	}

	known, err := d.Expect("(yes/no)?", "password:")
	if strings.Contains(string(known), "(yes/no)?") {
		d.WriteLine("yes")
		_, err := d.Expect("password:")
		if err != nil {
			return err
		}
	}

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err := d.Expect("#")
	if err != nil {
		return err
	}

	if !strings.Contains(string(data), "100%") {
		return fmt.Errorf("Upload file with error: %s", string(data))
	}

	return nil
}

func (d *Device) SCPGet(local, ip, user, pass, dir, file string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}

	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	_, err := d.WriteLine("scp " + user + "@" + ip + ":" + dir + "/" + file + " " + local)
	if err != nil {
		return err
	}

	known, err := d.Expect("(yes/no)?", "password:")
	if strings.Contains(string(known), "(yes/no)?") {
		d.WriteLine("yes")
		_, err := d.Expect("password:")
		if err != nil {
			return err
		}
	}

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err := d.Expect("#")
	if err != nil {
		return err
	}

	if !strings.Contains(string(data), "100%") {
		return fmt.Errorf("Upload file with error: %s", string(data))
	}

	return nil
}

func GetDeviceByConfig(c *Config) (*Device, error) {
	if !isValidDeviceConfig(c) {
		return nil, fmt.Errorf("Invalid config to create Device: %v", c)
	}
	newrut := &Device{
		Name:        "DUT" + strconv.Itoa(c.Index),
		Device:      c.Device,
		Username:    c.Username,
		Password:    c.Password,
		IP:          c.IP,
		Port:        c.Port,
		Hostname:    c.Hostname,
		BasePrompt:  c.BasePrompt,
		SessionID:   c.SessionID,
		SFU:         c.SFU,
		FileServers: configuration.DefaultFileServer,
	}

	log.Printf("%#v", newrut)

	if err := newrut.Init(); err != nil {
		return nil, fmt.Errorf("Cannot create new DUT with config :%v. Error Message: %s", c, err.Error())
	}

	return newrut, nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
