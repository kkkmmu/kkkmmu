The VPN Routing Process
The MPLS-VPN Routing process follows these steps:
1. Service Providers provide VPN services from PE routers that communicate directly with CE routers via an Ethernet Link.
2. Each PE router maintains a Routing and Forwarding table (VRF) for each customer. This guarantees isolation, and allows the usage of uncoordinated private addresses. When a packet is received from the CE, the VRF that is mapped to that site is used to determine the routing for the data. If a PE has multiple connections to the same site, a single VRF is mapped to all of those connections.
3. After the PE router learns of the IP prefix, it converts it into a VPN-IPv4 prefix by prepending it with an 8-byte Route Distinguisher (RD). The RD ensures that even if two customers have the same address, two separate routes to that address can be maintained. These VPN-IPv4 addresses are exchanged between the PE routers through MP-BGP.
4. A unique Router ID (usually the loopback address) is used to allocate a label, and enable VPN packet forwarding across the backbone.
5. Based on routing information stored in the VRF table, packets are forwarded to their destination using MPLS. Each PE router allocates a unique label to every route in each VRF (even if they have the same next hop), and propagates these labels, together with 12-byte VPN-IPv4 addresses, through Multi-Protocol BGP.
6. Ingress PE routers prepend a two-level label stack to the VPN packet, which is forwarded across the Provider network. This label stack contains a BGP-specific label from the VRF table (associated with the incoming interface), specifying the BGP next hop and an LDP-specific label from the global FTN table, specifying the IP next hop.
7. The Provider router in the network switches the VPN packet, based on the top label or the LDP-specific label in the stack. This top label is used as the key to lookup in the incoming interface’s Incoming Labels Mapping table (ILM). If there is an outbound label, the label is swapped, and the packet is forwarded to the next hop; if not, the router is the penultimate router, and it pops the LDP-specific label, and forwards the packet with only the BGP-specific label to the egress PE router.
8. The egress PE router pops the BGP-specific label, performs a single label lookup in the outbound interface, and sends the packet to the appropriate CE router.




Configure MPLS Layer-3 VPN
    The MPLS Layer-3 VPN configuration process can be divided into the following tasks
	1. Establish connection between PE routers.
	2. Configure PE1 and PE2 as iBGP neighbors.
	3. Create VRF.
	4. Associate interfaces to VRFs.
	5. Configure VRF Route Destination and Route Targets.
	6. Configure CE neighbor for the VPN.
	7. Verify the MPLS to VPN configuration.
 

