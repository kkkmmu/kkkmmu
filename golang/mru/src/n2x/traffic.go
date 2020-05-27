package n2x

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
 ListProtocols
 ARP ATM_Cell ATM_Payload ATM_OAM_Cell AAL1_Cell AAL5 AAL5_SNAP AAL5_SNAP_bridged AAL5_SNAP_bridged_PPPoE AAL5_NLPID AAL5_NLPID_PPP AAL5_CISCO AAL5_VCMux AAL5_VCMux_bridged AAL5_VCMux_bridged_PPPoE ATP BFD BGP4 BPDU CDP CGMP CiscoHDLC DDP Decnet DHCP DHCPv6 DVMRP EAP EAPOL EIGRP ELAP Ethernet Ethernet_SAP Ethernet_SNAP MAC_Control Slow_Control FrameRelay FrameRelay_NLPID FrameRelay_CISCO FrameRelay_SNAP FrameRelay_Q933 GMRP GRE GTP GVRP HSRP ICMP ICMP_v6 IGAP IGMP IGRP IPv4 raw_socket_header IPv6 IPX IS-IS ISL L2TPv2 L2TPv3 Control_Word PWE_ATM_1_1_Cell PWE_ATM_N_1_Cell PWE_AAL5_PDU PWE_AAL5_SDU PWE_FR_PAYLOAD LDP MDT MPLS MSDP OSPFv2 OSPFv3 PAgP PIMv4 POP3 PPPoHDLC PPP PPP-MP PPPoE_Session PPPoE_Discovery Raw_Shim Raw_Layer2 Raw_Payload RGMP RIP RIPng RSVP RTMP RTP TCP TCP_v6 UDP UDP_v6 VRRP VTP

 ListPacketTypes
 IPv4 IPv4/MPLS1 IPv4/MPLS2 TCPv4 UDPv4 IPv6 IPv6/MPLS1 IPv6/MPLS2 UDPv6 TCPv6 PING PINGv6 IGMPv3 MLD PWE_E/MPLS PWE_E/CW/MPLS PWE_VLAN/MPLS PWE_VLAN/CW/MPLS PWE_HDLC/MPLS PWE_HDLC/CW/MPLS PWE_PPP/MPLS PWE_PPP/CW/MPLS PWE_ATM_Cell(1:1)/MPLS PWE_ATM_Cell(N:1)/MPLS PWE_ATM_Cell(N:1)/CW/MPLS PWE_AAL5_PDU/CW/MPLS PWE_AAL5_SDU/CW/MPLS PWE_FR/MPLS PWE_FR/CW/MPLS PPPoE_Session PPPoE_Discovery
*/

var AddStreamWithProfileR = regexp.MustCompile("{(?P<stream>[[:alnum:][:space:]]+)}[[:space:]]+{(?P<pdu>[[:alnum:][:space:]]+)}")

//Should pay more attentaion to StreamGroupList object, We do not define this struct,
//At the same time we merge the function of this object into Traffic and StreamGroup.
type Traffic struct {
	Name       string
	Object     string
	Handler    string
	LoadObject string
	LoadMode   string
	Type       string
	DstPort    *Port
	*Port
	Unit               string
	Load               string
	StreamGroups       map[string]*StreamGroup
	DefaultStreamGroup *StreamGroup
	StreamCount        string
	StreamType         string
	Statistics         *Statistics
}

type StreamGroup struct {
	Name               string
	ID                 string
	Object             string
	Handler            string
	PDUObject          string
	PDUHandler         string
	Length             string
	LengthMode         string
	L2Protocol         string
	SourcePort         string
	DestinationPorts   map[string]string
	Type               string
	Enabled            bool
	NumberOfStream     string
	AllL2Protocol      string
	AllPacketType      string
	SourceEndpointType string
	SourceEndpoint     string
	*Traffic
}

type Ethernet struct {
	DMAC      string
	SMAC      string
	EtherType string
	VLANTag   string
}

/*
AGT_UNITS_PACKETS_PER_SEC
Specifies that the load is given in units of packets/s.

AGT_UNITS_MBITS_PER_SEC
Specifies that the load is given in units of Mb/s, expressed as an L2 load.

AGT_UNITS_PERCENTAGE_LINK_BANDWIDTH
Specifies that the load is given as a percentage of the link's available bandwidth.

AGT_UNITS_L3_MBITS_PER_SEC
Specifies that the load is given in units of Mb/s, expressed as an L3 load.

*/
const (
	PACKETS_PER_SEC = iota
	MBITS_PER_SEC
	PERCENTAGE_LINK_BANDWIDTH
	L3_MBITS_PER_SEC
	DEFAULT_TRAFFIC_UNIT = PACKETS_PER_SEC
)

var LoadUnitsNameMap = map[int]string{
	PACKETS_PER_SEC:           "AGT_UNITS_PACKETS_PER_SEC",
	MBITS_PER_SEC:             "AGT_UNITS_MBITS_PER_SEC",
	PERCENTAGE_LINK_BANDWIDTH: "AGT_UNITS_PERCENTAGE_LINK_BANDWIDTH",
	L3_MBITS_PER_SEC:          "AGT_UNITS_L3_MBITS_PER_SEC",
}

