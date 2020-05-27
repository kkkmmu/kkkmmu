package n2x

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
You can measure statistics for:

a test port
a stream group
a stream
a PDU header field (for example, VLAN, PVC, IP address)
*/

//var DefaultStatistics = "AGT_IP_DATAGRAMS_TRANSMITTED AGT_IP_DATAGRAMS_RECEIVED AGT_IPV6_PACKETS_TRANSMITTED AGT_IPV6_PACKETS_RECEIVED AGT_ETHERNET_FRAMES_TRANSMITTED AGT_ETHERNET_FRAMES_RECEIVED"
var PortStatistics = "AGT_TEST_PACKETS_TRANSMITTED AGT_TEST_PACKETS_RECEIVED"
var StreamStatistics = "AGT_STREAM_PACKETS_TRANSMITTED AGT_STREAM_PACKETS_RECEIVED AGT_STREAM_PACKET_LOSS"

type Statistics struct {
	Name    string
	Handler string
	*Port
	Type        string
	Traffic     *Traffic
	DstPorts    map[string]*Port
	TxIPv4Pkts  float64
	RxIPv4Pkts  float64
	TxIPv6Pkts  float64
	RxIPv6Pkts  float64
	TxEtherPkts float64
	RxEtherPkts float64
	TxPkts      float64
	RxPkts      float64
	LossPkts    float64

	Interval uint64

	pTxIPv4Pkts  float64
	pRxIPv4Pkts  float64
	pTxIPv6Pkts  float64
	pRxIPv6Pkts  float64
	pTxEtherPkts float64
	pRxEtherPkts float64
	pInterval    uint64
	Timer        <-chan time.Time
}

const (
	TxIPv4Pkts = iota
	RxIPv4Pkts
	TxIPv6Pkts
	RxIPv6Pkts
	TxEtherPkts
	RxEtherPkts
)

func (st *Statistics) Init() error {
	//Pay attantion to the function call method here for multiple parameters.
	if st.Type == "Port" {
		cmd := fmt.Sprintf("AgtStatistics SelectStatistics %s %s", st.Handler, "{"+PortStatistics+"}")
		_, err := st.Invoke(cmd)
		if err != nil {
			return fmt.Errorf("Cannot init statistics %s with: %s", st.Handler, err)
		}

		err = st.SelectPorts(st.Port.Handler)
		if err != nil {
			return fmt.Errorf("Cannot init statistics %s with: %s", st.Handler, err)
		}
		st.Name = st.Port.StatisticsName()
	} else if st.Type == "Stream" {
		cmd := fmt.Sprintf("AgtStatistics SelectStatistics %s %s", st.Handler, "{"+StreamStatistics+"}")
		_, err := st.Invoke(cmd)
		if err != nil {
			return fmt.Errorf("Cannot init statistics %s with: %s", st.Handler, err)
		}

		cmd = fmt.Sprintf("AgtStatistics SelectStreamGroups %s %s", st.Handler, st.Traffic.Handler)
		_, err = st.Invoke(cmd)
		if err != nil {
			return fmt.Errorf("Cannot init statistics %s with: %s", st.Handler, err)
		}
		st.Name = st.Traffic.StatisticsName()
	}

	err := st.SetName(st.Name)
	if err != nil {
		return fmt.Errorf("Cannot init statistics %s with: %s", st.Handler, err)
	}

	/*
		if st.Timer == nil {
			st.Timer = time.Tick(time.Duration(time.Second * 2))
			go st.Collect()
		}
	*/

	return nil
}

//Pay attantion to the function call method here for multiple parameters.
func (st *Statistics) SelectPorts(ports ...string) error {
	ps := "{"
	for _, port := range ports {
		ps += fmt.Sprintf("%s ", port)
	}
	ps += "}"

	cmd := fmt.Sprintf("AgtStatistics SelectPorts %s %s", st.Handler, ps)
	_, err := st.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot bind ports to statistics port %s with: %s", ports, err)
	}

	return nil
}

func (st *Statistics) SetName(name string) error {
	cmd := fmt.Sprintf("AgtStatisticsList SetName %s %s", st.Handler, name)
	_, err := st.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set statistics name %s with: %s", st.Handler, err)
	}

	return nil
}

func (st *Statistics) GetName() (string, error) {
	cmd := fmt.Sprintf("AgtStatistics GetName %s", st.Handler)
	res, err := st.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get statistics name with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	st.Name = res

	return res, nil
}

