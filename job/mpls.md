Baisc problem should be resolved in a L3 VPN technology.
1.Easy to extened.
2.Address overlapping.
3.Multi-homed site.

Label Swapping
    Label swapping is the use of the following procedures to forward a packet.
	In order to forward a labeled packet, a LSR examines the label at the top of the label stack.  It uses the ILM to map this label to an NHLFE.  Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet's label stack.  It then encodes the new label stack into the packet, and forwards the result.
	In order to forward an unlabeled packet, a LSR analyzes the network layer header, to determine the packet's FEC.  It then uses the FTN to map this to an NHLFE.  Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet's label stack.  (Popping the label stack would, of course, be illegal in this case.)  It then encodes the new label stack into the packet, and forwards the result.

	IT IS IMPORTANT TO NOTE THAT WHEN LABEL SWAPPING IS IN USE, THE NEXT HOP IS ALWAYS TAKEN FROM THE NHLFE; THIS MAY IN SOME CASES BE DIFFERENT FROM WHAT THE NEXT HOP WOULD BE IF MPLS WERE NOT IN USE.


MPLS Basic Sequence:
1. CE advertises routes to the local PE via some routing protocol.
2. The local PE marks these routes with a particular extended community (route target) and advertises them in BGP.
3. The routes are distributed to all remote PE by BGP.
4. Remote PE receives BGP routes, filters them based on the community and advertises them to the CE.

The P routers carry all VPN routes, so the addresses used in the VPNs need to be unique in the provider's network.

VPN V4 Address: Contrustructed by concatenating an IP address and an 8 byte unique identifier called the route distinguisher.
	VPN-V4 address is advertised in a special address family by BGP(MP-BGP)
	VPN-V4 address is used only in the provider's network.
	VPN-V4 address is used only in the control plane.
	VPN-V4 address: The translation from IP adresses to VPN-V4 addresses happens on the PE.
	VPN-V4 address cannot be used for route filtering.(use communities to do that)

Why MPLS ?
	VPN-IP addresses are used by the routing protocols, but do not appear in headers of IP packets.
	So we need a way to forward traffic along routes to VPN-IP addresses. MPLS decouples forwarding from the destination information.

VPN-LABEL:
	"VPN-IP addresses are used by the routing protocols, but do not appear in headers of IP packets."
	The IDEA: Use a label to identrify the next-hop at the remote PE. Also called VPN label.
	The label is distributed by BGP, along with the VPN-IP address.
	Traffic will carry two labels, the VPN label and the LSP label.
	The remote PE makes the forwarding decision based on the VPN label.

MPLS VPN model summary:
	1. P routers don't need to maintain VPN routes at all. Only need to maintain routes to other P and PE routers.
	2. PE routers maintain VPN routes, but only for VPNs that have sites attached to them.
	3. VPNs can have overlapping address spaces.

MPLS VPN basic concepts:
	1. Use MPLS to forward traffic across nodes that don't have routing information for the packet's final destination.
	2. Use a label to mark the traffic. Use this makring to determine the next-hop.
	3. The address of the next-hop in the BGP advertisement provides coupling between the VPN routes and the internal routing to the remote PE.

MPLS VPN Scaling properties:
	1. Only one routing peering (CE-PE), regardless of the number of sites in the VPN.
	2. The customer doesn't need routing skills. A customer doesn't need to operate its own backbone.
	3. Adding a new site requires configuration of one PE regardless of the number of site.
	4. PE has to maintain routes only for the VPNs to which it is connected.
	5. P routers don't have to maintain VPN routes at all.
	6. Can use overlapping address spaces - efficient use of private IP addresses.
	7. Route distinguishers are structured so that each service provider can manage its own number spsace.

ILM: Incoming Lable Map 
FEC: Fowarding Equivalence Class
NHLFE: Next-Hop Lable Forwarding Entry 
FTN: FEC-To-NHLFE
LSR: Label Switching Router
MPLS-TP: MPLS Transport Profile