func (t *Traffic) Init() error {
	if t.Name == "" {
		return fmt.Errorf("You must give the name for traffic")
	}

	err := t.SetName(t.Name)
	if err != nil {
		return fmt.Errorf("Cannot init traffic with: %s", err)
	}

	cmd := fmt.Sprintf("AgtStreamGroupList AddStreamGroupsWithExistingProfile %s AGT_PACKET_STREAM_GROUP %s", t.Handler, t.StreamCount)
	res, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Init traffic with: %s", err)
	}

	matches := AddStreamWithProfileR.FindStringSubmatch(res)
	if len(matches) != 3 {
		return fmt.Errorf("Cannot init traffic with invalid retrun: %s", res)
	}

	streams := strings.Split(matches[1], " ")
	pdus := strings.Split(matches[2], " ")
	if len(streams) != len(pdus) {
		return fmt.Errorf("Cannot init traffic with invalid retrun: %s", res)
	}

	t.StreamGroups = make(map[string]*StreamGroup, len(streams))

	for i, stream := range streams {
		stream = strings.TrimSpace(stream)
		if stream == "" {
			continue
		}

		nsg := &StreamGroup{
			Handler:    stream,
			Object:     "AgtStreamGroup",
			PDUHandler: strings.TrimSpace(pdus[i]),
			Traffic:    t,
			PDUObject:  "AgtPduHeader",
		}

		err = nsg.Sync()
		if err != nil {
			return fmt.Errorf("Cannot init traffic with: %s", err)
		}

		if i == 0 {
			t.DefaultStreamGroup = nsg
		}

		t.StreamGroups[nsg.Handler] = nsg
	}

	if t.DstPort != nil {
		cmd := fmt.Sprintf("AgtStreamGroup SetExpectedDestinationPorts %s %s", t.Handler, t.DstPort.Handler)
		_, err := t.Invoke(cmd)
		if err != nil {
			return fmt.Errorf("Cannot Set traffic dst port with: %s", err)
		}
	}

	return nil
}

func (t *Traffic) SetName(name string) error {
	cmd := fmt.Sprintf("AgtProfileList SetName %s %s", t.Handler, name)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannotg set traffic profile name with: %s", err)
	}

	t.Name = name

	return nil
}

func (t *Traffic) GetName() (string, error) {
	cmd := fmt.Sprintf("%s GetName %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get traffic profile name with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic profile name")
	}

	t.Name = res

	return res, nil
}

func (t *Traffic) GetType() (string, error) {
	cmd := fmt.Sprintf("%s GetType %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get traffic profile type with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic profile type")
	}

	t.Type = res

	return res, nil
}

func (t *Traffic) GetSourcePort() (string, error) {
	//AgtCustomProfile GetSourcePort
	cmd := fmt.Sprintf("%s GetSourcePort %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get source port with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic profile source port")
	}

	return res, nil
}

func (t *Traffic) GetProfileType() (string, error) {
	//AgtCustomProfile GetSourcePort
	cmd := fmt.Sprintf("%s GetProfileType %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get profile type with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic profile type")
	}

	return res, nil
}

//AgtCustomProfile SetProfileType

func (t *Traffic) SetProfileType(typ string) error {
	cmd := fmt.Sprintf("%s SetProfileType %s", t.Object, t.Handler, typ)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set profile type with: %s", err)
	}

	t.Type = typ

	return nil
}

func (t *Traffic) GetMode() (string, error) {
	//AgtCustomProfile GetSourcePort
	cmd := fmt.Sprintf("%s GetMode %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get traffic mode with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic mode")
	}

	t.LoadMode = res

	return res, nil
}

func (t *Traffic) SetMode(mode string) error {
	cmd := fmt.Sprintf("%s SetMode %s", t.Object, t.Handler, mode)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set traffic mode with: %s", err)
	}

	t.LoadMode = mode

	return nil
}

func (t *Traffic) SetAverageLoad(unit int, load string) error {
	uname, ok := LoadUnitsNameMap[unit]
	if !ok {
		return fmt.Errorf("Invalid load mode: %d", unit)
	}

	cmd := fmt.Sprintf("%s SetAverageLoad %s %s %s", t.Object, t.Handler, load, uname)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set traffic average load with: %s", err)
	}

	t.Unit = uname
	t.Load = load
	return nil
}

func (t *Traffic) GetAverageLoad() (string, error) {
	cmd := fmt.Sprintf("%s GetAverageLoad %s %s", t.Object, t.Handler, t.Unit)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get traffic average load with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid traffic average load mode")
	}

	t.Load = res

	return res, nil
}

func (t *Traffic) GetNumberOfPacketsToInject() (string, error) {
	cmd := fmt.Sprintf("%s GetNumberOfPacketsToInject %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get traffic number of packets to inject with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid return value when get injection count")
	}

	return res, nil
}

func (t *Traffic) SetNumberOfPacketsToInject(count string) error {
	cmd := fmt.Sprintf("%s SetNumberOfPacketsToInject %s", t.Object, t.Handler, count)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set traffic number of packets to inject with: %s", err)
	}

	return nil
}

//AgtConstantProfile ListStreamGroups
func (t *Traffic) ListStreamGroups() (string, error) {
	cmd := fmt.Sprintf("%s ListStreamGroups %s", t.Object, t.Handler)
	res, err := t.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get stream groups with: %s", err)
	}

	res = strings.Replace(res, "\"", "", -1)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	return res, nil
}

func (t *Traffic) GetAllStreamGroups() ([]*StreamGroup, error) {
	res, err := t.ListStreamGroups()
	if err != nil {
		return nil, fmt.Errorf("Cannot get stream groups of %s with: %s", t.Handler, err)
	}

	fields := strings.Split(res, " ")
	sgs := make([]*StreamGroup, 0, len(fields))
	t.StreamGroups = make(map[string]*StreamGroup, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		nsg := &StreamGroup{
			Handler:   field,
			Object:    "AgtStreamGroup",
			Traffic:   t,
			PDUObject: "AgtPduHeader",
		}

		err = nsg.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot Get All streamgroup with: %s", err)
		}

		sgs = append(sgs, nsg)
	}

	return sgs, nil
}

