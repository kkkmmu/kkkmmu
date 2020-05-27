IP Tunneling:
	Tunneling is used to provide a direct connection between two nodes across a network. Tunneling is done simply by encapsulation a packet supporting one protocol within another packet with the same or different protocol.
	Encapsulation in IP tunneling is achieved by adding an outer IP header before the original IP header of an IP packet. The destination and source IP addresses of the outer header spcify the endpoints or logical link of the tunnel. Hence a tunnel has two endpoints: entry point and exit point. 
		Entry point encapsulates the packet by adding outer IP header and decrements the TTL field of the inner header. Then it forwards and routes the tunneled packet based on the IP addresses of the outer header.
	    Exit point decapsulates the packet by removing the outer header and it forwards and routes the packet based on the IP addresses of the inner header to the recipient.	

	IP Tunneling provides connectivity across two routers or two hosts or a router and a host. Hence, each end point can be a router or a host with the following properties:
		1. Entry and exit points are two routers - the router header indicates the logical link and the inner header IP addresses specify local link. The outer source and destination IP addresses are not the same as inner source and destination IP addresses.
		2. Entry point is a outer and exist point is a host - The outer source IP address is not the same as inner source IP address. The outer destination IP address is the same as inner destination IP address.
		3. Entry point is a host and exit point is a router - The outer source IP address is the same as inner source IP address. The outer destination IP address is not the same as inner destination IP address.
		4. Entry point and exit point are two hosts - the outer source and destination IP addresses are the same as inner and destination IP addresses. (Logical link and local link are the same).

IP Node Types:
	There are three different type os IP nodes that support different protocols:
		1. IPv4-only node - A host or router that implements only IPv4. An IPv4 only node does not understand IPv6. Most hosts and routers installed today are IPv4-only nodes.
		2. IPv6-only node - A host or router that implements IPv6 and does not impelment IPv4.
		3. IPv4/IPv6 node or Dual Stack - They are able to transmit both IPv4 and IPv6 packets and thus interact with all IP systems in the network.

GRE tunnel的一点理解:
	配置GRE tunnel的时候，tunnel destination 和tunnel source必须是实际存在的物理接口的IP(loopback也可以).
	配置GRE tunnel的时候，tunnel destination 和tunnel source用来生成报文的外层IP报头.
	配置GRE tunnel的时候，tunnel 的IP地址主要用对数据流量进行路由，将数据流量的路由出接口指定为tunnel接口, 这样在出方向就可以对流量由chip进行GRE封装了.
	从这里可以看出来tunnel分为两端，两端对报文的操作也是不同的。入方向通过将报文的下一跳出接口指定为tunnel interface对报文进行GRE封装。出方向对GRE封装的报文进行解封装，并根据路由表项进行转发。
	对于转发芯片来说就应该能够鉴别合适对一个报文进行封装，如何封装，以及何时对一个报文进行解封装，如何解封装。

Tunnel/LSP 在出方向和入方向上与普通的接口（L3）对报文的处理是有不同的需求的。在出方向上如何封装一个报文，以及在入方向上如何甄别一个报文都是接口的属性。因此才会引入LIF的概念，从本质上来说，LSP/Tunnel Endpoint也是一种接口，只不过它们在入方向和出方向上与普通的三层接口具有不同的业务逻辑。将Tunnel/LSP Endpoint抽象成接口的好处是：可以将该接口作为三层路由下一跳的出接口，而具体如何封装和解封装报文则是该接口本身的属性。这样就简化了三层业务逻辑，同时也更容易扩展。

IPv4 and IPv6 packets are routed between RIFs, they enter the router through an InRIF and exit through an OutRIF

Each RIF Has the following properties:
	1. Routeing Enabled: The RIF enables IPv4 and IPv6 routing for Unicast and/or Multicast for incoming packets.
	2. My-MAC:  A MAC Address value is associated with eth Eth-RIF.
	3. - Incoming Unicast packets whose Ethernet-DA matches the My-MAC are forwarded through the RIF to the router.
	   - Outgoing Unicast packets have their SA set to My-MAC.
	4. Unicast RPF Enable: A flag inidicates whether to perform a Reverse Path Forwarding check on all incoming packets(loose or strict mode).
	5. VRF: Identify a Virtual Route Forwarding database.
	6. Cos-Profile: Indicate the CoS mode to apply to packets received through the Eth-RIF.
	7. InRIF and OutRIF are used for RPF and ICMP redirect tests.

Each Tunnel has the folloing atrributes:
	1. Tunnel(loopback) Interface IP Address(My-IP): This address is compared to incoming packet's DIP to determin if packet is addressed to the tunnel interface(used in termination).
	2. Detination IP Address: this address is used when transmiting packets out of the interface(used for encapsulation).
	3. Tunneling Protocol: for example, GRE encapsulation or IP in IP or VxLAN.
	4. Tunnel CoS Attribute: define how the DSCP is set upon encapsulation.
	5. Tunnle TTL: define how the TTL is set upon encapsulation.

MPLS:
	对于LER来说，使用LDP来构建LSP。每个LSP对应一个Tunnel，每个Tunnel又在逻辑上对应一个Tunnel Interface.
	因此只要将从某一VPN来的流量下一跳出接口指定为该Tunnel就可以在出方向上完成MPLS封装。这就是FTN所做的事情。
	而LSR在收到MPLS封装的流量时，根据Tunnel Terminator的结果再查找ILM来进行流量转发。
	也就是说L3 MPLS VPN的核心依然是 LSP所构建的Tunnel Interface.
