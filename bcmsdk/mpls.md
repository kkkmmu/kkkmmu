1. MPLS features can be classified into the following five groups.
    a. MPLS-Only. Basic MPLS label operations: MPLS-Only functions include label swapping, label pop and label push.
    b. MPLS PWE(VPWS). MPLS PWE provides point-to-point Ethernet transport services to end customers. It involves encapsulation of Ethernet frames with MPLS headers and decapsulation of MPLS headers. The customer packets are mapped to MPLS PW based on port or port + VID. Forwarding of Ethernet frames, after MPLS decapsulation, is based on the MPLS lookup. MPLS PWE is also called Virtual Private Wire Service (VPWS).
    c. VPLS. VPLS provides point-to-multipoint Ethernet transport services to end customers. This involves encapsulation of Ethernet frames with MPLS headers and decapsulation of MPLS headers. Forwarding of Ethernet frames after MPLS decapsulation is based on MAC DA lookup of customer Ethernet frame. It relies on MPLS PWE as part of the operation.
    d. L3 MPLS TE. L3 MPLS TE provides point-to-multipoint IP transport services to proiver edge routers. This involved encapsulation of IPv4 and IPv6 with MPLS headers and decapsulation of MPLS headers. Forwarding of IPv4 and IPv6 packets after MPLS decapsulation is based on DIP lookup of the IPv4 or IPv6 packets.
    e. L3 MPLS VPN. L3 MPLS provides point-to-multipoint IP transport services to end customers. This involves encapsulation of IPv4 and IPv6 with MPLS headers and decapsulation of MPLS headers. Forwarding of IPv4 and IPv6 packets after MPLS decapsulation is based on DIP lookup using virtual routing table. L3 MPLS VPN supports virtual routing tables for invidual VPN customers.

2. A core MPLS switch only requires MPLS-Only features. Any edge MPLS switch that originates and terminates MPLS Label Switched PSpaths (LSPs) may need at least one of the three edge functions: L2 MPLS PWE, L3 MPLS TE/VPN, or VPLS.

3. MPLS works by prefixing packets with MPLS shim header, containing one or more 'label' forming a lable stack. For Ethernet, PPP, FDDI, and other technologies, the shim header is located between the link layer and network layer headers. Each label stack entry contains four fields:
    a. A 20 bit label value.
    b. A 3 bit field for Quality of Service(QoS) prority (EXP - experimental).
    c. A 1 bit bottom of stack (BOS) flag.
    d. An 8 bit time to live (TTL) field.

4. The MPLS-labeld packets are switched based on Label Lookup. Routers that perform routing based only on the label are called Label Switch Rotuers (LSR). They function only as transit routers and are also referred to as Provider(P) routers (RFC2547). In the specific context of a MPLS-based Virtual Private Network(VPN), LSRs that function as ingress and/or egress routers to the VPN are often called PE (Provider Edge) routers or Label Edge Routers (LER), which, respectively, push an MPLS label onto the incoming packet and popi it off the outgoing packet.

5. Tunnels are constructed between the different PEs using targeted Label Switch Path(LSP) and the P routers are concerned only with transporting the traffic from PE to PE without getting invovled with the different services offered at the edge. The LSP tunnels, tunnle PseudoWires from one end of MPLS cloud to the other with opportunity of aggregating multiple PWs within a single LSP tunnel. The LSP tunnels themselves can be constructed manually, or via MPLS signalling using Label Distribution Protocol (LDP) or Resource reservation Protocol (RSVP) traffic Engineering(TE).

6. In general, LSPs can follow conventionally routed shorted paths or explicit paths that could be different from the shortest paths. A LSP tunnel is called Packet Switched Network (PSN) tunnel.

7. L2 MPLS:
    a.Virtual Private Wire Services (VPWS)
    b.Virtual Private LAN Services (VPLS)
    These technologies form an Ethernet PseudoWirez(PW) witch allows Ethernet Packet Date Units (PDUs) to be carried over a PSN (Packet Switched Network) like MPLS. This allows service provider or enterprise networks to leverage an existing MPLS network to offer Ethernet services.
    In L2 MPLS, customer Ethernet frames are encapsulated in two MPLS labels: Tunnel Lable and VC label. The VC label represents the PW itself, whereas the tunnel lable represents the MPLS tunnel that provides a tunnel to transport the PW. MPLS tunnel lables are used for aggregation and transportation of multiple PWs so that core routers inside the MPLS cloud see fewer number of lables.

    The main difference between the two L2 MPLS technologies is the VPWS is a point to point service aka E-line service, and VPLS is a point-to-multipoint service aka E-LAN service. Also the endpoints of VPWS PW are the external connected ports, whereas the VPLS endpoint are the switches themselves. This means that:
    a. VPWS packets are all unicast in nature, both at the start of the PW and at the end of the PW. VPLS packets may need to be multicast both across ports and within ports of the system, similar to IPMC replication, to reach several "equivalent" network destinations.
    b. VPWS pseudo-wires are completely "pre-configured", but VPLS requires L2 lookups, L2 learning, and L2 broadcast to support learning of customer MAC addresses.