//AgtConstantProfile AddStreamGroups
func (t *Traffic) AddStreamGroups(handler string) error {
	cmd := fmt.Sprintf("%s AddStreamGroups %s %s", t.Object, t.Handler, handler)
	_, err := t.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot add stream group %s to %s with: %s", handler, t.Handler, err)
	}

	return nil
}

/*
 * Each test port provides 15 customizable traffic profiles (128 profiles for all Fusion enabled load modules).
 * You assign traffic stream groups to one of these profiles.
 * For each profile, you can:
 *      select a constant (default), bursty, or custom transmission
 *      specify the transmit load and change it without having to restart a test
 *      indicate whether to send the defined traffic stream groups once or continuously
 *      send an exact number of packets
 *      determine and set the Interframe Departure Values (IDVs)
 *      dynamically disable and enable each profile and its stream groups
 *  You can either:
 *      create a stream group (and identify its traffic profile) first, then customize the profile
 *      customize the profile first, then create a stream group and assign it to the profile
 */

/* To define a constant profile

This example sends exactly 1000 layer 2 frames at 100 frames per second.

% # Reserve two test ports: 101/1 to send data, 101/2 to receive
% set hTxPort [AgtInvoke AgtPortSelector AddPort 101 1]
% set hRxPort [AgtInvoke AgtPortSelector AddPort 101 2]

% # Configure the first profile for constant traffic
% set hProfileConst [AgtInvoke AgtProfileList AddProfile $hTxPort AGT_CONSTANT_PROFILE]
% # Send all frames only once, in a single shot
% AgtInvoke AgtConstantProfile SetMode $hProfileConst AGT_TRAFFIC_PROFILE_MODE_ONE_SHOT
% # Set the average load to 100 layer 2 frames per second
% AgtInvoke AgtConstantProfile SetAverageLoad $hProfileConst 100 AGT_UNITS_PACKETS_PER_SEC
% # Send exactly 1000 layer 2 frames
% AgtInvoke AgtConstantProfile SetNumberOfPacketsToInject $hProfileConst 1000
*/

/*To define a custom profile

This example simulates jitter by varying the IDVs. It alternates between sending at one frame/s (one frame every 1,000,000,000 ns) and two frames/s (one frame every 500,000,000 ns). Note that the resolution of the IDVs depends on the port type; use the method shown below to determine a port's resolution (as a factor of 1 ns):

% # Configure the next profile to simulate jitter
% set hProfileCustom [AgtInvoke AgtProfileList AddProfile $hTxPort AGT_CUSTOM_PROFILE]
% # Set the IDVs in terms of ns
% set idvList [list 1000000000 500000000]
% # Get the IDV resolution of the port used by this profile
% set idvResolution [AgtInvoke AgtCustomProfile \
GetInterdepartureValueResolution $hProfileCustom]
% # Factor the port's resolution into the IDVs
% for {set i 0} {$i < [llength $idvList]} {incr i} {
    set idvList [lreplace $idvList $i $i \
    [expr round([expr [lindex $idvList $i] / $idvResolution])]]
}
% # Set the IDVs using the values adjusted for the port
% AgtInvoke AgtCustomProfile SetInterdepartureValues $hProfileCustom $idvList
*/

/*
For each stream group, you can assign a:
    traffic profile
    source (sending) and destination (receiving) test ports
    PDU template defining:
        the protocols encapsulated in the PDU header
        value in each PDU header field, which may be varied using field modifiers
        the variable field that is used to generate streams for per-stream statistics
        PDU length distribution
        PDU payload
*/

/*
To create a stream group

When you create a stream group, you can assign it to either:
    the default (first) traffic profile
    any existing traffic profile
    a new traffic profile
    Use the API object AgtStreamGroupList with one of the following methods.
    These examples illustrate commands entered interactively via the Tcl shell
    (showing the % prompt) so that you can see the values returned by different methods.

% # Create a stream group, assign it to the default traffic profile
% AgtInvoke AgtStreamGroupList AddStreamGroups $hTxPort \
AGT_PACKET_STREAM_GROUP 1
{1} 1 {1}
(returns handles to: new stream group, its profile, its PDU)
% # Create a stream group, assign it to a existing profile
% AgtInvoke AgtStreamGroupList AddStreamGroupsWithExistingProfile \
$hProfileBurst AGT_PACKET_STREAM_GROUP 1
{2} {2}
(returns handles to: new stream group, its PDU)
% # Create 3 stream groups, assign them to a new traffic profile
% AgtInvoke AgtStreamGroupList AddStreamGroupsWithNewProfile \
$hTxPort AGT_PACKET_STREAM_GROUP 3
{3 4 5} 2 {3 4 5}
(returns handles to: new stream groups, their profile, their PDUs)

*/

/*
To identify the test ports that are supposed to receive the traffic:

AgtInvoke AgtStreamGroup SetExpectedDestinationPorts $hStreamGroup $hRxPort
*/

