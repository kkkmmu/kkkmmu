There are two primary points in the device's pipeline where packets can accumulate during times of congestion.

1. The first point is just before packets are processed by the ingress pipeline. This is a small buffer referred to as the Oversubscription Bufer Moudle(OBM). Its purpose is primary to absorb packets from each of the port macros at line rate before they are admitted to the ingress pipeline. There is one instance of the OBM per port macro, so for a system that is using all of the SerDes cores, there are 32 OBM instances active.

2. The other primary point of buffering is located between the ingress and egress pipelines. This is a large, monolithic, highly configurable packet buffer. This is typically referred to as Memory Management Unit(MMU).

3. The Memory Management Unit (MMU) is the central block that connects the two ingress pipelines to the two egress pipelines. The MMU's primary purpose is to buffer packets during times of congestion and schedule the packets out to the appropriate port(s) in a user configurable and predictable manner. This is accomplished by receiving chunks of packet data (cells) from either ingress pipeline, linking the cells together to form a full packet, linking the packet to the appropriate queue(s), and then scheduling the packet out to the desired port(s). 

4. The MMU has a finite number of resources available to buffer packets. To manage these resources fairly and efficiently, the MMU provides many features to assit the user in meeting their desired buffering and fairness goals. These features can be categorized into the following general areas:
    a. Resource Management Mechanisms:
        1. Ingress Resource Accounting
        2. Egress Resource Accounting
        3. Dynamic Thresholding
        4. Hyseresis
        5. Flow Control Generation
        6. Weighted Random Early Detection
        7. Aging
    b. Packet Buffer Visibility
        1. Counters
        2. Buffer Statistics Tracking
        3. Transient Capture Buffer
        4. Packetized Statistics

5. Resource Types
    The MMU  has multiple finite resources that must be managed. These resource include:
        1. A Pool of cells used to buffer packet data.
        2. A pool of pointers used by the replication engine during spatial packet replication.
        3. A pool of pointers used by multicast packets.

6. The Memory Management Unit (MMU) tracks all finite resources from multiple perspectives to support multiple buffering goals. These goal typically include thing like asserting flow control to avoid packet loss, tail dropping when a resource becomes congested, or randomly dropping to allow statful protocols to degrade gracefully.
    The majority of the functionality within the MMU depends on using counters to track the usuage of the resources. When a cell arrives from the ingress pipeline, the state of various counters is checked to determine if the cell should be accepted. The decision to accpet the cell is based on user configurable thresholds for the various accouting sites.
    There are multiple blocks within the MMU that are responsible for resouces accouting and thresholding. These blocks are partitioned and named based on their location and the type of resource they account for. 

7. Ingress Accouting is primarily used to trigger flow control and ensure fair buffer usuage from and ingress port perpsective. The THDI block is the only ingress accouting resource in the device. It's purpose is to accout for cell ussage from the ingress port perspective.
   All of the other finite resources that are used by mulicast traffic are not tracked from the ingress port perspective. This means that it is not possible to trigger flow control if these resources become congested. As such, the device only guarantees lossless behavior for known unicast packets.
   For lossless operation, the ingerss accouting resources are typically configured such that they trigger flow control, while the egress accouting resources are configured to the maximum value allowed by the width of the field in the register or memory. This ensures that an ingress threshold is always hit before an egress threshold. Conversely, for lossy operation, the ingress accouting resources are typically configured to be the maxmimum allowable value and the egress thresholds are configured to a value that causes discards.

8. Cell Accounting
    Every cell that arrives at the MMU is subject to ingress accounting and thresholding checks. If this check is above a configured levle, the THDI block can signal to the ingress port that flow control should be generated. If this check is above a (typically) higer configured level, the cell(and associated packet) is discard.
    The THDI block is primaryily color-blind, in that all cells are treated equally regardless of the green/yellow/red color that was assigned in the ingress pipeline. The exception to this is if the Ingress Content Aware Processor(ICAP) performs a Shared Pool ID override. In that case, the cell is only subjected to resource checks in the specified shared pool, and all port level and priority group-level checks are ingnored. Addtionally, the ICAP can also assign a different Service Pool Acceptance Profile(SPAP), which cause the THDI block to use yellow or read thresholds assigned to the ingress shared pool.

9. Egress Accounting is primarily used to ensure fair buffer usuage from an egress port perspective by tail dropping. 
    For knwon unicast packets, only cell usuage is accouted for. This is done in the THDU block. 
    For multicast packets, the THDR_DB block is used to account for cell usuage while the packet is being replicated by the Replication and Queueing Engine(RQE). 
    The THDR_QE block is used to account for replication engine queue entry usage while the packet is being replicated by the RQE. The 
    THDM_DB block is used to account for cell usuage after the RQE has replicated the packet to the final destination queue(s). 
    Finally the THDM_MCQE block is used to account for multicast queue entry usage after the RQE has replicated the packet to the final destination queue(s).

10. Unicast Cell Accounting(THDU)
    Every known unicast cell that arrives at the MMU is subject to egress accouting and thresholding checks. If this check is above a configured level, the THDU block discards the cell(and associated packet).
    The THDU block is color awaer and has the ability to have separate thresholds depending on the green/yellow/red color that was assigned in the ingress pipeline.

11. Relication Engine Cell Accouting (THDR_DB) 
    Every multicast packet that arrives at the MMU is processed by the Replication and Queueing Engine(RQE). While the packet is assigend to the RQE for replication, the cells used by the packet are accounted for by the THDR_DB block. If the RQE becomes too congested, the THDR_DB block discards the cell(and associated packet). If this occurs, the packet is not replicated to any ports.
    The THDR_DB block is color aware and has the ability to have separate thresholds depending on the green/yello/red color that was assigned in the ingress pipeline.

12. Replication Engine Qeuue Entry Accouting(THDR_QE)
    Every multicast packet that arrives at the MMU is processed by the Replication and Queue Engine(RQE). While the packet is assigned to the RQE for replication, a RQE Queue Entry is used by the packet. This is accounted for by the THDR_QE block. If the RQE becomes too congested, the THDR_QE block discards the pakcet. If this occurs, the packet is not replicated to any ports.
    The THDR_QE block is color aware and has the ablility to have separate thresholds depending on the green/yellow/red color that was assigned in the ingress pipeline.

13. Multicast Cell Accouting (THDM_DB)
    Each time that the Replication and Queue Engine(RQE) performs a replication, the replicated cells are accounted for in the THDM_DB block. Unlike the previous cell-based checks, all the cells associated with the packet have already been stored in the MMU. As such, when the RQE performs the replication it provides the THDM_DB block with the total number of cells associated with the packet. The THDM_DB block then uses this to check available resources for the replicated copy. If this check fails, the packet is not replicated to that particular port.
    The THDM_DB block is color aware and has the ability to have separate thresholds depending on the green/yellow/red color that was assigned in the ingress pipeline.

14. Multicast Queue Entry Accounting(THDM_MCQE)
    Each time that the Replication and Queueing Engine(RQE) performs a replication, the replicated packet is accounted for in the THDM_MCQE block. If the egress port becomes too congested, the THDM_MCQE block does not accept the packet. If this occurs, the packet is not replicated to that particular port.
    The THDM_MCQE block is color aware and has the ability to have separate thresholds depending on the green/yellow/red color that was assigned in the ingress pipeline.