func (st *Statistics) GetHandle() (string, error) {
	cmd := fmt.Sprintf("AgtStatisticsList GetHandle %s", st.Name)
	res, err := st.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get statistics name with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	st.Handler = res

	return res, nil
}

func (st *Statistics) GetStatistics() (uint64, uint64, uint64, error) {
	if st.Type == "Port" || st.Type == "Stream" {
		cmd := fmt.Sprintf("AgtStatistics GetStatistics %s", st.Handler)
		res, err := st.Invoke(cmd)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("Cannot get statistics with: %s", err)
		}

		res = strings.Replace(res, "{", "", -1)
		res = strings.Replace(res, "}", "", -1)
		res = strings.TrimSpace(res)

		sts := strings.Split(res, " ")

		if res == "" || (len(sts) != 4 && len(sts) != 3) {
			return 0, 0, 0, fmt.Errorf("Cannot get statistics %d", len(sts))
		}

		for i, s := range sts {
			if i == 0 {
				value, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return 0, 0, 0, fmt.Errorf("Invalid statistics value %s", s)
				}
				st.Interval = value
				continue
			} else {
				value, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return 0, 0, 0, fmt.Errorf("Invalid statistics value %s", s)
				}

				switch i {
				case 1:
					st.TxPkts = value
				case 2:
					st.RxPkts = value
				case 3:
					st.LossPkts = value
				default:
					return 0, 0, 0, fmt.Errorf("Unknown statistics type")
				}
			}

			if st.Type == "Port" {
				st.LossPkts = st.TxPkts - st.RxPkts
			}
		}
	} else {
		return 0, 0, 0, fmt.Errorf("Unkown statistics type %s", st.Type)
	}

	return uint64(st.TxPkts), uint64(st.RxPkts), uint64(st.LossPkts), nil
}

func (st *Statistics) getStatistics() error {
	cmd := fmt.Sprintf("AgtStatistics GetStatistics %s", st.Handler)
	res, err := st.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get statistics with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	sts := strings.Split(res, " ")

	if res == "" || len(sts) != 7 {
		return fmt.Errorf("Cannot get statistics")
	}

	for i, s := range sts {
		if i == 0 {
			value, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid statistics value %s", s)
			}
			st.Interval = value
			continue
		} else {
			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return fmt.Errorf("Invalid statistics value %s", s)
			}

			switch i {
			case 1:
				st.TxIPv4Pkts = value
			case 2:
				st.RxIPv4Pkts = value
			case 3:
				st.TxIPv6Pkts = value
			case 4:
				st.RxIPv6Pkts = value
			case 5:
				st.TxEtherPkts = value
			case 6:
				st.RxEtherPkts = value
			default:
				return fmt.Errorf("Unknown statistics type")
			}
		}
	}

	return nil
}

func (st *Statistics) GetAll() (uint64, uint64, uint64, uint64) {
	return st.GetRxTotal(), st.GetTxTotal(), st.GetRxRate(), st.GetTxRate()
}
func (st *Statistics) GetRxRate() uint64 {
	return uint64(st.RxEtherPkts-st.pRxEtherPkts) / (st.Interval - st.pInterval)
}

func (st *Statistics) GetTxRate() uint64 {
	return uint64(st.TxEtherPkts-st.pTxEtherPkts) / (st.Interval - st.pInterval)
}

func (st *Statistics) GetRxTotal() uint64 {
	return uint64(st.RxEtherPkts)
}

func (st *Statistics) GetTxTotal() uint64 {
	return uint64(st.TxEtherPkts)
}

func (st *Statistics) Collect() {
	if st.Timer == nil {
		return
	}

	for range st.Timer {
		pInterval := st.Interval
		pTxIPv4Pkts := st.TxIPv4Pkts
		pRxIPv4Pkts := st.RxIPv4Pkts
		pTxIPv6Pkts := st.TxIPv6Pkts
		pRxIPv6Pkts := st.RxIPv6Pkts
		pTxEtherPkts := st.TxEtherPkts
		pRxEtherPkts := st.RxEtherPkts

		err := st.getStatistics()
		if err != nil {
			fmt.Println(err)
		}

		if pInterval != st.Interval {
			st.pInterval = pInterval
			st.pTxIPv4Pkts = pTxIPv4Pkts
			st.pRxIPv4Pkts = pRxIPv4Pkts
			st.pTxIPv6Pkts = pTxIPv6Pkts
			st.pRxIPv6Pkts = pRxIPv6Pkts
			st.pTxEtherPkts = pTxEtherPkts
			st.pRxEtherPkts = pRxEtherPkts
		}
	}
}