9. If all the sites in a VPN are owned by the same enterprise, the VPN   may be thought of as a corporate "intranet".  If the various sites in   a VPN are owned by different enterprises, the VPN may be thought of   as an "extranet".  A site can be in more than one VPN; e.g., in an   intranet and in several extranets.  In general, when we use the term   "VPN" we will not be distinguishing between intranets and extranets.
10. We refer to the owners of the sites as the "customers".  We refer to   the owners/operators of the backbone as the "Service Providers"   (SPs).  The customers obtain "VPN service" from the SPs.
11. Each VPN site must contain one or more Customer Edge (CE) devices.   Each CE device is attached, via some sort of attachment circuit, to   one or more Provider Edge (PE) routers.
12. Routers in the SP’s network that do not attach to CE devices are   known as "P routers".
13. CE devices are logically part of a customer’s VPN.  PE and P routers   are logically part of the SP’s network.
14. If every router in an SP’s backbone had to maintain routing   information for all the VPNs supported by the SP, there would be   severe scalability problems; the number of sites that could be   supported would be limited by the amount of routing information that   could be held in a single router.  It is important therefore that the   routing information about a particular VPN only needs to be present   in the PE routers that attach to that VPN.  In particular, the P   routers do not need to have ANY per-VPN routing information   whatsoever.
15. VRFs: Multiple Forwarding Tables in PEs   Each PE router maintains a number of separate forwarding tables.  One   of the forwarding tables is the "default forwarding table".  The   others are "VPN Routing and Forwarding tables", or "VRFs".
16. Every PE/CE attachment circuit is associated, by configuration, with   one or more VRFs.  An attachment circuit that is associated with a   VRF is known as a "VRF attachment circuit".
17. If an IP packet arrives over an attachment circuit that is not   associated with any VRF, the packet’s destination address is looked   up in the default forwarding table, and the packet is routed   accordingly.  Packets forwarded according to the default forwarding   table include packets from neighboring P or PE routers, as well as   packets from customer-facing attachment circuits that have not been   associated with VRFs.
18. Intuitively, one can think of the default forwarding table as   containing "public routes", and of the VRFs as containing "private   routes".  One can similarly think of VRF attachment circuits as being   "private", and of non-VRF attachment circuits as being "public".   
19. The BGP Multiprotocol Extensions [BGP-MP] allow BGP to carry routes   from multiple "address families".  We introduce the notion of the   "VPN-IPv4 address family".  A VPN-IPv4 address is a 12-byte quantity,   beginning with an 8-byte Route Distinguisher (RD) and ending with a 4-byte IPv4 address. If several VPNs use the same IPv4 address prefix, the PEs translate these into unique VPN-IPv4 address prefixes.  This ensures that if the same address is used in several different VPNs, it is possible for BGP to carry several completely different routes to that address, one for each VPN.
19. Communication between VPN sites and non-VPN sites is   prevented by keeping the routes to the VPN sites out of the default   forwarding table. (注意这里，三层VPN在PE上并不是将VPN内部流量“转发”到global, global路由本身是用来建立SP网络拓扑以及VPN tunnel的，各个site之间其实是通过LSP来转发流量的，而LSP是在global路由表的基础上构建起来的。也就是说不存在global到VRF的转发，也不存在VRF的转发，这种观念是概念性的错误。
20. An RD is simply a number, and it does not contain any inherent   information; it does not identify the origin of the route or the set   of VPNs to which the route is to be distributed.  The purpose of the   RD is solely to allow one to create distinct routes to a common IPv4   address prefix.  Other means are used to determine where to   redistribute the route.
21. A PE needs to be configured such that routes that lead to a   particular CE become associated with a particular RD.  The   configuration may cause all routes leading to the same CE to be   associated with the same RD, or it may cause different routes to be   associated with different RDs, even if they lead to the same CE.
22. a VPN-IPv4 address consists of an 8-byte Route   Distinguisher followed by a 4-byte IPv4 address.  The RDs are encoded   as follows:     - Type Field: 2 bytes     - Value Field: 6 bytes   The interpretation of the Value field depends on the value of the   type field. 
23. If a PE router is attached to a particular VPN (by being attached to   a particular CE in that VPN), it learns some of that VPN’s IP routes   from the attached CE router.  Routes learned from a CE routing peer   over a particular attachment circuit may be installed in the VRF   associated with that attachment circuit.  Exactly which routes are   installed in this manner is determined by the way in which the PE   learns routes from the CE.  In particular, when the PE and CE are   routing protocol peers, this is determined by the decision process of   the routing protocol; this is discussed in Section 7.   These routes are then converted to VPN-IP4 routes, and "exported" to   BGP.  If there is more than one route to a particular VPN-IP4 address   prefix, BGP chooses the "best" one, using the BGP decision process.   That route is then distributed by BGP to the set of other PEs that   need to know about it.  At these other PEs, BGP will again choose the   best route for a particular VPN-IP4 address prefix.  Then the chosen   VPN-IP4 routes are converted back into IP routes, and "imported" into   one or more VRFs.  Whether they are actually installed in the VRFs   depends on the decision process of the routing method used between   the PE and those CEs that are associated with the VRF in question.   Finally, any route installed in a VRF may be distributed to the   associated CE routers.
24. Every VRF is associated with one or more Route Target (RT)   attributes.
	When a VPN-IPv4 route is created (from an IPv4 route that the PE has   learned from a CE) by a PE router, it is associated with one or more Route Target attributes.  These are carried in BGP as attributes of the route.(RT是作为BGP commnity 传输的).
	Any route associated with Route Target T must be distributed to every   PE router that has a VRF associated with Route Target T.  When such a   route is received by a PE router, it is eligible to be installed in   those of the PE’s VRFs that are associated with Route Target T.
	A Route Target attribute can be thought of as identifying a set of   sites.  (Though it would be more precise to think of it as   identifying a set of VRFs.)  Associating a particular Route Target   attribute with a route allows that route to be placed in the VRFs   that are used for routing traffic that is received from the   corresponding sites.
	There is a set of Route Targets that a PE router attaches to a route   received from site S; these may be called the "Export Targets".  And   there is a set of Route Targets that a PE router uses to determine   whether a route received from another PE router could be placed in   the VRF associated with site S; these may be called the "Import   Targets".  The two sets are distinct, and need not be the same.  Note   that a particular VPN-IPv4 route is only eligible for installation in   a particular VRF if there is some Route Target that is both one of   the route’s Route Targets and one of the VRF’s Import Targets.

25. When a BGP speaker has received more than one route to the same VPN-   IPv4 prefix, the BGP rules for route preference are used to choose   which VPN-IPv4 route is installed by BGP.
26. Note that a route can only have one RD, but it can have multiple   Route Targets.  In BGP, scalability is improved if one has a single   route with multiple attributes, as opposed to multiple routes.
27. ow does a PE determine which Route Target attributes to associate   with a given route?  There are a number of different possible ways.   The PE might be configured to associate all routes that lead to a   specified site with a specified Route Target.  Or the PE might be   configured to associate certain routes leading to a specified site   with one Route Target, and certain with another.
28. When a PE router distributes a VPN-IPv4 route via BGP, it uses its   own address as the "BGP next hop".  This address is encoded as a   VPN-IPv4 address with an RD of 0.  ([BGP-MP] requires that the next   hop address be in the same address family as the Network Layer   Reachability Information (NLRI).)  It also assigns and distributes an   MPLS label.  (Essentially, PE routers distribute not VPN-IPv4 routes,   but Labeled VPN-IPv4 routes.  Cf. [MPLS-BGP].)  When the PE processes   a received packet that has this label at the top of the stack, the PE   will pop the stack, and process the packet appropriately.

29. Suppose that a PE has assigned label L to route R, and has   distributed this label mapping via BGP.  If R is an aggregate of a   set of routes in the VRF, the PE will know that packets from the   backbone that arrive with this label must have their destination   addresses looked up in a VRF.  When the PE looks up the label in its   Label Information Base, it learns which VRF must be used.  On the   other hand, if R is not an aggregate, then when the PE looks up the   label, it learns the egress attachment circuit, as well as the   encapsulation header for the packet.  In this case, no lookup in the   VRF is done.

30.    Whether or not each route has a distinct label is an implementation   matter.  There are a number of possible algorithms one could use to   determine whether two routes get assigned the same label:     
			- One may choose to have a single label for an entire VRF, so that       a single label is shared by all the routes from that VRF.  Then       when the egress PE receives a packet with that label, it must       look up the packet’s IP destination address in that VRF (the       packet’s "egress VRF"), in order to determine the packet’s egress       attachment circuit and the corresponding data link encapsulation.     
			- One may choose to have a single label for each attachment       circuit, so that a single label is shared by all the routes with       the same "outgoing attachment circuit".  This enables one to       avoid doing a lookup in the egress VRF, though some sort of       lookup may need to be done in order to determine the data link       encapsulation, e.g., an Address Resolution Protocol (ARP) lookup. 
			- One may choose to have a distinct label for each route.  Then if       a route is potentially reachable over more than one attachment       circuit, the PE/CE routing can switch the preferred path for a       route from one attachment circuit to another, without there being       any need to distribute new a label for that route.

31. The BGP Multiprotocol Extensions [BGP-MP] are used to encode the   NLRI.  If the Address Family Identifier (AFI) field is set to 1, and   the Subsequent Address Family Identifier (SAFI) field is set to 128,   the NLRI is an MPLS-labeled VPN-IPv4 address.  AFI 1 is used since   the network layer protocol associated with the NLRI is still IP.   Note that this VPN architecture does not require the capability to   distribute unlabeled VPN-IPv4 addresses.   In order for two BGP speakers to exchange labeled VPN-IPv4 NLRI, they   must use BGP Capabilities Advertisement to ensure that they both are   capable of properly processing such NLRI.  This is done as specified   in [BGP-MP], by using capability code 1 (multiprotocol BGP), with an   AFI of 1 and an SAFI of 128.   The labeled VPN-IPv4 NLRI itself is encoded as specified in   [MPLS-BGP], where the prefix consists of an 8-byte RD followed by an   IPv4 prefix.

32. By setting up the Import Targets and Export Targets properly, one can   construct different kinds of VPNs.

33. Route Distribution Among VRFs in a Single PE   
	It is possible to distribute routes from one VRF to another, even if   both VRFs are in the same PE, even though in this case one cannot say   that the route has been distributed by BGP.  Nevertheless, the   decision to distribute a particular route from one VRF to another   within a single PE is the same decision that would be made if the   VRFs were on different PEs.  That is, it depends on the Route Target   attribute that is assigned to the route (or would be assigned if the   route were distributed by BGP), and the import target of the second   VRF.