/*
To select pre-defined protocols for PDUs
Examples for Ethernet interfaces:

% # Check the default protocols selected for the interface
% AgtInvoke AgtPduHeader ListProtocolsInHeader $hPdu
Ethernet IPv4

% # Define IPv6 PDUs containing: UDP / IPv6 / Ethernet
% AgtInvoke AgtStreamGroup SetPduHeaders $hStreamGroup \
{ethernet ipv6 udp_v6}
% # note: equivalent to [list ethernet ipv6 udp]

% # Test GRE with: Ethernet / MPLS / GRE / IPv4 / Ethernet
% AgtInvoke AgtStreamGroup SetPduHeaders $hStreamGroup \
{ethernet ipv4 gre mpls ethernet}

*/
/*
To define the PDU header values
To determine the values currently used
The tester populates each field in a PDU with either a:

default value from the applicable protocol XML file—for example, an Ethernet frame's Ether Type field is defined by default to be 0x9000 (loopback)
derived value from encapsulated protocols—for example, adding an IPv4 packet to an Ethernet frame's payload resets the Ether Type field to 0x0800 (as instructed by the Ethernet XML file)
configured value for the test port—for example, addresses set for the test port during the link-layer configuration are automatically used in its transmitted PDUs, as described in About source and destination addresses. Thus, configure a test port's link layer before you create its stream groups—link-layer configuration done after you have created stream groups does not affect the stream groups.
computed value—for example, the Ethernet frame's FCS field is calculated based on the values in the other PDU fields
To check the values currently being used in the Ethernet frame header:

% # List the fields currently enabled for the Ethernet header
% AgtInvoke AgtPduHeader ListProtocolFieldsInHeader $hPdu ethernet 1
destination_address source_address ether_type

% # Get the destination MAC address currently being used
% AgtInvoke AgtPduHeader GetFieldFixedValue $hPdu ethernet 1 destination_address
11:22:33:44:55:66

% # Get the current length of the frame
% AgtInvoke AgtStreamGroup GetLength $hStreamGroup
AGT_PACKET_LENGTH_FIXED {64}

% # Check the bytes in the entire PDU
% AgtInvoke AgtRawPdu GetAllPduBytes $hPdu
0x112233445566778899aabbcc08004500002e00000000403df88ac0030102c00101020000000000

*/

/*
To select alternate or optional fields
The protocol XML files can define optional PDU fields that:

represent different formats possible for a field—For example, the IPv4 packet XML file defines a Priority field that be a ToS-formatted value, DS-formatted value, or raw free-form value. The Raw format is defined as the default.
are inserted between header fields—For example, VLAN tags and MPLS labels can be inserted into protocol headers. For examples, see To enable optional fields like VLAN tags and To enable additional instances of fields like MPLS tags.
The methods used to select alternate and optional fields are the same. To enable an alternate field like ToS:

% # List the fields currently enabled for the IPv4 header
% AgtInvoke AgtPduHeader ListProtocolFieldsInHeader $hPdu ipv4 1
version hlen Raw tot_len identification reserved_flag fragment_flag
last_fragment_flag fragment_offset ttl protocol header_checksum
source_address destination_address

% # List the optional fields defined for the IPv4 header
% AgtInvoke AgtPduHeader ListOptionalFields $hPdu ipv4
options Raw TOS DS

% # Enable the TOS (not Raw) format for the Priority field
% AgtInvoke AgtPduHeader EnableOptionalField $hPdu ipv4 1 TOS

% # List the fields now enabled for the IPv4 header
% AgtInvoke AgtPduHeader ListProtocolFieldsInHeader $hPdu ipv4 1
version hlen precedence delay throughput reliability unused
tot_len identification reserved_flag fragment_flag
last_fragment_flag fragment_offset ttl protocol
header_checksum source_address destination_address FCS
The ToS fields (precedence, delay, throughput, reliability) can now be set in the same way as the other IPv4 header fields.
*/

/*
To set field values
When defining PDU header fields, you can use:

the default value
a specific value
a value from a list of specific values *
incrementing or decrementing values *
random values *
* These options use a field modifier to generate different values for a PDU field. For complete details about field modifiers and the number you can use, see About field modifiers.

For the names of the PDU header fields and any pre-defined values for them, see the protocol XML files.

% # Set the source IP address to be an incrementing range of 10 addresses
% # between 100.1.1.1 and 100.1.1.10
% AgtInvoke AgtPduHeader SetFieldIncrementingValueRange $hPdu "ipv4" 1 \
"source_address" 0 "100.1.1.1" 10 1
% # Set the destination IP address to be a decrementing range of 20 addresses
% # between 200.2.2.40 and 200.2.2.1 (decrements by 0.0.0.2)
% AgtInvoke AgtPduHeader SetFieldDecrementingValueRange $hPdu "ipv4" 1 \
"destination_address" 0 "200.2.2.40" 20 2
% # Change the source IP addresses to be a random range between 100.1.1.1
% # and 200.1.1.1 changing the /8 range and keeping the host fixed at x.1.1.1
% AgtInvoke AgtPduHeader SetFieldRandomValueRange $hPdu "ipv4" 1 \
"source_address" 8 "100.1.1.1" "200.1.1.1"
% # Select the MPLS label values used from a set of possible values
% # {16, 19, 45, 46, 32, 30, 78, 99, 1550, 0 }
% AgtInvoke AgtPduHeader SetFieldValueList $hPdu "mpls" 1 "Label1" \
[list 16 19 45 46 32 30 78 99 1550 0]
% # Set the ether type to the (incorrect) specific value of 0x0801
% AgtInvoke AgtPduHeader SetFieldFixedValue $hPdu "ethernet" 1 "ether_type" 0x0801
*/

/*
To generate individual streams for modified fields
As described in About field modifiers, you can select one field modifier to generate unique stream IDs. This allows you to measure per-stream statistics on each of a series of generated header field values. You can measure statistics on each source address, destination address, VLAN ID, QoS level, etc.

The following example selects a ToS field in order to measure per- stream QoS statistics on individual service levels.

# Generate packets with the 8 different ToS Precedence values
AgtInvoke AgtPduHeader SetFieldIncrementingValueRange $hPdu \
"ipv4" 1 "precedence" 0 0 8 1
# Select the ToS field for per-stream statistics
AgtInvoke AgtStreamGroup SetStreamGenerationParameter $hStreamGroup \
"ipv4" 1 "precedence"
Stream IDs are assigned based on whether or not you have enabled them. If you have:

enabled stream IDs for an independently modified field—a stream ID is assigned to each distinct value in the selected field
enabled stream IDs for a linked (fully meshed) field—a stream ID is assigned to each distinct combination of values in the modified fields. When using linked fields, you cannot assign a stream ID to one modified field and not other modified fields.
disabled stream IDs for all field modifiers—a single stream ID is assigned to all packets from the stream group

*/