In conventional IP forwarding, a particular router will typically consider two packets to be in the same FEC if there is some address prefix X in that router's routing tables such that X is the "longest match" for each packet's destination address.  As the packet traverses the network, each hop in turn reexamines the packet and assigns it to a FEC.

In MPLS, the assignment of a particular packet to a particular FEC is done just once, as the packet enters the network.  The FEC to which the packet is assigned is encoded as a short fixed length value known as a "label".  When a packet is forwarded to its next hop, the label is sent along with it; that is, the packets are "labeled" before they are forwarded.

At subsequent hops, there is no further analysis of the packet's  network layer header.  Rather, the label is used as an index into a table which specifies the next hop, and a new label.  The old label is replaced with the new label, and the packet is forwarded to its next hop.

In the MPLS forwarding paradigm, once a packet is assigned to a FEC, no further header analysis is done by subsequent routers; all forwarding is driven by the labels.  This has a number of advantages over conventional network layer forwarding.
	-  MPLS forwarding can be done by switches which are capable of doing label lookup and replacement, but are either not capable of analyzing the network layer headers, or are not capable of analyzing the network layer headers at adequate speed.



Labels:
     A label is a short, fixed length, locally significant identifier which is used to identify a FEC.  The label which is put on a particular packet represents the Forwarding Equivalence Class to which that packet is assigned.
	 Most commonly, a packet is assigned to a FEC based (completely or partially) on its network layer destination address.  However, the label is never an encoding of that address.
	 If Ru and Rd are LSRs, they may agree that when Ru transmits a packet to Rd, Ru will label with packet with label value L if and only if the packet is a member of a particular FEC F.  That is, they can agree to a "binding" between label L and FEC F for packets moving from Ru to Rd.  As a result of such an agreement, L becomes Ru's "outgoing label" representing FEC F, and L becomes Rd's "incoming label" representing FEC F.
	 Note that L does not necessarily represent FEC F for any packets other than those which are being sent from Ru to Rd.  L is an arbitrary value whose binding to F is local to Ru and Rd.


Upstream and Downstream LSRs
   Suppose Ru and Rd have agreed to bind label L to FEC F, for packets  sent from Ru to Rd.  Then with respect to this binding, Ru is the "upstream LSR", and Rd is the "downstream LSR".
   To say that one node is upstream and one is downstream with respect to a given binding means only that a particular label represents a particular FEC in packets travelling from the upstream node to the downstream node.  This is NOT meant to imply that packets in that FEC would actually be routed from the upstream node to the downstream node.

Label Assignment and Distribution
   In the MPLS architecture, the decision to bind a particular label L to a particular FEC F is made by the LSR which is DOWNSTREAM with respect to that binding.  The downstream LSR then informs the upstream LSR of the binding.  Thus labels are "downstream-assigned", and label bindings are distributed in the "downstream to upstream"direction. If an LSR has been designed so that it can only look up labels that fall into a certain numeric range, then it merely needs to ensure that it only binds labels that are in that range.


Label Distribution Protocols
   A label distribution protocol is a set of procedures by which one LSR informs another of the label/FEC bindings it has made.  Two LSRs which use a label distribution protocol to exchange label/FEC binding information are known as "label distribution peers" with respect to the binding information they exchange.  If two LSRs are label distribution peers, we will speak of there being a "label distribution adjacency" between them.

