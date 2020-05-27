1. Basic:
    bcm debug bcmapi l3 +
    bcm debug bcmapi +
    bcm memwatch vlan read
    bcm memwatch vlan write
    bcm memwatch vlan off
    bcm lls
    bcm cos bandwith_show
    bcm cos discard_show
    bcm cos port show PBM=all
    bcm d sw cosq
    bcm cos show
    bcm show pmap
    bcm dump sw port
    bcm dump sw vlan
    bcm dump sw ifp
    bcm fp show 
    bcm show features
    bcm show st
    bcm show params
    bcm hash
    bcm HashDestination
    bcm pr
    bcm pp
    bcm rate
    bcm knetctrl netif show 
    bcm knetctrl filter show
    bcm egress show
    bcm tx 1 pbm=none
    bcm port xe1 encap=higi
    bcm vlan show 
    bcm vlan 1000
    bcm vlan port xe1
    bcm show portmap
	bcm soc 查看chip信息及目前状态
	bcm ser log 可以查看chip的ser log 信息
2.  TechSupport
    The "techsupport" utility internally maintains the following sets of data on a per chip per feature basis:

        diag shell commands (diag).
        Register list (reg).
        Memory table list (mem).
        Software maintained feature specific state (sw), where applicable. 
    The user can collect the data described above using the techsupport command. This command has the intelligence to identify the underlying chip on which it is being run and collects the chip specific data.

    The following is the format of techsupport command:

        techsupport basic - collects basic config/setup information from switch
        techsupport [diag] [reg] [mem] [list] [sw] [verbose]
    Command "techsupport basic" collects basic config/setup information from switch. It executes the following commands and is supported on all XGS chips.

        attach
        version
        config show
        show counters
        linkscan
        ps
        lls
        hsp
        soc
        fp show
        show pmap
        phy info
        show params
        show features
        dump soc diff
        dump pcic
   
    The command "techsupport " dumps all four sets of feature specific information. By default (i.e if verbose option is not mentioned), memory and register dumps are changes from defaults. If any other options (diag, reg, mem, sw) are specified, then only that specific information is dumped. 

    For example,
        Command "techsupport l3mc diag" executes the Layer 3 multicast related diag shell commands.
        Command "techsupport l3mc diag mem" executes the Layer3 multicast related diag shell commands and dumps the Layer3 multicast related memory tables. 
        Command "techsupport l3mc"  collects all the above 4 sets of data.
        Command "techsupport l3mc list" lists the diag shell commands, memory and register names.
    As of SDK 6.5.7, the techsupport command supports the following features. These features are supported for Trident2, Trident2+ and Tomahawk chipsets.
        l3uc - collects L3 unicast related debug information
        l3mc - collects L3 multicast related debug information
        mpls - collects mpls related debug information
        mmu - collects mmu related debug information
        niv - collects niv related debug information
        riot - collects riot related debug information
        oam - collects oam related debug information
        vxlan - collects vxlan related debug information
        vlan - collects vlan related debug information
        cos - collects cos related debug information
        load-balance - collects RTAG7,DLB,ECMP and trunk related debug information   
        efp - collects efp related debug information
        ifp - collects ifp related debug information
        vfp - collects vfp related debug information


3. Rule 相关的shell command
        fp show
        fp group get ID // 获取一个group支持的qualifier
        fp show entry ID
        fp entry enable ID
        fp entry disable ID
        fp entry remove ID
        fp entry destroy ID
        fp entry copy SID DID
        fp aset show 
        fp show brief
        fp show group ID
        fp list actions
        fp list qualifiers
        fp stat get statid=10 type=packets

4. How to debug packet drop on XGS
            RDGBC/TDGBC
        cstat ls
        cstat info rx
        cstat info tx
        cstat PORT

		通过cstat PORT set rx 0 TRIGGGER 来使能一个flag
		通过cstat PORT 来查看是否是该原因导致的丢包。
		通过cstat PORT set rx 0 TRIGGGER 来关闭一个flag
		通过cstat info 来查看所有可配置的flag。
		注意这里的set rx 0的0是counter,这个值可以从bcm cstat  get PORT rx来获取。同时改命令可以获取某一counter对应的trigger.这个数字其实就是rdbgc的索引。
		通过cstat PORT set tx 6 TRIGGER 来使能一个tx的trigger。
		注意rx和tx的counter值范围是不一样的。
		也可以通过cstat set PORT rx all TRIGGER来指定所有的counter。
		也可以通过cstat set PORT tx all TRIGGER来指定所有的tx counter。


5. 根据端口号快速计算Port Bit Map
        pbmp
        pbmp ge0
        pbmp ge
        pbmp xe
        pbmp all

6. 使用shell配置rule的最简单例子
    1. Count packet from ge0
        fp init
        fp qset clear
        fp qset add InPorts
        fp group create 0 0
        fp entry create 0 0
        fp qual 0 InPorts ge0 
        fp policer create PolId=0 mode=srtcm ColorBlind=0 ColorMergeOr=0 cbs=70000 Cir=600000 eir=200000    -----------------> create a meter
        fp policer attach entry=0 level=0 polid=0
        fp stat create group=0 type0=bytes type1=packets       --------------------> do counting on 'Bytes' and 'Packets'. here you can also count for 'GreenBytes' or 'YellowPackets' as type 
        fp stat attach entry=0
        fp entry install 0

    2. Redirect packet from xe0 destip 3.3.3.0/24 to 10003
        fp init
        fp qset clear
        fp qset add inport
        fp qset add DstIp
        fp group create 0 1
        fp entry create 1 2
        fp qual 2 inport xe0 0xffff
        fp qual 2 dstIp 0x03030300 0xffffff00
        fp action add 2 RedirectEgrNextHop 100003
        fp stat create group=1 type0=Packets
        fp stat attach entry=2 StatId=4
        fp entry install 2
    3. Drop packet which is send to cpu.
        fp init
        fp qset clear
        fp qset add stageegress
        fp qset add outport
        fp qset add SrcIp
        fp group create 0 0
        fp entry create 0 0
        fp qual 0 outport cpu0 0xffff
        fp action add 0 drop
        fp entry install 0

7. 遇到过的问题
        1. ING_DEVICE_PORT/PORT table的LPORT_PROFILE_IDX是用来索引LPORT_TAB的。而LPORT_TAB中包含了Port通用的控制寄存器.
        2. ModuleLoopback 特性默认是关闭的，也就是说如果从Higi port收到source modid 与本module_id 相同的报文时，默认时丢弃的。如果想要控制不丢弃，则可以使用LPORT_TAB的ALLOW_SRC_MOD field来控制。
        3. PORT_TABLE都对应一个switch_port_control的选项。
        4. 怎么让brodcom chip转发BPDU： 将L2_USER_ENTRY的CPU域设为0.将LPORT_TAB的DROP_BPDU域也设为0.

8. 通过sdkdebug 监控端口的状态变化
        1. debug bcm link
        2. 打开link的debug以后，在link状态改变的时候，就可以看到相应的sdkdebug消息了。

9. 通过shell获取gport
        1. gport port=mp(modid, port)
        3. gport port=local(modid, port)
        3. gport port=trunk(modid, port)
