MPLS Forwarding Overview
     A label switch router(LSR)is a router that supports MPLS.It is capable of understanding MPLS labels and of receiving and transmitting a labeled packet on a data link.An LSR can do the three operations:pop,push,or swap.There are three kindsof LSRs in an MPLS network:
	     • Ingress LSRs    —Ingress LSRs receive a packet that is not labeled yet,insert a label(stack) in front of the packet,and send it on a datalink.
		 • Egress  LSRs    —EgressLSRs receive labeled packets, remove the label(s),and send them on a datalink.Ingress and egress LSRs are edge LSRs.
		 • IntermediateLSRs—Intermediate LSRs receive an incoming labeled packet, perform an operation on it,switch the packet,and send the packet on the correct data link.
	The Ingress LSR(LER) adds the label header based on its label mapping database. Down stream LSRs swap the label based on their label mapping table.The last LSR(LER)in the LSP will remove the label and forward the packet as IP.

FTN(FEC-to-NHLFE)
	The FTN maps each FEC to a set of NHLFEs. It is usedwhen forwarding packets that arrive unlabeled, but which are to be labeled before being forwarded.

在MPLS体系中，数据转发涉及到3张表的维护：NHLFE、FTN和ILM。下面分别进行说明。
	1. NHLFE(Next Hop Label Forwarding Entry)：下一跳标签转发单元。这张表维护了MPLS数据包下一跳该怎样走的相关信息，在标签交换中，报文的下一跳总是取自于该表。NHLFE描述了对标签的具体操作、报文的出标签、报文出接口等相关信息。每个LSR上都要维护NHLFE，不管是LER节点还是中间的P节点。
	2. FTN(FEC to NHLFE)：FEC到NHLFE映射。这张表只在LSP的入节点，即LER上维护，用于完成对相关FEC到NHLFE的映射，对不带标签的报文打上相关标签。当一个不带标签的报文进来后，首先分析其IP头或者二层报文头，来决定其FEC（在IP报文中，FEC通常为IP地址，在二层报文中通常为vlan-id或者vc-id）。然后利用FTN将FEC映射到NHLFE中，用来进行报文的出标签、标签操作和决定往何处发送该报文。
	3. ILM(Incoming Loabel Map)：入标签映射。这张表只在LSP的中间P节点上维护，LSP的入节点不维护该表。当LSR收到带标签的报文后，要利用ILM将报文映射到NHLFE，来决定该标签报文的标签操作和下一跳地址。
	4. 在这三张表结构的中间还有一个结构XC(cross-connect)，用来将FTN和ILM链接到相关NHLFE中。在每个FTN和ILM节点中都要有一个XC指针指向对应的NHLFE，用来快速查找。

FEC: 
	MPLS将具有相同特征的报文归为一类，称为转发等价类FEC（Forwarding Equivalence Class）。属于相同FEC的报文在转发过程中被LSR以相同方式处理。
	FEC可以根据源地址、目的地址、源端口、目的端口、VPN等要素进行划分。例如，在传统的采用最长匹配算法的IP转发中，到同一条路由的所有报文就是一个转发等价类。

MPLS详细转发过程
基本概念

在MPLS详细转发过程中涉及的相关概念如下：
	Tunnel ID
		为了给使用隧道的上层应用（如VPN、路由管理）提供统一的接口，系统自动为隧道分配了一个ID，也称为Tunnel ID。该Tunnel ID的长度为32比特，只是本地有效。
	NHLFE
		下一跳标签转发表项NHLFE（Next Hop Label Forwarding Entry）用于指导MPLS报文的转发。
		NHLFE包括：Tunnel ID、出接口、下一跳、出标签、标签操作类型等信息。
		FEC到一组NHLFE的映射称为FTN（FEC-to-NHLFE）。通过查看FIB表中Tunnel ID值不为0x0的表项，能够获得FTN的详细信息。FTN只在Ingress存在。
	ILM
		入标签到一组下一跳标签转发表项的映射称为入标签映射ILM（Incoming Label Map）。
		ILM包括：Tunnel ID、入标签、入接口、标签操作类型等信息。
		ILM在Transit节点的作用是将标签和NHLFE绑定。通过标签索引ILM表，就相当于使用目的IP地址查询FIB，能够得到所有的标签转发信息。

当IP报文进入MPLS域时，首先查看FIB表，检查目的IP地址对应的Tunnel ID值是否为0x0。
	如果Tunnel ID值为0x0，则进入正常的IP转发流程。
	如果Tunnel ID值不为0x0，则进入MPLS转发流程。

在MPLS转发过程中，FIB、ILM和NHLFE表项是通过Tunnel ID关联的。
	Ingress的处理：通过查询FIB表和NHLFE表指导报文的转发。
		查看FIB表，根据目的IP地址找到对应的Tunnel ID。
		根据FIB表的Tunnel ID找到对应的NHLFE表项，将FIB表项和NHLFE表项关联起来。
		查看NHLFE表项，可以得到出接口、下一跳、出标签和标签操作类型。
		在IP报文中压入出标签，同时处理TTL，然后将封装好的MPLS报文发送给下一跳。

 		说明： 转发过程中对TTL的处理，请参见MPLS对TTL的处理。
	 Transit的处理：通过查询ILM表和NHLFE表指导MPLS报文的转发。
	 	根据MPLS的标签值查看对应的ILM表，可以得到Tunnel ID。
		根据ILM表的Tunnel ID找到对应的NHLFE表项。
		查看NHLFE表项，可以得到出接口、下一跳、出标签和标签操作类型。
		MPLS报文的处理方式根据不同的标签值而不同。
			如果标签值>＝16，则用新标签替换MPLS报文中的旧标签，同时处理TTL，然后将替换完标签的MPLS报文发送给下一跳。
			如果标签值为3，则直接弹出标签，同时处理TTL，然后进行IP转发或下一层标签转发。
 	 Egress的处理：通过查询ILM表指导MPLS报文的转发或查询路由表指导IP报文转发。
 		如果Egress收到IP报文，则查看路由表，进行IP转发。
 		如果Egress收到MPLS报文，则查看ILM表获得标签操作类型，同时处理TTL。
 		如果标签中的栈底标识S=1，表明该标签是栈底标签，直接进行IP转发。
 		如果标签中的栈底标识S=0，表明还有下一层标签，继续进行下一层标签转发。


Mapping packets to NHLFE
	    ILM (Incoming Label Map) Maps a labeled packet to a set of NHLFE
			Enables forwarding of labeled packet
		FTN (FEC-to-NHLFE)  Maps an FEC to a set of NHLFE
			Enables forwarding of unlabeled packets that are to be labeled before being forwarded
		Why mapping to a set of NHLFE? It might be useful for specific applicationsE.g., load balancing among alternate links/pathsCriteria for choice within set is not specified