//AgtStreamGroup
func (sg *StreamGroup) GetType() (string, error) {
	cmd := fmt.Sprintf("%s GetType %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup type with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup type")
	}

	sg.Type = res

	return res, nil
}

//Set name should use StreamGroupList object
func (sg *StreamGroup) GetName() (string, error) {
	cmd := fmt.Sprintf("%s GetName %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup name with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup name")
	}

	sg.Name = res

	return res, nil
}

func (sg *StreamGroup) GetLockCount() (string, error) {
	return "", nil
}

func (sg *StreamGroup) GetSourcePort() (string, error) {
	cmd := fmt.Sprintf("%s GetSourcePort %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup sport with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup sport")
	}

	sg.SourcePort = res

	return res, nil
}

func (sg *StreamGroup) SetSourceEndpointType() error {
	return nil
}

func (sg *StreamGroup) GetSourceEndpointType() (string, error) {
	cmd := fmt.Sprintf("%s GetSourceEndpointType %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup source endpoint type with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup source endpoint type")
	}

	sg.SourceEndpointType = res

	return res, nil
}

func (sg *StreamGroup) SetSourceEndpoint() error {
	return nil
}

func (sg *StreamGroup) GetSourceEndpoint() (string, error) {
	cmd := fmt.Sprintf("%s GetSourceEndpoint %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup source endpoint type with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup source endpoint type")
	}

	sg.SourceEndpoint = res

	return res, nil
}

func (sg *StreamGroup) Refresh() error {
	return nil
}

func (sg *StreamGroup) SetExpectedDestinationPorts() error {
	return nil
}

func (sg *StreamGroup) GetExpectedDestinationPorts() (string, error) {
	cmd := fmt.Sprintf("%s GetExpectedDestinationPorts %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup dst port with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup dst port")
	}

	fields := strings.Split(res, " ")

	sg.DestinationPorts = make(map[string]string, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		sg.DestinationPorts[field] = field
	}

	return res, nil
}

func (sg *StreamGroup) ListAllL2Protocols() (string, error) {
	cmd := fmt.Sprintf("%s ListAllL2Protocols %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup all packet type %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup all packet type")
	}

	sg.AllL2Protocol = res

	return res, nil

}

func (sg *StreamGroup) ListAllPacketTypes() (string, error) {
	cmd := fmt.Sprintf("%s ListAllPacketTypes", sg.Object)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup all packet type %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup all packet type")
	}

	sg.AllPacketType = res

	return res, nil
}

func (sg *StreamGroup) SetPduHeaders() error {
	return nil
}

func (sg *StreamGroup) SetPduHeadersByPacketType() error {
	return nil
}

func (sg *StreamGroup) AppendHeader() error {
	return nil
}

func (sg *StreamGroup) GetDefaultL2Protocol() (string, error) {
	cmd := fmt.Sprintf("%s GetDefaultL2Protocol %s", sg.Object, sg.Traffic.Port.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup default l2 protocol %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup default l2 protocol")
	}

	sg.L2Protocol = res

	return res, nil
}

func (sg *StreamGroup) SetLengthMode() error {
	return nil
}

func (sg *StreamGroup) GetLengthMode() (string, error) {
	cmd := fmt.Sprintf("%s GetLengthMode %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup length mode %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup length mode")
	}

	sg.LengthMode = res

	return res, nil
}

func (sg *StreamGroup) SetLength() error {
	return nil
}

func (sg *StreamGroup) GetLength() (string, error) {
	cmd := fmt.Sprintf("%s GetLength %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup length %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup length")
	}

	sg.Length = res

	return res, nil
}

func (sg *StreamGroup) GetStreamTag() (string, error) {
	return "", nil
}

func (sg *StreamGroup) GetStreamId(index string) (string, error) {
	cmd := fmt.Sprintf("%s GetStreamId %s %s", sg.Object, sg.Handler, index)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup id %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup id")
	}

	return res, nil
}

func (sg *StreamGroup) SetRepeatCount() error {
	return nil
}

func (sg *StreamGroup) GetRepeatCount() (string, error) {
	return "", nil
}

func (sg *StreamGroup) EnableTestPayload() error {
	return nil
}

func (sg *StreamGroup) DisableTestPayload() error {
	return nil
}

func (sg *StreamGroup) IsTestPayloadEnabled() error {
	return nil
}

func (sg *StreamGroup) SetProfile() error {
	return nil
}

func (sg *StreamGroup) GetProfile() (string, error) {
	return "", nil
}

func (sg *StreamGroup) GetPdu() (string, error) {
	cmd := fmt.Sprintf("%s GetPdu %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup pdu %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup pdu")
	}

	sg.PDUObject = "AgtPduHeader"
	sg.PDUHandler = res

	return res, nil
}

func (sg *StreamGroup) SetStreamGenerationParameter() error {
	return nil
}

func (sg *StreamGroup) GetStreamGenerationParameter() (string, error) {
	cmd := fmt.Sprintf("%s GetStreamGenerationParameter %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup nos %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup nos")
	}

	sg.NumberOfStream = res

	return res, nil
}

func (sg *StreamGroup) GetNumberOfStreams() (string, error) {
	cmd := fmt.Sprintf("%s GetNumberOfStreams %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup nos %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup nos")
	}

	sg.NumberOfStream = res

	return res, nil
}