The Label Stack
   So far, we have spoken as if a labeled packet carries only a single label.  As we shall see, it is useful to have a more general model in which a labeled packet carries a number of labels, organized as a last-in, first-out stack.  We refer to this as a "label stack".
   Although, as we shall see, MPLS supports a hierarchy, the processing of a labeled packet is completely independent of the level of hierarchy.  The processing is always based on the top label, without regard for the possibility that some number of other labels may have been "above it" in the past, or that some number of other labels may be below it at present.
   An unlabeled packet can be thought of as a packet whose label stack is empty (i.e., whose label stack has depth 0).
   If a packet's label stack is of depth m, we refer to the label at the bottom of the stack as the level 1 label, to the label above it (if such exists) as the level 2 label, and to the label at the top of the stack as the level m label.

 The Next Hop Label Forwarding Entry (NHLFE)
    The "Next Hop Label Forwarding Entry" (NHLFE) is used when forwarding a labeled packet.  It contains the following information:
	      1. the packet's next hop
		  2. the operation to perform on the packet's label stack; this is one of the following operations:
				 a) replace the label at the top of the label stack with a specified new label
				 b) pop the label stack
				 c) replace the label at the top of the label stack with a specified new label, and then push one or more specified new labels onto the label stack.
			 It may also contain:
				 d) the data link encapsulation to use when transmitting the packet
			     e) the way to encode the label stack when transmitting the packet
				 f) any other information needed in order to properly dispose of the packet.
    Note that at a given LSR, the packet's "next hop" might be that LSR itself.  In this case, the LSR would need to pop the top level label, and then "forward" the resulting packet to itself.  It would then make another forwarding decision, based on what remains after the label stacked is popped.  This may still be a labeled packet, or it may be the native IP packet.
	This implies that in some cases the LSR may need to operate on the IP header in order to forward the packet.
	If the packet's "next hop" is the current LSR, then the label stack operation MUST be to "pop the stack".

Incoming Label Map (ILM)
	The "Incoming Label Map" (ILM) maps each incoming label to a set of NHLFEs. It is used when forwarding packets that arrive as labeled packets.
	If the ILM maps a particular label to a set of NHLFEs that contains more than one element, exactly one element of the set must be chosen before the packet is forwarded.  The procedures for choosing an element from the set are beyond the scope of this document. Having the ILM map a label to a set containing more than one NHLFE may be useful if, e.g., it is desired to do load balancing over multiple equal-cost paths.

FEC-to-NHLFE Map (FTN)
	The "FEC-to-NHLFE" (FTN) maps each FEC to a set of NHLFEs.  It is used when forwarding packets that arrive unlabeled, but which are to be labeled before being forwarded.
	If the FTN maps a particular label to a set of NHLFEs that contains more than one element, exactly one element of the set must be chosen before the packet is forwarded.  The procedures for choosing an element from the set are beyond the scope of this document.  Having the FTN map a label to a set containing more than one NHLFE may be useful if, e.g., it is desired to do load balancing over multiple equal-cost paths.

Label Swapping
    Label swapping is the use of the following procedures to forward a packet.
	In order to forward a labeled packet, a LSR examines the label at the top of the label stack.  It uses the ILM to map this label to an NHLFE.  Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet's label stack.  It then encodes the new label stack into the packet, and forwards the result.
	In order to forward an unlabeled packet, a LSR analyzes the network layer header, to determine the packet's FEC.  It then uses the FTN to map this to an NHLFE.  Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet's label stack.  (Popping the label stack would, of course, be illegal in this case.)  It then encodes the new label stack into the packet, and forwards the result.

	IT IS IMPORTANT TO NOTE THAT WHEN LABEL SWAPPING IS IN USE, THE NEXT HOP IS ALWAYS TAKEN FROM THE NHLFE; THIS MAY IN SOME CASES BE DIFFERENT FROM WHAT THE NEXT HOP WOULD BE IF MPLS WERE NOT IN USE.










show mpls cross-connect-table
	Cross connect ix: Display the table index for the crosss-connect.
	in labelspace:    Indicates that all MPLS interfaces will use the platform wide label space("0")
	in label:         Displayes the ingress (incoming interface) lable for this segment.
	out-segment ix:   Displays the outbound segment index.
	Owner:            Displays the creator of this segment, typically a protocol such as BGP.
	Persistent:       Displays whether the tunnel is persitent - Yes or No.
	Admin Status:     Indicates whether the user can administratively disable a peer while still preserving its configuration.
	Oper Status:      Displays the current status of the cross-connect segment- Up or Down.

