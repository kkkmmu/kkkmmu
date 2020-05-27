Logical Interface:
	Logical interfaces (LIFs) emulate L2 tunnels (for example, AC, PWE, and VXLAN) or L3 tunnels (for example, MPLS/IP-Tunnels) and have incoming and/or outgoing interfaces. A logical tunnel ID is known as the Global-LIF-ID. For the incoming interface, the ID is the In-Global-LIF, and for the outgoint interface the ID is the Out-Glocal-LIF.
	The Global-LIF-ID is used as the software object handling ID as well as the hardware-representation over the system headers and the different forwarding/learning/ACL databases. The software object handling ID is used as an input to configure the BCM LIF APIs and to wrap with GPORT entity. For the hardware-representer, the object-ID is used to point to L2/L3 tunnels from the forwarding tables by the ACL databases, multicast APIs. 

	如何来理解LIF？
	对于交换设备来说，我们有二层端口和三层接口。我们需要重新审视着两个概念。
	二层端口本身隐藏了入方向收ETH报文出方向发ETH报文的概念，由于普通的二层转发本省不会对报文内容进行编辑，因此一般我们不会注意这一概念。
	三层接口本身隐藏了入方向收目的MAC为本地的报文，出方向上发源MAC为本地目的MAC为下一跳MAC的IP报文报文的概念。这一概念在软件上可以抽象为同一Object的收发callback函数，但是在硬件上是不可能的，我们必须在硬件上创建该接口入方向的实例用以解析报文，同时创建该接口在出方向上的实例用以封装报文。(当然这里的描述是高层次抽象，具体到实现上我们看到的入方向和出方向实例可能是指一个标识符，而相应的功能由其他部分实现）.
	交换设备的物理端口可能因为为业务需要而支持多种协议（ETH/IPoEth/MPLSoEth/VXLAN/PWE/GRE等等），此时就需要通过硬件同时来实现对这些协议报文的封装和解封装。为此我们就需要在入方向上能够解析该协议的封装，在出方向上又能进行该协议的封装。不同的协议其对应的封装和解封装方法也是不同的。而这种封装和解封装的实现在逻辑上就是该协议对应的接口在入方向和出方向上的功能。也就是说不同的协议对应一种不同的接口类型。（ETH对应简单的二层端口，IPoETH对应简单的三层接口，MPLSoETH对应MPLS Tunnel接口，VXLANoETH对应VXLAN Tunnel接口）。LIF就是软硬件上对于这些逻辑接口的表示。
