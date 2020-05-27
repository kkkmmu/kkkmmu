Communication between VPN sites and non-VPN sites is   prevented by keeping the routes to the VPN sites out of the default   forwarding table. (注意这里，三层VPN在
PE上并不是将VPN内部流量“转发”到global, global路由本身是用来建立SP网络拓扑以及VPN tunnel的，各个site之间其实是通过LSP来转发流量的，而LSP是在global路由表的基础上构建起来
的。也就是说不存在global到VRF的转发，也不存在VRF的转发，这种观念是概念性的错误。

MPLS L3 VPN VRF内部是L3 转发而backbone是MPLS转发。

A BGP speaker may not use BGP to   send labels to a particular BGP peer unless that peer indicates,   through BGP Capability Advertisement, that it can process Update   messages with the specified SAFI field.

Label mapping information is carried as part of the Network Layer   Reachability Information (NLRI) in the Multiprotocol Extensions   attributes.  The AFI indicates, as usual, the address family of the   associated route.  The fact that the NLRI contains a label is   indicated by using SAFI value 4.

The label(s) specified for a particular route (and associated with   its address prefix) must be assigned by the LSR which is identified   by the value of the Next Hop attribute of the route.

When a BGP speaker redistributes a route, the label(s) assigned to   that route must not be changed (except by omission), unless the   speaker changes the value of the Next Hop attribute of the route.

A BGP speaker can withdraw a previously advertised route (as well as   the binding between this route and a label) by either (a) advertising   a new route (and a label) with the same NLRI as the previously   advertised route, or (b) listing the NLRI of the previously   advertised route in the Withdrawn Routes field of an Update message.   The label information carried (as part of NLRI) in the Withdrawn   Routes field should be set to 0x800000.  (Of course, terminating the   BGP session also withdraws all the previously advertised routes.)

Consider the following LSR topology: A--B--C--D.  Suppose that D   distributes a label L to A.  In this topology, A cannot simply push L   onto a packet’s label stack, and then send the resulting packet to B.   D must be the only LSR that sees L at the top of the stack.  Before A   sends the packet to B, it must push on another label, which was   distributed by B.  B must replace this label with yet another label,   which was distributed by C.  In other words, there must be an LSP   between A and D.  If there is no such LSP, A cannot make use of label   L.  This is true any time labels are distributed between non-adjacent   LSRs, whether that distribution is done by BGP or by some other   method.
