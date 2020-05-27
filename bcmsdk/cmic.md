CMIC: CPU Management Interface Controller
	this is a gateway to the host CPU.


Packet DMA operations are controlled by DMA Control Blocks(DCBs).
The DMA controller uses chains of DMA Control Blocks(DCBs). DCBs are used by the Broadcom SDK software for 
all packet transmit and receive transactions. The DCB contains all of the inforamtion required for q packet
data transfer. DCBs can be changed together to form a physically contiguous array, and/or linked to permit
gathering of descriptor arrays from multiple independent locations in memory.
There are different DCBs for different devices and for the same device there are different DCBs for packet transmit and 
packet receive operations.

DCB is also called DMA descriptor sometime.

broadcom设备发送和接收报文分别使用不同DMA Channel. 对SWITCH设备来说，一般RX用的CHANNEL会比TX用的CHANNEL要少，因为一般入方向上流量会比较大，
经过CPU过滤，出方向流量会相应的变小。比如说如果一共有四个Channel的话，可以用一个CHANNEL来TX，用另外三个CHANNEL来RX。

对于用作RX的CHANNEL来说，需要配置与CPU队列的映射关系(将特定类型，特定优先级的报文映射进特定的DMA Channel。
一般来说CPU Port拥有比普通Port更多的队列。
没有映射到RX DMA Channel的CPU Queue的报文，CPU不会处理。 即，如果我们将流量映射进了某个CPU队列，但是该队列没有映射到DMA RX Channel的话，这部分报文会被丢弃掉。

bcmsdk中的寄存器带0和1后缀的是同一个配置LSB和MSB.
