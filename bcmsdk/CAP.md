1. Terminology
    1. Physical Slice/Physical Lookup Table
        ICAP supports multiple parallel lookups. Each parallel lookup requires dedicated resources which constitute a physical slice of ICAP.
        Each physical slice includes a lookup table, a key selector, a key generator, logical table resolution, a policy table, and other resources.
        The search in a physical slice may implemented using TCAMs, or hash tables.
    2. Multi-wide Mode
        Each TCAM has a native width for input keys. Adjacent TCAMs can be connected together to get a multiwide mode with a larger key. A max of three adjacent TCAMs can be connected in the 56870 ICAP design.
        The terms pairing, intra-slice, inter-slice, double wide, and quad wide were used previously to describle multi-wide modes. However, there are different connotations now depending on the chip considered, so those usages are avoided.
    3. Rule Set
        A set of rules and associated actions as defined by a user in an implementation-agnostic manner. Rules within a set are ordered by explicit priority and check in decreasing order of priority. A hit in the rule set matches the highest priority matching rule.
    4. Logical Table/Logical Table Partition/Logical Table Partition Selection
        A rule set corresponds to a logical table in hardware. The logical table may span multiple TCAMs. The portion of a logical table in one physical TCAM is termed a logical table partition. The search of a logical table searches all of its logical partitions. For multiple partition hits, the partition with the highest priority is used.
        When a logcial table requires a multi-wide of operation, each logical partition also spans the multiwide TCAMs. The portion of a logical partion within each TCAM is referred to as a logical table partion section.
    5. Mutulally Exclusive Rule Sets
        Two rule sets are mutually exclusive if they can never apply to the same packet. Examples are rules qualified by EtherType. A rule set define IPv4 rules is mutually exclusive to a rule set defining IPv6 or FCoE rules.
        Mutually exclusive rule sets, each implemented as a logical table, can share the same search resources such as TCAMs if the lookup key is qualified with a rule set id or a logical table id. A search in the physical resource only searches within one logical table. This can change on a per packet basis.
    6. Logical Partition Priority
        Software assigns a priority to each partition within a logical table. Since logical paritions are not visible to the user, the priority assignment is an implmentation internal issue and is managed by software.
    7. Logical Table Action Priority
        Each logical table is assigned an action priority. It is used as the strenght for all actions from the logical table lookup.
2. The Major blocks that make up a single slice of a ContentAware processor:
    1. Protocol-Aware Selector
    2. Lookup Engine
    3. Meter Engine
    4. Counter Engine
    5. Policy Engine
    6. Action Resolution Engine