show mpls ftn
	Prefix/mask:      Displays the IP address and mask stored in for this FEC-to-NHLFE table entry.
	NHLFE ix:         Displays the index number for the Next-Hop Label Forwarding Etnry
	opcode:           PUSH - Replace the top label with another and the push one or more additional labels onto the label stack.
		              SET  - Set the next hop label.
	label/ifindex:    Displays the label associated with the interface.
	nh-addr:          Displays the IP address of the next-hop.

show mpls ilm
	Label:            Displays the label ID for this enty in the Incoming Label Map table.
	Opcode:           POP            - Remove label from packet.
					  CONTEXT-CHANGE - 
					  DELIVER        -
	nhlfe-ix/contex-id: Displays the Next-Hop Lable Forwarding Entry (NHLFE) index or context ID for this entry.


The MPLS/VPN architecture uses BGP routing protocol in two different ways: 
    VPNv4 routes are propagated across a MPLS/VPN backbone using multi-protocol BGP between the PE routers 
    BGP can be used as the PE-CE routing protocol to exchange VPN routes between the provider edge routers and the customer edge routers 

The address-family router configuration command is used to select the routing context that you’d like to configure:     
    Internet routing (global IP routing table) is the default address family that you configure when you start configuring the BGP routing process;     
	To configure multi-protocol BGP sessions between the PE routers, use the vpnv4 address family     
	To configure BGP between the PE routers and the CE routes within individual VRF, use the ipv4 vrf name address family 

MPLS/VPN architecture defines two types of BGP neighbors: 
	Global BGP neighbors (other PE routers), with which the PE router can exchange multiple types of routes. 
		These neighbors are defined in the global BGP definition and only have to be activated for individual address families 
	Per-VRF BGP neighbors (the CE routers) which are configured and activated within the ipv4 vrf name address family 

BGP connectivity between two PE routers is configured in four steps: 	
	The remote PE router is configured as global BGP neighbor under BGP router configuration mode 	
	Parameters that affect the BGP session itself (for example, source address for the TCP session) are defined on the global BGP neighbor 	
	VPNv4 address family is selected and the BGP neighbor is activated for VPNv4 route exchange 	
	Additional VPNv4-specific BGP parameters that affect the VPNv4 routing updates (filters, next-hop processing, route-maps) are configured within the VPNv4 address family

MPLS/VPN architecture has introduced the extended community BGP attribute. BGP still supports the standard community attribute, which has not been superseded with the extended communities. The default community propagation behavior for standard BGP communities has not changed – community propagation still needs to be configured manually. Extended BGP communities are propagated by default, because their propagation is mandatory for successful MPLS/VPN operation. The neighbor send-community command was extended to support standard and extended communities. You should use this command to configure propagation of standard and extended communities if your BGP design relies on usage of standard communities (for example, to propagate Quality of Service information across the network). 


The LDP (or TDP) protocol is used to establish label switched paths between all PE routers. 
The P and PE routers must exchange the LDP (or TDP) protocol. Each IGP route will be assigned a label. This label is propagated to each neighbor. The labels used for forwarding are stored in the LFIB. 


The import function into vrf A is verified by the command show ip bgp vpnv4 vrf A.


DASAN Debug Command:
	show ip bgp vpnv4 vrf VRF label 查看PE对每个VRF的每个PREFIX分配的LABEL.
	show mpls forwarding-table
	show mpls ilm-table
	show mpls ftn-table
	show mpls ftn statistics 
	show mpls out-segment-table
	show mpls in-segment-table
	show ip vrf VRF 可以查看每个VRF的LABEL分配模式(Per VRF/Per Interface/Per Prefix)

