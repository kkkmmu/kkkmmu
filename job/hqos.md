1. 普通的QOS为不同提供不同服务。
2. HQOS为不同的用户的不同业务提供不同的服务。
3. CAR: Commited Access Rate.

4. 传统流量管理是基于端口带宽进行调度的，这样产生的结果就是对用户不敏感，只对服务等级敏感，适合网络侧，但不适合业务接入侧。
   传统流量管理很难做到同时对多个用户的多个业务进行控制
   为解决以上的问题，提供更好的QoS能力，迫切需要一种既能控制用户的流量，又能同时对用户业务的优先级进行调度的QoS技术。HQoS采用多级调度的方式，采用全新的硬件设计，使设备具有内部资源的控制策略，既能够为高级用户提供质量保证，又能够从整体上节约网络构造成本。
   普通QoS调度算法是基于端口的，只对业务等级敏感，对用户不敏感，会产生混乱的抢占结果，HQoS通过多级调度来解决此问题。
   HQoS解决多用户的多业务带宽保证原理: 接口就像一个大的水管，每个用户从大管子中分一个小管子，大管子就能对小管子进行流量管理，而每个用户又对其自己的小管子再细分出不同的流，对不同的流进行不同的处理。
   这种原理在设备上通过分级调度来实现，第一级完成对用户带宽的保证，第二级完成对每个用户各业务的带宽保证。

5. 目前IP组网典型方案是分为接入层、汇聚层、骨干层三层结构。这种宽带网络的分层结构使网络层次清晰。接入层为不同用户提供各种接入手段，汇聚层对接入层业务流汇聚，城域网核心层保证快速转发，网络拓扑有下列显著特征：
	(1 网络拓扑结构通常是树形拓扑结构，即使有些场合，物理上的拓扑结构是环形结构，但链路层还是树形拓扑结构；
	(2)各种业务的转发路径基本固定，不会动态改变；
	(3)业务路由器是树形拓扑结构的根节点。
6. 为了更好的对流量进行管理，HQoS有着完善的流量统计功能。通过流量统计功能，使用者可以看到各种业务的带宽使用情况，并通过分析流量，合理的划分各业务的带宽分配。

7. 在分层调度模型中，共包括端口层、用户组层、用户层、用户业务层共四个层次。HQoS处理引擎完成拥塞避免、多级调度、流量限速、流量统计等功能。
8. 调度器的层次主要影响到HQoS是否可以应用到拓扑结构更复杂的场合，层次越多，可以应用的场合越多，（一般分四个层次描述，分别称为用户业务队列->用户->用户组->端口），体现为支持业务类型的丰富性和用户个数。
	队列主要用于实现丰富的QoS功能，例如流量限速、拥塞避免、流量整形和队列调度等，用户业务队列数目决定了同时最多可以有多少流参与HQoS处理，队列的最大深度则决定了最多可以缓存的报文数目。