func (sg *StreamGroup) GetStreamGenerationFieldValue() (string, error) {
	return "", nil
}

func (sg *StreamGroup) SetFieldModifiersRelation() error {
	return nil
}

func (sg *StreamGroup) GetFieldModifiersRelation() (string, error) {
	cmd := fmt.Sprintf("%s GetFieldModifiersRelation %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup FMR %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup FMR")
	}

	return res, nil
}

func (sg *StreamGroup) SetLinkedFieldModifiers() error {
	return nil
}

func (sg *StreamGroup) GetLinkedFieldModifiers() (string, error) {
	cmd := fmt.Sprintf("%s GetLinkedFieldModifiers %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup LFM %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup LFM")
	}

	return res, nil
}

func (sg *StreamGroup) Enable() error {
	cmd := fmt.Sprintf("%s Enable %s", sg.Object, sg.Handler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Enable streamgroup with: %s", err)
	}

	sg.Enabled = false

	return nil
}

func (tr *Traffic) Enable() error {
	return tr.DefaultStreamGroup.Enable()
}

func (sg *StreamGroup) Disable() error {
	cmd := fmt.Sprintf("%s Disable %s", sg.Object, sg.Handler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable streamgroup with: %s", err)
	}

	sg.Enabled = false

	return nil
}

func (tr *Traffic) Disable() error {
	return tr.DefaultStreamGroup.Disable()
}

func (sg *StreamGroup) IsEnabled() (bool, error) {
	cmd := fmt.Sprintf("%s IsEnabled %s", sg.Object, sg.Handler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return false, fmt.Errorf("Cannot get streamgroup state with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "\"", "", -1)

	if res == "" {
		return false, fmt.Errorf("Invalid streamgroup state")
	}

	if res == "1" {
		sg.Enabled = true
	} else {
		sg.Enabled = false
	}

	return sg.Enabled, nil
}

func (sg *StreamGroup) SetL2Error() error {
	return nil
}

func (sg *StreamGroup) GetL2Error() (string, error) {
	return "", nil
}

/*This function get all the layers in a packet */
func (sg *StreamGroup) ListProtocolsInHeader() (string, error) {
	cmd := fmt.Sprintf("%s ListProtocolsInHeader %s", sg.PDUObject, sg.PDUHandler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get streamgroup state with: %s", err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid streamgroup state")
	}

	return res, nil
}

/*This function get the fields name in a particular protocol header */
func (sg *StreamGroup) ListProtocolFieldsInHeader(proto string) (string, error) {
	cmd := fmt.Sprintf("%s ListProtocolFieldsInHeader %s %s 1", sg.PDUObject, sg.PDUHandler, proto)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get protocol field in %s with: %s", proto, err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Invalid protocol field")
	}

	return res, nil
}

/*This function get the fields value */
//AgtInvoke AgtPduHeader GetFieldFixedValue $hPdu ethernet 1 destination_address
func (sg *StreamGroup) GetFieldFixedValue(proto, field string) (string, error) {
	cmd := fmt.Sprintf("%s GetFieldFixedValue %s %s 1 %s", sg.PDUObject, sg.PDUHandler, proto, field)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get  protocol %s field %s's value with: %s", proto, field, err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Cannot get  protocol %s field %s's value with: %s", proto, field, err)
	}

	return res, nil
}

//AgtInvoke AgtPduHeader ListOptionalFields $hPdu ipv4
func (sg *StreamGroup) ListOptionalFields(proto string) (string, error) {
	cmd := fmt.Sprintf("%s ListOptionalFields %s %s", sg.PDUObject, sg.PDUHandler, proto)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get optional fileds in %s with: %s", proto, err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Cannot get optional fileds in %s with: %s", proto, err)
	}

	return res, nil
}

//AgtInvoke AgtPduHeader EnableOptionalField $hPdu ipv4 1 TOS
func (sg *StreamGroup) EnableOptionalField(proto, field string) error {
	cmd := fmt.Sprintf("%s EnableOptionalField %s %s 1 %s", sg.PDUObject, sg.PDUHandler, proto, field)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable optional fileds in %s with: %s", proto, err)
	}

	return nil
}

//AgtInvoke AgtPduHeader DisableOptionalField $hPdu ipv4 1 TOS
func (sg *StreamGroup) DisableOptionalField(proto, field string) error {
	cmd := fmt.Sprintf("%s DisableOptionalField %s %s 1 %s", sg.PDUObject, sg.PDUHandler, proto, field)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable optional fileds in %s with: %s", proto, err)
	}

	return nil
}

//AgtInvoke AgtPduHeader GetFieldLength $hPdu2 BGP4 1 prefix
func (sg *StreamGroup) GetFieldLength(proto, field string) (string, error) {
	cmd := fmt.Sprintf("%s GetFieldLength %s %s 1 %s", sg.PDUObject, sg.PDUHandler, proto, field)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get field length in %s with: %s", proto, err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Cannot get field length in %s with: %s", proto, err)
	}

	return res, nil
}

func (sg *StreamGroup) SetFieldLength(proto, field, length string) error {
	cmd := fmt.Sprintf("%s SetFieldLength %s %s 1 %s %s", sg.PDUObject, sg.PDUHandler, proto, field, length)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set field length in %s with: %s", proto, err)
	}

	return nil
}

//{PduHandle Protocol ProtocolInstance Field InFieldOffset StartValue NumValues StepSize}
func (sg *StreamGroup) SetFieldIncrementingValueRange(proto, field, offset, start, count, step string) error {
	cmd := fmt.Sprintf("%s SetFieldIncrementingValueRange %s %s 1 %s %s %s %s %s", sg.PDUObject, sg.PDUHandler, proto, field, offset, start, count, step)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Set field incremental value range in %s with: %s", proto, err)
	}

	return nil
}