FTN Entry:
	FTN represents FEC-to-NHLFE, a logical table that implements the MPLS architecture defined in RFC 3031: 
		Each FTN entry is used to map incoming traffic to an MPLS LSP. This information is assigned to NHLFE at the edge of the MPLS cloud. NHLFE represents the Next-Hop-Label-Forwarding-Entry, which specifies the MPLS properties for egressing a packet onto LSP. MPLS properties include lable, nexthop IP address and outgoing interface.
		FTN is an entity that is present only on LSP ingress. It is used for Pushing lable information to native packets and tunneling them via LSP. Howerver, in MPLS-TP, incoming packets are not directly mapped to MPLS-TP LSP. Instead, service such as PW map, indirectly route map traffic to the LSP. 
		For MPLS-TP, FTN is only used to specify NHLFE information the MPLS-TP LSP at the LSP ingress. FTN implements the forward-componet of an LSP on ingress and reserse component of LSP on egress.

ILM Entry and Reverse ILM Entry:
	ILM represents the Incoming Label Map. ILM is a logical table the is indexed by the incoming interface and label. An ILM entry specifies behavior for processing labeled packets that arrive on MPLS core and egress nodes.
		1. On core nodes, the behavior is to swap the lables.
		2. On egress nodes, the behavior is to POP the label and process the native packet.
		3. For bidirectional LSPs, ILM entries performing label-POP implement the forwarding component of LSPs on the egress components of LSPs on ingress.
		4. ILM entries performing label-SWAP implement the forward and revers components of LSPs on transit/core nodes.


MPLS is a lable swapping and forwarding technology in which every packet contains a short, fixed-length label. Routers use these label values to forward incoming packets to their destination.

MPLS divides its functions into two distinct categories:
		1. Assigning and exchanging labels between Label Switching Routers (LSRs) through the Label Distribution Protocol (LDP).
		2. Forwarding labeled packets. The MPLS forwarder swaps incoming labels with outgoing labels and forwards the resulting packets to the outgoing interface. MPLS supports the Ethernet interface.


Internal Architecture
    The MPLS Forwarder processes incoming packets from all the network interfaces. Its primary function is to forward traffic over data paths created by label distribution protocols such as LDP and RSVP-TE.
	The MPLS Forwarder handles the following functions:
	   • Receiving IP/MPLS unicast packets from the interface queues.
	   • Receiving Layer-2 packets from the interface queues.
	   • Forwarding labeled/unlabeled unicast packets based on IP address or MPLS labels.
	   • Forwarding labeled/unlabeled Layer-2 frames based on incoming interface or MPLS labels.
	   • Fragmenting packets exceeding the MTU size of the outgoing interface.
    The MPLS Forwarder is made up of the following separate entities:
       • A global FTN (FEC to Next-Hop-Label-Forwarding-Entry) table. The kernel interfaces use this table when processing non-labeled packets (IP packets) and when the kernel interface is not bound to a VRF table. Multiple kernel interfaces may use one global FTN table; or each interface may use its own table.
	   • One or more ILM (Incoming Label Map) tables. The kernel interfaces use these tables to process labeled packets and the label space in the interface contains a positive integer. An ILM table is created per label-space used. Therefore, if all interfaces in the system are using the same label-space, only one ILM table is created.
	   • One or more VRF tables. The kernel interfaces use these tables to process non-labeled packets and the interface is bound to a VRF table. Many kernel interfaces may use one VRF table; or each interface may use its own table.
	   • One or more interfaces with the flexibility to be enabled for either MPLS or VRF, or both.
	   If an interface is not enabled for MPLS or VRF, all labeled packets are dropped.



Valid Label Ranges
In the scheme of labels, there are several reserved labels and reserved label ranges. For the detailed explanation of these see RFC 3032 section 2.1.
The range 0-3 is reserved for the following uses:
      • 0 is the IPv4 Explicit NULL Label. When the Forwarder receives a packet with this label value at the bottom of the label stack, the stack pops the label value; it bases all forwarding of this packet on the IPv4 header.
	  • 1 is the Router Alert Label. When the Forwarder receives a packet with this label value at the top of the label stack, the forwarder uses the next label beneath it in the stack. The Router Alert Label is pushed on top of the stack if further forwarding is required.
	  • 2 is the IPv6 Explicit NULL Label. When the Forwarder receives a packet with this label value at the bottom of the label stack, it pops the label stack and forwards the packet based on the IPv6 header.
	  • 3 is the Implicit NULL Label. When the Forwarder receives a packet with this label value at the top of the label stack, the LSR pops the stack.
