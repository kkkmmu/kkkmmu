1. PORT.MPLS_ENABLE， VLAN_MPLS.MPLS_ENABLE
2. 如果入方向报文命中MY_STATION_TCAM:
	1. 如果该报文包含两个及以上LABEL的话查找L3_TUNNEL表, 如果命中的话继续查找MPLS_ENTRY表。
	2. 否则查找MPLS_ENTRY表。
3. 从L3_TUNNEL/MPLS_ENTRY表我们可以获取到下一跳index：NEXT_HOP_INDEX。
	NEXT_HOP_INDEX一般用来同时索引两张表(ING_L3_NEXT_HOP, EGR_L3_NEXT_HOP)。
	ING_L3_NEXT_HOP: 我们从该表可以获取出端口, 出MODID。
	EGR_L3_NEXT_HOP: 
		1. EGR_MAC_DA_PROFILE： 可以获取到报文接下来的目的MAC。
		2. EGR_MPLS_VC_AND_SWAP_LABEL_TABLE: 可以获取报文接下来的LABEL。
		3. EGR_L3_INTF: 可以获取报文接下来的源MAC以及VID.
		4. EGR_IP_TUNNEL_MPLS: 用在LER上来为报文出方向加多个LABEL.

4. The device performs the first MPLS label lookup into the L3_TUNNEL table if the incoming packet has three or more lables. If an entry is found, device pops the outmost lable. Next device performs two lookups into the MPLS_ENTRY tables and removes pop two MPLS labels. The MPLS_ENTRY table provides one of the following:
	1. MPLS_ACTION_IF_BOS = PHP_NHI for the BOS label. Logic terminates the MPLS header and forwards the L3 packet to the next hop (No Routing).
	2. MPLS_ACTION_IF_BOS=L3_INTF for the BOS label. Logic terminates the MPLS header, perform L3 Loopkups and route the packet.
	3. MPLS_ACTION_IF_BOS=L3_ECMP for the BOS lable. Logic terminates the MPLS header, obtains the destination port using the ECMP logic and route the pakcets.


5. L3 MPLS deals with the estabilishment of MPLS tunnels. At the MPLS tunnel initiation point, packets entering the tunnel are encapsulated with an MPLS header (label PUSH operation). IPv4 and IPv6 packets can be directed into an MPLS tunnel based on the DIP lookup. Along the path of the MPLS tunnel, MPLS label SWAP, PUSH, PHP, and POP operations are supported. At the end of an MPLS tunnel, forwarding of IPv4 and IPv6 packets is based on DIP loopkup. Virtual routing tables are supported for the DIP lookup.

6. In the MPLS APIs, the concept of VPN is introduced. A VPN can be of type VPLPS, VPWS or L3:
	L3 VPNs are used to indentify a virtual routing table. When creating an L3 VPN, a VRF must be specified. At the end of MPLS tunnel, the IPv4 or IPv6 DIP lookup is qualified by VRF associated with the L3 VPN.


7. Ingress PE (Provider Edge) refers to the PE netowrk entity or node where IP packets or Layer 2 frames are encapsulated with an MPLS header. There may be more than one MPLS label (header) on the outging packet. Tunnel Midpoint or Provider(P) node typically swaps the outer MPLS label with antother MPLS label. Egress PE refers to the PE (Provider Edge) node where the outer MPLS label(s) is popped and the inner payload is forwarded optionally based on L2 （MAC address) lookup or IP lookup.
	这里针对不同的设备 （PE/P），需要的API也不同.
	1. Ingress PE: API used is mpls_tunnel_initiator_create or mpls_tunnel_initiator_set
	2. Tunnel Midpoint: API used is mpls_tunnel_switch_create or mpls_tunnel_switch_add.
	3. Penultimate to Egress PE (PHP): API used is mpls_tunnel_switch_create or mpls_tunnel_switch_add.
	4. Egress PE: API used is mpls_tunnel_switch_create or mpls_tunnel_switch_add.
  也就是说，ILM和FTN需要调用不同的SDK API. 
		ILM主要用来给为带标签的报文在出方向添加标签，所以用mpls_tunnel_initiator_create/add.
		FTN主要用来给入方向带标签的报文根据规则处理它所带的标签（也就是标签交换），所以用mpls_tunnel_switch_create/add.

8. mpls_tunnel_initiator_create/add
	This API is called to set the tunnel initiator parameters for an L3 interface. The label_arrary contains information for forming the MPLS labels to be pushed onto a packet entering the tunnel. The destination port information for the tunnel is specified by calling l3_egress_create API to create an egress object.
	The steps to compeletly setup an MPLS tunnel initiator are:
		1. l3_intf_create
		2. mpls_tunnel_initiator_set
		3. l3_egress_create.
	The egress object handle returned by l3_egress_create can be used in l3_route_add to direct IPv4 packets into an MPLS tunnel.