func (sg *StreamGroup) SetFieldDecrementingValueRange(proto, field, offset, start, count, step string) error {
	cmd := fmt.Sprintf("%s  SetFieldDecrementingValueRange %s %s 1 %s 0 %s %s %s", sg.PDUObject, sg.PDUHandler, proto, field, offset, start, count, step)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Set field decremental value range in %s with: %s", proto, err)
	}

	return nil
}

//AgtInvoke AgtPduHeader SetFieldRandomValueRange $hPdu "ipv4" 1 "source_address" 8 "100.1.1.1" "200.1.1.1"
func (sg *StreamGroup) SetFieldRandomValueRange(proto, field, start, count, end string) error {
	cmd := fmt.Sprintf("%s SetFieldRandomValueRange %s %s 1 %s %s %s %s", sg.PDUObject, sg.PDUHandler, proto, field, count, start, end)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Set field random value range in %s with: %s", proto, err)
	}

	return nil
}

//AgtInvoke AgtPduHeader SetFieldFixedValue $hPdu "ethernet" 1 "ether_type" 0x0801
func (sg *StreamGroup) SetFieldFixedValue(proto, field, value string) error {
	cmd := fmt.Sprintf("%s SetFieldFixedValue %s %s 1 %s %s", sg.PDUObject, sg.PDUHandler, proto, field, value)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Set field fixed value range in %s with: %s", proto, err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv6SourceAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv6", "source_address", plen, start, count, step)
}

func (sg *StreamGroup) SetIPv6DestinationAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv6", "destination_address", plen, start, count, step)
}