The range 4-15 is reserved for future use.



Ethernet Interface
    Upon receipt of a frame by the Data Link device driver, the frame is passed to the MPLS module, which handles the processing of all packets of type IP or MPLS. The following steps are then taken:
        1. Determine if packet is labeled: if the incoming packet is labeled, this is an MPLS packet; go to 2. Otherwise, the packet is an IP packet; go to 3
		2. Use the top label in the packet to look up the destination in ILM table that this interface is bound to. If the lookup finds an outgoing label, go to 4. Else, drop the packet and exit this function.
		3. By employing the best-match principle, use the destination IP address to determine the FEC that this destination address belongs to. Using this FEC as the key, lookup in the FTN table for a valid outgoing label. If we have a valid label for the destination address, push the outgoing label found in the FTN table onto the packet; continue with 6. Otherwise go to 5.
		4. 	If there is no outgoing label for the incoming label, the LSR is an egress for the current LSP, and the packet shall be routed using traditional, native routing; continue with 5. Otherwise, push the mapped label onto the packet; continue with 8.
		5. 	If the opcode associated with this ILM entry was POP_FOR_VC, go to 6. Else go to 7.
		6. 	Remove the MPLS shim and pass the Ethernet frame that had been encapsulated with the MPLS shim to the outgoing interface's controller.
		7. 	Decrement the TTL fields of the labeled or unlabeled IP packet, use the kernel URF to route the packet using conventional routing, then exit this function.
		8. Decrement the TTL fields of the packet and label-switch the packet; then exit this function.
    The following steps are taken in an environment where LDP, BGP and the MPLS Forwarder all exist together:
        1. Determine if packet is labeled: if the incoming packet is labeled, this is an MPLS packet; go to 2. Otherwise, the packet is an IP packet; go to 3
		2. Determine which ILM table to use to look up outgoing labels: if the interface that accepts incoming packets is bound to an ILM table (meaning that the label-space-value is not zero), use that ILM table. Else, because no ILM table is bound to this interface, use the global ILM table.Use the top label in the packet as a key to lookup an outgoing label in either the bound or global ILM table.If the lookup finds an outgoing label, go to 4. Else, drop the packet and exit this function.
		3. By employing the best-match principle, use the destination IP address to figure out the FEC that this destination address belongs to. Using this FEC as the key, lookup in the FTN table for a valid outgoing label. The FTN table to be used is decided as follows: If there is a VRF table bound to this interface, lookup in that VRF table. If not, lookup in the global FTN table. If we have a valid label for the destination address, push the outgoing label found in the FTN table onto the packet; continue with 6. Otherwise go to 5.
		4. If there is no outgoing label for the incoming label, the LSR is an egress for the current LSP, and the packet shall be routed using traditional native routing; continue with 5. Otherwise push the mapped label onto the packet; continue with 6.
		5. If the lookup was done in a VRF table, drop the packet and exit this function. Otherwise, decrement the TTL fields of packet and use the kernel URF to route the packet using conventional routing. Exit this function.
		6. Decrement the TTL fields of the and label-switch the packet. Exit this function.
   The following steps are taken in an environment where LDP is being used to set up a Layer 2 Virtual Circuit between remote nodes:
        1. Determine if the interface on which the frame is received is bound to a Virtual Circuit. If not, the packet is handled as per the steps enumerated earlier. If the incoming interface is bound to a Virtual Circuit, go to 2.
		2. Check whether an FTN entry is associated with this interface. If yes, go to 3. If no, drop the packet and exit.
		3. Determine whether the opcode associated with the FTN entry is either PUSH_AND_LOOKUP_FOR_VC or PUSH_FOR_VC. If it is the former, go to 4. Else go to 5.
		4. Using the FTN entry bound to the interface, add an MPLS shim on top of the Ethernet frame received. Using the Virtual Circuit endpoint - the nexthop in the FTN entry - specified, carry out a lookup in the global FTN table for an FTN entry identifying an LSP from this node to the nexthop specified. If no entry is found, drop the packet. Else go to 6.
		5. Using the FTN entry bound to the interface, add an MPLS shim on top of the Ethernet frame received, and pass on the shim+frame to the outgoing interface's Ethernet controller, so that this shim+frame can be encapsulated inside an Ethernet frame. Go to 7.
		6. Using the newly found FTN entry, add to the existing shim that has been added on top of the received Ethernet frame, and then pass this shim+frame to the outgoing interface's Ethernet controller, so that this shim+frame can be encapsulated inside an Ethernet frame. Go to 7.
		7. Forward the newly generated frame.


	In conventional IP forwarding, a particular router will typically   consider two packets to be in the same FEC if there is some address   prefix X in that router’s routing tables such that X is the "longest   match" for each packet’s destination address.  As the packet   traverses the network, each hop in turn reexamines the packet and   assigns it to a FEC.   In MPLS, the assignment of a particular packet to a particular FEC is   done just once, as the packet enters the network.  The FEC to which   the packet is assigned is encoded as a short fixed length value known   as a "label".  When a packet is forwarded to its next hop, the label   is sent along with it; that is, the packets are "labeled" before they   are forwarded.   
	At subsequent hops, there is no further analysis of the packet’s   network layer header.  Rather, the label is used as an index into a   table which specifies the next hop, and a new label.  The old label   is replaced with the new label, and the packet is forwarded to its   next hop.
	A label is a short, fixed length, locally significant identifier   which is used to identify a FEC.  The label which is put on a   particular packet represents the Forwarding Equivalence Class to   which that packet is assigned.