9. mpls_tunnel_switch_create/add
	The mpls_tunnel_switch_t structure is used to specify an entry in the MPLS label table. The ingress parameters specify the entry key, action (SWAP, POP, PHP), and QOS settings. The entry key is the MPLS label and optionally the incoming port. The action_if_bos and action_if_not_bos parameters can be used on certain devices to specify different lable actions for the cases of bottom-of-stack(BOS) labels and non-BOS labels. The vpn parameter is used only if the specified action is POP. This is an l3 VPN used to get the VRF to used for IPv4 and IPv6 DIP lookup.
	There are 2 modes for SWAP/LSR configuration:
		Mode-1: The egress paramter egress_label is only used if the action is SWAP. This parameter is used to specify the SWAP label inforation (Label, EXP, ttL). The egress_intf paramter is only used for SWAP and PHP operations. It points to an egress object which was created useing l3_egress_create. The egress opbject can be an MPLS tunnel, for example, if the desired operation is SWAP and PUSH. The MPLS information (Label, EXP, TTL) within the Egress object must be initialized to 0.
		Mode-2: The egress_label must be initialized to 0, and egress_if point to an egress object which must contain the SWAP lable inforation (label, EXP, TTL).
	Although Mode-1 is maintained, Mode-2 is recommanded for several reasons. The nexthop allocation happens during the creation of Egress object. No new next hops are allocated during tunnel switch creation. Multiple tunnel switch entries can point to the same egress object.

10. An MPLS domain consists of a group of MPLS-enabled routers, called Label Switching Routers(LSRs). In an MPLS doamin, packets are forwarded from on MPLS-enabled router to another along a predetermined path, called an LSP. LSPs are one-way paths between MPLS-enabled routers on a network. To provide two-way traffic, the use configures LSPs in each direction.
	The LSRs at the headend and tailend of an LSP are known as Label Edge Routers(LERs). The LER at the headend, where packets enter the LSP, is known as the ingress LER, the LER at the tailed, where packets exit the LSP, is known as the egress LER. Each LSP has one ingress LER and on egress LER. Packets in an LSP flow in one direction; from the ingress LER torwards the egress LER. In between the ingress and egress LERs there may be zero or more transit LSRs. A device enabled for MPLS can perform the role of ingress LER, transit LSR, or egress LER in an LSP. Furter, a device can serve simultaneously as an ingress LER for one LSP, transit LSR for another LSP, and egress LER for some other LSP.
	Label switching in an MPLS domain depicts an MPLS domain with a single LSP consisting of three LSRs: an ingress LSR, a transite LSR, and an egress LSR.

11. Label switching in an MPLS domain works as described below.
	1. The ingress LER receives a packet and pushes a label onto it.
		When a packet arrives on an MPLS-enabled interfae, the device determines to which LSP(if any) the packet are assigned. Specifically, the device determins to which Forwarding Equiavlence Class(FEC) the packet belongs. An FEC is simply a group of packets that are all forwarded in the same way. For example, a FEC could be defined as all packets from a given Virtual Leased Line(VLL). FECs are mapped to LSPs. When a packet belongs to a FEC, and an LSP is mapped to that FEC, the packet is assigned to the LSP.
		When a packet is assigned to an LSP, the device, acting as an ingress LER, applies(pushes) a tunnel label onto the packet. A label is a 32-bit, fixed-length identifier that is significant only to MPLS. From this point until the packet reaches the egress LER at the end of the path, the packet is forwarded using information in its label, not information in its IP header. The packet's IP header is not examined aging as long as the packet traverses the LSP. The ingress LER may also apply a VC label onto the packet based on the VPN application.
		On the ingress LER, the label is associated with an outbound interface. After receiving a label, the packet is forwarded over the outboud interface to the next router in the LSP.
	2. A transit LSR receives the labled packet, swaps the lable, and forwards the packet to the next LSR.
		In an LSP, zero or more transit LSRs can exist between the ingress and egress LSRs. A transit LSR swaps labes on an MPLS packet and forwards the packt to the next router in the LSP.
		When a transit LSR receive an MPLS packet, it looks up the labe in its MPLS forwarding table. This table maps the label and inbound interface to a new label and outboud interface. The transit LSP replaces the ould lable with the new lable and sends the packet out the outbound interface specifed in the table. This process repeats at each transit LSR until the packet reaches the next-to-last LSR in the LSP.
	3. The egress LER receives labeled packet, pops lable, and forwards IP packet.
		When a packets reaches the egress LER, the MPLS lable is removed(called popping the label), and the packet can then be forwarded to its destination using standard hop-by-hop routing protocols.