func (sg *StreamGroup) SetIPv6NextHeader(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv6", "next_header", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6TrafficClass(start, count, step string) error {
	err := sg.EnableOptionalField("ipv6", "traffic_class")
	if err != nil {
		return fmt.Errorf("Cannot enable traffic class with %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ipv6", "traffic_class", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6DSCP(start, count, step string) error {
	err := sg.EnableOptionalField("ipv6", "ds")
	if err != nil {
		return fmt.Errorf("Cannot enable traffic class with %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ipv6", "default_dsfc", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6TcpSourcePort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("tcp_v6", "source_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6TcpDestinationPort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("tcp_v6", "destination_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6UdpSourcePort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("udp_v6", "source_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv6UdpDestinationPort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("udp_v6", "destination_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4SourceAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv4", "source_address", plen, start, count, step)
}

func (sg *StreamGroup) SetIPv4DestinationAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv4", "destination_address", plen, start, count, step)
}

func (sg *StreamGroup) SetIPv4Protocol(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ipv4", "protocol", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4Tos(start, count, step string) error {
	err := sg.EnableOptionalField("ipv4", "tos")
	if err != nil {
		return fmt.Errorf("Cannot set ipv4 tos with %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ipv4", "precedence", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4DSCP(start, count, step string) error {
	err := sg.EnableOptionalField("ipv4", "ds")
	if err != nil {
		return fmt.Errorf("Cannot enable traffic class with %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ipv4", "default_phb", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4TcpSourcePort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("tcp", "source_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4TcpDestinationPort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("tcp", "destination_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4UdpSourcePort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("udp", "source_port", "0", start, count, step)
}

func (sg *StreamGroup) SetIPv4UdpDestinationPort(start, count, step string) error {
	return sg.SetFieldIncrementingValueRange("udp", "destination_port", "0", start, count, step)
}

func (sg *StreamGroup) SetDestinationMacAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ethernet", "destination_address", plen, start, count, step)
}

func (sg *StreamGroup) SetSourceMacAddress(start, plen, count, step string) error {
	return sg.SetFieldIncrementingValueRange("ethernet", "source_address", plen, start, count, step)
}

func (sg *StreamGroup) SetVlan(start, count, step string) error {
	err := sg.EnableOptionalField("ethernet", "vlan_tag1")
	if err != nil {
		return fmt.Errorf("Cannot set vlan with: %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ethernet", "vlan_id1", "0", start, count, step)
}

func (sg *StreamGroup) UnsetVlan() error {
	err := sg.SetFieldIncrementingValueRange("ethernet", "vlan_id1", "0", "0", "1", "0")
	if err != nil {
		return fmt.Errorf("Cannot unset vlan with: %s", err)
	}
	return sg.DisableOptionalField("ethernet", "vlan_tag1")
}

func (sg *StreamGroup) SetCos(start, count, step string) error {
	err := sg.EnableOptionalField("ethernet", "vlan_tag1")
	if err != nil {
		return fmt.Errorf("Cannot set cos with: %s", err)
	}
	return sg.SetFieldIncrementingValueRange("ethernet", "vlan_user_priority1", "0", start, count, step)
}

func (sg *StreamGroup) UnsetCos() error {
	err := sg.SetFieldIncrementingValueRange("ethernet", "vlan_user_priority1", "0", "0", "1", "0")
	if err != nil {
		return fmt.Errorf("Cannot unset cos with: %s", err)
	}
	return sg.DisableOptionalField("ethernet", "vlan_tag1")
}

//AgtInvoke AgtRawPdu GetAllPduBytes $hPdu
func (sg *StreamGroup) GetAllPduBytes() (string, error) {
	cmd := fmt.Sprintf("AgtRawPdu GetAllPduBytes %s", sg.PDUHandler)
	res, err := sg.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get all byte of %s with: %s", sg.PDUHandler, err)
	}

	res = strings.TrimSpace(res)
	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	if res == "" {
		return "", fmt.Errorf("Cannot get all byte of %s with: %s", sg.PDUHandler, err)
	}

	return res, nil
}

func (sg *StreamGroup) SetIPv4UDP() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet ipv4 udp}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv4 udp: %s", err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv4TCP() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet ipv4 tcp}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv4 tcp: %s", err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv6TCP() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet ipv6 tcp_v6}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv6 tcp: %s", err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv6UDP() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet ipv6 udp_v6}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv6 udp: %s", err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv6ND() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet ipv6 icmp_v6}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv6 tcp: %s", err)
	}

	return nil
}

func (sg *StreamGroup) SetIPv4ARP() error {
	cmd := fmt.Sprintf("%s SetPduHeaders %s {ethernet arp}", sg.Object, sg.PDUHandler)
	_, err := sg.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set stream as ipv6 tcp: %s", err)
	}

	return nil
}

func (tr *Traffic) Sync() error {
	if tr.Handler == "" {
		return fmt.Errorf("Cannot sync traffic with not initialized")
	}

	_, err := tr.GetName()
	if err != nil {
		return fmt.Errorf("CAnnot sync traffic with: %s", err)
	}

	_, err = tr.GetType()
	if err != nil {
		return fmt.Errorf("CAnnot sync traffic with: %s", err)
	}

	_, err = tr.GetMode()
	if err != nil {
		return fmt.Errorf("CAnnot sync traffic with: %s", err)
	}

	_, err = tr.GetAverageLoad()
	if err != nil {
		return fmt.Errorf("CAnnot sync traffic with: %s", err)
	}

	return nil
}

func (sg *StreamGroup) Sync() error {
	_, err := sg.GetName()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetType()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetPdu()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetLength()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetLengthMode()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetDefaultL2Protocol()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetSourcePort()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetNumberOfStreams()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetSourceEndpointType()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetSourceEndpoint()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.ListAllPacketTypes()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.ListAllL2Protocols()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetExpectedDestinationPorts()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.IsEnabled()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	_, err = sg.GetAllPduBytes()
	if err != nil {
		return fmt.Errorf("Cannot sync sg %s with: %s", sg.Handler, err)
	}

	return nil
}

func (tr *Traffic) SetPPS(pps string) error {
	return tr.SetAverageLoad(PACKETS_PER_SEC, pps)
}

func (tr *Traffic) SetMPS(mps string) error {
	return tr.SetAverageLoad(MBITS_PER_SEC, mps)
}

func verifyIPv4PrefixLen(plen string) string {
	lent, _ := strconv.ParseInt(plen, 10, 64)
	return fmt.Sprintf("%d", 32-lent)
}

func verifyIPv6PrefixLen(plen string) string {
	lent, _ := strconv.ParseInt(plen, 10, 64)
	return fmt.Sprintf("%d", 128-lent)
}

func (tr *Traffic) SetStreamsDstIP(ip, plen, step, count string) error {
	plen = verifyIPv4PrefixLen(plen)
	return tr.DefaultStreamGroup.SetFieldIncrementingValueRange("ipv4", "destination_address", plen, ip, count, step)
}

func (tr *Traffic) SetStreamsSrcIP(ip, plen, step, count string) error {
	plen = verifyIPv4PrefixLen(plen)
	return tr.DefaultStreamGroup.SetFieldIncrementingValueRange("ipv4", "source_address", plen, ip, count, step)
}

func (tr *Traffic) SetStreamsDstIP6(ip, plen, step, count string) error {
	plen = verifyIPv6PrefixLen(plen)
	return tr.DefaultStreamGroup.SetFieldIncrementingValueRange("ipv6", "destination_address", plen, ip, count, step)
}

func (tr *Traffic) SetStreamsSrcIP6(ip, plen, step, count string) error {
	plen = verifyIPv6PrefixLen(plen)
	return tr.DefaultStreamGroup.SetFieldIncrementingValueRange("ipv6", "source_address", plen, ip, count, step)
}

func (tr *Traffic) SetStreamsDstMAC(mac, count string) error {
	return tr.DefaultStreamGroup.SetDestinationMacAddress(mac, "0", count, "1")
}

func (tr *Traffic) SetStreamsSrcMAC(mac, count string) error {
	return tr.DefaultStreamGroup.SetSourceMacAddress(mac, "0", count, "1")
}

func (tr *Traffic) SetStreamsVLAN(vid, count string) error {
	return tr.DefaultStreamGroup.SetVlan(vid, count, "1")
}

func (tr *Traffic) SetStreamsIPProtocol(proto, count string) error {
	return tr.DefaultStreamGroup.SetIPv4Protocol(proto, count, "1")
}

func (tr *Traffic) SetStreamsIPv6NextHeader(nh, count string) error {
	return tr.DefaultStreamGroup.SetIPv6NextHeader(nh, count, "1")
}

func (tr *Traffic) AddStatistics() error {
	if tr.DstPort == nil {
		return fmt.Errorf("Traffic destination port is not set, cannot init statistics")
	}

	cmd := fmt.Sprintf("AgtStatisticsList Add AGT_STATISTICS")
	res, err := tr.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot add statistics with for traffic %s with %s", tr.Name, err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	if res == "" {
		return fmt.Errorf("Cannot add statistics for traffic %s with cannot add statistics", tr.Name)
	}

	st := &Statistics{
		Name:    tr.StatisticsName(),
		Handler: res,
		Port:    tr.Port,
		Traffic: tr,
		Type:    "Stream",
	}

	err = st.Init()
	if err != nil {
		return fmt.Errorf("Cannot add statistics %s with %s", err)
	}

	tr.Statistics = st

	return nil
}

func (tr *Traffic) StatisticsName() string {
	return fmt.Sprintf("ATSN2X_STREAM_%s", tr.Name)
}

func (tr *Traffic) GetStatistics() (uint64, uint64, uint64, error) {
	return tr.Statistics.GetStatistics()
}