Suppose Ru and Rd have agreed to bind label L to FEC F, for packets sent from Ru to Rd. Then with respect to this binding, Ru is the "upstream LSR", and Rd is the "downstream LSR".   
To say that one node is upstream and one is downstream with respect to a given binding means only that a particular label represents a particular FEC in packets travelling from the upstream node to the downstream node. This is NOT meant to imply that packets in that FEC would actually be routed from the upstream node to the downstream   node.
In the MPLS architecture, the decision to bind a particular label L to a particular FEC F is made by the LSR which is DOWNSTREAM with respect to that binding. The downstream LSR then informs the upstream LSR of the binding. Thus labels are "downstream-assigned", and label bindings are distributed in the "downstream to upstream" direction.

An unlabeled packet can be thought of as a packet whose label stack is empty.


The Next Hop Label Forwarding Entry (NHLFE):
	The "Next Hop Label Forwarding Entry" (NHLFE) is used when forwarding a labeled packet.  It contains the following information:
		1. the packet’s next hop
		2. the operation to perform on the packet’s label stack; this is one of the following operations:      
			a) replace the label at the top of the label stack with a specified new label      
			b) pop the label stack
			c) replace the label at the top of the label stack with a specified new label, and then push one or more specified new labels onto the label stack.
			d) the data link encapsulation to use when transmitting the packet      
			e) the way to encode the label stack when transmitting the packet      
			f) any other information needed in order to properly dispose of the packet.
		If the packet’s "next hop" is the current LSR, then the label stack   operation MUST be to "pop the stack".

Incoming Label Map (ILM)   
	The "Incoming Label Map" (ILM) maps each incoming label to a set of NHLFEs.  It is used when forwarding packets that arrive as labeled packets.   
	If the ILM maps a particular label to a set of NHLFEs that contains  more than one element, exactly one element of the set must be chosen before the packet is forwarded.  The procedures for choosing an element from the set are beyond the scope of this document. Having  the ILM map a label to a set containing more than one NHLFE may be useful if, e.g., it is desired to do load balancing over multiple equal-cost paths.

