1. Each network application may come with a unique set of network QoS requirements for packet latency, jitter, drop sensitivity and bandwidth. To ensure the system Quality of Service(QoS) requirements, several processes-including determining Per Hop Behaviopr(PHB), applying PHB, and marking outgoing packets -- are supported in the StrataXGS systems.

2. PHB (Per Hob Behavior)
    The PHB is composed of the traffic class and color. The PHB is used internally for switching decisions, such as selectiong the appropriate egress queue. As such, it may influence the packet latency or rate of the packet drop.
    The traffic class, also called internal priority, is usually related to the level of importance of delivery or time sensitivity of the traffic. 
    The traffic color also referred internally as CNG, usually represents temporal network congestion conditions that determine the packet drop sensitivity.

3. Outgoint Packet Marking
    The outgoing packet marking refers to modifications to outgoing packet headers that are intended to influence the PHB of the packet elsewhere in the network. The Outgoint Packet Marking does not necessarily impact the internal switching decisions, but impacts the form of the outgoing packet.

4. Determining PHB
    The goal of determining PHB is to classify traffic based on service requirements(low latency, jitter sensitive, bandwidth guarantees etc.) so that it can be handled appropriately within the node. The PHB is often determined by the QoS markings in the packet such as PCP/CFI and DSCP. For traffic entering the node without valid QoS markings, deep packet inspection might be required to assign PHB. This is typically done through ACLs(ICAP/VCAP).

5. Applying PHB
    When PHB is determined for an incoming packet, there are several mechanisms which use PHB, such as buffer management, queueing, sheduling, and shapping.

6. Marking outgoint Packets
    The goal of marking outgoint packet is to encode the PHB derived at this node into the packet so that a downstream node can use it to derive its PHB.