FEC-to-NHLFE Map (FTN)   
	The "FEC-to-NHLFE" (FTN) maps each FEC to a set of NHLFEs.  It is used when forwarding packets that arrive unlabeled, but which are to be labeled before being forwarded.
	If the FTN maps a particular label to a set of NHLFEs that contains more than one element, exactly one element of the set must be chosen before the packet is forwarded.  Having the FTN map a label to a set containing more than one NHLFE may be useful if, e.g., it is desired to do load balancing over multiple   equal-cost paths.


Label Swapping   
	Label swapping is the use of the following procedures to forward a packet. In order to forward a labeled packet, a LSR examines the label at the top of the label stack.  It uses the ILM to map this label to an NHLFE.  Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet’s label stack.  It then encodes the new label stack into the packet, and forwards the result.
	In order to forward an unlabeled packet, a LSR analyzes the network layer header, to determine the packet’s FEC. It then uses the FTN to map this to an NHLFE. Using the information in the NHLFE, it determines where to forward the packet, and performs an operation on the packet’s label stack.  (Popping the label stack would, of course,  be illegal in this case.)  It then encodes the new label stack into the packet, and forwards the result.


LSP Control: Ordered versus Independent   
	Some FECs correspond to address prefixes which are distributed via a dynamic routing algorithm. The setup of the LSPs for these FECs can be done in one of two ways: Independent LSP Control or Ordered LSP Control.
	In Independent LSP Control, each LSR, upon noting that it recognizes a particular FEC, makes an independent decision to bind a label to that FEC and to distribute that binding to its label distribution peers. This corresponds to the way that conventional IP datagram routing works; each node makes an independent decision as to how to treat each packet, and relies on the routing algorithm to converge rapidly so as to ensure that each datagram is correctly delivered.
	In Ordered LSP Control, an LSR only binds a label to a particular FEC if it is the egress LSR for that FEC, or if it has already received a label binding for that FEC from its next hop for that FEC.
	If one wants to ensure that traffic in a particular FEC follows a   path with some specified set of properties (e.g., that the traffic   does not traverse any node twice, that a specified amount of   resources are available to the traffic, that the traffic follows an   explicitly specified path, etc.)  ordered control must be used.  With   independent control, some LSRs may begin label switching a traffic in   the FEC before the LSP is completely set up, and thus some traffic in   the FEC may follow a path which does not have the specified set of properties.  Ordered control also needs to be used if the recognition of the FEC is a consequence of the setting up of the corresponding LSP.

Aggregation:
	One way of partitioning traffic into FECs is to create a separate FEC   for each address prefix which appears in the routing table.  However,   within a particular MPLS domain, this may result in a set of FECs   such that all traffic in all those FECs follows the same route.  For   example, a set of distinct address prefixes might all have the same   egress node, and label swapping might be used only to get the the   traffic to the egress node.  In this case, within the MPLS domain,   the union of those FECs is itself a FEC.  This creates a choice:   should a distinct label be bound to each component FEC, or should a single label be bound to the union, and that label applied to all traffic in the union?


In order to use MPLS for the forwarding of packets according to the   hop-by-hop route corresponding to any address prefix, each LSR MUST:      
	1. bind one or more labels to each address prefix that appears in  its routing table;      

 The Implicit NULL label is a label with special semantics which an   LSR can bind to an address prefix.  If LSR Ru, by consulting its ILM,   sees that labeled packet P must be forwarded next to Rd, but that Rd   has distributed a binding of Implicit NULL to the corresponding   address prefix, then instead of replacing the value of the label on   top of the label stack, Ru pops the label stack, and then forwards   the resulting packet to Rd.
	2. for each such address prefix X, use a label distribution protocol to distribute the binding of a label to X to each of its label distribution peers for X.
